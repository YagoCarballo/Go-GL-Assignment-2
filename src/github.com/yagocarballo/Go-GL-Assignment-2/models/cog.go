package models

import (
    "github.com/go-gl/gl/all-core/gl"
    "github.com/go-gl/mathgl/mgl32"
    "fmt"
    "math"

    "github.com/yagocarballo/Go-GL-Assignment-2/wrapper"
)

// Define buffer object indices
type Cog struct {
    Name                                                    string      // Name

    cogBufferObject, cogNormals, cogColours  uint32      // Model Buffers
    elementBuffer                                           uint32      // Indices Buffer

    DrawMode                                                DrawMode    // Defines drawing mode of cog as points, lines or filled polygons

    VerticesPerDisk                                         uint32      // Defines the Vertices Per Disk for the cog
    Height                                                  float32     // Defines the Height for the cog
    Radius                                                  float32     // Defines the Radius for the cog
    ToothSize                                               float32     // Defines the Size of the Tooth

    Model                                                   mgl32.Mat4  // Model

    numCogVertices                                          uint32      // Total count of vertices

    ShaderManager                                           *wrapper.ShaderManager // Pointer to the Shader Manager
}

func NewCog (name string, vertices uint32, height, radius, toothSize float32, shaderManager *wrapper.ShaderManager) *Cog {
    return &Cog{
        name,           // Name
        0, 0, 0,        // cogBufferObject, cogNormals, cogColours
        0,              // elementBuffer
        DRAW_POLYGONS,  // DrawMode
        vertices,       // VerticesPerDisk
        height,         // Height
        radius,         // Radius
        toothSize,      // ToothSize
        mgl32.Ident4(), // Model
        0,              // numcogVertices
        shaderManager,  // Pointer to the Shader Manager
    }
}

