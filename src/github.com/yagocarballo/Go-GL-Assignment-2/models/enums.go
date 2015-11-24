package models

type DrawMode int32
type ColorMode int32
type EmitMode uint32

const (
	_ = iota // ignore first value by assigning to blank identifier
	DRAW_POINTS DrawMode = 0 + iota
	DRAW_LINES
	DRAW_POLYGONS
)

const (
	_ = iota // ignore first value by assigning to blank identifier
	COLOR_PER_SIDE ColorMode = -1 + iota
	COLOR_SOLID
)

const (
	_ = iota // ignore first value by assigning to blank identifier
    EMIT_BRIGHT EmitMode = 0 + iota
    EMIT_COLORED
)


var drawModeNames = [...]string{
	"_",
	"[ Draw Points ]",
	"[ Draw Lines ]",
	"[ Draw Polygons ]",
}

var colorModeNames = [...]string{
	"[ Color per side ]",
	"[ Solid Color ]",
}

var emitModeNames = [...]string{
	"[ Emit Colored ]",
	"[ Emit Bright ]",
}

func (drawMode DrawMode) String() string {
	return drawModeNames[drawMode]
}

func (colorMode ColorMode) String() string {
	return colorModeNames[colorMode]
}

func (emitMode EmitMode) String() string {
    return emitModeNames[emitMode]
}

func (emitMode EmitMode) AsUint32() uint32 {
    return uint32(emitMode)
}
