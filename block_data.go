package main

import "github.com/go-gl/mathgl/mgl32"

type BlockData struct {
	side0UV  mgl32.Vec2
	side1UV  mgl32.Vec2
	side2UV  mgl32.Vec2
	side3UV  mgl32.Vec2
	topUV    mgl32.Vec2
	bottomUV mgl32.Vec2
}

const BLOCK_AIR = 0
const BLOCK_DIRT = 1
const BLOCK_GRASS = 2
const BLOCK_STONE = 3

const BLOCK_DATA_UV_SPACE = 64.0 / 1024.0

var blockDirtData = BlockData{
	mgl32.Vec2{0, 0},
	mgl32.Vec2{0, 0},
	mgl32.Vec2{0, 0},
	mgl32.Vec2{0, 0},
	mgl32.Vec2{0, 0},
	mgl32.Vec2{0, 0},
}

var blockGrassData = BlockData{
	mgl32.Vec2{1, 0},
	mgl32.Vec2{1, 0},
	mgl32.Vec2{1, 0},
	mgl32.Vec2{1, 0},
	mgl32.Vec2{2, 0},
	mgl32.Vec2{0, 0},
}

var blockStoneData = BlockData{
	mgl32.Vec2{3, 0},
	mgl32.Vec2{3, 0},
	mgl32.Vec2{3, 0},
	mgl32.Vec2{3, 0},
	mgl32.Vec2{3, 0},
	mgl32.Vec2{3, 0},
}

var blockData = map[int]BlockData{
	BLOCK_DIRT:  blockDirtData,
	BLOCK_GRASS: blockGrassData,
	BLOCK_STONE: blockStoneData,
}
