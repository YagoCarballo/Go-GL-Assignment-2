//
// Wavefront Object Parser adapted from the VU 3D Engine (BSD License)
// https://github.com/gazed/vu/blob/master/load/obj.go
//
// - Heavily adapted to work with a few more use cases and models than the original version.
// - Removed dependencies to the VU engine.
//

package loader

import (
	"bufio"
	"log"
	"os"
	"strings"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/yagocarballo/Go-GL-Assignment-2/wrapper"
)

type ObjectData struct {
	Name                       string     // Data name from .obj file.
	Vertex                     []float32  // Vertex positions.    Arranged as [][3]float32
	Normals                    []float32  // Vertex normals.      Arranged as [][3]float32
	Coordinates                []float32  // Texture coordinates. Arranged as [][2]float32
	Faces                      []uint16   // Triangle faces.      Arranged as [][3]uint16

	VertexBufferObjectVertices 		uint32     // Vertex Buffer Object (Vertices)
	VertexBufferObjectNormals  		uint32     // Vertex Buffer Object (Normals)
	VertexBufferObjectFaces    		uint32     // Vertex Buffer Object (Faces)
	VertexBufferObjectTextureCoords	uint32     // Texture Coordinates Buffer Object (Texture Coordinates)

	Model                      mgl32.Mat4 // Transformation Info
	Material                   *MtlData   // Material Info
}

type Loader struct {
	Materials map[string]*MtlData
}

//
// NewLoader
// Constructor, Creates a new Loader
//
// @return loader (*Loader) a pointer to the new Loader.
//
func NewLoader () *Loader {
	return &Loader{ map[string]*MtlData{} }
}

// objStrings is an intermediate data structure used in parsing.
type objectStrings struct {
	name	string
	lines	[]string
}

// objectData is an intermediate data structure used in parsing.
// Each .obj file keeps a global count of the data below.  This is referenced
// from the face data.
type objectData struct {
	vertices	[]dataPoint // vertices
	normals		[]dataPoint // normals
	texture		[]uvPoint   // texture coordinates
	material	string		// material name
}

// dataPoint is an internal structure for passing vertices or normals.
type dataPoint struct {
	x, y, z float32
}

// uvPoint is an internal structure for passing texture coordinates.
type uvPoint struct {
	u, v float32
}

// face is an internal structure for passing face indexes.
type face struct {
	s []string // each point is a "x/y/z" value.
}

//
// Load
// Loads a .obj file into an array of objects.
//
// @param filename (string) the path to the .obj file
//
// @return objectsData ([]*ObjectData) an array of objects.
// @return error (error) the error (if any)
//
func (loader *Loader) Load (filename string) (objectsData []*ObjectData, err error) {
	objectsData = []*ObjectData{}

	file, err := os.Open(filename)

	if err != nil {
		log.Println(err)
		return objectsData, err
	}

	defer file.Close()

	objects := loader.objectToStrings(file)

	// parse each wavefront object into a mesh.
	object_data := &objectData{}
	for _, object := range objects {
		if faces, derr := loader.objectToData(object.lines, object_data); derr == nil {
			if objectData, merr := loader.objectToObjectData(object.name, object_data, faces); merr == nil {
				objectData.Material = loader.Materials[object_data.material]

				if objectData.Material != nil {
					if objectData.Material.MapBump != "" {
						var bumperr error
						objectData.Material.NormalMap, bumperr = loader.LoadTexture("resources/models/" + objectData.Material.MapBump)
						if bumperr != nil {
							return objectsData, fmt.Errorf("Bump Map %s: %s", objectData.Material.MapBump, bumperr)
						}
					}

					if objectData.Material.MapKD != "" {
						var texErr error
						// Load the texture
						objectData.Material.Texture, texErr = loader.LoadTexture("resources/models/" + objectData.Material.MapKD)
						if texErr != nil {
							return objectsData, fmt.Errorf("Texture %s: %s", objectData.Material.MapKD, texErr)
						}
					}

					if objectData.Material.MapKS != "" {
						var specErr error
						// Load the texture
						objectData.Material.SpecularMap, specErr = loader.LoadTexture("resources/models/" + objectData.Material.MapKS)
						if specErr != nil {
							return objectsData, fmt.Errorf("Specular Map %s: %s", objectData.Material.MapKS, specErr)
						}
					}
				}

				objectsData = append([]*ObjectData{ objectData }, objectsData...) // prepend

			} else {
				return objectsData, fmt.Errorf("Object To Object Data %s: %s", filename, merr)
			}
		} else {
			return objectsData, fmt.Errorf("Object To Data %s: %s", filename, derr)
		}
	}

	return
}

