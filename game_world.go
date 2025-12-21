// Implements world management for a voxel terrain generator.
// The GameWorld struct handles chunk generation, loading, and rendering
// based on player position with dynamic view distance management.

package main

import (
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

// GameWorld manages all chunks in the game world and handles dynamic
// chunk loading/unloading based on camera position.
type GameWorld struct {
	chunks                     map[mgl32.Vec2]*Chunk // Map of all loaded chunks keyed by their position
	renderChunks               []*Chunk              // Subset of chunks currently within render distance
	currentCamera              *Camera               // Reference to the active camera for position tracking
	renderDistance             int                   // Number of chunks to render in each direction from camera
	closeCameraMovementRoutine chan bool             // Channel to signal shutdown of the camera tracking goroutine
}

// Initialize sets up the game world with default values and starts
// the camera tracking goroutine for dynamic chunk loading.
func (gameWorld *GameWorld) Initialize() {
	gameWorld.renderDistance = 16 // Render 16 chunks in each direction (32x32 chunk area)
	gameWorld.chunks = make(map[mgl32.Vec2]*Chunk)

	// Start goroutine that monitors camera position and loads/unloads chunks
	gameWorld.closeCameraMovementRoutine = gameWorld.ProcessCameraMovementRoutine()
}

// ProcessCameraMovementRoutine starts a background goroutine that periodically
// checks camera position and updates which chunks should be rendered.
// Returns a channel that can be closed to stop the goroutine.
func (gameWorld *GameWorld) ProcessCameraMovementRoutine() chan bool {
	closeChan := make(chan bool)

	go func() {
		// Check camera position every 300ms (balances responsiveness with performance)
		ticker := time.NewTicker(time.Millisecond * 300)

		// Track previous camera chunk position to detect movement
		prevXPos := 1110 // Initial sentinel value that won't match any real position
		prevYPos := 1110

		for {
			select {
			case <-ticker.C:
				// Convert camera world position to chunk coordinates
				xPos := int(math.Round(float64(gameWorld.currentCamera.position[0] / 16.0)))
				yPos := int(math.Round(float64(gameWorld.currentCamera.position[2] / 16.0)))

				// Only update render list if camera moved to a new chunk
				if xPos != prevXPos || yPos != prevYPos {
					prevXPos = xPos
					prevYPos = yPos

					newRenderChunks := []*Chunk{}

					// Calculate bounding box of chunks to render based on render distance
					for x := xPos + (-gameWorld.renderDistance); x < xPos+gameWorld.renderDistance; x++ {
						for y := yPos + (-gameWorld.renderDistance); y < yPos+gameWorld.renderDistance; y++ {
							position := mgl32.Vec2{float32(x), float32(y)}

							// Check if chunk already exists in memory
							chunk, exists := gameWorld.chunks[position]
							if !exists {
								// Create and generate new chunk
								chunk = &Chunk{}
								chunk.position = position
								chunk.Generate() // Starts async generation
								gameWorld.chunks[position] = chunk
							}

							// Add chunk to render list
							newRenderChunks = append(newRenderChunks, chunk)
						}
					}

					// Update render list atomically
					gameWorld.renderChunks = newRenderChunks
				}

			case <-closeChan:
				// Stop the goroutine when close channel receives signal
				ticker.Stop()
				return
			}
		}
	}()

	return closeChan
}

// Render draws all chunks currently within render distance.
// Called each frame from the main game loop.
func (gameWorld *GameWorld) Render() {
	for _, chunk := range gameWorld.renderChunks {
		chunk.Render()
	}
}
