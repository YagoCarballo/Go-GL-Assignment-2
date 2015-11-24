package models

import (
    "github.com/go-gl/gl/all-core/gl"
    "github.com/go-gl/mathgl/mgl32"
    "fmt"
    "math"
    "github.com/yagocarballo/Go-GL-Assignment-2/wrapper"
)

// Define buffer object indices
type Cylinder struct {
    Name                                                    string      // Name

    cylinderBufferObject, cylinderNormals, cylinderColours  uint32      // Model Buffers
    elementBuffer                                           uint32      // Indices Buffer

    DrawMode                                                DrawMode    // Defines drawing mode of Cylinder as points, lines or filled polygons

    VerticesPerDisk                                         uint32      // Defines the Vertices Per Disk for the Cylinder
    Height                                                  float32     // Defines the Height for the Cylinder
    Radius                                                  float32     // Defines the Radius for the Cylinder

    Model                                                   mgl32.Mat4  // Model

    numCylinderVertices                                     uint32      // Total count of vertices


    ShaderManager                                           *wrapper.ShaderManager // Pointer to the Shader Manager
}

func NewCylinder (name string, vertices uint32, height, radius float32, shaderManager *wrapper.ShaderManager) *Cylinder {
    return &Cylinder{
        name,           // Name
        0, 0, 0,        // cylinderBufferObject, cylinderNormals, cylinderColours
        0,              // elementBuffer
        DRAW_POLYGONS,  // DrawMode
        vertices,       // VerticesPerDisk
        height,         // Height
        radius,         // Radius
        mgl32.Ident4(), // Model
        0,              // numCylinderVertices
        shaderManager,  // Pointer to the Shader Manager
    }
}

