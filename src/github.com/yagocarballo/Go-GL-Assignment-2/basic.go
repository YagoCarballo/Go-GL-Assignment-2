package main
import (
	"fmt"
	"runtime"

	"github.com/yagocarballo/Go-GL-Assignment-2/wrapper"
    "github.com/yagocarballo/Go-GL-Assignment-2/models"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
)

const windowWidth = 1024
const windowHeight = 768
const windowFPS = 60

// The Window Wrapper
var glw *wrapper.Glw

// Vertex array (Container) object.
var vertexArrayObject uint32

// Position and view globals
var speed float64 = 10

// Animation progress
var animationProgress float32 = 0

// Aspect ratio of the window defined in the reshape callback
var aspect_ratio float32

// Shader Manager
var shaderManager *wrapper.ShaderManager

// Selection
var selected_model models.Model

// Models
var view					*models.Camera
var lightPoint				*models.Sphere

// Index of a uniform to switch the colour mode in the vertex shader
var colorMode models.ColorMode

// Light point emit mode
var emitMode models.EmitMode

/////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////// Initialization ///////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////

//
// init
// This function is called by go as soon as this class is opened
//
func init() {
	// Locks the Execution in the main Thread as OpenGL is not thread Safe
	runtime.LockOSThread()
}

//
// main
// Entry point of program
//
func main() {
	// Creates the Window Wrapper
	glw = wrapper.NewWrapper(windowWidth, windowHeight, "Lab4")
	glw.SetFPS(windowFPS)

	// Creates the Window
	glw.CreateWindow()

	// Sets the Event Callbacks
	glw.SetRenderCallback(drawLoop)
	glw.SetKeyCallBack(keyCallback)
	glw.SetReshapeCallback(reshape)

	// Initializes the App
	InitApp(glw)

	// Prints the Keyboard Instructions
	printKeyboardMappings()

	// Sets the Viewport (Important !!, this has to run after the loop!!)
	defer gl.Viewport(0, 0, windowWidth, windowHeight)

	// Starts the Rendering Loop
	glw.StartLoop()
}

//
// InitShaders
// Loads the Shaders into the Shader Manager
//
func InitShaders () {
	// Creates the Shader Program
	shaderManager = wrapper.NewShaderManager()
	var err error; err = shaderManager.LoadShader("default", "./resources/shaders/basic.vert", "./resources/shaders/basic.frag")

	// If there is any error loading the shaders, it panics
	if err != nil {
		panic(err)
	}

	// Creates the Shader Program
	err = shaderManager.LoadShader("shiny", "./resources/shaders/shiny.vert", "./resources/shaders/shiny.frag")

	// If there is any error loading the shaders, it panics
	if err != nil {
		panic(err)
	}


	err = shaderManager.LoadShader("fragment_light", "./resources/shaders/fragment_light.vert", "./resources/shaders/fragment_light.frag")

	// If there is any error loading the shaders, it panics
	if err != nil {
		panic(err)
	}

	// Define uniforms to send to the shaders
	for name, _ := range shaderManager.Shaders {
		shaderManager.CreateUniform(name, "model")
		shaderManager.CreateUniform(name, "colourmode")
		shaderManager.CreateUniform(name, "emitmode")
		shaderManager.CreateUniform(name, "view")
		shaderManager.CreateUniform(name, "projection")
		shaderManager.CreateUniform(name, "lightpos")
	}
}

//
// Init App
// This function initializes the variables and sets up the environment.
//
// @param wrapper (*wrapper.Glw) the window wrapper
//
func InitApp(glw *wrapper.Glw) {
	/* Set the object transformation controls to their initial values */
	aspect_ratio = 1.3333
	colorMode = models.COLOR_SOLID
	emitMode = models.EMIT_COLORED

	// Generate index (name) for one vertex array object
	gl.GenVertexArrays(1, &vertexArrayObject);

	// Create the vertex array object and make it current
	gl.BindVertexArray(vertexArrayObject);

	InitShaders();

	// Create the Camera / View
	view = models.NewCamera("View", mgl32.LookAtV(
		mgl32.Vec3{0, 0, 2.5},	// Camera is at (0,0,4), in World Space
		mgl32.Vec3{0, 0, 0},	// and looks at the origin
		mgl32.Vec3{0, 1, 0},	// Head is up (set to 0,-1,0 to look upside-down)
	))

    // Creates the Lightpoint indicator
    lightPoint = models.NewSphere(
		"Light Point", 	// Name
		20, 			// Latitude
		20, 			// Longitude
		shaderManager,	// A pointer to the ShaderManager
	)
    lightPoint.MakeSphereVBO() // Creates the Light Point Buffer Object

	// Applies Initial Transforms to the Models
	InitialModelTransforms()

	// Sets the Light Point as the Default Selected Model
	selected_model = lightPoint

	// Changes the Title of the Window to display the Selected Model
	glw.Window.SetTitle(fmt.Sprintf("Selected Model --> %s", selected_model.GetName()))
}

