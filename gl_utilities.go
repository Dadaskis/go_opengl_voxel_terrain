package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func GetTriangleMesh() Mesh {
	mesh := Mesh{}

	mesh.AddVertex(mgl32.Vec3{-0.5, -0.5, 0.0}, mgl32.Vec3{1.0, 0.0, 0.0}, mgl32.Vec3{1.0, 0.0, 0.0}, mgl32.Vec2{0.0, 0.0})
	mesh.AddVertex(mgl32.Vec3{0.5, -0.5, 0.0}, mgl32.Vec3{0.0, 1.0, 0.0}, mgl32.Vec3{1.0, 0.0, 0.0}, mgl32.Vec2{0.0, 0.0})
	mesh.AddVertex(mgl32.Vec3{0.0, 0.5, 0.0}, mgl32.Vec3{0.0, 0.0, 1.0}, mgl32.Vec3{1.0, 0.0, 0.0}, mgl32.Vec2{0.0, 0.0})
	mesh.UpdateVAO()

	return mesh
}

func GLString(str string) *uint8 {
	return gl.Str(str + "\x00")
}

func NewTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
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