func (cog *Cog) MakeCogVBO () {
    var i uint32

    // Calculate the number of vertices required in sphere
    cog.numCogVertices = ((cog.VerticesPerDisk + 2) * 4)
    pVertices, pNormals := cog.MakeUnitcog()
    pColours := make([]float32, ((cog.VerticesPerDisk * 4) * 4))

    // Define colours as the x,y,z components of the cog vertices
    for i = 0; i < (cog.VerticesPerDisk * 4); i++ {
        pColours[i * 4] = 0.3 + pVertices[i * 2]
        pColours[i * 4 + 1] = 0.5 + pVertices[i * 2 + 1]
        pColours[i * 4 + 2] = 0.3 + pVertices[i * 2 + 2]
        pColours[i * 4 + 3] = 1.0
    }

    /* Generate the vertex buffer object */
    gl.GenBuffers(1, &cog.cogBufferObject)
    gl.BindBuffer(gl.ARRAY_BUFFER, cog.cogBufferObject)
    gl.BufferData(gl.ARRAY_BUFFER, int(8 * len(pVertices) * 3), gl.Ptr(pVertices), gl.STATIC_DRAW)
    gl.BindBuffer(gl.ARRAY_BUFFER, 0)

    /* Store the normals in a buffer object */
    gl.GenBuffers(1, &cog.cogNormals)
    gl.BindBuffer(gl.ARRAY_BUFFER, cog.cogNormals)
    gl.BufferData(gl.ARRAY_BUFFER, int(8 * len(pNormals) * 3), gl.Ptr(pNormals), gl.STATIC_DRAW)
    gl.BindBuffer(gl.ARRAY_BUFFER, 0)

    /* Store the colours in a buffer object */
    gl.GenBuffers(1, &cog.cogColours)
    gl.BindBuffer(gl.ARRAY_BUFFER, cog.cogColours)
    gl.BufferData(gl.ARRAY_BUFFER, int(8 * len(pColours) * 4), gl.Ptr(pColours), gl.STATIC_DRAW)
    gl.BindBuffer(gl.ARRAY_BUFFER, 0)

    /* Calculate the number of indices in our index array and allocate memory for it */
    numIndices := (2 * (cog.VerticesPerDisk + 4)) * 2
    pIndices := make([]uint32, numIndices)

    // fill "indices" to define triangle strips
    var index int = 0 // Current index

    // Define indices for the first triangle fan for one pole
    for i = 0; i < cog.VerticesPerDisk + 1; i++ {
        pIndices[index] = i
        index++
    }

    // Join last triangle in the triangle fan
    pIndices[index] = 1
    index++

    // Creates the Sides
    for i = 1; i < (cog.VerticesPerDisk * 2) + 1; i++ {
        pIndices[index] = i + cog.VerticesPerDisk
        index++
    }

    // Join last triangle in the triangle fan
    pIndices[index] = 1 + cog.VerticesPerDisk
    index++

    // Define indices for the last triangle fan for the south pole region

    // Start on a corner to avoid breaking the model
    pIndices[index] = ((cog.VerticesPerDisk + 3) - 2) + (cog.VerticesPerDisk * 3)
    index++

    // Go to Center and keep lopping till the end
    for i=(cog.VerticesPerDisk + 3) - 1; i >= 1; i-- {
        pIndices[index] = i + (cog.VerticesPerDisk * 3)
        index++
    }

    // Generate a buffer for the indices
    gl.GenBuffers(1, &cog.elementBuffer)
    gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, cog.elementBuffer)
    gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(numIndices * 4), gl.Ptr(pIndices), gl.STATIC_DRAW)
    gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (cog *Cog) MakeUnitcog () ([]float32, []float32) {
    var i float32
    var pair float64 = 0
    var step float32 = 360.0 / float32(cog.VerticesPerDisk)

    var pos int32 = 0
    pVertices := make([]float32, (((2 + cog.VerticesPerDisk) * 3) * 4))
    pNormals := make([]float32, (((2 + cog.VerticesPerDisk) * 3) * 4))

    var x float32 = 0
    var y float32 = cog.Height * 0.5
    var z float32 = 0

    pVertices[pos] = x; pNormals[pos] = 0; pos++
    pVertices[pos] = y; pNormals[pos] = 1; pos++
    pVertices[pos] = z; pNormals[pos] = 0; pos++

    // Top Lid
    for i=-180; i<180; i+=step {
        radius := cog.Radius
        if math.Mod(pair, 2) == 0 {
            radius = cog.ToothSize
        }

        var x float32 = radius * float32(math.Cos(float64(i * DEG_TO_RADIANS)))
        var y float32 = cog.Height * 0.5
        var z float32 = radius * float32(math.Sin(float64(i * DEG_TO_RADIANS)))

        pVertices[pos] = x; pNormals[pos] = 0; pos++
        pVertices[pos] = y; pNormals[pos] = 1; pos++
        pVertices[pos] = z; pNormals[pos] = 0; pos++

        pair++
    }

    // Sides
    for i=-180; i<180; i+=step {
        radius := cog.Radius
        if math.Mod(pair, 2) == 0 {
            radius = cog.ToothSize
        }

        var x float32 = radius * float32(math.Cos(float64(i * DEG_TO_RADIANS)))
        var y float32 = cog.Height * 0.5
        var z float32 = radius * float32(math.Sin(float64(i * DEG_TO_RADIANS)))

        // Circle at the top
        pVertices[pos] = x; pNormals[pos] = x; pos++
        pVertices[pos] = y; pNormals[pos] = 0; pos++
        pVertices[pos] = z; pNormals[pos] = z; pos++

        // Updates the Y position, the rest stays the same
        y = cog.Height * -0.5

        // Circle at the Bottom
        pVertices[pos] = x; pNormals[pos] = x; pos++
        pVertices[pos] = y; pNormals[pos] = 0; pos++
        pVertices[pos] = z; pNormals[pos] = z; pos++

        pair++
    }

    // Bottom Lid
    for i=180; i>=-180; i-=step {
        radius := cog.Radius
        if math.Mod(pair, 2) == 0 {
            radius = cog.ToothSize
        }

        var x float32 = radius * float32(math.Cos(float64(i * DEG_TO_RADIANS)))
        var y float32 = cog.Height * -0.5
        var z float32 = radius * float32(math.Sin(float64(i * DEG_TO_RADIANS)))

        pVertices[pos] = x; pNormals[pos] = 0; pos++
        pVertices[pos] = y; pNormals[pos] = -1; pos++
        pVertices[pos] = z; pNormals[pos] = 0; pos++

        pair++
    }


    x = 0
    y = cog.Height * -0.5
    z = 0

    pVertices[pos] = x; pNormals[pos] = 0; pos++
    pVertices[pos] = y; pNormals[pos] = 1; pos++
    pVertices[pos] = z; pNormals[pos] = 0; pos++

    return pVertices, pNormals
}

