# mfch - Make FLIR Crosshair

**mfch** is a lightweight, high-performance crosshair overlay application written in Go. It utilizes OpenGL and GLFW to render a highly customizable, persistent crosshair that stays on top of all other windows, making it the perfect tool for FLIR overlay.

## Key Features

- **Always on Top:** The overlay window floats above all other applications, ensuring the crosshair is always visible.
- **Click-Through:** The window is fully transparent to mouse events, allowing you to interact with applications underneath without obstruction.
- **Transparent Background:** Using a transparent framebuffer, only the crosshair pixels are drawn, keeping your view the clear.
- **Smart Multi-Monitor Support:** 
  - Automatically detects multiple monitors.
  - Defaults to the **secondary monitor**, ideal for dual-screen setups.
  - Falls back gracefully to the primary monitor if only one is detected.
- **Legacy Compatibility:** Verified support for **Windows 7** and **Windows XP** (requires Go 1.20/1.21 for build).
- **Custom Rendering:** Built on OpenGL 2.1 legacy profile to ensure maximum hardware compatibility across older and newer machines.

## Installation & Building

### Prerequisites

To build `mfch` from source, you need:

1.  **Go:** Version 1.21 or 1.20 (Required for Windows 7/XP compatibility).
2.  **C Compiler:**
    - **macOS/Linux:** standard `gcc` (via Xcode Command Line Tools or build-essential).
    - **Windows Cross-Compile:** `mingw-w64` is required on macOS/Linux.

### Quick Build (macOS/Linux)

A convenience script `build.sh` is provided to handle dependencies, check versions, and cross-compile.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/LeyllProst/mfch.git
    cd mfch
    ```

2.  **Run the build script:**
    ```bash
    ./build.sh
    ```

    The script will output:
    - `mfch`: Native binary for your current OS (macOS).
    - `mfch.exe`: Windows 64-bit binary (if `mingw-w64` is installed).

### Manual Build

**For macOS (Native):**
```bash
go build -o mfch main.go
```

**For Windows (Cross-compile from macOS):**
```bash
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o mfch.exe main.go
```

## Running the Application

### macOS
Open your terminal and run:
```bash
./mfch
```
*Note: You may need to grant "Screen Recording" or "Accessibility" permissions to your terminal or the app to allow it to draw over other full-screen apps.*

### Windows
Simply double-click `mfch.exe`. 
- The application implies no console window (via `-ldflags -H=windowsgui` if built manually, currently the build script uses `-ldflags "-s -w"` which might strip symbols but leave console; modify build flags to `-ldflags -H=windowsgui` to hide console if desired).
- **Windows 7/XP Note:** If you encounter missing DLL errors, ensure you have the Visual C++ Redistributable or appropriate drivers installed, though Go binaries are mostly static.

## Configuration

Currently, configuration is handled directly in the code for maximum performance and simplicity. 

**Adjusting Position & Monitor:**
Edit `main.go` around line 40:
```go
if len(monitors) > 1 {
    monitor := monitors[1] // Selects the 2nd monitor
    // ...
}
```

**Adjusting Crosshair Style:**
Edit `drawCrosshair` function in `main.go`:
```go
func drawCrosshair() {
    gl.LineWidth(5.0)              // Thickness
    gl.Color4f(0.0, 0.0, 0.0, 1.0) // Color (R, G, B, Alpha)
    // ...
}
```

## Troubleshooting

- **"DLL not found" on Windows:** Ensure your graphics drivers support OpenGL 2.1.
- **Crosshair not visible:**
  - Check if the application is running (Task Manager).
  - On macOS, ensure permissions are granted.
  - On Windows, some exclusive-fullscreen games might override the overlay. Run games in "Bordered Windowless" or "Windowed" mode for best results.

## License

MIT License
