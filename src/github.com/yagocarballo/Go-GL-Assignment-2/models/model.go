package models

import (
    "github.com/go-gl/mathgl/mgl32"
)

const DEG_TO_RADIANS = 3.141592 / 180.0

type Model interface {
    ResetModel()
    Translate(float32, float32, float32)
    Scale(float32, float32, float32)
    Rotate(float32, mgl32.Vec3)

    GetName() string
    GetDrawMode() DrawMode
    SetDrawMode(DrawMode)
}