// Draws the cog form the previously defined vertex and index buffers
func (cog *Cog) Draw() {
    // Adds the Sphere Model to the Active Shader
    cog.ShaderManager.SetUniformMatrix4fv(cog.ShaderManager.ActiveShader, "model", 1, false, &cog.Model[0])

    /* Draw the vertices as GL_POINTS */
    gl.BindBuffer(gl.ARRAY_BUFFER, cog.cogBufferObject)
    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
    gl.EnableVertexAttribArray(0)

    /* Bind the sphere colours */
    gl.BindBuffer(gl.ARRAY_BUFFER, cog.cogColours)
    gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, nil)
    gl.EnableVertexAttribArray(1)

    /* Bind the sphere normals */
    gl.BindBuffer(gl.ARRAY_BUFFER, cog.cogNormals)
    gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 0, nil)
    gl.EnableVertexAttribArray(2)

    gl.PointSize(3.0)

    // Enable this line to show model in wireframe
    if cog.DrawMode == DRAW_LINES {
        gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
    } else {
        gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
    }

    if cog.DrawMode == DRAW_POINTS {
        gl.DrawArrays(gl.POINTS, 0, int32(cog.numCogVertices))
    } else {
        // Bind the indexed vertex buffer
        gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, cog.elementBuffer)

        // Draw the north pole regions as a triangle
        gl.DrawElements(gl.TRIANGLE_FAN, int32(cog.VerticesPerDisk + 2), gl.UNSIGNED_INT, nil)

        // Calculate offsets into the indexed array. Note that we multiply offsets by 4
        // because it is a memory offset the indices are type GLuint which is 4-bytes
        var lat_offset_jump int = int((cog.VerticesPerDisk * 2) + 2)
        var lat_offset_start int = int(cog.VerticesPerDisk + 2)
        var lat_offset_current int = lat_offset_start * 4

        // Draw the triangle strips of Sides
        gl.DrawElements(gl.TRIANGLE_STRIP, int32(cog.VerticesPerDisk * 2 + 2), gl.UNSIGNED_INT, gl.PtrOffset(lat_offset_current))
        lat_offset_current += (lat_offset_jump * 4)

        // Draw the south pole as a triangle fan
        gl.DrawElements(gl.TRIANGLE_FAN, int32(cog.VerticesPerDisk + 2), gl.UNSIGNED_INT, gl.PtrOffset(lat_offset_current))
    }
}

func (cog *Cog) ResetModel () {
    cog.Model = mgl32.Ident4()
}

func (cog *Cog) Translate (Tx, Ty, Tz float32) {
    cog.Model = cog.Model.Mul4(mgl32.Translate3D(Tx, Ty, Tz))
}

func (cog *Cog) Scale (scaleX, scaleY, scaleZ float32) {
    cog.Model = cog.Model.Mul4(mgl32.Scale3D(scaleX, scaleY, scaleZ))
}

func (cog *Cog) Rotate (angle float32, axis mgl32.Vec3) {
    cog.Model = cog.Model.Mul4(mgl32.HomogRotate3D(angle, axis))
}

func (cog *Cog) GetDrawMode () DrawMode {
    return cog.DrawMode
}

func (cog *Cog) SetDrawMode (drawMode DrawMode) {
    cog.DrawMode = drawMode
}

func (cog *Cog) GetName () string {
    return cog.Name
}

func (cog *Cog) String () string {
    return fmt.Sprintf(`
           Cog --> %s
    -------------------------------------
    %s
    -------------------------------------
    `, cog.Name, cog.Model)
}

