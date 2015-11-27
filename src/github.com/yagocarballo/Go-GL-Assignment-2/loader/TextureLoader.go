package loader

import (
	"fmt"
	"os"
	"image"
	"image/draw"
	_ "image/png"
	_ "image/jpeg"
	_ "image/gif"
	_ "golang.org/x/image/bmp"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/kardianos/osext"
)

func (loader *Loader) LoadTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		// Get the Folder of the current Executable
		dir, err := osext.ExecutableFolder()
		if err != nil {
			return 0, err
		}

		// Read the file and return content or error
		var secondErr error
		imgFile, secondErr = os.Open(fmt.Sprintf("%s/%s", dir, file))
		if secondErr != nil {
			return 0, secondErr
		}
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}