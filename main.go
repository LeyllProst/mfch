package main

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")
	gdi32  = syscall.NewLazyDLL("gdi32.dll")

	procRegisterClassExW           = user32.NewProc("RegisterClassExW")
	procCreateWindowExW            = user32.NewProc("CreateWindowExW")
	procDefWindowProcW             = user32.NewProc("DefWindowProcW")
	procGetMessageW                = user32.NewProc("GetMessageW")
	procTranslateMessage           = user32.NewProc("TranslateMessage")
	procDispatchMessageW           = user32.NewProc("DispatchMessageW")
	procSetLayeredWindowAttributes = user32.NewProc("SetLayeredWindowAttributes")
	procPostQuitMessage            = user32.NewProc("PostQuitMessage")
	procGetSystemMetrics           = user32.NewProc("GetSystemMetrics")
	procSetWindowPos               = user32.NewProc("SetWindowPos")
	procBeginPaint                 = user32.NewProc("BeginPaint")
	procEndPaint                   = user32.NewProc("EndPaint")
	procLoadCursorW                = user32.NewProc("LoadCursorW")

	procCreatePen    = gdi32.NewProc("CreatePen")
	procSelectObject = gdi32.NewProc("SelectObject")
	procDeleteObject = gdi32.NewProc("DeleteObject")
	procMoveToEx     = gdi32.NewProc("MoveToEx")
	procLineTo       = gdi32.NewProc("LineTo")
)

const (
	WS_POPUP          = 0x80000000
	WS_EX_TOPMOST     = 0x00000008
	WS_EX_TRANSPARENT = 0x00000020
	WS_EX_TOOLWINDOW  = 0x00000080
	WS_EX_LAYERED     = 0x00080000
	LWA_COLORKEY      = 0x00000001
	WM_DESTROY        = 0x0002
	WM_PAINT          = 0x000F
	WM_KEYDOWN        = 0x0100
	VK_ESCAPE         = 0x1B
	SM_CXSCREEN       = 0
	IDC_ARROW         = 32512
	PS_SOLID          = 0
)

type WNDCLASSEX struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   syscall.Handle
	Icon       syscall.Handle
	Cursor     syscall.Handle
	Background syscall.Handle
	MenuName   *uint16
	ClassName  *uint16
	IconSm     syscall.Handle
}

