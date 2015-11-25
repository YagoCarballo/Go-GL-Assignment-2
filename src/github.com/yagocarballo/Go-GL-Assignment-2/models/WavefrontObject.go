package models

import (
	"log"
	"fmt"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/yagocarballo/Go-GL-Assignment-2/loader"
	"github.com/yagocarballo/Go-GL-Assignment-2/wrapper"
)

type WavefrontObject struct {
	Name						string

	VertexCoordinates 			uint32
	VertexNormals 				uint32

	Objects						[]*loader.ObjectData

	DrawMode					DrawMode
}

func NewObjectLoader () *WavefrontObject {
	return &WavefrontObject{
		"Obj", // Name

		0, // VertexCoordinates
		1, // VertexNormals

		[]*loader.ObjectData{}, // Objects

		DRAW_POLYGONS, // Draw Mode
	}
}

func (objectLoader *WavefrontObject) LoadObject (filename string) {
	load := loader.NewLoader()
	objects, err := load.Load(filename)

	log.Printf("Loaded %d Objects. \n", len(objects))

	if err != nil {
		log.Println(err)
		return
	}

	objectLoader.Objects = objects
}

func (objectLoader *WavefrontObject) CreateObject () {
	for _, object := range objectLoader.Objects {
		// Sets the Model in the Initial position
		object.Model = mgl32.Ident4()

		if wrapper.DEBUG {
			// Print the object
			fmt.Println(object)
		}

		// Generate the vertex buffer object
		gl.GenBuffers(1, &object.VertexBufferObjectVertices)
		gl.BindBuffer(gl.ARRAY_BUFFER, object.VertexBufferObjectVertices)
		gl.BufferData(gl.ARRAY_BUFFER, int(len(object.Vertex)*4), gl.Ptr(&(object.Vertex[0])), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0);

		// Obj might not have normals
		if len(object.Normals) != 0 {
			// Store the normals in a buffer object
			gl.GenBuffers(1, &object.VertexBufferObjectNormals)
			gl.BindBuffer(gl.ARRAY_BUFFER, object.VertexBufferObjectNormals)
			gl.BufferData(gl.ARRAY_BUFFER, int(len(object.Normals) * 4), gl.Ptr(&(object.Normals[0])), gl.STATIC_DRAW)
			gl.BindBuffer(gl.ARRAY_BUFFER, 0);
		}

		// Generate a buffer for the indices
		gl.GenBuffers(1, &object.VertexBufferObjectFaces)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, object.VertexBufferObjectFaces)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(len(object.Faces)*3), gl.Ptr(&(object.Faces[0])), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0);

		if len(object.Coordinates) != 0 {
			// Generate a buffer for the Texture Coordinates
			gl.GenBuffers(1, &object.VertexBufferObjectTextureCoords)
			gl.BindBuffer(gl.ARRAY_BUFFER, object.VertexBufferObjectTextureCoords)
			gl.BufferData(gl.ARRAY_BUFFER, int(len(object.Coordinates) * 2) * 2, gl.Ptr(&(object.Coordinates[0])), gl.STATIC_DRAW)
			gl.BindBuffer(gl.ARRAY_BUFFER, 0);
		}
	}
}

