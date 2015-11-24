package models

import (
    "github.com/go-gl/mathgl/mgl32"
    "fmt"
)

type Camera struct {
    Name    string
    Model   mgl32.Mat4
}

func NewCamera (name string, initialModel mgl32.Mat4) *Camera {
    return &Camera{ name, initialModel }
}

func (camera *Camera) Draw () {

}

func (camera *Camera) ResetModel () {
    camera.Model = mgl32.Ident4()
}

func (camera *Camera) Translate (Tx, Ty, Tz float32) {
    camera.Model = camera.Model.Mul4(mgl32.Translate3D(Tx, Ty, Tz))
}

func (camera *Camera) Scale (scaleX, scaleY, scaleZ float32) {
    camera.Model = camera.Model.Mul4(mgl32.Scale3D(scaleX, scaleY, scaleZ))
}

func (camera *Camera) Rotate (angle float32, axis mgl32.Vec3) {
    camera.Model = camera.Model.Mul4(mgl32.HomogRotate3D(angle, axis))
}

func (camera *Camera) GetName () string {
    return camera.Name
}

func (camera *Camera) GetDrawMode () DrawMode { return DRAW_POLYGONS }
func (camera *Camera) SetDrawMode (drawMode DrawMode) {}

func (camera *Camera) String () string {
    return fmt.Sprintf(`
             Camera --> %s
    -------------------------------------
    %s
    -------------------------------------
    `, camera.Name, camera.Model)
}

