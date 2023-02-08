package gui

import mgl "github.com/go-gl/mathgl/mgl32"

type ViewInfo struct {
	BarsMat    mgl.Mat4
	VolumesMat mgl.Mat4
	MinPrice   float32
	MaxPrice   float32
	MaxVolume  int32
}
