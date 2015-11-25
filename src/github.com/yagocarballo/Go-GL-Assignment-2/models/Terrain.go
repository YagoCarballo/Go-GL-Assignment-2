//	terrain_object.cpp
//	Example class to show how to render a height map
//	Iain Martin November 2014
package models

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/ojrac/opensimplex-go"
	"fmt"
)

const (
	SizeOfUint16 = 2
	SizeOfInt32 = 4
	SizeOfFloat32 = 4
	SizeOfVec3 = SizeOfFloat32 * 3
)

type Terrain struct  {
	XSize, ZSize		uint32
	PerlinOctaves		uint32
	HeightScale			float32

	VBOVertices			uint32
	VBOColors			uint32
	VBONormals			uint32
	VBOIndices			uint32

	Vertices			[]mgl32.Vec3
	Normals				[]mgl32.Vec3
	Colors				[]mgl32.Vec3
	Indices				[]uint16

	Noise				[]float32

	Model				mgl32.Mat4

	Name				string
	DrawMode            DrawMode // Defines drawing mode of cube as points, lines or filled polygons
}

//	Define the vertex attributes for vertex positions and normals.
//	Make these match your application and vertex shader
//	You might also want to add colours and texture coordinates
func NewTerrain () *Terrain {
	return &Terrain{
		0,				// XSize: Set to zero because we haven't created the heightField array yet
		0,				// ZSize
		4,				// PerlinOctaves
		1.0,			// HeightScale

		0,				// VBOVertices
		0,				// VBOColors
		0,				// VBONormals
		0,				// VBOIndices

		[]mgl32.Vec3{},	// Vertices
		[]mgl32.Vec3{},	// Normals
		[]mgl32.Vec3{},	// Colors
		[]uint16{},		// Indices

		[]float32{},	// Noise

		mgl32.Ident4(), // Model

		"Terrain",
		DRAW_POLYGONS,
	}
}

//
// Copy the vertices, normals and element indices into vertex buffers
//
func (terrain *Terrain) CreateObject() {
	// Generate the vertex buffer object
	gl.GenBuffers(1, &terrain.VBOVertices)
	gl.BindBuffer(gl.ARRAY_BUFFER, terrain.VBOVertices)
	gl.BufferData(gl.ARRAY_BUFFER, int(len(terrain.Vertices) * 3 * 4), gl.Ptr(terrain.Vertices), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	/* Store the normals in a buffer object */
	gl.GenBuffers(1, &terrain.VBONormals)
	gl.BindBuffer(gl.ARRAY_BUFFER, terrain.VBONormals)
	gl.BufferData(gl.ARRAY_BUFFER, int(len(terrain.Normals) * 3 * 4), gl.Ptr(terrain.Normals), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	/* Store the Colors in a buffer object */
	gl.GenBuffers(1, &terrain.VBOColors)
	gl.BindBuffer(gl.ARRAY_BUFFER, terrain.VBOColors)
	gl.BufferData(gl.ARRAY_BUFFER, int(len(terrain.Colors) * 5 * 4), gl.Ptr(terrain.Colors), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// Generate a buffer for the indices
	gl.GenBuffers(1, &terrain.VBOIndices)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, terrain.VBOIndices)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(len(terrain.Indices) * 2), gl.Ptr(terrain.Indices), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

/* Enable vertex attributes and draw object
Could improve efficiency by moving the vertex attribute pointer functions to the
create object but this method is more general
This code is almost untouched fomr the tutorial code except that I changed the
number of elements per vertex from 4 to 3*/
func (terrain *Terrain) DrawObject(shaderProgram uint32) {
	// Reads the uniform Locations
	modelUniform := gl.GetUniformLocation(shaderProgram, gl.Str("model\x00"));

	// Send our uniforms variables to the currently bound shader
	gl.UniformMatrix4fv(modelUniform, 1, false, &terrain.Model[0]);

	var size int32	// Used to get the byte size of the element (vertex index) array

	// Get the vertices uniform position
	verticesUniform := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("position\x00")))
	normalsUniform := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("normal\x00")))
	colorsUniform := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("colour\x00")))

	// Describe our vertices array to OpenGL (it can't guess its format automatically)
	gl.EnableVertexAttribArray(verticesUniform)
	gl.BindBuffer(gl.ARRAY_BUFFER, terrain.VBOVertices)
	gl.VertexAttribPointer(
		verticesUniform, // attribute index
		3,               // number of elements per vertex, here (x,y,z)
		gl.FLOAT,        // the type of each element
		false,           // take our values as-is
		0,               // no extra data between each position
		nil,             // offset of first element
	)

	gl.EnableVertexAttribArray(normalsUniform)
	gl.BindBuffer(gl.ARRAY_BUFFER, terrain.VBONormals)
	gl.VertexAttribPointer(
		normalsUniform, 	// attribute
		3,                  // number of elements per vertex, here (x,y,z)
		gl.FLOAT,           // the type of each element
		false,           	// take our values as-is
		0,                  // no extra data between each position
		nil,                // offset of first element
	)

	gl.EnableVertexAttribArray(colorsUniform)
	gl.BindBuffer(gl.ARRAY_BUFFER, terrain.VBOColors)
	gl.VertexAttribPointer(
		colorsUniform, 	// attribute
		3,                  // number of elements per vertex, here (x,y,z)
		gl.FLOAT,           // the type of each element
		false,           	// take our values as-is
		0,                  // no extra data between each position
		nil,                // offset of first element
	)

	size = int32(len(terrain.Vertices))

	gl.PointSize(3.0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, terrain.VBOIndices);
	gl.GetBufferParameteriv(gl.ELEMENT_ARRAY_BUFFER, gl.BUFFER_SIZE, &size);

	// Enable this line to show model in wireframe
	switch terrain.DrawMode {
	case 1:
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	case 2:
		gl.DrawArrays(gl.POINTS, 0, int32(len(terrain.Vertices)))
//		gl.DrawElements(gl.POINTS, int32(len(terrain.Indices)), gl.UNSIGNED_SHORT, nil)
		return
	default:
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}

	/* Draw the triangle strips */
	for i := uint32(0); i < terrain.XSize - 1; i++ {
		var location int = SizeOfUint16 * int(i * terrain.ZSize * 2)
		gl.DrawElements(gl.TRIANGLE_STRIP, int32(terrain.ZSize * 2), gl.UNSIGNED_SHORT, gl.PtrOffset(location))
	}
}