//
// objectToStrings
// Reads in all the file data grouped by object name. This is needed
// because a single wavefront file can hold many objects. Separating the objects
// makes parsing easier.
//
// @param file (*os.File) the file to be parsed
//
// @return objects ([]*objectStrings) an array of object strings.
//
func (loader *Loader) objectToStrings(file *os.File) (objects []*objectStrings) {
	objects = []*objectStrings{}
	name := ""
	var current *objectStrings

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, " ")

		if len(tokens) == 2 && tokens[0] == "mtllib" {
			var mtlPath string

			if _, e := fmt.Sscanf(line, "mtllib %s", &mtlPath); e != nil {
				log.Printf("Bad Materials: %s\n", line)
				log.Printf("could not parse materials %s", e)
			}

			mtlData, merr := loader.LoadMTL("resources/models/" + mtlPath)
			if merr != nil {
				log.Printf("Loading Material %s: %s", mtlPath, merr)
			}

			for _, material := range mtlData {
				loader.Materials[material.Name] = material
			}

		} else if len(tokens) == 2 && tokens[0] == "o" {
			name = strings.TrimSpace(tokens[1])
			current = &objectStrings{name, []string{}}

			objects = append(objects, current)

		} else if len(name) > 0 {
			current.lines = append(current.lines, strings.TrimSpace(line))
		}
	}

	return objects
}

//
// objectToData
// Turns a wavefront object into numbers and temporary data structures.
//
// @param lines ([]string) An array of lines to be parsed.
// @param odata (*objectData) A temporary object data pointer.
//
// @return faces ([]face) An array of faces / indices
// @return error (error) the error (if any)
//
func (loader *Loader) objectToData(lines []string, odata *objectData) (faces []face, err error) {
	for _, line := range lines {
		tokens := strings.Split(line, " ")
		var f1, f2, f3 float32
		var s1 string
//		var s1, s2, s3, s4 string
		switch tokens[0] {
		case "v":
			if _, e := fmt.Sscanf(line, "v %f %f %f", &f1, &f2, &f3); e != nil {
				log.Printf("- Obj - Bad vertex: %s\n", line)
			}
			odata.vertices = append(odata.vertices, dataPoint{f1, f2, f3})
		case "vn":
			if _, e := fmt.Sscanf(line, "vn %f %f %f", &f1, &f2, &f3); e != nil {
				log.Printf("- Obj - Bad normal: %s\n", line)
			}
			odata.normals = append(odata.normals, dataPoint{f1, f2, f3})
		case "vt":
			if _, e := fmt.Sscanf(line, "vt %f %f", &f1, &f2); e != nil {
				log.Printf("- Obj - Bad texture coord: %s\n", line)
			}
			odata.texture = append(odata.texture, uvPoint{f1, 1 - f2})
		case "f":
			faceTokens := strings.Split(line, " ")
			faces = append(faces, face{faceTokens[1:]})

//			// There might be 4 points
//			if _, e := fmt.Sscanf(line, "f %s %s %s %s", &s1, &s2, &s3, &s4); e != nil {
//				// If not there are 3
//				if _, e = fmt.Sscanf(line, "f %s %s %s", &s1, &s2, &s3); e != nil {
//					log.Printf("Bad face: %s\n", line)
//					return faces, fmt.Errorf("could not parse face %s", e)
//				}
//				faces = append(faces, face{[]string{s1, s2, s3}})
//				continue
//			} else {
//				faces = append(faces, face{[]string{s1, s2, s3, s4}})
//			}

		case "o": 		// mesh name is processed before this method is called.
		case "mtllib": 	// materials loaded separately and explicitly.
		case "usemtl": 	// material name - ignored, see above.
			if _, e := fmt.Sscanf(line, "usemtl %s", &s1); e != nil {
				log.Printf("- Obj - Bad material name: %s\n", line)
				log.Printf("- Obj - could not parse material name %s", e)
			}else {
				odata.material = s1
			}
		case "s": 		// smoothing group - ignored for now.
			if wrapper.DEBUG {
				log.Printf("- Obj - Smoothing Group not implemented: %s\n", line)
			}
		default:
			if wrapper.DEBUG {
				log.Printf("- Obj - Feature not implemented: %s\n", line)
			}
		}
	}
	return
}

