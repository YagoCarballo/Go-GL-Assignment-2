package models

import (
	"math"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
    "fmt"
    "math/rand"
    "time"

    "github.com/yagocarballo/Go-GL-Assignment-2/wrapper"
)

// Define buffer object indices
type Sphere struct {
    Name                                             string

	sphereBufferObject, sphereNormals, sphereColours uint32
	elementBuffer                                    uint32

	DrawMode                                         DrawMode // Defines drawing mode of sphere as points, lines or filled polygons
	numLats, numLongs                                uint32      //Define the resolution of the sphere object

	numSphereVertices                                uint32

	Model                                            mgl32.Mat4

    Position                                         mgl32.Vec4

    ShaderManager                                    *wrapper.ShaderManager // Pointer to the Shader Manager
}

func NewSphere(name string, numLats, numLongs uint32, shaderManager *wrapper.ShaderManager) *Sphere {
	return &Sphere{
        name,               // Name
		0, 0, 0,            // sphereBufferObject, sphereNormals, sphereColours
		0,                  // elementBuffer
		DRAW_POLYGONS,      // drawmode
		numLats, numLongs,  // numLats, numLongs
		0,                  // numSphereVertices
		mgl32.Ident4(),     // model
        mgl32.Vec4{},       // Position
        shaderManager,      // Pointer to the Shader Manager
	}
}

