package main

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func setBit(flags *uint8, flag uint8) {
	if (*flags)&flag != flag {
		(*flags) |= flag
	}
}

func remBit(flags *uint8, flag uint8) {
	if (*flags)&flag == flag {
		(*flags) ^= flag
	}
}

func hasBit(flags uint8, flag uint8) bool {
	return (flags & flag) == flag
}

func getDistance(player mgl32.Vec3, chunk Vec2i) float64 {
	dx := (player[0] * chunkSizeInv) - float32(chunk[0])
	dz := (player[2] * chunkSizeInv) - float32(chunk[1])

	return math.Sqrt(float64(dx*dx + dz*dz))
}

func copyRemeshData(in V3_Remesh) V3_Remesh {
	var out V3_Remesh

	out.done = in.done
	out.glBufferLen = in.glBufferLen
	out.glBuffer = make([]int32, in.glBufferLen)

	for k := range in.glBuffer {
		out.glBuffer[k] = in.glBuffer[k]
	}

	return out
}

func copyChunkData(in V3_ChunksList) V3_ChunksList {
	var out V3_ChunksList

	out.voxel = make(map[V3_Vec3i16]V3_VoxelData, len(in.voxel))
	out.flag = in.flag
	out.maxY = in.maxY
	for key := range in.voxel {
		out.voxel[key] = in.voxel[key]
	}
	return out
}

func copyVoxMap(in map[V3_Vec3i16]V3_VoxelData) map[V3_Vec3i16]V3_VoxelData {
	result := make(map[V3_Vec3i16]V3_VoxelData, len(in))

	for key := range in {
		result[key] = in[key]
	}
	return result
}

func getTheoricalMaxCircleVision(position, rotation mgl32.Vec3) float64 {

	ray := mgl32.Vec3{0, 0, 1}
	rotationMat := mgl32.HomogRotate3D(rotation[1], mgl32.Vec3{0, -1, 0})
	rotationMat = rotationMat.Mul4((mgl32.HomogRotate3D(rotation[0], mgl32.Vec3{-1, 0, 0})))
	ray = rotationMat.Mul4x1(ray.Vec4(1)).Vec3()

	adaptiveDistance := float64(((float32(MAX_CUBE_VISION) / 1.2) / 10) - 4)

	maxVisionCircle := ((1 - math.Abs(float64(ray[1]))) * adaptiveDistance) + 4
	heightCompensation := math.Abs(math.Abs(float64(position[1]))-70) / (160)

	return maxVisionCircle + (heightCompensation * 2)
}

func ResetWorld(store *V3_ChunkStore) {
	store.mutexChunks.Lock()
	store.chunks = make(map[Vec2i]V3_ChunksList)
	store.mutexChunks.Unlock()

	store.mutexRemesh.Lock()
	store.remesh = make(map[Vec2i]V3_Remesh)
	store.mutexRemesh.Unlock()

	for c := range store.vbo {
		vboValue := store.vbo[c].glVbo
		gl.DeleteBuffers(1, &vboValue)
	}
	store.vbo = make(map[Vec2i]V3_Vbo)
}
