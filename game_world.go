package main

import (
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type GameWorld struct {
	chunks                     map[mgl32.Vec2]*Chunk
	renderChunks               []*Chunk
	currentCamera              *Camera
	renderDistance             int
	closeCameraMovementRoutine chan bool
}

func (gameWorld *GameWorld) Initialize() {
	gameWorld.renderDistance = 16
	gameWorld.chunks = make(map[mgl32.Vec2]*Chunk)
	gameWorld.closeCameraMovementRoutine = gameWorld.ProcessCameraMovementRoutine()
}

func (gameWorld *GameWorld) ProcessCameraMovementRoutine() chan bool {
	closeChan := make(chan bool)
	go func() {
		ticker := time.NewTicker(time.Millisecond * 300)
		prevXPos := 1110
		prevYPos := 1110
		for {
			select {
			case <-ticker.C:
				xPos := int(math.Round(float64(gameWorld.currentCamera.position[0] / 16.0)))
				yPos := int(math.Round(float64(gameWorld.currentCamera.position[2] / 16.0)))
				if xPos != prevXPos || yPos != prevYPos {
					prevXPos = xPos
					prevYPos = yPos
					newRenderChunks := []*Chunk{}
					for x := xPos + (-gameWorld.renderDistance); x < xPos+gameWorld.renderDistance; x++ {
						for y := yPos + (-gameWorld.renderDistance); y < yPos+gameWorld.renderDistance; y++ {
							position := mgl32.Vec2{float32(x), float32(y)}
							chunk, exists := gameWorld.chunks[position]
							if !exists {
								chunk = &Chunk{}
								chunk.position = position
								chunk.Generate()
								gameWorld.chunks[position] = chunk
							}
							newRenderChunks = append(newRenderChunks, chunk)
						}
					}
					gameWorld.renderChunks = newRenderChunks
				}
			case <-closeChan:
				return
			}
		}
	}()

	return closeChan
}

func (gameWorld *GameWorld) Render() {
	for _, chunk := range gameWorld.renderChunks {
		chunk.Render()
	}
}
