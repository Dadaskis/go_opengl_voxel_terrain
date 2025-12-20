package main

import "github.com/go-gl/mathgl/mgl32"

type GameWorld struct {
	chunks map[mgl32.Vec2]Chunk
}

func (gameWorld *GameWorld) Initialize() {
	gameWorld.chunks = make(map[mgl32.Vec2]Chunk)
	for x := -4; x < 4; x++ {
		for y := -4; y < 4; y++ {
			position := mgl32.Vec2{float32(x), float32(y)}
			chunk := Chunk{}
			chunk.position = position
			chunk.Generate()
			chunk.UpdateMesh()
			gameWorld.chunks[position] = chunk
		}
	}
}

func (gameWorld *GameWorld) Render() {
	for x := -4; x < 4; x++ {
		for y := -4; y < 4; y++ {
			position := mgl32.Vec2{float32(x), float32(y)}
			chunk := gameWorld.chunks[position]
			chunk.Render()
		}
	}
}