//
// InitialModelTransforms
// Applies the Initial transforms for the Models
//
func InitialModelTransforms () {
	// Define the model transformations for the LightPoint
	lightPoint.ResetModel()
	lightPoint.Translate(0.0, 0.0, 1.0)
	lightPoint.Scale(0.05, 0.05, 0.05)

	// Apply rotations to the view position
	view.Rotate(0, mgl32.Vec3{1, 0, 0}) //rotating in clockwise direction around x-axis
	view.Rotate(0, mgl32.Vec3{0, 1, 0}) //rotating in clockwise direction around y-axis
	view.Rotate(0, mgl32.Vec3{0, 0, 1})
}

/////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////// Callbacks /////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////

//
// Draw Loop Function
// This function gets called on every update.
//
func drawLoop(glw *wrapper.Glw, delta float64) {
	// Sets the Clear Color (Background Color)
	gl.ClearColor(0.0, 0.8274509804, 0.9490196078, 1.0)

	// Clears the Window
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Enables Depth
	gl.Enable(gl.DEPTH_TEST)

	// Applies the Animations
	applyAnimations(delta)

	// Fov / Aspect / Near / Far
	// Projection matrix : 30° Field of View, 4:3 ratio, display range : 0.1 unit <-> 100 units
	var Projection mgl32.Mat4 = mgl32.Perspective(30.0, aspect_ratio, 0.1, 100.0)

    // Define the light position and transform by the view matrix
	var lightpos mgl32.Vec4 = view.Model.Mul4x1(mgl32.Vec4{lightPoint.Position.X(), lightPoint.Position.Y(), lightPoint.Position.Z(), 1.0})

	for name, _ := range shaderManager.Shaders {
		// Sets the Shader program to Use
		shaderManager.EnableShader(name)

		// Send our uniforms variables to the shader
		shaderManager.SetUniform1ui(name, "colourmode", uint32(colorMode))
		shaderManager.SetUniform1ui(name, "emitmode", emitMode.AsUint32())
		shaderManager.SetUniformMatrix4fv(name, "view", 1, false, &view.Model[0])
		shaderManager.SetUniformMatrix4fv(name, "projection", 1, false, &Projection[0])
		shaderManager.SetUniform4f(name, "lightpos", lightpos.X(), lightpos.Y(), lightpos.Z(), lightpos.W())
	}

	// Sets the Shader program to Use
	shaderManager.EnableShader("shiny")
	shaderManager.EnableShader("default")
	shaderManager.EnableShader("shiny")
	shaderManager.EnableShader("default")

    // Draw our light Position sphere
    emitMode = models.EMIT_BRIGHT
	shaderManager.SetUniform1ui(shaderManager.ActiveShader, "emitmode", emitMode.AsUint32())

    lightPoint.Draw() // Draws the Light Point

	gl.DisableVertexAttribArray(0);
	shaderManager.DisableShader()
    emitMode = models.EMIT_COLORED
}

//
// applyAnimations
// Applies animations (called once every loop)
//
// @param delta (float64) delta time of the update
//
func applyAnimations (delta float64) {

}

