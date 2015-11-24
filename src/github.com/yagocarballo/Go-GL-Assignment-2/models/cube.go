package models

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
    "fmt"

    "github.com/yagocarballo/Go-GL-Assignment-2/wrapper"
)

type Cube struct {
    Name                                       string

	bufferObject, normalsObject, coloursObject uint32
	elementBuffer                              uint32

	DrawMode                                   DrawMode // Defines drawing mode of cube as points, lines or filled polygons

	vertexPositions, vertexColours, normals    *[]float32

	Model                                      mgl32.Mat4


    ShaderManager                              *wrapper.ShaderManager // Pointer to the Shader Manager
}

func NewCube(name string, vertexPositions, vertexColours, normals *[]float32, shaderManager *wrapper.ShaderManager) *Cube {
	return &Cube{
        name,               // Name
		0, 0, 0,            // bufferObject, normals, colours
		0,                  // elementBuffer
		DRAW_POLYGONS,      // drawmode
		vertexPositions, vertexColours, normals, // vertexPositions, vertexColours, normals
		mgl32.Ident4(),     // model
        shaderManager,      // Pointer to the Shader Manager
	}
}

func (cube *Cube) MakeVBO() {
	// Create a vertex buffer object to store vertices for the cube
	gl.GenBuffers(1, &cube.bufferObject);
	gl.BindBuffer(gl.ARRAY_BUFFER, cube.bufferObject);
	gl.BufferData(gl.ARRAY_BUFFER, len(*cube.vertexPositions) * 4, gl.Ptr(*cube.vertexPositions), gl.STATIC_DRAW);
	gl.BindBuffer(gl.ARRAY_BUFFER, 0);

	// Create a vertex buffer object to store vertex colours for the cube
	gl.GenBuffers(1, &cube.coloursObject);
	gl.BindBuffer(gl.ARRAY_BUFFER, cube.coloursObject);
	gl.BufferData(gl.ARRAY_BUFFER, len(*cube.vertexColours) * 4, gl.Ptr(*cube.vertexColours), gl.STATIC_DRAW);
	gl.BindBuffer(gl.ARRAY_BUFFER, 0);

	// Create the normals buffer for the cube
	gl.GenBuffers(1, &cube.normalsObject);
	gl.BindBuffer(gl.ARRAY_BUFFER, cube.normalsObject);
	gl.BufferData(gl.ARRAY_BUFFER, len(*cube.normals) * 4, gl.Ptr(*cube.normals), gl.STATIC_DRAW);
	gl.BindBuffer(gl.ARRAY_BUFFER, 0);
}

func (cube *Cube) Draw() {
    // Adds the Sphere Model to the Active Shader
    cube.ShaderManager.SetUniformMatrix4fv(cube.ShaderManager.ActiveShader, "model", 1, false, &cube.Model[0])

    /* Bind cube vertices. Note that this is in attribute index 0 */
	gl.BindBuffer(gl.ARRAY_BUFFER, cube.bufferObject)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	/* Bind cube colours. Note that this is in attribute index 1 */
	gl.BindBuffer(gl.ARRAY_BUFFER, cube.coloursObject)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, nil)

	/* Bind cube normals. Note that this is in attribute index 2 */
	gl.EnableVertexAttribArray(2)
	gl.BindBuffer(gl.ARRAY_BUFFER, cube.normalsObject)
	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 0, nil)

	// Enable this line to show model in wireframe
	if cube.DrawMode == DRAW_LINES {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}

	/* Draw our cube*/
	if cube.DrawMode == DRAW_POINTS {
		gl.DrawArrays(gl.POINTS, 0, int32(32))
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	}
}

func (cube *Cube) ResetModel () {
	cube.Model = mgl32.Ident4()
}

func (cube *Cube) Translate (Tx, Ty, Tz float32) {
	cube.Model = cube.Model.Mul4(mgl32.Translate3D(Tx, Ty, Tz))
}

func (cube *Cube) Scale (scaleX, scaleY, scaleZ float32) {
	cube.Model = cube.Model.Mul4(mgl32.Scale3D(scaleX, scaleY, scaleZ))
}

func (cube *Cube) Rotate (angle float32, axis mgl32.Vec3) {
	cube.Model = cube.Model.Mul4(mgl32.HomogRotate3D(angle, axis))
}

func (cube *Cube) GetDrawMode () DrawMode {
    return cube.DrawMode
}

func (cube *Cube) SetDrawMode (drawMode DrawMode) {
    cube.DrawMode = drawMode
}

func (cube *Cube) GetName () string {
    return cube.Name
}

func (cube *Cube) String () string {
    return fmt.Sprintf(`
               Cube --> %s
    -------------------------------------
    %s
    -------------------------------------
    `, cube.Name, cube.Model)
}