func (objectLoader *WavefrontObject) DrawObject(shaderProgram uint32) {
	for _, object := range objectLoader.Objects {
		// Reads the uniform Locations
		modelUniform := gl.GetUniformLocation(shaderProgram, gl.Str("model\x00"));
		ambientUniform := gl.GetUniformLocation(shaderProgram, gl.Str("ambient\x00"));
		diffuseUniform := gl.GetUniformLocation(shaderProgram, gl.Str("diffuse\x00"));
		specularUniform := gl.GetUniformLocation(shaderProgram, gl.Str("specular\x00"));
		emissiveUniform := gl.GetUniformLocation(shaderProgram, gl.Str("emissive\x00"));

		// Send our uniforms variables to the currently bound shader
		if object.Material != nil {
			gl.Uniform4f(ambientUniform, object.Material.KaR, object.Material.KaG, object.Material.KaB, object.Material.Tr); // Ambient colour.
			gl.Uniform4f(diffuseUniform, object.Material.KdR, object.Material.KdG, object.Material.KdB, object.Material.Tr); // Diffuse colour.
			gl.Uniform4f(specularUniform, object.Material.KsR, object.Material.KsG, object.Material.KsB, object.Material.Tr); // Specular colour.
			gl.Uniform4f(emissiveUniform, object.Material.KeR, object.Material.KeG, object.Material.KeB, object.Material.Tr); // Emissive colour.

			if object.Texture != 0 {
				textureUniform := gl.GetUniformLocation(shaderProgram, gl.Str("DiffuseTextureSampler\x00"))
				gl.Uniform1i(textureUniform, 0)

				normalTextureUniform := gl.GetUniformLocation(shaderProgram, gl.Str("NormalTextureSampler\x00"))
				gl.Uniform1i(normalTextureUniform, 1)

				specularTextureUniform := gl.GetUniformLocation(shaderProgram, gl.Str("SpecularTextureSampler\x00"))
				gl.Uniform1i(specularTextureUniform, 2)

				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(gl.TEXTURE_2D, object.Texture)

				gl.ActiveTexture(gl.TEXTURE1)
				gl.BindTexture(gl.TEXTURE_2D, object.NormalMap)

				gl.ActiveTexture(gl.TEXTURE2)
				gl.BindTexture(gl.TEXTURE_2D, object.SpecularMap)
			}
		}

		// Geometry
		var size int32    // Used to get the byte size of the element (vertex index) array

		gl.UniformMatrix4fv(modelUniform, 1, false, &object.Model[0]);

		// Get the vertices uniform position
		verticesUniform := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("position\x00")))
		normalsUniform := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("normal\x00")))
		textureCoordinatesUniform := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("texcoord\x00")))

		// Describe our vertices array to OpenGL (it can't guess its format automatically)

		gl.BindBuffer(gl.ARRAY_BUFFER, object.VertexBufferObjectVertices);
		gl.VertexAttribPointer(
			verticesUniform,				// attribute index
			3,								// number of elements per vertex, here (x,y,z)
			gl.FLOAT,						// the type of each element
			false,							// take our values as-is
			0,								// no extra data between each position
			nil,							// offset of first element
		)

		gl.EnableVertexAttribArray(normalsUniform)
		gl.BindBuffer(gl.ARRAY_BUFFER, object.VertexBufferObjectNormals);
		gl.VertexAttribPointer(
			normalsUniform,				// attribute
			3, 							// number of elements per vertex, here (x,y,z)
			gl.FLOAT,					// the type of each element
			false,						// take our values as-is
			0,							// no extra data between each position
			nil,						// offset of first element
		)

		gl.EnableVertexAttribArray(textureCoordinatesUniform)
		gl.BindBuffer(gl.ARRAY_BUFFER, object.VertexBufferObjectTextureCoords);
		gl.VertexAttribPointer(
			textureCoordinatesUniform,	// attribute
			2, 							// number of elements per vertex, here (u,v)
			gl.FLOAT,					// the type of each element
			false,						// take our values as-is
			0,							// no extra data between each position
			nil,						// offset of first element
		)

		size = int32(len(object.Vertex))

		gl.PointSize(3.0)

		// Enable this line to show model in wireframe
		switch objectLoader.DrawMode {
		case 1:
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		default:
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		}

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, object.VertexBufferObjectFaces);
		gl.GetBufferParameteriv(gl.ELEMENT_ARRAY_BUFFER, gl.BUFFER_SIZE, &size);
		gl.DrawElements(gl.TRIANGLES, int32(len(object.Faces)), gl.UNSIGNED_SHORT, nil)
//		gl.DrawElements(gl.POINTS, int32(len(object.Faces)), gl.UNSIGNED_SHORT, nil)
	}
}


// Individual Objects

func (objectLoader *WavefrontObject) ResetChildModel(index int) {
	objectLoader.Objects[index].Model = mgl32.Ident4()
}

func (objectLoader *WavefrontObject) TranslateChild(index int, Tx, Ty, Tz float32) {
	objectLoader.Objects[index].Model = objectLoader.Objects[index].Model.Mul4(mgl32.Translate3D(Tx, Ty, Tz))
}

func (objectLoader *WavefrontObject) ScaleChild(index int, scaleX, scaleY, scaleZ float32) {
	objectLoader.Objects[index].Model = objectLoader.Objects[index].Model.Mul4(mgl32.Scale3D(scaleX, scaleY, scaleZ))
}

func (objectLoader *WavefrontObject) RotateChild(index int, angle float32, axis mgl32.Vec3) {
	objectLoader.Objects[index].Model = objectLoader.Objects[index].Model.Mul4(mgl32.HomogRotate3D(angle, axis))
}

// All Objects

func (objectLoader *WavefrontObject) ResetModel() {
	for index, _ := range objectLoader.Objects {
		objectLoader.Objects[index].Model = mgl32.Ident4()
	}
}

func (objectLoader *WavefrontObject) Translate(Tx, Ty, Tz float32) {
	for index, _ := range objectLoader.Objects {
		objectLoader.Objects[index].Model = objectLoader.Objects[index].Model.Mul4(mgl32.Translate3D(Tx, Ty, Tz))
	}
}

func (objectLoader *WavefrontObject) Scale(scaleX, scaleY, scaleZ float32) {
	for index, _ := range objectLoader.Objects {
		objectLoader.Objects[index].Model = objectLoader.Objects[index].Model.Mul4(mgl32.Scale3D(scaleX, scaleY, scaleZ))
	}
}

func (objectLoader *WavefrontObject) Rotate(angle float32, axis mgl32.Vec3) {
	for index, _ := range objectLoader.Objects {
		objectLoader.Objects[index].Model = objectLoader.Objects[index].Model.Mul4(mgl32.HomogRotate3D(angle, axis))
	}
}

func (objectLoader *WavefrontObject) GetDrawMode () DrawMode {
	return objectLoader.DrawMode
}

func (objectLoader *WavefrontObject) SetDrawMode (drawMode DrawMode) {
	objectLoader.DrawMode = drawMode
}

func (objectLoader *WavefrontObject) GetName () string {
	return objectLoader.Name
}

func (objectLoader *WavefrontObject) String () string {
	return fmt.Sprintf(``)
}

