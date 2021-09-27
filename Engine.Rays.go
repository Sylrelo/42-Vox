package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

func (store *V3_ChunkStore) SendDeleteRay(vox Vox, camera, projection mgl32.Mat4) {
	store.mutexPosition.RLock()
	originPos := vox.pos

	ray := mgl32.Vec3{0, 0, 1}

	rotationMat := mgl32.HomogRotate3D(vox.rot.Y(), mgl32.Vec3{0, -1, 0})
	rotationMat = rotationMat.Mul4((mgl32.HomogRotate3D(vox.rot.X(), mgl32.Vec3{-1, 0, 0})))
	ray = rotationMat.Mul4x1(ray.Vec4(1)).Vec3()

	if originPos[0] > 0 {
		originPos[0] += chunkSize
	}
	if originPos[2] > 0 {
		originPos[2] += chunkSize
	}

	store.mutexPosition.RUnlock()
	for f := 1; f < 4; f++ {
		pos := originPos.Add(ray.Mul(float32(f)))

		chunkId := Vec2i{
			int(pos[0] / chunkSize),
			int(pos[2] / chunkSize),
		}

		voxelPos := V3_Vec3i16{
			int16(pos[0]) - (int16(pos[0]/chunkSize) * chunkSize),
			int16(pos[1]) - (int16(pos[1]/chunkSizeYF) * chunkSizeYF),
			int16(pos[2]) - (int16(pos[2]/chunkSize) * chunkSize),
		}
		if voxelPos[0] > 0 {
			voxelPos[0] -= chunkSize
		}
		if voxelPos[1] > 0 {
			voxelPos[1] -= chunkSizeYF
		}
		if voxelPos[2] > 0 {
			voxelPos[2] -= chunkSize
		}

		store.mutexChunks.RLock()
		_, chunkExists := store.chunks[chunkId]
		store.mutexChunks.RUnlock()
		if chunkExists {
			store.mutexChunks.RLock()
			chunkTested := copyChunkData(store.chunks[chunkId])
			store.mutexChunks.RUnlock()

			if _, voxExists := chunkTested.voxel[V3_Vec3i16{int16(voxelPos[0] * -1), int16(voxelPos[1] * -1), int16(voxelPos[2] * -1)}]; voxExists {
				delete(chunkTested.voxel, V3_Vec3i16{int16(voxelPos[0] * -1), int16(voxelPos[1] * -1), int16(voxelPos[2] * -1)})
				remBit(&chunkTested.flag, V3_C_MESHED)
				remBit(&chunkTested.flag, V3_C_MESHING)
				setBit(&chunkTested.flag, V3_C_READY)
				setBit(&chunkTested.flag, V3_C_IN_VIEW)
				setBit(&chunkTested.flag, V3_C_REMESH)
				store.mutexChunks.Lock()
				store.chunks[chunkId] = chunkTested
				store.mutexChunks.Unlock()
				break
			}
			chunkTested.voxel = make(map[V3_Vec3i16]V3_VoxelData)
		}
	}
}