//
// key Callback
// This function gets called when a key is pressed
//
// @param window (*glfw.Window) a pointer to the window
// @param key (glfw.Key) the pressed key
// @param scancode (int) the scancode
// @param action (glfw.Action) the state of the key
// @param mods (glfw.ModifierKey) the pressed modified keys.
//
func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	var keySpeed float32 = 0.05
	var position mgl32.Vec4 = mgl32.Vec4{0, 0, 0, 0}
	var rotation mgl32.Vec4 = mgl32.Vec4{0, 0, 0, 0}
	var zoom float32 = 0.0

	// Increases the Speed of the Light Point
	if selected_model.GetName() == "Light Point" {
		keySpeed = 0.5
	}

	switch key {
	// If the Key Excape is pressed, it closes the App
	case glfw.KeyEscape:
		if action == glfw.Press {
			window.SetShouldClose(true)
		}

	// Changes the Selected Model
	case glfw.Key1:
		selected_model = lightPoint

	case glfw.Key2:
		selected_model = lightPoint

	case glfw.Key3:
		selected_model = lightPoint

	case glfw.Key4:
		selected_model = lightPoint

	case glfw.Key5:
		selected_model = lightPoint

	case glfw.Key6:
		selected_model = lightPoint

	case glfw.Key7:
		selected_model = lightPoint

	case glfw.Key8:
		selected_model = lightPoint

	case glfw.Key9:
		selected_model = view

	case glfw.Key0:
		selected_model = lightPoint

	// Applies Movement
	case glfw.KeyQ:
		position = mgl32.Vec4{0, 0, keySpeed, 0}

	case glfw.KeyW:
		position = mgl32.Vec4{0, -keySpeed, 0, 0}

	case glfw.KeyE:
		position = mgl32.Vec4{0, 0, -keySpeed, 0}

	case glfw.KeyA:
		position = mgl32.Vec4{keySpeed, 0, 0, 0}

	case glfw.KeyS:
		position = mgl32.Vec4{0, keySpeed, 0, 0}

	case glfw.KeyD:
		position = mgl32.Vec4{-keySpeed, 0, 0, 0}

	case glfw.KeyTab:
//		position = mgl32.Vec4{0, 0, 0, speed}

	case glfw.KeyR:
//		position = mgl32.Vec4{0, 0, 0, -speed}

	// Rotates
	case glfw.KeyI:
		rotation = mgl32.Vec4{1, 0, 0, 0}

	case glfw.KeyK:
		rotation = mgl32.Vec4{-1, 0, 0, 0}

	case glfw.KeyJ:
		rotation = mgl32.Vec4{0, 1, 0, 0}

	case glfw.KeyL:
		rotation = mgl32.Vec4{0, -1, 0, 0}

	case glfw.KeyU:
		rotation = mgl32.Vec4{0, 0, 1, 0}

	case glfw.KeyO:
		rotation = mgl32.Vec4{0, 0, -1, 0}

	case glfw.KeyY:
//		rotation = mgl32.Vec4{0, 0, 0, 1}

	case glfw.KeyP:
//		rotation = mgl32.Vec4{0, 0, 0, -1}

	// Zooms In / Out
	case glfw.KeyZ:
		zoom = -0.02

	case glfw.KeyX:
		zoom = 0.02

	// Speed Up / Down
	case glfw.KeyC:
		speed -= 1

	case glfw.KeyV:
		speed += 1
	}

	// Changes the Title of the Window to display the Selected Model
	glw.Window.SetTitle(fmt.Sprintf("Selected Model --> %s", selected_model.GetName()))

	// Applies the Transformations to the Selected Model
	ApplyTransformations(selected_model, rotation, position, zoom, keySpeed)

	// React only if the key was just pressed
	if action != glfw.Press {
		return;
	}

	switch key {
	case glfw.KeyM:
		if colorMode == models.COLOR_PER_SIDE {
			colorMode = models.COLOR_SOLID
		} else {
			colorMode = models.COLOR_PER_SIDE
		}
		fmt.Printf("Color Mode: %s \n", colorMode)

	// Cycle between drawing vertices, mesh and filled polygons
	case glfw.KeyN:
		selected_model.SetDrawMode(selected_model.GetDrawMode() + 1)
		if selected_model.GetDrawMode() > models.DRAW_POLYGONS {
			selected_model.SetDrawMode(models.DRAW_POINTS)
		}
		fmt.Printf("%s Draw Mode: %s \n", selected_model.GetName(), selected_model.GetDrawMode())

	// Prints the Keyboard Mappings
	case glfw.KeyB:
		printKeyboardMappings()

	case glfw.KeyEnter:
		// Creates the Shader Program
		err := shaderManager.LoadShader("default", "./resources/shaders/basic.vert", "./resources/shaders/basic.frag")

		// If there is any error loading the shaders, it panics
		if err != nil {
			log.Println(err)
		}

		// Creates the Shader Program
		err = shaderManager.LoadShader("shiny", "./resources/shaders/shiny.vert", "./resources/shaders/shiny.frag")

		// If there is any error loading the shaders, it panics
		if err != nil {
			log.Println(err)
		}

		err = shaderManager.LoadShader("fragment_light", "./resources/shaders/fragment_light.vert", "./resources/shaders/fragment_light.frag")

		// If there is any error loading the shaders, it panics
		if err != nil {
			log.Println(err)
		}
	case glfw.KeySpace:
		fmt.Println(selected_model)
	}
}

