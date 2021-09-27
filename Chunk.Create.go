package main

import (
	"math"
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

func (store *V3_ChunkStore) CreateChunkAroundIfVisible(frustrum FrustrumPlanes, circleVision float64, position mgl32.Vec3, chunkId Vec2i) {
	px := chunkId[0] * chunkSize
	pz := chunkId[1] * chunkSize

	inFrustrum := InsideFrustrum(frustrum, -px, -px+chunkSize, -pz, -pz+chunkSize)
	_, existsInQueue := store.chunks[chunkId]

	if inFrustrum && getDistance(position, chunkId) <= circleVision && !existsInQueue {
		store.chunks[chunkId] = V3_ChunksList{flag: V3_C_WAITING}
	}
}

func (store *V3_ChunkStore) CreateNewChunks(chunkId Vec2i) {

	store.mutexPosition.RLock()
	frustrum := store.frustrum
	position := store.position
	circleVision := getTheoricalMaxCircleVision(store.position, store.rotation)
	store.mutexPosition.RUnlock()

	store.mutexChunks.Lock()
	store.CreateChunkAroundIfVisible(frustrum, circleVision, position, Vec2i{chunkId[0] + 1, chunkId[1]})
	store.CreateChunkAroundIfVisible(frustrum, circleVision, position, Vec2i{chunkId[0], chunkId[1] + 1})
	store.CreateChunkAroundIfVisible(frustrum, circleVision, position, Vec2i{chunkId[0] - 1, chunkId[1]})
	store.CreateChunkAroundIfVisible(frustrum, circleVision, position, Vec2i{chunkId[0], chunkId[1] - 1})
	store.mutexChunks.Unlock()
}

func CreateFilledChunk(chunkId Vec2i, chunk *V3_ChunksList) {
	var blockType uint8
	px := chunkId[0] * chunkSize
	pz := chunkId[1] * chunkSize

	for x := uint8(0); x < chunkSize; x++ {
		for z := uint8(0); z < chunkSize; z++ {

			voxelValueY := Biome(
				chunk,
				float64(px-int(x)),
				float64(pz-int(z)),
				float64(x),
				float64(z))

			chunk.maxY[x][z] = uint8(voxelValueY)

			for cy := voxelValueY - 1; cy > 0; cy-- {
				cavernNoise := Noise3dSimplexCavern(float64(px)-float64(x), float64(cy), float64(pz)-float64(z), 0)

				if cavernNoise > 0.50 {
					mineralGold := Noise3dMinerals(float64(px)-float64(x), float64(cy), float64(pz)-float64(z), 1)
					mineralLapis := Noise3dMinerals(float64(px)-float64(x), float64(cy), float64(pz)-float64(z), 2)

					if mineralGold < 0.3275 && cavernNoise <= 0.575 && cavernNoise < 0.557 {
						probaTexture := Schneider(Vector2{float64(x + z*2), float64(cy) - float64(z)})
						blockType = BT_CAVERN_GOLD_2
						if probaTexture < 0.5 {
							blockType = BT_CAVERN_GOLD_1
						}
					} else if mineralLapis < 0.325 && cavernNoise <= 0.575 {
						probaTexture := Schneider(Vector2{float64(x + z*2), float64(cy) - float64(z)})
						blockType = BT_CAVERN_LAPIS_2
						if probaTexture < 0.5 {
							blockType = BT_CAVERN_LAPIS_1
						}
					} else {
						blockType = BT_CAVERN_2
						if cavernNoise < 0.557 {
							blockType = BT_CAVERN
						}
					}
					chunk.voxel[V3_Vec3i16{int16(x), int16(cy), int16(z)}] = V3_VoxelData{
						V32_NO_FILL, 0xdddddd, blockType,
					}
				}
			}

			nTerrain := Noise2dSimplex(float64(px-int(x)), float64(pz-int(z)), 0.0, 0.75, 0.0145, 3, 2)
			nTerrain = math.Pow(nTerrain*0.75, 2) * 15
			chunk.voxel[V3_Vec3i16{int16(x), int16(nTerrain), int16(z)}] = V3_VoxelData{
				V32_NO_FILL, 0xdddddd, BT_CAVERN,
			}

			FillChunkWithTrees(chunk, chunkId, float64(px-int(x)), float64(pz-int(z)), float64(x), float64(z))
		}
	}
	FillInside(chunk, chunkId)

}

func (store *V3_ChunkStore) GenerateChunk(chunkId Vec2i, wg *sync.WaitGroup) {
	defer wg.Done()
	var chunk V3_ChunksList

	chunk.voxel = make(map[V3_Vec3i16]V3_VoxelData)

	CreateFilledChunk(chunkId, &chunk)

	store.mutexChunks.Lock()
	remBit(&chunk.flag, V3_C_WAITING)
	setBit(&chunk.flag, V3_C_READY)
	store.chunks[chunkId] = copyChunkData(chunk)
	store.mutexChunks.Unlock()

}
