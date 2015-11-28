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
	"math/rand"
	"time"
	"math"
	"github.com/davecheney/profile"
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
var fishAnimationProgress []float32 = []float32{}

// Aspect ratio of the window defined in the reshape callback
var aspect_ratio float32

// Shader Manager
var shaderManager *wrapper.ShaderManager

// Selection
var selected_model models.Model

// Models
var view					*models.Camera
var lightPoint				*models.Sphere

var terrain 				*models.Terrain
var water 					*models.Terrain
var gopher 					*models.WavefrontObject
var gingerbreadHouse        *models.WavefrontObject
var dragon       		 	*models.WavefrontObject
var wall                    *models.WavefrontObject
var car       		 		*models.WavefrontObject
var seaCreatures       		[]*models.WavefrontObject

// Index of a uniform to switch the colour mode in the vertex shader
var colorMode models.ColorMode

// Light point emit mode
var emitMode models.EmitMode

// List of shaders
var shaderList = []string{
	"basic",
	"textureMaterial",
	"bumpMapMaterial",
	"terrain",
	"colorMaterial",
}


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
	if wrapper.DEBUG {
		cfg := profile.Config{
			MemProfile:     true,
			ProfilePath:    ".", // store profiles in current directory
			NoShutdownHook: true, // do not hook SIGINT
		}

		defer profile.Start(&cfg).Stop()
	}

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

	// Loads In the list of shaders
	for _, shaderName := range shaderList {
		var err error; err = shaderManager.LoadShader(
			shaderName,
			fmt.Sprintf("./resources/shaders/%s.vert", shaderName),
			fmt.Sprintf("./resources/shaders/%s.frag", shaderName),
		)

		// If there is any error loading the shaders, it panics
		if err != nil {
			panic(err)
		}
	}

	// Define uniforms to send to the shaders
	for name, _ := range shaderManager.Shaders {
		shaderManager.CreateUniform(name, "model")
		shaderManager.CreateUniform(name, "view")
		shaderManager.CreateUniform(name, "projection")
		shaderManager.CreateUniform(name, "lightpos")

		shaderManager.CreateUniform(name, "colourmode")
		shaderManager.CreateUniform(name, "emitmode")
		shaderManager.CreateUniform(name, "tone")
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
		mgl32.Vec3{0, -20, 6},	// Camera is at (0,0,4), in World Space
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

	// Creates the Terrain
	terrain = models.NewTerrainWithSeed(
//		time.Now().UnixNano(),
		999999,
		4.0,
		5.0,
		mgl32.Vec4{ 0.662, 0.405, 0.022, 1 },
	)
//	terrain.CreateTerrain(200, 200, 150.0, 150.0)
	terrain.CreateTerrain(200, 200, 350.0, 350.0)
//	terrain.CreateTerrain(300, 300, 150.0, 150.0)

	// Creates the Water
	water = models.NewTerrainWithSeed(
		incSeed,
		4.0,
		50.0,
		mgl32.Vec4{ 0, 0.618, 1, 0.9 },
	)
	water.CreateTerrain(250, 250, 350.0, 350.0)

	// Creates the Gopher
	gopher = models.NewObjectLoader()
	gopher.LoadObject("./resources/models/gopher/gopher.obj")
	gopher.CreateObject()

	// Creates the Gingerbread House
	gingerbreadHouse = models.NewObjectLoader()
	gingerbreadHouse.LoadObject("./resources/models/gingebreadHouse/gingebreadHouse.obj")
	gingerbreadHouse.CreateObject()

	// Creates the Dragon
	dragon = models.NewObjectLoader()
	dragon.LoadObject("./resources/models/dragon/dragon.obj")
	dragon.CreateObject()

	// Creates the Tree
	wall = models.NewObjectLoader()
	wall.LoadObject("./resources/models/wall/wall.obj")
	wall.CreateObject()


	// Creates the Car
	car = models.NewObjectLoader()
	car.LoadObject("./resources/models/car/car.obj")
	car.CreateObject()

	// Creates Sea Creatures
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	for i:=0; i<30; i++ {
		creature := models.NewObjectLoader()
		if (i % 2) == 0 {
			creature.LoadObject("./resources/models/clownFish/clownFish.obj")
		} else {
			creature.LoadObject("./resources/models/fish/fish.obj")
		}

		creature.CreateObject()
		seaCreatures = append(seaCreatures, creature)
		fishAnimationProgress = append(fishAnimationProgress, (random.Float32() * 10.0))
	}

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
	lightPoint.Translate(-10.0, -15.0, -20.0)
	lightPoint.Scale(1, 1, 1)

	// Apply rotations to the view position
	view.Rotate(0.7, mgl32.Vec3{1, 0, 0}) //rotating in clockwise direction around x-axis
	view.Rotate(0, mgl32.Vec3{0, 1, 0}) //rotating in clockwise direction around y-axis
	view.Rotate(0, mgl32.Vec3{0, 0, 1})
	view.Scale(0.7, 0.7, 0.7)

	// Apply the initial terrain Position
	terrain.ResetModel()
	terrain.Translate(0.0, 10.0, -20.0)
	terrain.Scale(0.5, 0.5, 0.5)
	terrain.Rotate(180, mgl32.Vec3{0, 1, 0})

	// Apply the initial terrain Position
	water.ResetModel()
	water.Translate(0.0, 10.0, -20.0)
	water.Scale(0.5, 0.5, 0.5)
	water.Rotate(180, mgl32.Vec3{0, 1, 0})

	// Apply the initial Gopher Position
	gopher.ResetModel()
	gopher.Translate(0.0, 1.8, -30.0)
	gopher.Scale(3.0, 3.0, 3.0)
	gopher.Rotate(1.7, mgl32.Vec3{1, 0, 0})
	gopher.Rotate(1.5, mgl32.Vec3{0, 1, 0})
	gopher.Rotate(1.5, mgl32.Vec3{0, 0, 1})

	// Moves the House
	gingerbreadHouse.ResetModel()
	gingerbreadHouse.Translate(10.5, 1.0, 20.0)
	gingerbreadHouse.Scale(3.0, 3.0, 3.0)
	gingerbreadHouse.Rotate(3.5, mgl32.Vec3{1, 0, 0})
	gingerbreadHouse.Rotate(-0.15, mgl32.Vec3{0, 0, 1})

	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	for _, creature := range seaCreatures {
		creature.ResetModel()
		creature.Translate((random.Float32() * 50.0), 12, (random.Float32() * 50.0))
		creature.Scale(0.2, 0.2, 0.2)
		creature.Rotate(180, mgl32.Vec3{1, 0, 0})
	}

	// Moves the dragon
	dragon.ResetModel()
	dragon.Translate(-18.0, -3.0, -30.0)
	dragon.Scale(3.0, 3.0, 3.0)
	dragon.Rotate(1.7, mgl32.Vec3{1, 0, 0})
	dragon.Rotate(1.5, mgl32.Vec3{0, 1, 0})
	dragon.Rotate(1.5, mgl32.Vec3{0, 0, 1})

	// Moves the tree
	wall.ResetModel()
	wall.Translate(-5.0, 4.0, 0.0)
//	tree.Scale(3.0, 3.0, 3.0)
	wall.Rotate(1.7, mgl32.Vec3{1, 0, 0})
	wall.Rotate(1.5, mgl32.Vec3{0, 1, 0})
	wall.Rotate(1.5, mgl32.Vec3{0, 0, 1})

	// Moves the car
	car.ResetModel()
	car.Translate(0.0, 6.0, -10.0)
	car.Scale(3.0, 3.0, 3.0)
	car.Rotate(1.5, mgl32.Vec3{1, 0, 0})
	car.Rotate(1.5, mgl32.Vec3{0, 1, 0})
	car.Rotate(1.5, mgl32.Vec3{0, 0, 1})
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
	gl.ClearColor(0.028, 0.156, 0.348, 1)

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
	for name, _ := range shaderManager.Shaders {
		// Sets the Shader program to Use
		shaderManager.EnableShader(name)
		shaderManager.SetUniform4f(name, "tone", terrain.ColorTone.X(), terrain.ColorTone.Y(), terrain.ColorTone.Z(), terrain.ColorTone.W())
	}
	shaderManager.EnableShader("terrain")
	terrain.DrawObject(shaderManager.Shaders["terrain"].Shader)

	// Sets the Shader program to Use
	for name, _ := range shaderManager.Shaders {
		// Sets the Shader program to Use
		shaderManager.EnableShader(name)
		shaderManager.SetUniform4f(name, "tone", water.ColorTone.X(), water.ColorTone.Y(), water.ColorTone.Z(), water.ColorTone.W())
	}

	shaderManager.EnableShader("bumpMapMaterial")
	gingerbreadHouse.DrawObject(shaderManager.CurrentShader())
	wall.DrawObject(shaderManager.CurrentShader())

	for _, creature := range seaCreatures {
		creature.DrawObject(shaderManager.CurrentShader())
	}

	shaderManager.EnableShader("textureMaterial")
	dragon.DrawObject(shaderManager.CurrentShader())

	shaderManager.EnableShader("colorMaterial")
	gopher.DrawObject(shaderManager.CurrentShader())
	car.DrawObject(shaderManager.CurrentShader())

	// Enables Transparencies
	gl.Enable(gl.BLEND)
//		gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	//	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)


	gl.BlendFunc(gl.SRC_COLOR, gl.ONE)

	shaderManager.EnableShader("terrain")
	water.DrawObject(shaderManager.CurrentShader())

	// Disables transparencies
	gl.Disable(gl.BLEND)


    // Draw our light Position sphere
    emitMode = models.EMIT_BRIGHT
	shaderManager.EnableShader("basic")
	shaderManager.SetUniform1ui(shaderManager.ActiveShader, "emitmode", emitMode.AsUint32())

    lightPoint.Draw() // Draws the Light Point

	gl.DisableVertexAttribArray(0);
	shaderManager.DisableShader()
    emitMode = models.EMIT_COLORED
}