func (cylinder *Cylinder) MakeCylinderVBO () {
    var i uint32

    // Calculate the number of vertices required in sphere
    cylinder.numCylinderVertices = ((cylinder.VerticesPerDisk + 2) * 4)
    pVertices, pNormals := cylinder.MakeUnitCylinder()
    pColours := make([]float32, ((cylinder.VerticesPerDisk * 4) * 4))

    // Define colours as the x,y,z components of the cylinder vertices
    for i = 0; i < (cylinder.VerticesPerDisk * 4); i++ {
        pColours[i * 4] = 0.3 + pVertices[i * 2]
        pColours[i * 4 + 1] = 0.3 + pVertices[i * 2 + 1]
        pColours[i * 4 + 2] = 0.5 + pVertices[i * 2 + 2]
        pColours[i * 4 + 3] = 1.0
    }

    /* Generate the vertex buffer object */
    gl.GenBuffers(1, &cylinder.cylinderBufferObject)
    gl.BindBuffer(gl.ARRAY_BUFFER, cylinder.cylinderBufferObject)
    gl.BufferData(gl.ARRAY_BUFFER, int(8 * len(pVertices) * 3), gl.Ptr(pVertices), gl.STATIC_DRAW)
    gl.BindBuffer(gl.ARRAY_BUFFER, 0)

    /* Store the normals in a buffer object */
    gl.GenBuffers(1, &cylinder.cylinderNormals)
    gl.BindBuffer(gl.ARRAY_BUFFER, cylinder.cylinderNormals)
    gl.BufferData(gl.ARRAY_BUFFER, int(8 * len(pNormals) * 3), gl.Ptr(pNormals), gl.STATIC_DRAW)
    gl.BindBuffer(gl.ARRAY_BUFFER, 0)

    /* Store the colours in a buffer object */
    gl.GenBuffers(1, &cylinder.cylinderColours)
    gl.BindBuffer(gl.ARRAY_BUFFER, cylinder.cylinderColours)
    gl.BufferData(gl.ARRAY_BUFFER, int(8 * len(pColours) * 4), gl.Ptr(pColours), gl.STATIC_DRAW)
    gl.BindBuffer(gl.ARRAY_BUFFER, 0)

    /* Calculate the number of indices in our index array and allocate memory for it */
    numIndices := (2 * (cylinder.VerticesPerDisk + 4)) * 2
    pIndices := make([]uint32, numIndices)

    // fill "indices" to define triangle strips
    var index int = 0 // Current index

    // Define indices for the first triangle fan for one pole
    for i = 0; i < cylinder.VerticesPerDisk + 1; i++ {
        pIndices[index] = i
        index++
    }

    // Join last triangle in the triangle fan
    pIndices[index] = 1
    index++

    // Creates the Sides
    for i = 1; i < (cylinder.VerticesPerDisk * 2) + 1; i++ {
        pIndices[index] = i + cylinder.VerticesPerDisk
        index++
    }

    // Join last triangle in the triangle fan
    pIndices[index] = 1 + cylinder.VerticesPerDisk
    index++

    // Define indices for the last triangle fan for the south pole region

    // Start on a corner to avoid breaking the model
    pIndices[index] = ((cylinder.VerticesPerDisk + 3) - 2) + (cylinder.VerticesPerDisk * 3)
    index++

    // Go to Center and keep lopping till the end
    for i=(cylinder.VerticesPerDisk + 3) - 1; i >= 1; i-- {
        pIndices[index] = i + (cylinder.VerticesPerDisk * 3)
        index++
    }

    // Generate a buffer for the indices
    gl.GenBuffers(1, &cylinder.elementBuffer)
    gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, cylinder.elementBuffer)
    gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(numIndices * 4), gl.Ptr(pIndices), gl.STATIC_DRAW)
    gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (cylinder *Cylinder) MakeUnitCylinder () ([]float32, []float32) {
    var i float32
    var step float32 = 360.0 / float32(cylinder.VerticesPerDisk)

    var pos int32 = 0
    pVertices := make([]float32, (((2 + cylinder.VerticesPerDisk) * 3) * 4))
    pNormals := make([]float32, (((2 + cylinder.VerticesPerDisk) * 3) * 4))

    var x float32 = 0
    var y float32 = cylinder.Height * 0.5
    var z float32 = 0

    pVertices[pos] = x; pNormals[pos] = 0; pos++
    pVertices[pos] = y; pNormals[pos] = 1; pos++
    pVertices[pos] = z; pNormals[pos] = 0; pos++

    // Top Lid
    for i=-180; i<180; i+=step {
        var x float32 = cylinder.Radius * float32(math.Cos(float64(i * DEG_TO_RADIANS)))
        var y float32 = cylinder.Height * 0.5
        var z float32 = cylinder.Radius * float32(math.Sin(float64(i * DEG_TO_RADIANS)))

        pVertices[pos] = x; pNormals[pos] = 0; pos++
        pVertices[pos] = y; pNormals[pos] = 1; pos++
        pVertices[pos] = z; pNormals[pos] = 0; pos++
    }

    // Sides
    for i=-180; i<180; i+=step {
        var x float32 = cylinder.Radius * float32(math.Cos(float64(i * DEG_TO_RADIANS)))
        var y float32 = cylinder.Height * 0.5
        var z float32 = cylinder.Radius * float32(math.Sin(float64(i * DEG_TO_RADIANS)))

        // Circle at the top
        pVertices[pos] = x; pNormals[pos] = x; pos++
        pVertices[pos] = y; pNormals[pos] = 0; pos++
        pVertices[pos] = z; pNormals[pos] = z; pos++

        // Updates the Y position, the rest stays the same
        y = cylinder.Height * -0.5

        // Circle at the Bottom
        pVertices[pos] = x; pNormals[pos] = x; pos++
        pVertices[pos] = y; pNormals[pos] = 0; pos++
        pVertices[pos] = z; pNormals[pos] = z; pos++
    }

    // Bottom Lid
    for i=180; i>=-180; i-=step {
        var x float32 = cylinder.Radius * float32(math.Cos(float64(i * DEG_TO_RADIANS)))
        var y float32 = cylinder.Height * -0.5
        var z float32 = cylinder.Radius * float32(math.Sin(float64(i * DEG_TO_RADIANS)))

        pVertices[pos] = x; pNormals[pos] = 0; pos++
        pVertices[pos] = y; pNormals[pos] = -1; pos++
        pVertices[pos] = z; pNormals[pos] = 0; pos++
    }


    x = 0
    y = cylinder.Height * -0.5
    z = 0

    pVertices[pos] = x; pNormals[pos] = 0; pos++
    pVertices[pos] = y; pNormals[pos] = 1; pos++
    pVertices[pos] = z; pNormals[pos] = 0; pos++

    return pVertices, pNormals
}

