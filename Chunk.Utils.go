package main

import (
	"image/color"
)

func _i16min(a, b int16) int16 {
	if a > b {
		return b
	}
	return a
}

func _getMaximumValue(diff *int16, chunkId Vec2i, voxelId V3_Vec3i16) {
	cx := chunkId[0] * chunkSize
	cz := chunkId[1] * chunkSize

	voxelValueY := Biome(nil, float64(cx-int(voxelId[0])), float64(cz-int(voxelId[2])), float64(voxelId[0]), float64(voxelId[2]))

	*diff = _i16min(*diff, voxelValueY)

}

func FillInside(chunk *V3_ChunksList, chunkId Vec2i) {
	for voxel := range (*chunk).voxel {
		currVoxel := (*chunk).voxel[voxel]

		if hasBit(currVoxel.flags, V32_FILLER) || hasBit(currVoxel.flags, V32_NO_FILL) || hasBit(currVoxel.flags, V32_SPRITE) {
			continue
		}

		diff := voxel[1]

		if voxel[0] == chunkSize-1 {
			_getMaximumValue(&diff, Vec2i{chunkId[0] - 1, chunkId[1]}, V3_Vec3i16{0, voxel[1], voxel[2]})
		}
		if voxel[0] > 0 {
			diff = _i16min(int16((*chunk).maxY[voxel[0]-1][voxel[2]]), diff)
		}
		if voxel[0] == 0 {
			_getMaximumValue(&diff, Vec2i{chunkId[0] + 1, chunkId[1]}, V3_Vec3i16{chunkSize - 1, voxel[1], voxel[2]})
		}

		if voxel[0] < chunkSize-1 {
			diff = _i16min(int16((*chunk).maxY[voxel[0]+1][voxel[2]]), diff)
		}

		if voxel[2] == chunkSize-1 {
			_getMaximumValue(&diff, Vec2i{chunkId[0], chunkId[1] - 1}, V3_Vec3i16{voxel[0], voxel[1], 0})
		}
		if voxel[2] > 0 {
			diff = _i16min(int16((*chunk).maxY[voxel[0]][voxel[2]-1]), diff)
		}
		if voxel[2] == 0 {
			_getMaximumValue(&diff, Vec2i{chunkId[0], chunkId[1] + 1}, V3_Vec3i16{voxel[0], voxel[1], chunkSize - 1})
		}
		if voxel[2] < chunkSize-1 {
			diff = _i16min(int16((*chunk).maxY[voxel[0]][voxel[2]+1]), diff)
		}

		for cy := int8(1); cy < int8(voxel[1]-diff); cy++ {
			if currVoxel.blockType == BT_CANYON || currVoxel.blockType == BT_CANYON_TOP || currVoxel.blockType == BT_CANYON_WHITE || currVoxel.blockType == BT_CANYON_GRANITE {
				currY := voxel[1] - int16(cy)

				if currY == 43 || currY == 44 || currY == 45 {
					currVoxel.blockType = BT_CANYON_WHITE
				} else if currY == 54 || currY == 55 {
					currVoxel.blockType = BT_CANYON_WHITE
				} else if currY < 54 && currY > 45 {
					currVoxel.blockType = BT_CANYON_GRANITE
				} else if currY > 45 {
					currVoxel.blockType = BT_CANYON_TOP
				} else {
					currVoxel.blockType = BT_CANYON
				}
			}
			(*chunk).voxel[V3_Vec3i16{voxel[0], voxel[1] - int16(cy), voxel[2]}] = V3_VoxelData{
				currVoxel.flags | V32_FILLER, currVoxel.color, currVoxel.blockType,
			}
		}
	}
}

func CreateTreeVoxel(chunks *V3_ChunksList, x, y, z int16, blockType uint8, color color.RGBA, block uint8) {

	idVoxel := V3_Vec3i16{int16(x), int16(y), int16(z)}

	if _, exists := (*chunks).voxel[idVoxel]; exists {
		return
	}

	(*chunks).voxel[V3_Vec3i16{int16(x), int16(y), int16(z)}] = V3_VoxelData{
		blockType | V32_NO_OCCLUDE, 0xffffff, block,
	}
}