var incSeed int64

//
// applyAnimations
// Applies animations (called once every loop)
//
// @param delta (float64) delta time of the update
//
func applyAnimations (delta float64) {
	animationProgress += float32(delta) * float32(1000)
	if animationProgress > 10 {
		model := water.Model
		water = models.NewTerrainWithSeed(
			incSeed,
			4.0,
			50.0,
			mgl32.Vec4{ 0, 0.618, 1, 0.9 },
		)
		water.CreateTerrain(250, 250, 350.0, 350.0)
		water.Model = model
		animationProgress = 0.0
		incSeed++
	}

	incStep := (float32(delta) * 20.0)

	// Animate the Planets
	for index, creature := range seaCreatures {
		pair := float32(1.0)
		// Gets the current Progress
		if (index % 2) == 0 {
			pair = -1.0
		}

		// Inverts direction
		fishAnimationProgress[index] += (pair * incStep)

		// Calculates the Next Step
		step := (360.0 / 1000.0) * float64(len(seaCreatures) - index) * float64(fishAnimationProgress[index])

		// Calculates the Next position
		var x float32 = (float32(index) * 3.0) * float32(math.Cos(step * float64(models.DEG_TO_RADIANS)))
		var y float32 = 12
		var z float32 = (float32(index) * 3.0) * float32(math.Sin(step * float64(models.DEG_TO_RADIANS)))

		// Resets model and applies transformations
		creature.ResetModel()
		creature.Translate(x, y, z)
		creature.Scale(0.2, 0.2, 0.2)
		creature.Rotate(180, mgl32.Vec3{1, 0, 0})
		creature.Rotate(-90, mgl32.Vec3{0, 0, pair})
	}
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
	var moveSpeed float32 = 0.5
	var position mgl32.Vec4 = mgl32.Vec4{0, 0, 0, 0}
	var rotation mgl32.Vec4 = mgl32.Vec4{0, 0, 0, 0}
	var zoom float32 = 0.0

	// Increases the Speed of the Light Point
	if selected_model.GetName() != "Terrain" && selected_model.GetName() != "View" {
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
		selected_model = gopher

	case glfw.Key2:
		selected_model = gingerbreadHouse

	case glfw.Key3:
		selected_model = dragon

	case glfw.Key4:
		selected_model = car

	case glfw.Key5:
		selected_model = wall

	case glfw.Key6:
		selected_model = lightPoint

	case glfw.Key7:
		selected_model = lightPoint

	case glfw.Key8:
		selected_model = terrain

	case glfw.Key9:
		selected_model = view

	case glfw.Key0:
		selected_model = lightPoint

	// Applies Movement
	case glfw.KeyQ:
		position = mgl32.Vec4{0, 0, moveSpeed, 0}

	case glfw.KeyW:
		position = mgl32.Vec4{0, -moveSpeed, 0, 0}

	case glfw.KeyE:
		position = mgl32.Vec4{0, 0, -moveSpeed, 0}

	case glfw.KeyA:
		position = mgl32.Vec4{moveSpeed, 0, 0, 0}

	case glfw.KeyS:
		position = mgl32.Vec4{0, moveSpeed, 0, 0}

	case glfw.KeyD:
		position = mgl32.Vec4{-moveSpeed, 0, 0, 0}

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
		// Loads In the list of shaders
		for _, shaderName := range shaderList {
			var err error; err = shaderManager.LoadShader(
				shaderName,
				fmt.Sprintf("./resources/shaders/%s.vert", shaderName),
				fmt.Sprintf("./resources/shaders/%s.frag", shaderName),
			)

			// If there is any error loading the shaders, it panics
			if err != nil {
				log.Println(err)
			}
		}
	case glfw.KeySpace:
		model := water.Model
		water = models.NewTerrainWithSeed(
			incSeed,
			4.0,
			50.0,
			mgl32.Vec4{ 0, 0.618, 1, 0.9 },
		)
		water.CreateTerrain(250, 250, 350.0, 350.0)
		water.Model = model
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

------------------------------------
 |-----||-----||-----||-----||-----|
 |  1  ||  2  ||  3  ||  4  ||  5  |
 |-----||-----||-----||-----||-----|
 |-----||-----||-----||-----||-----|
 |  6  ||  7  ||  8  ||  9  ||  0  |
 |-----||-----||-----||-----||-----|

1. Gopher				| 6.
2. Gingerbread House	| 7.
3. Dragon				| 8. Terrain
4. Car					| 9. Camera / View
5. Wall					| 0. Light Point


		  Position (Move)
-------------------------------------

|-------||-----||-----||-----||-----|
|  Tab  ||  Q  ||  W  ||  E  ||  R  |
|-------||-----||-----||-----||-----|
		 |-----||-----||-----|
		 |  A  ||  S  ||  D  |
		 |-----||-----||-----|

		 Rotation (Rotate)
-----------------------------------
|-----||-----||-----||-----||-----|
|  Y  ||  U  ||  I  ||  O  ||  P  |
|-----||-----||-----||-----||-----|
 	   |-----||-----||-----|
 	   |  J  ||  K  ||  L  |
 	   |-----||-----||-----|


 Zoom (-/+)  | Speed Up/Down
----------------------------
|-----||-----||-----||-----|
|  Z  ||  X  ||  Z  ||  X  |
|-----||-----||-----||-----|



Instructions / Draw Mode / Color
	---------------------
	|-----||-----||-----|
	|  B  ||  N  ||  M  |
	|-----||-----||-----|

----------------------------

DEBUG Options

- The Enter Key will reload the Shaders.
- The Space Key will print the current selected model matrix.

-------------------------------------------------------------


	`)
}