/* Define the terrian heights */
/* Uses code adapted from OpenGL Shading Language Cookbook: Chapter 8 */
func (terrain *Terrain) CalculateNoise(freq, scale float32) {
	//	Create the array to store the noise values
	//	The size is the number of vertices * number of octaves
	terrain.Noise = []float32{}
	for i := uint32(0); i < (terrain.XSize * terrain.ZSize * terrain.PerlinOctaves); i++ {
		terrain.Noise = append(terrain.Noise, 0)
	}

	xFactor := 1. / float32(terrain.XSize - 1)
	zFactor := 1. / float32(terrain.ZSize - 1)

	for row := uint32(0); row < terrain.ZSize; row++ {
		for col := uint32(0); col < terrain.XSize; col++ {
			x := xFactor * float32(col)
			z := zFactor * float32(row)
			sum := float32(0.0);
			current_scale := scale;
			current_freq := freq;

			// Compute the sum for each octave
			for oct := uint32(0); oct < 4; oct++ {
				noiseGenerator := opensimplex.New()
				p := noiseGenerator.Eval2(float64(x) * float64(current_freq), float64(z) * float64(current_freq))
				val := float32(p) / current_scale
				sum += val;
				result := (sum + 1.0) / 2.0

				// Store the noise value in our noise array
				terrain.Noise[(row * terrain.XSize + col) * 4 + oct] = result

				// Move to the next frequency and scale
				current_freq *= 2.0
				current_scale *= scale
			}
		}
	}
}

