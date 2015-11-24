package wrapper

import (
	"runtime"
	"log"
	"fmt"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type Glw struct  {
	// Properties
	Width, Height int
	Title string

	// State
	fps int
	running bool
	Window *glfw.Window

	// Callbacks
	renderer func(glw *Glw, elapsed float64)
	keyCallBack glfw.KeyCallback
	reshape glfw.FramebufferSizeCallback
}

// This function is called by go as soon as this library is imported
func init () {
	// Locks the Execution in the main Thread as OpenGL is not thread safe
	runtime.LockOSThread()
}

//
// New Wrapper
// This is the Constructor, Creates a wrapper instance and returns the pointer to it.
//
// @param width (int) the window width
// @param heigh (int) the window height
// @param title (string) the window title
//
// @return wrapper (*Glw) a pointer to the wrapper.
//
func NewWrapper(width, height int, title string) *Glw {
	return &Glw{ width, height, title, 60, true, nil, nil, nil, nil }
}

// Public Functions

//
// Create Window
// This creates a window, initiates GL and sets the event listeners
//
// @return window (*glfw.Window) pointer to the window
//
func (glw *Glw) CreateWindow () *glfw.Window {
	// Init GLFW
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	// Sets the OpenGL Version
	setOpenGlVersion()

	// Creates the Window
	win, err := glfw.CreateWindow(glw.Width, glw.Height, glw.Title, nil, nil)
	if err != nil {
		panic(err)
	}

	// Prints the OpenGL Versions at the end
	defer printOpenGlVersionInfo()

	// Sets this context as the current context
	win.MakeContextCurrent()

	// Initiates GL
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// Enables Depth
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	win.SetInputMode(glfw.StickyKeysMode, 1)

	// Sets the Window to the Wrapper
	glw.SetWindow(win)
	return win
}

//
// Start Loop
// this starts the event loop which runs until the program ends
//
func (glw *Glw) StartLoop () {
	optimalTime := 1000.0 / float64(glw.fps)
	previousTime := glfw.GetTime()

	// If the Window is open keep looping
	for !glw.GetWindow().ShouldClose() {
		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		delta	:= elapsed / optimalTime
		previousTime = time

		// Calls the Render Callback
		glw.renderer(glw, delta)

		// Triggers window refresh
		glw.GetWindow().SwapBuffers()

		// Triggers events
		glfw.PollEvents()
	}

	// Called at the end of the program, and terminates the window system
	glw.Terminate()
}

//
// Terminate
// When this is called, it destroys the window and terminates glfw
//
func (glw *Glw) Terminate () {
	// Clean up
	glw.Window.Destroy()
	glfw.Terminate()
}

// Private Functions

//
// set OpenGl Version
// Sets the openGL version to the window
//
func setOpenGlVersion() {
	glfw.WindowHint(glfw.Samples, 4) // Anti Aliasing (16 for nice Screenshots)
	glfw.WindowHint(glfw.ContextVersionMajor, 3) // Mac will use the latest available, even if 3.3 is selected
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)    // Necessary for OS X (This removes any deprecated API in 4.1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile) // Necessary for OS X
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.False)

	glfw.WindowHint(glfw.Resizable, glfw.True)
}

//
// print OpenGl Version Info
// Prints the OpenGL Version to the Console
//
func printOpenGlVersionInfo() {
	version := gl.GoStr(gl.GetString(gl.VERSION))
	renderer := gl.GoStr(gl.GetString(gl.RENDERER))
	shaderVersion := gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION))
	fmt.Println("OpenGL version", version)
	fmt.Println("OpenGL renderer", renderer)
	fmt.Println("Max supported shader language version", shaderVersion)
}

// Getters & Setters

//
// SetFPS
// Sets the FPS for the window
//
func (glw *Glw) SetFPS (fps int) {
	glw.fps = fps
}

//
// GetFPS
// Gets the FPS for the Window
//
func (glw *Glw) GetFPS () int {
	return glw.fps
}

//
// SetWindow
// Sets the Window Object
//
func (glw *Glw) SetWindow (window *glfw.Window) {
	glw.Window = window
}

//
// GetWindow
// Gets the Window Object
//
func (glw *Glw) GetWindow () *glfw.Window {
	return glw.Window
}

//
// SetRenderCallback
// Sets the Render Callback (loop)
//
func (glw *Glw) SetRenderCallback (callback func(glw *Glw, elapsed float64)) {
	glw.renderer = callback
}

//
// SetKeyCallBack
// Sets the Key Callback
//
func (glw *Glw) SetKeyCallBack (callback glfw.KeyCallback) {
	glw.keyCallBack = callback
	glw.Window.SetKeyCallback(callback)
}

//
// SetReshapeCallback
// Sets the window resize callback
//
func (glw *Glw) SetReshapeCallback (callback glfw.FramebufferSizeCallback) {
	glw.reshape = callback
	glw.Window.SetFramebufferSizeCallback(callback)
}