// Make a sphere from two triangle fans (one at each pole) and triangle strips along latitudes
// This version uses indexed vertex buffers for both the fans at the poles and the latitude strips
func (sphere *Sphere) MakeSphereVBO() {
	var i uint32

	// Calculate the number of vertices required in sphere
	sphere.numSphereVertices = 2 + ((sphere.numLats - 1) * sphere.numLongs)
	pVertices, pNormals := sphere.MakeUnitSphere()
	pColours := sphere.GenerateColors()

	/* Generate the vertex buffer object */
	gl.GenBuffers(1, &sphere.sphereBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, sphere.sphereBufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, int(4 * sphere.numSphereVertices * 3), gl.Ptr(pVertices), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	/* Store the normals in a buffer object */
	gl.GenBuffers(1, &sphere.sphereNormals)
	gl.BindBuffer(gl.ARRAY_BUFFER, sphere.sphereNormals)
	gl.BufferData(gl.ARRAY_BUFFER, int(4 * sphere.numSphereVertices * 3), gl.Ptr(pNormals), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	/* Store the colours in a buffer object */
	gl.GenBuffers(1, &sphere.sphereColours)
	gl.BindBuffer(gl.ARRAY_BUFFER, sphere.sphereColours)
	gl.BufferData(gl.ARRAY_BUFFER, int(4 * sphere.numSphereVertices * 4), gl.Ptr(pColours), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	/* Calculate the number of indices in our index array and allocate memory for it */
	numIndices := ((sphere.numLongs * 2) + 2) * (sphere.numLats - 1) + ((sphere.numLongs + 2) * 2)
	pIndices := make([]uint32, numIndices)

	// fill "indices" to define triangle strips
	var index int = 0 // Current index

	// Define indices for the first triangle fan for one pole
	for i = 0; i < sphere.numLongs + 1; i++ {
		pIndices[index] = i
		index++
	}

	pIndices[index] = 1    // Join last triangle in the triangle fan
	index++

	var j uint32
	var start uint32 = 1        // Start index for each latitude row
	for j = 0; j < sphere.numLats - 2; j++ {
		for i = 0; i < sphere.numLongs; i++ {
			pIndices[index] = start + i
			index++

			pIndices[index] = start + i + sphere.numLongs
			index++
		}

		// close the triangle strip loop by going back to the first vertex in the loop
		pIndices[index] = start
		index++

		// close the triangle strip loop by going back to the first vertex in the loop
		pIndices[index] = start + sphere.numLongs
		index++

		start += sphere.numLongs
	}

	// Define indices for the last triangle fan for the south pole region
	for i = sphere.numSphereVertices - 1; i > sphere.numSphereVertices - sphere.numLongs - 2; i-- {
		pIndices[index] = i
		index++
	}
	pIndices[index] = sphere.numSphereVertices - 2 // Tie up last triangle in fan
	index++

	// Generate a buffer for the indices
	gl.GenBuffers(1, &sphere.elementBuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, sphere.elementBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(numIndices * 4), gl.Ptr(pIndices), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (sphere *Sphere) GenerateColors() []float32 {
    pColours := make([]float32, (sphere.numSphereVertices * 4))

    // Seeds the Random Number
    rand.Seed(time.Now().UnixNano())
    var colorLine uint32 = 0
    var red float32 = 0.5
    var green float32 = 0.5
    var blue float32 = 0.5

    var lockColor int = rand.Intn(3)

    if lockColor == 0 {
        red = rand.Float32() * 0.6
    }
    if lockColor == 1 {
        green = rand.Float32() * 0.6
    }
    if lockColor == 2 {
        blue = rand.Float32() * 0.6
    }

    var i uint32 = 0

    // Define colours as the x,y,z components of the sphere vertices
    for i = 0; i < sphere.numSphereVertices; i++ {
        colorLine++
        if colorLine >= sphere.numLongs {
            colorLine = 0
            if lockColor == 0 {
                red = rand.Float32() * 0.6
            }
            if lockColor == 1 {
                green = rand.Float32() * 0.6
            }
            if lockColor == 2 {
                blue = rand.Float32() * 0.6
            }
        }

        pColours[i * 4] = red
        pColours[i * 4 + 1] = green
        pColours[i * 4 + 2] = blue
        pColours[i * 4 + 3] = 1.0
    }

    return pColours
}

// Define the vertex positions for a sphere. The array of vertices must have previosuly been created.
func (sphere *Sphere) MakeUnitSphere() ([]float32, []float32) {
	var vnum int32 = 0
	var x, y, z, lat_radians, lon_radians float32
	var lat, lon float32

	pVertices := make([]float32, (sphere.numSphereVertices * 3))
	pNormals := make([]float32, (sphere.numSphereVertices * 3))

	// Define north pole
	pVertices[0] = 0.0; pNormals[0] = 0
	pVertices[1] = 0.0; pNormals[1] = 1
	pVertices[2] = 1.0; pNormals[2] = 0
	vnum++

	latStep := 180.0 / float32(sphere.numLats)
	longStep := 360.0 / float32(sphere.numLongs)

	/* Define vertices along latitude lines */
	for lat = 90.0 - latStep; lat > -90.0; lat -= latStep {
		lat_radians = lat * DEG_TO_RADIANS
		for lon = -180.0; lon < 180.0; lon += longStep {
			lon_radians = lon * DEG_TO_RADIANS

			x = float32(math.Cos(float64(lat_radians)) * math.Cos(float64(lon_radians)))
			y = float32(math.Cos(float64(lat_radians)) * math.Sin(float64(lon_radians)))
			z = float32(math.Sin(float64(lat_radians)))

			/* Define the vertex */
			pVertices[vnum * 3] = x
			pVertices[vnum * 3 + 1] = y
			pVertices[vnum * 3 + 2] = z
			vnum++
		}
	}

	/* Define south pole */
	pVertices[vnum * 3] = 0
	pVertices[vnum * 3 + 1] = 0
	pVertices[vnum * 3 + 2] = -1.0

	return pVertices, pNormals
}

// Draws the sphere form the previously defined vertex and index buffers
func (sphere *Sphere) Draw() {
    // Adds the Sphere Model to the Active Shader
    sphere.ShaderManager.SetUniformMatrix4fv(sphere.ShaderManager.ActiveShader, "model", 1, false, &sphere.Model[0])

	/* Draw the vertices as GL_POINTS */
	gl.BindBuffer(gl.ARRAY_BUFFER, sphere.sphereBufferObject)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	/* Bind the sphere colours */
	gl.BindBuffer(gl.ARRAY_BUFFER, sphere.sphereColours)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(1)

	/* Bind the sphere normals */
	gl.BindBuffer(gl.ARRAY_BUFFER, sphere.sphereNormals)
	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(2)

	gl.PointSize(3.0)

	// Enable this line to show model in wireframe
	if sphere.DrawMode == DRAW_LINES {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}

	if sphere.DrawMode == DRAW_POINTS {
		gl.DrawArrays(gl.POINTS, 0, int32(sphere.numSphereVertices))
	} else {
		/* Bind the indexed vertex buffer */
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, sphere.elementBuffer)

		/* Draw the north pole regions as a triangle  */
		gl.DrawElements(gl.TRIANGLE_FAN, int32(sphere.numLongs + 2), gl.UNSIGNED_INT, nil)

		/* Calculate offsets into the indexed array. Note that we multiply offsets by 4
		   because it is a memory offset the indices are type GLuint which is 4-bytes */
		var lat_offset_jump int = int((sphere.numLongs * 2) + 2)
		var lat_offset_start int = int(sphere.numLongs + 2)
		var lat_offset_current int = lat_offset_start * 4

		var i uint32

		/* Draw the triangle strips of latitudes */
		for i = 0; i < sphere.numLats - 2; i++ {
			gl.DrawElements(gl.TRIANGLE_STRIP, int32(sphere.numLongs * 2 + 2), gl.UNSIGNED_INT, gl.PtrOffset(lat_offset_current))
			lat_offset_current += (lat_offset_jump * 4)
		}
		/* Draw the south pole as a triangle fan */
		gl.DrawElements(gl.TRIANGLE_FAN, int32(sphere.numLongs + 2), gl.UNSIGNED_INT, gl.PtrOffset(lat_offset_current))
	}
}

func (sphere *Sphere) ResetModel() {
	sphere.Model = mgl32.Ident4()
}

func (sphere *Sphere) Translate(Tx, Ty, Tz float32) {
    sphere.Position = sphere.Position.Add(mgl32.Vec4{Tx, Ty, Tz, 0})
	sphere.Model = sphere.Model.Mul4(mgl32.Translate3D(Tx, Ty, Tz))
}

func (sphere *Sphere) Scale(scaleX, scaleY, scaleZ float32) {
	sphere.Model = sphere.Model.Mul4(mgl32.Scale3D(scaleX, scaleY, scaleZ))
}

func (sphere *Sphere) Rotate(angle float32, axis mgl32.Vec3) {
	sphere.Model = sphere.Model.Mul4(mgl32.HomogRotate3D(angle, axis))
}

func (sphere *Sphere) GetDrawMode () DrawMode {
    return sphere.DrawMode
}

func (sphere *Sphere) SetDrawMode (drawMode DrawMode) {
    sphere.DrawMode = drawMode
}

func (sphere *Sphere) GetName () string {
    return sphere.Name
}

func (sphere *Sphere) String () string {
    return fmt.Sprintf(`
             Sphere --> %s
    -------------------------------------
    %s
    -------------------------------------
    `, sphere.Name, sphere.Model)
}
