package main

const (
	DIR_X = uint8(iota)
	DIR_Y = uint8(iota)
	DIR_Z = uint8(iota)
	ZPOS  = uint8(0)
	XPOS  = uint8(1)
	ZNEG  = uint8(2)
	XNEG  = uint8(3)
)

func _isOccluded(voxId V3_Vec3i16, voxels *map[V3_Vec3i16]V3_VoxelData, occluded V3_Vec3i16) bool {
	voxel, exists := (*voxels)[V3_Vec3i16{int16(int16(voxId[0]) + occluded[0]), int16(int16(voxId[1]) + occluded[1]), int16(int16(voxId[2]) + occluded[2])}]

	if exists && !hasBit(voxel.flags, V32_SPRITE) && !hasBit(voxel.flags, V32_NO_OCCLUDE) {
		return true
	}
	return false
}

func _search(voxId V3_Vec3i16, voxels *map[V3_Vec3i16]V3_VoxelData, visited *map[V3_Vec3i16]bool, dir uint8, occluded V3_Vec3i16) (bool, int16, int16) {
	max := int16(voxId[dir])
	min := int16(voxId[dir])
	value := int16(voxId[dir])

	if _isOccluded(voxId, voxels, occluded) {
		return false, 0, 0
	}

	if hasVisited, exists := (*visited)[voxId]; exists && hasVisited {
		return false, min, max
	} else {
		(*visited)[voxId] = true
	}

	for i := value + 1; i < chunkSize; i++ {
		id := voxId
		id[dir] = i

		if _isOccluded(id, voxels, occluded) {
			break
		}
		if hasVisited, exists := (*visited)[id]; exists && hasVisited {
			break
		}
		if chunk, exists := (*voxels)[id]; exists && chunk.color == (*voxels)[voxId].color && chunk.blockType == (*voxels)[voxId].blockType {
			(*visited)[id] = true
		} else {
			break
		}
		max = i
	}

	if value == 0 {
		min = 0
		id := voxId
		id[dir] = 0
		(*visited)[id] = true
	} else {
		for i := int16(voxId[dir] - 1); i > 0; i-- {
			id := voxId
			id[dir] = int16(i)

			if _isOccluded(id, voxels, occluded) {
				break
			}
			if hasVisited, exists := (*visited)[id]; exists && hasVisited {
				break
			}
			if chunk, exists := (*voxels)[id]; exists && chunk.color == (*voxels)[voxId].color && chunk.blockType == (*voxels)[voxId].blockType {
				(*visited)[id] = true
			} else {
				break
			}
			min = int16(i)
		}
	}

	return true, min, max
}

func (store *V3_ChunkStore) GreedyFaces(meshing *V3_Remesh, chunkId Vec2i, chunk *V3_ChunksList) {
	store.mutexChunks.RLock()
	voxels := copyVoxMap((*chunk).voxel)
	store.mutexChunks.RUnlock()

	visitedTop := make(map[V3_Vec3i16]bool)
	visitedBottom := make(map[V3_Vec3i16]bool)
	visitedLeft := make(map[V3_Vec3i16]bool)
	visitedRight := make(map[V3_Vec3i16]bool)
	visitedFront := make(map[V3_Vec3i16]bool)
	visitedBack := make(map[V3_Vec3i16]bool)

	for v := range voxels {
		var vertexes []int32

		if hasBit(voxels[v].flags, V32_SPRITE) {
			(*meshing).glBuffer = append((*meshing).glBuffer, GenerateSpriteBlock(v[0], v[1], v[2], int(voxels[v].color), voxels[v].blockType)...)
			continue
		}
		if processed, min, max := _search(v, &voxels, &visitedTop, DIR_X, V3_Vec3i16{0, 1, 0}); processed {
			(*meshing).glBuffer = append((*meshing).glBuffer,
				GenerateTopFace(min, v[1], v[2], max-min, 0, 0, int(voxels[v].color), voxels[v].blockType)...)
		}

		if processed, min, max := _search(v, &voxels, &visitedBottom, DIR_X, V3_Vec3i16{0, -1, 0}); processed {
			(*meshing).glBuffer = append((*meshing).glBuffer,
				GenerateBottomFace(min, v[1], v[2], max-min, 0, 0, int(voxels[v].color), voxels[v].blockType)...)
		}

		if processed, min, max := _search(v, &voxels, &visitedLeft, DIR_Z, V3_Vec3i16{-1, 0, 0}); processed {
			(*meshing).glBuffer = append((*meshing).glBuffer,
				GenerateLeftFace(v[0], v[1], min, 0, 0, max-min, int(voxels[v].color), voxels[v].blockType)...)
		}

		if processed, min, max := _search(v, &voxels, &visitedRight, DIR_Z, V3_Vec3i16{1, 0, 0}); processed {
			(*meshing).glBuffer = append((*meshing).glBuffer,
				GenerateRightFace(v[0], v[1], min, 0, 0, max-min, int(voxels[v].color), voxels[v].blockType)...)
		}

		if processed, min, max := _search(v, &voxels, &visitedFront, DIR_X, V3_Vec3i16{0, 0, 1}); processed {
			(*meshing).glBuffer = append((*meshing).glBuffer,
				GenerateFrontFace(min, v[1], v[2], max-min, 0, 0, int(voxels[v].color), voxels[v].blockType)...)
		}

		if processed, min, max := _search(v, &voxels, &visitedBack, DIR_X, V3_Vec3i16{0, 0, -1}); processed {
			(*meshing).glBuffer = append((*meshing).glBuffer,
				GenerateBackFace(min, v[1], v[2], max-min, 0, 0, int(voxels[v].color), voxels[v].blockType)...)
		}

		(*meshing).glBuffer = append((*meshing).glBuffer, vertexes...)
	}

	(*meshing).glBufferLen = len(meshing.glBuffer)
	(*meshing).done = true

	visitedTop = make(map[V3_Vec3i16]bool)
	visitedLeft = make(map[V3_Vec3i16]bool)
	visitedFront = make(map[V3_Vec3i16]bool)
	visitedRight = make(map[V3_Vec3i16]bool)
	visitedBottom = make(map[V3_Vec3i16]bool)
	visitedBack = make(map[V3_Vec3i16]bool)
	voxels = make(map[V3_Vec3i16]V3_VoxelData)
}