//
// ApplyTransformations
// Applies transformations to the a model
//
// @param model (models.Model): The Model
// @param rotation (mgl32.Vec4): Rotation Vector (Which side rotates)
// @param position (mgl32.Vec4): Position Increments
// @param zoom (float32): Zoom Increments
// @param speed (float32): Rotation Increment
//
func ApplyTransformations (model models.Model, rotation mgl32.Vec4, position mgl32.Vec4, zoom float32, speed float32) {
	if rotation.X() != 0 || rotation.Y() != 0 || rotation.Z() != 0 || rotation.W() != 0 {
		model.Rotate(speed, mgl32.Vec3{rotation.X(), rotation.Y(), rotation.Z()})
	}
	if position.X() != 0 || position.Y() != 0 || position.Z() != 0 || rotation.W() != 0 {
		model.Translate(position.X(), position.Y(), position.Z())
	}
	if zoom != 0 {
		var scaleAmount float32 = (1.0 + zoom)
		model.Scale(scaleAmount, scaleAmount, scaleAmount)
	}
}

//
// Reshape
// This gets called when the window changes its size
//
// @param window (*glfw.Window) a pointer to the window
// @param width (int) the width of the window
// @param height (int) the height of the window
//
func reshape(window *glfw.Window, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height));
	aspect_ratio = (float32(width) / 640.0 * 4.0) / (float32(height) / 480.0 * 3.0);
}

//
// printKeyboardMappings
// Prints the Keyboard instructions to console
//
func printKeyboardMappings () {
	fmt.Println(`

							Keyboard Instructions
							---------------------

Select A Model with a number
- The name of the selected model is displayed on the title of the window.
- Once a model is selected, it can be moved using the Keys Below
- Move and Rotate work like a cross: ( Example )

						| W   \ S 	-->  Up  \ Down	 |
						| A   \ D 	--> Left \ Right |
						| Q   \ E 	--> Back \ Front |
						| Tab \ R 	-->   -W \ +W	 |

-------------------------------------------------------------------------------
	 |-----||-----||-----||-----||-----||-----||-----||-----||-----||-----|
	 |  1  ||  2  ||  3  ||  4  ||  5  ||  6  ||  7  ||  8  ||  9  ||  0  |
	 |-----||-----||-----||-----||-----||-----||-----||-----||-----||-----|

					1. Cube			| 6.
					2. Sphere		| 7.
					3. Cylinder		| 8.
					4. Cog			| 9. Camera / View
					5. Small Cog	| 0. Light Point


			Position (Move)							   Rotation (Rotate)
--------------------------------------------------------------------------

|---------||-----||-----||-----||-----|  |-----||-----||-----||-----||-----|
|   Tab   ||  Q  ||  W  ||  E  ||  R  |  |  Y  ||  U  ||  I  ||  O  ||  P  |
|---------||-----||-----||-----||-----|  |-----||-----||-----||-----||-----|
		   |-----||-----||-----|	      	    |-----||-----||-----|
		   |  A  ||  S  ||  D  |	      	    |  J  ||  K  ||  L  |
		   |-----||-----||-----|	      	  	|-----||-----||-----|


 Zoom (-/+)  | Speed Up/Down				Instructions / Draw Mode / Color
----------------------------					---------------------

|-----||-----||-----||-----|					|-----||-----||-----|
|  Z  ||  X  ||  Z  ||  X  |					|  B  ||  N  ||  M  |
|-----||-----||-----||-----|					|-----||-----||-----|
-------------------------------------------------------------------------------

								DEBUG Options

- The Enter Key will reload the Shaders.
- The Space Key will print the current selected model matrix.

-------------------------------------------------------------------------------


	`)
}

