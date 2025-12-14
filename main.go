package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 800
	height = 600
	// MousePassthrough is the definition for GLFW_MOUSE_PASSTHROUGH which might be missing in older bindings.
	// Value: 0x0002000D
	MousePassthrough glfw.Hint = 0x0002000D
)

func main() {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Floating, glfw.True)               // Always on top
	glfw.WindowHint(glfw.Decorated, glfw.False)             // No border/title bar
	glfw.WindowHint(glfw.TransparentFramebuffer, glfw.True) // Transparent background
	glfw.WindowHint(MousePassthrough, glfw.True)            // Click-through

	window, err := glfw.CreateWindow(width, height, "Crosshair", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	var xPos, yPos int

	monitors := glfw.GetMonitors()
	// Target the second monitor if available (User request: "second NVIDIA DVI port")
	if len(monitors) > 1 {
		monitor := monitors[1]
		// For the second monitor, we often want to start at its top-left.
		// GetPos returns the screen coordinates of the upper-left corner of the monitor's viewport.
		xPos, yPos = monitor.GetPos()
	} else {
		// Fallback to primary monitor
		monitor := glfw.GetPrimaryMonitor()
		xPos, yPos = monitor.GetPos()
	}

	window.SetPos(xPos, yPos)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		drawCrosshair()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func drawCrosshair() {
	// Draw a black crosshair
	gl.LineWidth(5.0)
	gl.Begin(gl.LINES)
	gl.Color4f(0.0, 0.0, 0.0, 1.0) // Black color

	const length = 0.07
	const ratio = float32(height) / float32(width)

	// Horizontal line
	gl.Vertex2f(-length*ratio, 0.0)
	gl.Vertex2f(length*ratio, 0.0)

	// Vertical line
	gl.Vertex2f(0.0, -length)
	gl.Vertex2f(0.0, length)

	gl.End()
}
