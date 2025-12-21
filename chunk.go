// Each Chunk represents a 16x16x256 region of blocks with procedurally generated terrain
// and optimized mesh generation using face culling.

package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/ojrac/opensimplex-go"
)

// Chunk represents a 16x16x256 block region in the world.
// It contains block data, a renderable mesh, and manages mesh generation.
type Chunk struct {
	position    mgl32.Vec2       // Chunk position in chunk coordinates (X,Z)
	blocks      [16][16][256]int // 3D array of block IDs (X, Y, Z) where Y is vertical
	mesh        Mesh             // Renderable mesh data for this chunk
	isMeshDirty bool             // Flag indicating if mesh needs to be regenerated
}

// Generate creates procedural terrain for the chunk using OpenSimplex noise.
// Runs asynchronously in a goroutine to prevent blocking the main thread.
// Terrain features include height-based layering (stone, dirt, grass) and caves.
func (chunk *Chunk) Generate() {
	go func() {
		// Initialize noise generator with seed 0
		noise := opensimplex.New(0)

		// Convert chunk position to world coordinates (chunks are 16 blocks wide)
		blockPos := chunk.position.Mul(16)

		// Noise scales for different terrain features
		scale := 0.01     // Scale for terrain height variation
		caveScale := 0.04 // Scale for cave generation (smaller = larger caves)

		// Generate blocks for each column in the chunk
		for x := range 16 {
			for y := range 16 {
				// Calculate height using 2D noise (creates rolling hills)
				height := 50 + noise.Eval2(
					float64(int(blockPos[0])+x)*scale,
					float64(int(blockPos[1])+y)*scale,
				)*30.0

				// Fill blocks from bottom up to calculated height
				for z := range int(height) {
					// Default to stone
					chunk.blocks[x][y][z] = BLOCK_STONE

					// Create dirt layer on top of stone
					if (int(height) - z) < 5 {
						chunk.blocks[x][y][z] = BLOCK_DIRT
					}

					// Create grass layer on very top
					if (int(height) - z) <= 1 {
						chunk.blocks[x][y][z] = BLOCK_GRASS
					}

					// Generate caves using 3D noise
					caveValue := noise.Eval3(
						float64(int(blockPos[0])+x)*caveScale,
						float64(int(blockPos[1])+y)*caveScale,
						float64(z)*caveScale,
					)

					// Create air blocks where cave noise exceeds threshold
					if caveValue > 0.6 {
						chunk.blocks[x][y][z] = BLOCK_AIR
					}
				}

				// Ensure bedrock layer at bottom (z=0)
				chunk.blocks[x][y][0] = BLOCK_STONE
			}
		}

		// Update mesh after generation completes
		chunk.UpdateMesh()
	}()
}