// Draws the Cylinder form the previously defined vertex and index buffers
func (cylinder *Cylinder) Draw() {
    // Adds the Sphere Model to the Active Shader
    cylinder.ShaderManager.SetUniformMatrix4fv(cylinder.ShaderManager.ActiveShader, "model", 1, false, &cylinder.Model[0])

    /* Draw the vertices as GL_POINTS */
    gl.BindBuffer(gl.ARRAY_BUFFER, cylinder.cylinderBufferObject)
    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
    gl.EnableVertexAttribArray(0)

    /* Bind the sphere colours */
    gl.BindBuffer(gl.ARRAY_BUFFER, cylinder.cylinderColours)
    gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, nil)
    gl.EnableVertexAttribArray(1)

    /* Bind the sphere normals */
    gl.BindBuffer(gl.ARRAY_BUFFER, cylinder.cylinderNormals)
    gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 0, nil)
    gl.EnableVertexAttribArray(2)

    gl.PointSize(3.0)

    // Enable this line to show model in wireframe
    if cylinder.DrawMode == DRAW_LINES {
        gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
    } else {
        gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
    }

    if cylinder.DrawMode == DRAW_POINTS {
        gl.DrawArrays(gl.POINTS, 0, int32(cylinder.numCylinderVertices))
    } else {
        // Bind the indexed vertex buffer
        gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, cylinder.elementBuffer)

        // Draw the north pole regions as a triangle
        gl.DrawElements(gl.TRIANGLE_FAN, int32(cylinder.VerticesPerDisk + 2), gl.UNSIGNED_INT, nil)

        // Calculate offsets into the indexed array. Note that we multiply offsets by 4
        // because it is a memory offset the indices are type GLuint which is 4-bytes
        var lat_offset_jump int = int((cylinder.VerticesPerDisk * 2) + 2)
        var lat_offset_start int = int(cylinder.VerticesPerDisk + 2)
        var lat_offset_current int = lat_offset_start * 4

        // Draw the triangle strips of Sides
        gl.DrawElements(gl.TRIANGLE_STRIP, int32(cylinder.VerticesPerDisk * 2 + 2), gl.UNSIGNED_INT, gl.PtrOffset(lat_offset_current))
        lat_offset_current += (lat_offset_jump * 4)

        // Draw the south pole as a triangle fan
        gl.DrawElements(gl.TRIANGLE_FAN, int32(cylinder.VerticesPerDisk + 2), gl.UNSIGNED_INT, gl.PtrOffset(lat_offset_current))
    }
}

func (cylinder *Cylinder) ResetModel () {
    cylinder.Model = mgl32.Ident4()
}

func (cylinder *Cylinder) Translate (Tx, Ty, Tz float32) {
    cylinder.Model = cylinder.Model.Mul4(mgl32.Translate3D(Tx, Ty, Tz))
}

func (cylinder *Cylinder) Scale (scaleX, scaleY, scaleZ float32) {
    cylinder.Model = cylinder.Model.Mul4(mgl32.Scale3D(scaleX, scaleY, scaleZ))
}

func (cylinder *Cylinder) Rotate (angle float32, axis mgl32.Vec3) {
    cylinder.Model = cylinder.Model.Mul4(mgl32.HomogRotate3D(angle, axis))
}

func (cylinder *Cylinder) GetDrawMode () DrawMode {
    return cylinder.DrawMode
}

func (cylinder *Cylinder) SetDrawMode (drawMode DrawMode) {
    cylinder.DrawMode = drawMode
}

func (cylinder *Cylinder) GetName () string {
    return cylinder.Name
}

func (cylinder *Cylinder) String () string {
    return fmt.Sprintf(`
           Cylinder --> %s
    -------------------------------------
    %s
    -------------------------------------
    `, cylinder.Name, cylinder.Model)
}

