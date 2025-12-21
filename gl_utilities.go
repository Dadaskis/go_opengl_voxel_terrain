// Provides utility functions for OpenGL graphics operations.
// Includes mesh creation helpers, OpenGL string utilities, and texture loading.

package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"os"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// GetTriangleMesh creates a simple test triangle mesh for debugging and visualization.
// The triangle is centered at origin with vertices at (-0.5,-0.5), (0.5,-0.5), and (0.0,0.5).
// Each vertex has a different RGB color for easy identification.
func GetTriangleMesh() Mesh {
	mesh := Mesh{}

	// Vertex 1: Bottom-left, red color
	mesh.AddVertex(
		mgl32.Vec3{-0.5, -0.5, 0.0}, // Position
		mgl32.Vec3{1.0, 0.0, 0.0},   // Color (red)
		mgl32.Vec3{1.0, 0.0, 0.0},   // Normal (facing +X)
		mgl32.Vec2{0.0, 0.0},        // Texture coordinates
	)

	// Vertex 2: Bottom-right, green color
	mesh.AddVertex(
		mgl32.Vec3{0.5, -0.5, 0.0}, // Position
		mgl32.Vec3{0.0, 1.0, 0.0},  // Color (green)
		mgl32.Vec3{1.0, 0.0, 0.0},  // Normal (facing +X)
		mgl32.Vec2{0.0, 0.0},       // Texture coordinates
	)

	// Vertex 3: Top-center, blue color
	mesh.AddVertex(
		mgl32.Vec3{0.0, 0.5, 0.0}, // Position
		mgl32.Vec3{0.0, 0.0, 1.0}, // Color (blue)
		mgl32.Vec3{1.0, 0.0, 0.0}, // Normal (facing +X)
		mgl32.Vec2{0.0, 0.0},      // Texture coordinates
	)

	// Prepare the mesh for OpenGL rendering
	mesh.PrepareArrayData() // Organizes vertex data into arrays
	mesh.UpdateVAO()        // Creates Vertex Array Object and buffers

	return mesh
}

// GLString converts a Go string to a C-style null-terminated string
// for use with OpenGL functions that expect *uint8 pointers.
// str: The Go string to convert
// Returns: Pointer to null-terminated C string
func GLString(str string) *uint8 {
	return gl.Str(str + "\x00")
}

// NewTexture loads an image file from disk and creates an OpenGL texture.
// file: Path to the image file (supports PNG and JPEG formats)
// Returns: OpenGL texture ID and any error encountered
func NewTexture(file string) (uint32, error) {
	// Open image file
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	defer imgFile.Close()

	// Decode image using registered decoders (PNG/JPEG)
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, fmt.Errorf("failed to decode image %q: %v", file, err)
	}

	// Convert image to RGBA format (required by OpenGL)
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride: got %d, expected %d",
			rgba.Stride, rgba.Rect.Size().X*4)
	}

	// Draw source image onto RGBA canvas
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	// Generate OpenGL texture
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	// Set texture parameters for pixelated/minecraft-style look
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)   // Nearest-neighbor filtering
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)   // No interpolation
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE) // Clamp texture coordinates
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	// Upload texture data to GPU
	gl.TexImage2D(
		gl.TEXTURE_2D,             // Target
		0,                         // Mipmap level (0 = base)
		gl.RGBA,                   // Internal format
		int32(rgba.Rect.Size().X), // Width
		int32(rgba.Rect.Size().Y), // Height
		0,                         // Border (must be 0)
		gl.RGBA,                   // Format of pixel data
		gl.UNSIGNED_BYTE,          // Data type
		gl.Ptr(rgba.Pix),          // Pointer to pixel data
	)

	return texture, nil
}