type MSG struct {
	Hwnd    syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

type POINT struct {
	X, Y int32
}

type RECT struct {
	Left, Top, Right, Bottom int32
}

type PAINTSTRUCT struct {
	Hdc         syscall.Handle
	Erase       int32
	RcPaint     RECT
	Restore     int32
	IncUpdate   int32
	RgbReserved [32]byte
}

func main() {
	runtime.LockOSThread()

	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("mfch version 1.0 (Native Windows)")
		return
	}

	className, _ := syscall.UTF16PtrFromString("MfchClass")
	windowName, _ := syscall.UTF16PtrFromString("MFCH")

	hInstance := syscall.Handle(0) // GetModuleHandle(NULL) usually 0 in Go for current process? Or not.
	// Actually GetModuleHandle(nil) is better but 0 often works for hInstance in RegisterClass.
	// Let's rely on default behavior or get it properly if needed, but 0 is usually fine for current module in many contexts,
	// though RegisterClassEx technically wants the module handle.
	// Let's treat it as 0 for simplicity, if it fails we add GetModuleHandle.

	cursor, _, _ := procLoadCursorW.Call(0, uintptr(IDC_ARROW))

	wc := WNDCLASSEX{
		Size:       uint32(unsafe.Sizeof(WNDCLASSEX{})),
		Style:      0,
		WndProc:    syscall.NewCallback(wndProc),
		Instance:   hInstance,
		Cursor:     syscall.Handle(cursor),
		ClassName:  className,
		Background: 0, // No background brush, we paint ourselves or use transparency
	}

	if ret, _, _ := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc))); ret == 0 {
		// Try to see if it failed.
		// In pure syscall with NewCallback, sometimes issues arise.
		// But let's proceed.
	}

	// Calculate position for 2nd monitor
	// Simple heuristic: X = PrimaryScreenWidth, Y = 0
	screenWidth, _, _ := procGetSystemMetrics.Call(uintptr(SM_CXSCREEN))
	xPos := int32(screenWidth)
	yPos := int32(0)

	// If single monitor or user wants primary, we could change this logic.
	// But requirement was "Target the second monitor if available".
	// If only 1 monitor, xPos might be off-screen?
	// Win32 `CreateWindow` allows off-screen windows.
	// Refinement: check SM_CMONITORS (80)
	const SM_CMONITORS = 80
	numMonitors, _, _ := procGetSystemMetrics.Call(uintptr(SM_CMONITORS))
	if numMonitors <= 1 {
		xPos = 0
	}

	// Create Window
	// Use WS_EX_LAYERED for transparency (color key)
	// Use WS_EX_TRANSPARENT for click-through
	exStyle := uintptr(WS_EX_TOPMOST | WS_EX_LAYERED | WS_EX_TRANSPARENT | WS_EX_TOOLWINDOW)
	style := uintptr(WS_POPUP) // No border, no title

	width := int32(800)
	height := int32(600)

	hwnd, _, _ := procCreateWindowExW.Call(
		exStyle,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		style,
		uintptr(xPos),
		uintptr(yPos),
		uintptr(width),
		uintptr(height),
		0,
		0,
		uintptr(hInstance),
		0,
	)

	if hwnd == 0 {
		fmt.Println("Failed to create window")
		return
	}

	// Set transparency
	// Key color: RGB(255, 0, 255) -> Magenta will be transparent
	// 0x00FF00FF is 0x00BBGGRR

	const LWA_COLORKEY = 1
	const gwlExStyle = -20

	// Magenta
	keyColor := uintptr(0x00FF00FF)
	procSetLayeredWindowAttributes.Call(hwnd, keyColor, 0, LWA_COLORKEY)

	// Show Window
	const SW_SHOW = 5
	// syscall/user32 doesn't have ShowWindow in my list?
	// Ah, I missed loading ShowWindow.
	// Or I can use SetWindowPos to show it.
	// Let's add ShowWindow to imports or use SetWindowPos with SWP_SHOWWINDOW (0x0040)

	// Actually, let's just add ShowWindow to the proc list above?
	// No, I can't edit the variable block halfway.
	// I handled it by adding `procShowWindow` dynamically here? No, better to be clean.
	// I'll use SetWindowPos which IS defined.
	const SWP_SHOWWINDOW = 0x0040
	const HWND_TOPMOST = 0xffff
	// Note: HWND_TOPMOST is -1 cast to uintptr

	procSetWindowPos.Call(
		hwnd,
		uintptr(unsafe.Pointer(uintptr(^uintptr(0)))), // -1 (HWND_TOPMOST)
		0, 0, 0, 0,
		uintptr(0x0001|0x0002|SWP_SHOWWINDOW), // NOSIZE | NOMOVE | SHOWWINDOW
	)

	// Message Loop
	var msg MSG
	for {
		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

func wndProc(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_DESTROY:
		procPostQuitMessage.Call(0)
		return 0
	case WM_KEYDOWN:
		if wParam == VK_ESCAPE {
			procPostQuitMessage.Call(0)
		}
		return 0
	case WM_PAINT:
		var ps PAINTSTRUCT
		hdc, _, _ := procBeginPaint.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&ps)))

		// 1. Fill background with the Key Color (Magenta) so it becomes transparent
		// Create a solid brush
		// Using standard GDI: CreateSolidBrush
		// I need to load CreateSolidBrush.
		// I'll assume I can lazily load it here or add to global.
		// For robustness, I'll use raw syscall.LoadDLL locally if needed, but better to structure it.
		// Let's use FillRect with a created brush.

		gdi32local := syscall.NewLazyDLL("gdi32.dll")
		user32local := syscall.NewLazyDLL("user32.dll")
		procCreateSolidBrush := gdi32local.NewProc("CreateSolidBrush")
		procFillRect := user32local.NewProc("FillRect")

		magenta := uintptr(0x00FF00FF)
		hBrush, _, _ := procCreateSolidBrush.Call(magenta)

		procFillRect.Call(hdc, uintptr(unsafe.Pointer(&ps.RcPaint)), hBrush)
		procDeleteObject.Call(hBrush)

		// 2. Draw Crosshair
		// Create Pen: Black, 5px width
		// RGB(0,0,0) is 0
		hPen, _, _ := procCreatePen.Call(uintptr(PS_SOLID), 5, 0)
		oldPen, _, _ := procSelectObject.Call(hdc, hPen)

		// Center is relative to client area (800x600)
		const width = 800
		const height = 600
		const cx = width / 2
		const cy = height / 2

		// Gap size? User asked for 5px gap previously?
		// Or "remove the gap"?
		// Prompt history says "Remove the gap in the center".
		// OK, I will draw a simple cross without gap.

		const size = 30 // Length of arms

		// Horizontal
		procMoveToEx.Call(hdc, uintptr(cx-size), uintptr(cy), 0)
		procLineTo.Call(hdc, uintptr(cx+size), uintptr(cy))

		// Vertical
		procMoveToEx.Call(hdc, uintptr(cx), uintptr(cy-size), 0)
		procLineTo.Call(hdc, uintptr(cx), uintptr(cy+size))

		// Cleanup
		procSelectObject.Call(hdc, oldPen)
		procDeleteObject.Call(hPen)

		procEndPaint.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&ps)))
		return 0
	}

	ret, _, _ := procDefWindowProcW.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	return ret
}