// Define vertices for a cube in 12 triangles
var vertexPositions = []float32{
	-0.25, 0.25, -0.25,
	-0.25, -0.25, -0.25,
	0.25, -0.25, -0.25,

	0.25, -0.25, -0.25,
	0.25, 0.25, -0.25,
	-0.25, 0.25, -0.25,

	0.25, -0.25, -0.25,
	0.25, -0.25, 0.25,
	0.25, 0.25, -0.25,

	0.25, -0.25, 0.25,
	0.25, 0.25, 0.25,
	0.25, 0.25, -0.25,

	0.25, -0.25, 0.25,
	-0.25, -0.25, 0.25,
	0.25, 0.25, 0.25,

	-0.25, -0.25, 0.25,
	-0.25, 0.25, 0.25,
	0.25, 0.25, 0.25,

	-0.25, -0.25, 0.25,
	-0.25, -0.25, -0.25,
	-0.25, 0.25, 0.25,

	-0.25, -0.25, -0.25,
	-0.25, 0.25, -0.25,
	-0.25, 0.25, 0.25,

	-0.25, -0.25, 0.25,
	0.25, -0.25, 0.25,
	0.25, -0.25, -0.25,

	0.25, -0.25, -0.25,
	-0.25, -0.25, -0.25,
	-0.25, -0.25, 0.25,

	-0.25, 0.25, -0.25,
	0.25, 0.25, -0.25,
	0.25, 0.25, 0.25,

	0.25, 0.25, 0.25,
	-0.25, 0.25, 0.25,
	-0.25, 0.25, -0.25,
}

// Define an array of colours
var vertexColours = []float32{
	0.0, 0.0, 1.0, 1.0,
	0.0, 0.0, 1.0, 1.0,
	0.0, 0.0, 1.0, 1.0,
	0.0, 0.0, 1.0, 1.0,
	0.0, 0.0, 1.0, 1.0,
	0.0, 0.0, 1.0, 1.0,

	0.0, 1.0, 0.0, 1.0,
	0.0, 1.0, 0.0, 1.0,
	0.0, 1.0, 0.0, 1.0,
	0.0, 1.0, 0.0, 1.0,
	0.0, 1.0, 0.0, 1.0,
	0.0, 1.0, 0.0, 1.0,

	1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 0.0, 1.0,

	1.0, 0.0, 0.0, 1.0,
	1.0, 0.0, 0.0, 1.0,
	1.0, 0.0, 0.0, 1.0,
	1.0, 0.0, 0.0, 1.0,
	1.0, 0.0, 0.0, 1.0,
	1.0, 0.0, 0.0, 1.0,

	1.0, 0.0, 1.0, 1.0,
	1.0, 0.0, 1.0, 1.0,
	1.0, 0.0, 1.0, 1.0,
	1.0, 0.0, 1.0, 1.0,
	1.0, 0.0, 1.0, 1.0,
	1.0, 0.0, 1.0, 1.0,

	0.0, 1.0, 1.0, 1.0,
	0.0, 1.0, 1.0, 1.0,
	0.0, 1.0, 1.0, 1.0,
	0.0, 1.0, 1.0, 1.0,
	0.0, 1.0, 1.0, 1.0,
	0.0, 1.0, 1.0, 1.0,
}

/* Manually specified normals for our cube */
var normals = []float32{
	0, 0, -1,
	0, 0, -1,
	0, 0, -1,
	0, 0, -1,
	0, 0, -1,
	0, 0, -1,
	1, 0, 0,
	1, 0, 0,
	1, 0, 0,
	1, 0, 0,
	1, 0, 0,
	1, 0, 0,
	0, 0, 1,
	0, 0, 1,
	0, 0, 1,
	0, 0, 1,
	0, 0, 1,
	0, 0, 1,
	-1, 0, 0,
	-1, 0, 0,
	-1, 0, 0,
	-1, 0, 0,
	-1, 0, 0,
	-1, 0, 0,
	0, -1, 0,
	0, -1, 0,
	0, -1, 0,
	0, -1, 0,
	0, -1, 0,
	0, -1, 0,
	0, 1, 0,
	0, 1, 0,
	0, 1, 0,
	0, 1, 0,
	0, 1, 0,
	0, 1, 0,
}
