//
// Wavefront Mtl Parser adapted from the VU 3D Engine (BSD License)
// https://github.com/gazed/vu/blob/master/load/mtl.go
//
// - Added Support for Textures, Normal Bump Maping, Specular Maping,
// 	 Illumination, Optical Density, Emissive Color and Name
// - Removed dependencies to VU project
//

package loader

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"log"
	"os"
	"github.com/yagocarballo/Go-GL-Assignment-2/wrapper"
	"github.com/kardianos/osext"
)

// MtlData holds colour and alpha information.
// It is intended for populating rendered models.
type MtlData struct {
	Name		  string	// Material Name
	KaR, KaG, KaB float32	// Ambient colour.
	KdR, KdG, KdB float32	// Diffuse colour.
	KsR, KsG, KsB float32	// Specular colour.
	KeR, KeG, KeB float32	// Emissive color.
	Tr            float32	// Transparency
	Ni			  float32	// Optical Density (Scaler)
	Illum		  int32		// Illumination model
	MapKD		  string	// Map Texture
	MapKS		  string	// Map Specular
	MapBump		  string	// Map Normals

	Texture		  uint32	  // Texture Pointer
	NormalMap	  uint32	  // Normal Map Texture Pointer
	SpecularMap	  uint32	  // Specular Map Texture Pointer
}

// Load a Wavefront .mtl file which is a text representation of one
// or more material descriptions.  See the file format specification at:
//    https://en.wikipedia.org/wiki/Wavefront_.obj_file#File_format
//    http://web.archive.org/web/20080813073052/
//    http://paulbourke.net/dataformats/mtl/
func (loader *Loader) LoadMTL(filename string) (data []*MtlData, err error) {
	log.Printf("Loading material: '%s'", filename)

	pos := -1
	materials := []*MtlData{}

	file, err := os.Open(filename)
	if err != nil {
		// Get the Folder of the current Executable
		dir, err := osext.ExecutableFolder()
		if err != nil {
			log.Println(err)
			return materials, fmt.Errorf("could not open %s %s", filename, err)
		}

		// Read the file and return content or error
		var secondErr error
		file, secondErr = os.Open(fmt.Sprintf("%s/%s", dir, file))
		if secondErr != nil {
			log.Println(secondErr)
			return materials, fmt.Errorf("could not open %s %s", filename, secondErr)
		}
	}

	var f1, f2, f3 float32
	var ni float32
	var illum int32

	reader := bufio.NewReader(file)
	line, e1 := reader.ReadString('\n')

	for ; e1 == nil; line, e1 = reader.ReadString('\n') {
		tokens := strings.Split(line, " ")

		// If line is empty, Ignore
		if line == "\n" {
			continue
		}

		switch tokens[0] {
		case "Ka": // ambient
			if _, e := fmt.Sscanf(line, "Ka %f %f %f", &f1, &f2, &f3); e != nil {
				log.Printf("- Mtl - could not parse ambient values %s", e)
			}
			materials[pos].KaR, materials[pos].KaG, materials[pos].KaB = f1, f2, f3
		case "Kd": // diffuse
			if _, e := fmt.Sscanf(line, "Kd %f %f %f", &f1, &f2, &f3); e != nil {
				log.Printf("- Mtl - could not parse diffuse values %s", e)
			}
			materials[pos].KdR, materials[pos].KdG, materials[pos].KdB = f1, f2, f3
		case "Ks": // specular
			if _, e := fmt.Sscanf(line, "Ks %f %f %f", &f1, &f2, &f3); e != nil {
				log.Printf("- Mtl - could not parse specular values %s", e)
			}
			materials[pos].KsR, materials[pos].KsG, materials[pos].KsB = f1, f2, f3
		case "Ke": // Emissive color
			if _, e := fmt.Sscanf(line, "Ke %f %f %f", &f1, &f2, &f3); e != nil {
				log.Printf("- Mtl - could not parse the emissive color values %s", e)
			}
			materials[pos].KeR, materials[pos].KeG, materials[pos].KeB = f1, f2, f3
		case "d": // transparency
			a, _ := strconv.ParseFloat(strings.TrimSpace(tokens[1]), 32)
			materials[pos].Tr = float32(a)
		case "newmtl": // material name
			pos++
			materials = append(materials, &MtlData{
				"",
				0.0, 0.0, 0.0,
				0.0, 0.0, 0.0,
				0.0, 0.0, 0.0,
				0.0, 0.0, 0.0,
				1.0,
				1.0,
				1,
				"",
				"",
				"",
				0,
				0,
				0,
			})


			if _, e := fmt.Sscanf(line, "newmtl %s", &materials[pos].Name); e != nil {
				log.Printf("- Mtl - could not parse the material name %s \n", e)
			}
		case "Ns": // specular exponent - scaler. Ignored for now.
		case "Ni": // optical density - scaler.
			if _, e := fmt.Sscanf(line, "Ni %f", &ni); e != nil {
				log.Printf("- Mtl - could not parse the optical density %s", e)
			}
			materials[pos].Ni = ni
		case "illum": // illumination model - int.
			if _, e := fmt.Sscanf(line, "illum %d", &illum); e != nil {
				log.Printf("- Mtl - could not parse the illumination model %s", e)
			}
			materials[pos].Illum = illum
		case "map_Kd": // Map Texture
			if _, e := fmt.Sscanf(line, "map_Kd %s", &materials[pos].MapKD); e != nil {
				log.Printf("- Mtl - could not parse the map kd %s", e)
			}
		case "map_Ks": // Specular color texture map
			if _, e := fmt.Sscanf(line, "map_Ks %s", &materials[pos].MapKS); e != nil {
				log.Printf("- Mtl - could not parse the map kd %s", e)
			}
		case "map_Bump": // Map Normals
			if _, e := fmt.Sscanf(line, "map_Bump %s", &materials[pos].MapBump); e != nil {
				log.Printf("- Mtl - could not parse the map bump %s", e)
			}
		case "#": // Comment, Ignore
		default:
			if wrapper.DEBUG {
				log.Printf("- Mtl - Feature not implemented: %s \n", line)
			}
		}
	}

	log.Printf("Loaded %d materials.", len(materials))

	return materials, nil
}

//
// String
// Implements the String function for pretty printing
//
// @return string (string) The representation of this Object as String
//
func (material *MtlData) String () string {
	return fmt.Sprintf(`
	Name: %s
	Ambient colour: { %f, %f, %f }
	Diffuse colour: { %f, %f, %f }
	Specular colour: { %f, %f, %f }
	Emissive colour: { %f, %f, %f }
	Transparency: %f
	Optical Density: %f
	Illumination: %d
	Map KD: %s
	Map Bump: %s
	Map Specular: %s
	`, material.Name,
		material.KaR, material.KaG, material.KaB,
		material.KdR, material.KdG, material.KdB,
		material.KsR, material.KsG, material.KsB,
		material.KeR, material.KeG, material.KeB,
		material.Tr,
		material.Ni,
		material.Illum,
		material.MapKD,
		material.MapBump,
		material.MapKS,
	)
}