// UpdateMesh generates a renderable mesh from the chunk's block data.
// Implements face culling by only generating faces between air and solid blocks.
func (chunk *Chunk) UpdateMesh() {
	// Start with empty mesh
	chunk.mesh = Mesh{}

	// Convert chunk position to world coordinates for vertex positioning
	blockPos := chunk.position.Mul(16)

	// Iterate through all blocks in the chunk
	for x := range 16 {
		for y := range 16 {
			for z := range 256 {
				blockID := chunk.blocks[x][y][z]

				// Skip air blocks (no faces to render)
				if blockID == BLOCK_AIR {
					continue
				}

				// Calculate world position of this block
				vertexPos := mgl32.Vec3{
					float32(int(blockPos[0]) + x),
					float32(z),
					float32(int(blockPos[1]) + y),
				}

				// Default vertex color (white - actual coloring from textures)
				color := mgl32.Vec3{1.0, 1.0, 1.0}

				// Helper function to check if a neighboring block is occupied (non-air)
				isBlockOccupied := func(x, y, z int) bool {
					// Check bounds - if outside chunk, treat as unoccupied
					// (allows faces on chunk edges to always render)
					if x < 0 {
						return false
					}
					if x >= 16 {
						return false
					}
					if y < 0 {
						return false
					}
					if y >= 16 {
						return false
					}
					if z < 0 {
						return false
					}
					if z >= 256 {
						return false
					}

					// Block is occupied if it's not air
					return chunk.blocks[x][y][z] != BLOCK_AIR
				}

				// SIDE 0 (-Z face - typically "north" side)
				// Only generate face if neighboring block in -Z direction is air/unoccupied
				if !isBlockOccupied(x, y-1, z) {
					// Add two triangles forming a quad for this face
					// Triangle 1
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 0.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, -1.0},
						mgl32.Vec2{
							(blockData[blockID].side0UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side0UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, -1.0},
						mgl32.Vec2{
							(blockData[blockID].side0UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side0UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, -1.0},
						mgl32.Vec2{
							(blockData[blockID].side0UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side0UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)

					// Triangle 2
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, -1.0},
						mgl32.Vec2{
							(blockData[blockID].side0UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side0UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, -1.0},
						mgl32.Vec2{
							(blockData[blockID].side0UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side0UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, -1.0},
						mgl32.Vec2{
							(blockData[blockID].side0UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side0UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
				}

				// SIDE 1 (-X face - typically "west" side)
				if !isBlockOccupied(x-1, y, z) {
					// Similar pattern for -X face
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 0.0, 0.0}), color, mgl32.Vec3{-1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side1UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side1UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 0.0, 1.0}), color, mgl32.Vec3{-1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side1UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side1UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 0.0}), color, mgl32.Vec3{-1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side1UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side1UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 1.0}), color, mgl32.Vec3{-1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side1UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side1UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 0.0, 1.0}), color, mgl32.Vec3{-1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side1UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side1UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 0.0}), color, mgl32.Vec3{-1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side1UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side1UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
				}

				// SIDE 2 (+X face - typically "east" side)
				if !isBlockOccupied(x+1, y, z) {
					// Similar pattern for +X face
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, 1.0},
						mgl32.Vec2{
							(blockData[blockID].side2UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side2UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 1.0}), color, mgl32.Vec3{0.0, 0.0, 1.0},
						mgl32.Vec2{
							(blockData[blockID].side2UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side2UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, 1.0},
						mgl32.Vec2{
							(blockData[blockID].side2UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side2UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 1.0, 1.0}), color, mgl32.Vec3{0.0, 0.0, 1.0},
						mgl32.Vec2{
							(blockData[blockID].side2UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side2UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 1.0}), color, mgl32.Vec3{0.0, 0.0, 1.0},
						mgl32.Vec2{
							(blockData[blockID].side2UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side2UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 0.0, 1.0},
						mgl32.Vec2{
							(blockData[blockID].side2UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side2UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
				}

				// SIDE 3 (+Z face - typically "south" side)
				if !isBlockOccupied(x, y+1, z) {
					// Similar pattern for +Z face
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 0.0, 1.0}), color, mgl32.Vec3{1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side3UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side3UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 1.0}), color, mgl32.Vec3{1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side3UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side3UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 1.0}), color, mgl32.Vec3{1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side3UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side3UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 1.0, 1.0}), color, mgl32.Vec3{1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side3UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side3UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 1.0}), color, mgl32.Vec3{1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side3UV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side3UV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 1.0}), color, mgl32.Vec3{1.0, 0.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].side3UV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].side3UV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
				}

				// TOP face (+Y direction)
				if !isBlockOccupied(x, y, z+1) {
					// Add two triangles for top face (different winding order for top)
					// Triangle 1
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 1.0}), color, mgl32.Vec3{0.0, 1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)

					// Triangle 2
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 1.0, 0.0}), color, mgl32.Vec3{0.0, 1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 1.0, 1.0}), color, mgl32.Vec3{0.0, 1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 1.0, 1.0}), color, mgl32.Vec3{0.0, 1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
				}

				// BOTTOM face (-Y direction)
				if !isBlockOccupied(x, y, z-1) {
					// Similar pattern for bottom face (uses topUV coordinates for simplicity)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 0.0, 0.0}), color, mgl32.Vec3{0.0, -1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 0.0, 1.0}), color, mgl32.Vec3{0.0, -1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 0.0}), color, mgl32.Vec3{0.0, -1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 0.0}), color, mgl32.Vec3{0.0, -1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 0) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{0.0, 0.0, 1.0}), color, mgl32.Vec3{0.0, -1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 0) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
					chunk.mesh.AddVertex(
						vertexPos.Add(mgl32.Vec3{1.0, 0.0, 1.0}), color, mgl32.Vec3{0.0, -1.0, 0.0},
						mgl32.Vec2{
							(blockData[blockID].topUV[0] + 1) * BLOCK_DATA_UV_SPACE,
							(blockData[blockID].topUV[1] + 1) * BLOCK_DATA_UV_SPACE,
						},
					)
				}
			}
		}
	}

	// Prepare the mesh data for OpenGL rendering
	chunk.mesh.PrepareArrayData()

	// Mark mesh as dirty so VAO gets updated before next render
	chunk.isMeshDirty = true
}

// Render draws the chunk's mesh to the screen.
// Updates the VAO if the mesh has changed since last render.
func (chunk *Chunk) Render() {
	// Update VAO if mesh data has changed
	if chunk.isMeshDirty == true {
		chunk.mesh.UpdateVAO()
		chunk.isMeshDirty = false
	}

	// Render the mesh
	chunk.mesh.Render()
}
