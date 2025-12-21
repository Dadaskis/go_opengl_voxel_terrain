// Basic block data definitions for a voxel terrain generator.
// It defines block types, their UV texture coordinates for different faces,
// and a lookup map for accessing block data by type ID.
package main

import "github.com/go-gl/mathgl/mgl32"

// BlockData stores UV texture coordinates for each face of a block.
// UV coordinates are represented as multiples of BLOCK_DATA_UV_SPACE
// (i.e., tile indices in a texture atlas).
type BlockData struct {
	side0UV  mgl32.Vec2 // Texture coordinates for side 0 (typically -X or West face)
	side1UV  mgl32.Vec2 // Texture coordinates for side 1 (typically +X or East face)
	side2UV  mgl32.Vec2 // Texture coordinates for side 2 (typically -Z or North face)
	side3UV  mgl32.Vec2 // Texture coordinates for side 3 (typically +Z or South face)
	topUV    mgl32.Vec2 // Texture coordinates for the top face (+Y)
	bottomUV mgl32.Vec2 // Texture coordinates for the bottom face (-Y)
}

// Block type constants representing different materials in the game world.
const (
	BLOCK_AIR   = 0 // Invisible, non-collidable block
	BLOCK_DIRT  = 1 // Standard dirt block
	BLOCK_GRASS = 2 // Grass block with dirt sides and bottom
	BLOCK_STONE = 3 // Stone block
)

// BLOCK_DATA_UV_SPACE defines the normalized size of a single texture tile
// in the texture atlas (64px tile / 1024px atlas = 0.0625).
const BLOCK_DATA_UV_SPACE = 64.0 / 1024.0

// blockDirtData defines the texture coordinates for a dirt block.
// All faces use the same dirt texture (tile 0,0 in the atlas).
var blockDirtData = BlockData{
	side0UV:  mgl32.Vec2{0, 0},
	side1UV:  mgl32.Vec2{0, 0},
	side2UV:  mgl32.Vec2{0, 0},
	side3UV:  mgl32.Vec2{0, 0},
	topUV:    mgl32.Vec2{0, 0},
	bottomUV: mgl32.Vec2{0, 0},
}

// blockGrassData defines the texture coordinates for a grass block.
// Sides use grass side texture (tile 1,0), top uses grass top texture (tile 2,0),
// and bottom uses dirt texture (tile 0,0).
var blockGrassData = BlockData{
	side0UV:  mgl32.Vec2{1, 0},
	side1UV:  mgl32.Vec2{1, 0},
	side2UV:  mgl32.Vec2{1, 0},
	side3UV:  mgl32.Vec2{1, 0},
	topUV:    mgl32.Vec2{2, 0},
	bottomUV: mgl32.Vec2{0, 0},
}

// blockStoneData defines the texture coordinates for a stone block.
// All faces use the same stone texture (tile 3,0 in the atlas).
var blockStoneData = BlockData{
	side0UV:  mgl32.Vec2{3, 0},
	side1UV:  mgl32.Vec2{3, 0},
	side2UV:  mgl32.Vec2{3, 0},
	side3UV:  mgl32.Vec2{3, 0},
	topUV:    mgl32.Vec2{3, 0},
	bottomUV: mgl32.Vec2{3, 0},
}

// blockData is a lookup map that associates block type IDs with their
// corresponding BlockData. Note: BLOCK_AIR is intentionally omitted as
// air blocks have no visual representation.
var blockData = map[int]BlockData{
	BLOCK_DIRT:  blockDirtData,
	BLOCK_GRASS: blockGrassData,
	BLOCK_STONE: blockStoneData,
}