//
// objectToObjectData
// Turns the data from .obj format into an internal OpenGL friendly
// format. The following information needs to be created for each mesh.
//
//    mesh.V = append(mesh.V, ...4-float32) - indexed from 0
//    mesh.N = append(mesh.N, ...3-float32) - indexed from 0
//    mesh.T = append(mesh.T, ...2-float32)	- indexed from 0
//    mesh.F = append(mesh.F, ...3-uint16)	- refers to above zero indexed values
//
// objectData holds the global vertex, texture, and normal point information.
// faces are the indexes for this mesh.
//
// Additionally the normals at each vertex are generated as the sum of the
// normals for each face that shares that vertex.
//
// @param name (string) The Name of the Object.
// @param objectData (*objectData) A temporary object data pointer.
// @param objectData (*objectData) The array of Faces / Indices.
//
// @return data (*ObjectData) A pointer to the Object.
// @return error (error) the error (if any)
//
func (loader *Loader) objectToObjectData(name string, objectData *objectData, faces []face) (data *ObjectData, err error) {
	data = &ObjectData{}
	data.Name = name
	vmap := make(map[string]int) // the unique vertex data points for this face.
	vcnt := -1

	// process each vertex of each face.  Each one represents a combination vertex,
	// texture coordinate, and normal.
	for _, face := range faces {
		for pi := 0; pi < 3; pi ++ { // Load only triangles
//		for pi, _ := range face.s {
			v, t, n := -1, -1, -1

			if len(face.s) > pi {
				faceIndex := face.s[pi]

				if v, t, n, err = parseFaceIndices(faceIndex); err != nil {
					return data, fmt.Errorf("Could not parse face data %s", err)
				}

				// cut down the amount of information passed around by reusing points
				// where the vertex and the texture coordinate information is the same.
				vertexIndex := fmt.Sprintf("%d/%d/%d", v, t, n)
				if _, ok := vmap[vertexIndex]; !ok {

					// add a new data point.
					vcnt++
					vmap[vertexIndex] = vcnt
					data.Vertex = append(data.Vertex, objectData.vertices[v].x, objectData.vertices[v].y, objectData.vertices[v].z)

					// Object might not have normals
					if n != -1 {
						data.Normals = append(data.Normals, objectData.normals[n].x, objectData.normals[n].y, objectData.normals[n].z)
					}

					// Object might not have texture information
					if t != -1 {
						data.Coordinates = append(data.Coordinates, objectData.texture[t].u, objectData.texture[t].v)
					}
				} else {
					// update the normal at the vertex to be a combination of
					// all the normals of each face that shares the vertex.
					ni := vmap[vertexIndex] * 3

					// Obj might not have normals
					if n != -1 && len(data.Normals) > (ni + 2) {
						var n1 mgl32.Vec3 = mgl32.Vec3{
							float32(data.Normals[ni]),
							float32(data.Normals[ni + 1]),
							float32(data.Normals[ni + 2]),
						}

						var n2 mgl32.Vec3 = mgl32.Vec3{
							float32(objectData.normals[n].x),
							float32(objectData.normals[n].y),
							float32(objectData.normals[n].z),
						}

						n2 = n2.Add(n1).Normalize()

						data.Normals[ni], data.Normals[ni + 1], data.Normals[ni + 2] = float32(n2.X()), float32(n2.Y()), float32(n2.Z())
					}
				}

				data.Faces = append(data.Faces, uint16(vmap[vertexIndex]))
			}
		}
	}
	return data, err
}

//
// parseFaceIndices
// Turns a face index point string (representing multiple indices)
// into 3 integer indices. The texture index is optional and is returned with
// a -1 value if it is not there.
//
// @param faceIndex (string) A string with the Face Index Line to be parsed.
//
// @return v (int) V is the reference number for a vertex in the face element. A
// minimum of three vertices are required.

// @return t (int) T is the reference number for a texture vertex in the face
// element. It always follows the first slash. (Optional)

// @return n (int) N is the reference number for a vertex normal in the face element.
// It must always follow the second slash. (Optional)

// @return error (error) the error (if any)
//
func parseFaceIndices(faceIndex string) (v, t, n int, err error) {
	v, t, n = -1, -1, -1

	// If there are TWO // then the T is not given
	if _, err = fmt.Sscanf(faceIndex, "%d//%d", &v, &n); err != nil {

		// If there are three values v/vt/vn all values are given
		if _, err = fmt.Sscanf(faceIndex, "%d/%d/%d", &v, &t, &n); err != nil {

			// If there is ONE / then the N is not given
			if _, err = fmt.Sscanf(faceIndex, "%d/%d", &v, &t); err != nil {
				return v, t, n, fmt.Errorf("Bad face (%s)\n", faceIndex)
			}
		}
	}
	v = int(v - 1)

	if n != -1 {
		n = int(n - 1) // should all have the same value.
	}

	if t != -1 {
		t = int(t - 1)
	}

	return
}

//
// String
// Implements the String function for pretty printing
//
// @return string (string) The representation of this Object as String
//
func (objectData *ObjectData) String () string {
	return fmt.Sprintf(`
	Name: %s

	Vertex Count: %d
	Normals Count: %d
	Texture Count: %d
	Faces Count: %d

	Material
	--------------------
	%s
	`, objectData.Name,
		len(objectData.Vertex),
		len(objectData.Normals),
		len(objectData.Coordinates),
		len(objectData.Faces),
		objectData.Material,
	)
}