//	Define the vertex array that specifies the terrain
//	(x, y) specifies the pixel dimensions of the heightfield (x * y) vertices
//	(xs, ys) specifies the size of the heightfield region
func (terrain *Terrain) CreateTerrain(xp, zp uint32, xs, zs float32) {
	terrain.XSize = xp
	terrain.ZSize = zp
	width := xs
	height := zs

	/* Create array of vertices */
	numVertices := terrain.XSize * terrain.ZSize;
	terrain.Vertices = make([]mgl32.Vec3, numVertices);
	terrain.Colors   = make([]mgl32.Vec3, numVertices);
	terrain.Normals  = make([]mgl32.Vec3, numVertices);

	/* Scale heights in relation to the terrain size */
	terrain.HeightScale = xs;

	/* First calculate the noise array which we'll use for our vertex height values */
	terrain.CalculateNoise(4.0, 5.0)

//	if wrapper.DEBUG {
//		// Debug code to check that noise values are sensible
//		for i := uint32(0); i < (terrain.XSize * terrain.ZSize * terrain.PerlinOctaves); i++ {
//			log.Printf("\n noise[%d] = %f", i, terrain.Noise[i]);
//		}
//	}

	/* Define starting (x,z) positions and the step changes */
	xpos := -width / 2.0;
	xpos_step := width / float32(xp);
	zpos_step := height / float32(zp);
	zpos_start := -height / 2.0;

	/* Define the vertex positions and the initial normals for a flat surface */
	for x := uint32(0); x < terrain.XSize; x++ {
		zpos := zpos_start;
		for z := uint32(0); z < terrain.ZSize; z++ {
			height := terrain.Noise[(x * terrain.ZSize + z) * 4 + 3]
			terrain.Vertices[x * terrain.XSize + z]	= mgl32.Vec3{ xpos, (height - 0.5) * terrain.HeightScale, zpos }
			terrain.Normals[x * terrain.XSize + z]	= mgl32.Vec3{ 0, 1.0, 0 } // Normals for a flat surface

			terrain.Colors[x * terrain.XSize + z]	= mgl32.Vec3{
				((1.0 * height) / 1.0),
				((1.0 * height) / 1.0),
				((1.0 * height) / 1.0),
			}

			zpos += zpos_step;
		}
		xpos += xpos_step;
	}

	/* Define vertices for triangle strips */
	for x := uint32(0); x < terrain.XSize - 1; x++ {
		top    := uint16(x * terrain.ZSize);
		bottom := uint16(top + uint16(terrain.ZSize));
		for z := uint32(0); z < terrain.ZSize; z++ {
			terrain.Indices = append(terrain.Indices, top, bottom)
			top ++
			bottom ++
		}
	}

	terrain.CalculateNormals()
	terrain.CreateObject()
}

//	Calculate normals by using cross products along the triangle strips
//	and averaging the normals for each vertex
func (terrain *Terrain) CalculateNormals() {
	var element_pos uint32 = 0
	var AB, AC, cross_product mgl32.Vec3

	// Loop through each triangle strip
	for x := uint32(0); x < terrain.XSize - 1; x++ {
		// Loop along the strip
		for tri := uint32(0); tri < terrain.ZSize * 2 - 2; tri++ {
			// Extract the vertex indices from the element array
			v1 := terrain.Indices[element_pos];
			v2 := terrain.Indices[element_pos + 1];
			v3 := terrain.Indices[element_pos + 2];

			// Define the two vectors for the triangle
			AB = terrain.Vertices[v2].Sub(terrain.Vertices[v1])
			AC = terrain.Vertices[v3].Sub(terrain.Vertices[v1])

			// Calculate the cross product
			cross_product = AB.Cross(AC)

			// Add this normal to the vertex normal for all three vertices in the triangle
			terrain.Normals[v1] = terrain.Normals[v1].Add(cross_product)
			terrain.Normals[v2] = terrain.Normals[v2].Add(cross_product)
			terrain.Normals[v3] = terrain.Normals[v3].Add(cross_product)

			// Move on to the next vertex along the strip
			element_pos++
		}

		// Jump past the lat two element positions to reach the start of the strip
		element_pos += 2
	}

	// Normalise the normals
	for v := uint32(0); v < terrain.XSize * terrain.ZSize; v ++ {
		terrain.Normals[v] = terrain.Normals[v].Normalize()
	}
}

func (terrain *Terrain) ResetModel() {
	terrain.Model = mgl32.Ident4()
}

func (terrain *Terrain) Translate(Tx, Ty, Tz float32) {
	terrain.Model = terrain.Model.Mul4(mgl32.Translate3D(Tx, Ty, Tz))
}

func (terrain *Terrain) Scale(scaleX, scaleY, scaleZ float32) {
	terrain.Model = terrain.Model.Mul4(mgl32.Scale3D(scaleX, scaleY, scaleZ))
}

func (terrain *Terrain) Rotate(angle float32, axis mgl32.Vec3) {
	terrain.Model = terrain.Model.Mul4(mgl32.HomogRotate3D(angle, axis))
}

func (terrain *Terrain) GetDrawMode () DrawMode {
	return terrain.DrawMode
}

func (terrain *Terrain) SetDrawMode (drawMode DrawMode) {
	terrain.DrawMode = drawMode
}

func (terrain *Terrain) GetName () string {
	return terrain.Name
}

func (terrain *Terrain) String () string {
	return fmt.Sprintf(`
               Terrain --> %s
    -------------------------------------
    %s
    -------------------------------------
    `, terrain.Name, terrain.Model)
}

