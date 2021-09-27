package main

import (
	"image/color"
	"math"
)

func generateVarianteColorBillowNoise(coco color.RGBA, x, z float64) color.RGBA {
	colorVariationPerlin := Noise2dSimplex(x, z, 0.0, 1.0, 0.035, 5, 1) * 13
	var newColorVariationR uint8
	var newColorVariationG uint8
	var newColorVariationB uint8
	if float64(coco.R)-colorVariationPerlin*2.5 < 0 {
		newColorVariationR = 0
	} else {
		newColorVariationR = coco.R - uint8(colorVariationPerlin*2.5)
	}
	if float64(coco.G)-colorVariationPerlin*2.5 < 0 {
		newColorVariationG = 0
	} else {
		newColorVariationG = coco.G - uint8(colorVariationPerlin*2.5)
	}
	if float64(coco.B)-colorVariationPerlin*2.5 < 0 {
		newColorVariationB = 0
	} else {
		newColorVariationB = coco.B - uint8(colorVariationPerlin*2.5)
	}

	return color.RGBA{newColorVariationR, newColorVariationG, newColorVariationB, 0xff}
}

func generateVarianteColorLEAVESchneider(coco color.RGBA, x, z, ix, iy, iz float64) color.RGBA {
	colorVariationSchneiderTree := Schneider(Vector2{(x), (z)}) * 50
	colorVariationSchneiderLeaf := Schneider(Vector2{(ix - iy), (iz + iy)}) * 5

	var newColorVariationR uint8
	var newColorVariationG uint8
	var newColorVariationB uint8
	if float64(coco.R)-colorVariationSchneiderLeaf*2-colorVariationSchneiderTree < 0 {
		newColorVariationR = 0
	} else {
		newColorVariationR = coco.R - uint8(colorVariationSchneiderLeaf*2) - uint8(colorVariationSchneiderTree)
	}
	if float64(coco.G)-colorVariationSchneiderLeaf*2-colorVariationSchneiderTree < 0 {
		newColorVariationG = 0
	} else {
		newColorVariationG = coco.G - uint8(colorVariationSchneiderLeaf*2) - uint8(colorVariationSchneiderTree)
	}
	if float64(coco.B)-colorVariationSchneiderLeaf*2-colorVariationSchneiderTree < 0 {
		newColorVariationB = 0
	} else {
		newColorVariationB = coco.B - uint8(colorVariationSchneiderLeaf*2) - uint8(colorVariationSchneiderTree)
	}
	return color.RGBA{newColorVariationR, newColorVariationG, newColorVariationB, 0xff}
}

func _generateBasicLeaf(chunks *V3_ChunksList, chunkId Vec2i, cx, nInt, cz int16, trunkHeight int, blockType uint8, color color.RGBA) {
	leafLength := trunkHeight
	if leafLength%2 == 0 {
		leafLength++
	}
	if leafLength < 3 {
		leafLength = 3
	}
	var bt uint8

	for ix := 1; ix <= leafLength; ix++ {
		for iy := 1; iy <= leafLength; iy++ {
			for iz := 1; iz <= leafLength; iz++ {
				if !(ix == 1 && iy == 1 && iz == 1) &&
					!(ix == 1 && iy == leafLength && iz == 1) &&
					!(ix == leafLength && iy == 1 && iz == 1) &&
					!(ix == leafLength && iy == leafLength && iz == 1) &&
					!(ix == 1 && iy == 1 && iz == leafLength) &&
					!(ix == 1 && iy == leafLength && iz == leafLength) &&
					!(ix == leafLength && iy == 1 && iz == leafLength) &&
					!(ix == leafLength && iy == leafLength && iz == leafLength) {

					if Schneider(Vector2{float64(ix*4 + int(cx)), float64(iz + iy*3 - int(cz))}) > 0.35 {
						bt = BT_TREE_LEAVES_1
						if Schneider(Vector2{float64(ix - int(cx)), float64(iz - iy*int(cz))}) > 0.5 {
							bt = BT_TREE_LEAVES_2
						}
						nCol := generateVarianteColorLEAVESchneider(color, float64(cx), float64(cz), float64(ix), float64(iy), float64(iz))
						CreateTreeVoxel(chunks, cx+int16(ix)-int16(leafLength)/2-1, nInt+int16(trunkHeight+iy), cz+int16(iz)-int16(leafLength)/2-1, blockType, nCol, bt)
					}
				}
			}
		}
	}
}

func _generateMeinLeaf(chunks *V3_ChunksList, cx, nInt, cz int16, trunkHeight int, blockType uint8, color color.RGBA) {
	leafLength := trunkHeight
	if leafLength%2 == 0 {
		leafLength++
	}
	var oldY int = 0
	var bt uint8

	for ix := 1; ix <= leafLength; ix++ {
		for iy := 1; iy <= leafLength/2; iy++ {
			for iz := 1; iz <= leafLength; iz++ {
				if !(ix == leafLength && iy == leafLength && iz == leafLength) {
					if Schneider(Vector2{float64(ix + int(cx)), float64(iz + iy - int(cz))}) > 0.1 {
						bt = BT_TREE_LEAVES_1
						if Schneider(Vector2{float64(ix - int(cx)), float64(iz - iy*int(cz))}) > 0.5 {
							bt = BT_TREE_LEAVES_2
						}
						newColor := generateVarianteColorLEAVESchneider(color, float64(cx), float64(cz), float64(ix), float64(iy), float64(iz))
						CreateTreeVoxel(chunks,
							int16(cx+int16(ix)-int16(leafLength)/2-1),
							int16(nInt+int16(trunkHeight+iy)),
							int16(cz+int16(iz)-int16(leafLength)/2-1),
							blockType, newColor, bt)
					}
				}

			}
			oldY = iy
		}
	}
	for ix := 1; ix <= leafLength-2; ix++ {
		for iy := 1; iy <= leafLength/2; iy++ {
			for iz := 1; iz <= leafLength-2; iz++ {
				if !(ix == leafLength && iy == leafLength && iz == leafLength) {
					if Schneider(Vector2{float64(ix - iy*5 - int(cx)), float64(iz - iy + int(cz))}) > 0.1 {
						bt = BT_TREE_LEAVES_1
						if Schneider(Vector2{float64(ix - int(cx)), float64(iz - iy*int(cz))}) > 0.5 {
							bt = BT_TREE_LEAVES_2
						}
						newColor := generateVarianteColorLEAVESchneider(color, float64(cx), float64(cz), float64(ix), float64(iy), float64(iz))
						CreateTreeVoxel(chunks,
							int16(cx+int16(ix)-int16(leafLength)/2),
							int16(nInt+int16(trunkHeight+iy+oldY)),
							int16(cz+int16(iz)-int16(leafLength)/2),
							blockType, newColor, bt)
					}
				}
			}
		}
	}
}

func _generateFirLeaf(chunks *V3_ChunksList, cx, nInt, cz int16, trunkHeight int, blockType uint8, coco color.RGBA) {
	leafLength := trunkHeight * 2
	randomBlocktype := uint8(blockType)

	if leafLength%2 != 0 {
		leafLength++
	}

	if leafLength < 2 {
		leafLength = 2
	}
	ixStart := -leafLength / 2
	izStart := -leafLength / 2
	ixEnd := -leafLength / 2
	izEnd := -leafLength / 2
	for iy := 1; iy <= leafLength*2; iy += 2 {
		for ix := ixStart; ix <= leafLength+ixEnd; ix++ {
			for iz := izStart; iz <= leafLength+izEnd; iz++ {
				if !(ix == ixStart && iz == izStart) &&
					!(ix == ixStart && iz == leafLength+izEnd) &&
					!(ix == leafLength+ixEnd && iz == izStart) &&
					!(ix == leafLength+ixEnd && iz == leafLength+izEnd) &&
					!(ix == leafLength && iy == leafLength && iz == leafLength) {
					if Schneider(Vector2{float64(ix - iy*5 - int(cx)), float64(iz - iy + int(cz))}) > 0.15 {
						fillerType := V32_NO_FILL | V32_SPRITE
						if blockType == BT_TREE_LEAVES_HIGH_DENSITY_1 {
							if Schneider(Vector2{float64(ix + iy), float64(iz)}) > 0.5 {
								randomBlocktype = BT_TREE_LEAVES_HIGH_DENSITY_1
							} else {
								randomBlocktype = BT_TREE_LEAVES_HIGH_DENSITY_2
							}
						} else if blockType == BT_TREE_LEAVES_1 {
							if Schneider(Vector2{float64(ix + iy), float64(iz)}) > 0.5 {
								randomBlocktype = BT_TREE_LEAVES_1
							} else {
								randomBlocktype = BT_TREE_LEAVES_2
							}
						} else if blockType == BT_TREE_LEAVES_PINK_1 {
							fillerType = V32_NO_FILL
							if Schneider(Vector2{float64(ix + iy), float64(iz)}) > 0.5 {
								randomBlocktype = BT_TREE_LEAVES_PINK_1
							} else {
								randomBlocktype = BT_TREE_LEAVES_PINK_2
							}
						} else if blockType == BT_TREE_LEAVES_SNOW_1 {
							fillerType = V32_NO_FILL
							if Schneider(Vector2{float64(ix + iy), float64(iz)}) > 0.5 {
								randomBlocktype = BT_TREE_LEAVES_SNOW_1
							} else {
								randomBlocktype = BT_TREE_LEAVES_SNOW_2
							}
						}
						newColor := generateVarianteColorLEAVESchneider(coco, float64(cx), float64(cz), float64(ix), float64(iy), float64(iz))
						CreateTreeVoxel(chunks,
							int16((cx + int16(ix))),
							int16((nInt + int16(trunkHeight+iy))),
							int16((cz + int16(iz))),
							uint8(fillerType), newColor, randomBlocktype)
					}
				}
			}
		}
		ixStart++
		izStart++
		ixEnd--
		izEnd--
	}
	ixStart = -leafLength / 2
	izStart = -leafLength / 2
	ixEnd = -leafLength / 2
	izEnd = -leafLength / 2
	for iy := 1; iy <= leafLength*2; iy += 2 {
		for ix := ixStart + 2; ix <= leafLength+ixEnd-2; ix++ {
			for iz := izStart + 2; iz <= leafLength+izEnd-2; iz++ {
				if !(ix == 0 && iz == 0) {
					if Schneider(Vector2{float64(ix - iy*5 - int(cx)), float64(iz - iy + int(cz))}) > 0.35 {
						fillerType := V32_NO_FILL | V32_SPRITE
						if blockType == BT_TREE_LEAVES_HIGH_DENSITY_1 {
							if Schneider(Vector2{float64(ix + iy), float64(iz)}) > 0.5 {
								randomBlocktype = BT_TREE_LEAVES_HIGH_DENSITY_1
							} else {
								randomBlocktype = BT_TREE_LEAVES_HIGH_DENSITY_2
							}
						} else if blockType == BT_TREE_LEAVES_1 {
							if Schneider(Vector2{float64(ix + iy), float64(iz)}) > 0.5 {
								randomBlocktype = BT_TREE_LEAVES_1
							} else {
								randomBlocktype = BT_TREE_LEAVES_2
							}
						} else if blockType == BT_TREE_LEAVES_SNOW_1 {
							fillerType = V32_NO_FILL
							if Schneider(Vector2{float64(ix + iy), float64(iz)}) > 0.5 {
								randomBlocktype = BT_TREE_LEAVES_SNOW_1
							} else {
								randomBlocktype = BT_TREE_LEAVES_SNOW_2
							}
						}
						newColor := generateVarianteColorLEAVESchneider(coco, float64(cx), float64(cz), float64(ix), float64(iy), float64(iz))
						CreateTreeVoxel(chunks,
							int16((cx + int16(ix))),
							int16((nInt + int16(trunkHeight+iy) + 1)),
							int16((cz + int16(iz))),
							uint8(fillerType), newColor, randomBlocktype)
					}
				}
			}
		}
		ixStart++
		izStart++
		ixEnd--
		izEnd--
	}
}

/*
func _generateFirSnowLeaf(chunks *V3_ChunksList, cx, nInt, cz int16, trunkHeight int, blockType uint8, coco color.RGBA) {
	leafLength := trunkHeight * 3

	if leafLength%2 != 0 {
		leafLength++
	}

	if leafLength < 2 {
		leafLength = 2
	}
	ixStart := -leafLength / 2
	izStart := -leafLength / 2
	ixEnd := -leafLength / 2
	izEnd := -leafLength / 2
	for iy := 1; iy <= leafLength; iy++ {
		size := int(Schneider(Vector2{float64(iy) + float64(cz)*2, float64(cz) * 3}) * 3)
		for ix := -size; ix <= size; ix++ {
			for iz := -size; iz <= size; iz++ {
				if iy < leafLength/2 {
					probaNot := Schneider(Vector2{float64(ix), float64(iz + iy)})
					if probaNot > 0.3 {
						var randomBlocktype uint8
						if Schneider(Vector2{float64(ix + iy), float64(iz)}) > 0.5 {
							randomBlocktype = BT_TREE_LEAVES_SNOW_1
						} else {
							randomBlocktype = BT_TREE_LEAVES_SNOW_2
						}
						newColor := generateVarianteColorLEAVESchneider(coco, float64(cx), float64(cz), float64(ix), float64(iy), float64(iz))
						CreateTreeVoxel(chunks,
							int16((cx + int16(ix))),
							int16((nInt + int16(trunkHeight+iy) + 1)),
							int16((cz + int16(iz))),
							uint8(V32_NO_FILL), newColor, int(randomBlocktype))
					}
				}
			}
		}
		ixStart++
		izStart++
		ixEnd--
		izEnd--
	}
	// (*chunk).glBuffer = append(chunk.glBuffer, generateCube(x+float32(cx)*2, (y+float32(nInt+trunkHeight*3)*2), z+float32(cz)*2, faceBits, float32(coco.R)/256.0, float32(coco.G)/256.0, float32(coco.B)/256.0)...)
}
*/

func _randomPalmLeaf(x, z float64) uint8 {
	if Schneider(Vector2{x, z}) > 0.5 {
		return BT_TREE_LEAVES_PALM_1
	} else {
		return BT_TREE_LEAVES_PALM_2
	}
}

func _generatePalmTreeLeaf(chunks *V3_ChunksList, x, z float64, cx, cz, nInt int16, trunkHeight, xDir, zDir, i int16, color color.RGBA) {
	probaDirLeaf1 := Schneider(Vector2{float64(x) + float64(cz)*2, float64(z) * 3})
	probaDirLeaf2 := Schneider(Vector2{float64(x) + float64(cz), float64(z) * 2})
	probaDirLeaf3 := Schneider(Vector2{float64(x+z)*2 + float64(cx), float64(z) * 4})
	probaDirLeaf4 := Schneider(Vector2{float64(x+x) + float64(cx)*5, float64(z + x)})
	randomReduceSize := int16(0)
	for leafGen := int16(1); leafGen <= trunkHeight/2+randomReduceSize; leafGen++ {
		probaReduceSize := Schneider(Vector2{float64(float32(cx)*2) + float64(leafGen), float64(z + float64(cz))})
		if probaReduceSize > 0.85 {
			randomReduceSize -= 2
		}
		newColor := generateVarianteColorLEAVESchneider(color, float64(cx), float64(cz), float64(cx), float64(leafGen), float64(cz))
		if probaDirLeaf1 > 0.5 {
			CreateTreeVoxel(chunks, int16(cx+xDir+leafGen), int16(nInt+i+1-leafGen), int16(cz+zDir), V32_NO_FILL|V32_SPRITE, newColor, _randomPalmLeaf(x, z))
			CreateTreeVoxel(chunks, int16(cx+xDir+leafGen), int16(nInt+i+1-leafGen+1), int16(cz+zDir), V32_NO_FILL|V32_SPRITE, newColor, _randomPalmLeaf(x, z))
		} else {
			CreateTreeVoxel(chunks, int16(cx+xDir+leafGen), int16(nInt+i+1-leafGen), int16(cz+zDir+leafGen), V32_NO_FILL|V32_SPRITE, newColor, _randomPalmLeaf(x, z))
			CreateTreeVoxel(chunks, int16(cx+xDir+leafGen), int16(nInt+i+1-leafGen+1), int16(cz+zDir+leafGen), V32_NO_FILL|V32_SPRITE, newColor, _randomPalmLeaf(x, z))
		}
	}
	randomReduceSize = 0
	for leafGen := int16(1); leafGen <= trunkHeight/2+randomReduceSize; leafGen++ {
		probaReduceSize := Schneider(Vector2{float64(x+float64(cx)) + float64(leafGen), float64(z + float64(cz)*2)})
		if probaReduceSize > 0.85 {
			randomReduceSize -= 2
		}
		newColor := generateVarianteColorLEAVESchneider(color, float64(cx), float64(cz), float64(cx), float64(leafGen), float64(cz))
		if probaDirLeaf2 > 0.5 {
			CreateTreeVoxel(chunks, int16(cx+xDir-leafGen), int16(nInt+i+1-leafGen), int16(cz+zDir), V32_NO_FILL|V32_SPRITE, newColor, _randomPalmLeaf(x, z))
			CreateTreeVoxel(chunks, int16(cx+xDir-leafGen), int16(nInt+i+1-leafGen+1), int16(cz+zDir), V32_NO_FILL|V32_SPRITE, newColor, _randomPalmLeaf(x, z))
		} else {
			CreateTreeVoxel(chunks, int16(cx+xDir-leafGen), int16(nInt+i+1-leafGen), int16(cz+zDir-leafGen), V32_NO_FILL|V32_SPRITE, newColor, _randomPalmLeaf(x, z))
			CreateTreeVoxel(chunks, int16(cx+xDir-leafGen), int16(nInt+i+1-leafGen+1), int16(cz+zDir-leafGen), V32_NO_FILL|V32_SPRITE, newColor, _randomPalmLeaf(x, z))
		}
	}
	randomReduceSize = 0
	for leafGen := int16(1); leafGen <= trunkHeight/2+randomReduceSize; leafGen++ {
		probaReduceSize := Schneider(Vector2{float64(x+float64(cz)*3) + float64(leafGen), float64(z + float64(cz) + 5)})
		if probaReduceSize > 0.85 {
			randomReduceSize -= 2
		}
		newColor := generateVarianteColorLEAVESchneider(color, float64(cx), float64(cz), float64(cx), float64(leafGen), float64(cz))
		if probaDirLeaf3 > 0.5 {
			CreateTreeVoxel(chunks, int16(cx+xDir), int16(nInt+i+1-leafGen), int16(cz+zDir-leafGen), V32_NO_FILL|V32_SPRITE, newColor, _randomPalmLeaf(x, z))
			CreateTreeVoxel(chunks, int16(cx+xDir), int16(nInt+i+1-leafGen+1), int16(cz+zDir-leafGen), V32_NO_FILL|V32_SPRITE, newColor, _randomPalmLeaf(x, z))
		} else {
			CreateTreeVoxel(chunks, int16(cx+xDir+leafGen), int16(nInt+i+1-leafGen), int16(cz+zDir-leafGen), V32_NO_FILL|V32_SPRITE, newColor, _randomPalmLeaf(x, z))
			CreateTreeVoxel(chunks, int16(cx+xDir+leafGen), int16(nInt+i+1-leafGen+1), int16(cz+zDir-leafGen), V32_NO_FILL|V32_SPRITE, newColor, _randomPalmLeaf(x, z))
		}
	}
	randomReduceSize = 0
	for leafGen := int16(1); leafGen <= trunkHeight/2+randomReduceSize; leafGen++ {
		probaReduceSize := Schneider(Vector2{float64(x*2) + float64(leafGen)*0.5, float64(z+float64(cz*2)) * 2})
		if probaReduceSize > 0.85 {
			randomReduceSize -= 2
		}
		newColor := generateVarianteColorLEAVESchneider(color, float64(cx), float64(cz), float64(cx), float64(leafGen), float64(cz))
		if probaDirLeaf4 > 0.5 {
			CreateTreeVoxel(chunks, int16(cx+xDir), int16(nInt+i+1-leafGen), int16(cz+zDir+leafGen), V32_NO_FILL|V32_SPRITE, newColor, _randomPalmLeaf(x, z))
			CreateTreeVoxel(chunks, int16(cx+xDir), int16(nInt+i+1-leafGen+1), int16(cz+zDir+leafGen), V32_NO_FILL|V32_SPRITE, newColor, _randomPalmLeaf(x, z))
		} else {
			CreateTreeVoxel(chunks, int16(cx+xDir-leafGen), int16(nInt+i+1-leafGen), int16(cz+zDir+leafGen), V32_NO_FILL|V32_SPRITE, newColor, _randomPalmLeaf(x, z))
			CreateTreeVoxel(chunks, int16(cx+xDir-leafGen), int16(nInt+i+1-leafGen+1), int16(cz+zDir+leafGen), V32_NO_FILL|V32_SPRITE, newColor, _randomPalmLeaf(x, z))
		}
	}
}

func _generateRiver(nRiverNoise float64, nInt *int) bool {
	minRiver := 0.34
	maxRiver := 0.35
	if nRiverNoise > minRiver && nRiverNoise < maxRiver {
		(*nInt)--
		if nRiverNoise > minRiver+(maxRiver-minRiver)*0.25 && nRiverNoise < minRiver+(maxRiver-minRiver)*0.75 {
			(*nInt)--
		}
		if nRiverNoise > minRiver+(maxRiver-minRiver)*0.35 && nRiverNoise < minRiver+(maxRiver-minRiver)*0.65 {
			(*nInt)--
		}
		return true
	}
	return false
}

func FillChunkWithTrees(chunks *V3_ChunksList, chunkId Vec2i, x, z, cx, cz float64) {
	nTerrain := Noise2dSimplex(x, (z), 0.5, 0.55, 0.0006, 8, 0)
	nTerrain = math.Pow(nTerrain*0.75, 2)
	nMountainBiome := Noise2dSimplex(x, (z), 0.0, 0.75, 0.00025, 1, 1)
	nDesertSnowBiome := Noise2dSimplex(x, (z), 0.0, 0.85, 0.00025, 5, 2)
	nHoleNoise := Noise2dSimplex(float64(x), float64(z), 0.0, 0.65, 0.0045, 2, 4)
	nHoleNoise = math.Pow(nHoleNoise, 2.0)
	nRiverNoise := Noise2dSimplex(x, (z), 0.2, 1.0, 0.00025, 1, 3)
	nInt := int(nTerrain * chunkSizeYF * 0.5)

	colorgb := color.RGBA{0, 0, 0, 1}
	var folliageType uint8
	var bt uint8

	if nInt <= 32 || nTerrain < 0.2625 || nHoleNoise > .825 || nRiverNoise > 0.34 && nRiverNoise < 0.35 {
		// no tree cuz in sea || cuz in hole || cuz in river //
	} else {
		if nMountainBiome > 0.75 {
		} else if nMountainBiome < 0.2 {
		} else if nDesertSnowBiome >= 0.705 && nTerrain >= 0.25 {
			if nTerrain > 0.55 {
				if Schneider(Vector2{float64(x), float64(z)}) > 0.99 {
					// Tree high-snow generation //
					var trunkHeight int = int(math.Round(Schneider(Vector2{float64(x + (cx)*2), float64(z + (cz)*2)}) * 4))
					_generateFirLeaf(chunks, int16(cx), int16(nInt), int16(cz), trunkHeight, BT_TREE_LEAVES_PINK_1, colorgb)
					trunkHeight = trunkHeight*3 - 1
					for trunkHeight > 0 {
						chunks.voxel[V3_Vec3i16{int16(cx), int16(nInt + trunkHeight), int16(cz)}] = V3_VoxelData{
							V32_NO_FILL, 0xffffff, BT_TREE_TRUNK_DARK,
						}
						trunkHeight--
					}
				} else if Schneider(Vector2{float64(z + x), float64(x - z)}) > 0.95 {
					// Folliage high-snow generation //
					typeFolliageProba := Schneider(Vector2{float64(x + 10 + cx*2), float64(z + cz*2)})

					if typeFolliageProba < 0.6 {
						folliageType = BT_FOLLIAGE_ORCHID
					} else {
						folliageType = BT_FOLLIAGE_LILY
					}
					chunks.voxel[V3_Vec3i16{int16(cx), int16(nInt + 1), int16(cz)}] = V3_VoxelData{
						V32_NO_FILL | V32_SPRITE, 0xffffff, folliageType,
					}
				}
			} else {
				if Schneider(Vector2{float64(x), float64(z)}) > 0.9975 {
					// Tree snow generation //
					var trunkHeight int = int(math.Round(Schneider(Vector2{float64(x + (cx)*2), float64(z + (cz)*2)}) * 4))
					_generateFirLeaf(chunks, int16(cx), int16(nInt), int16(cz), trunkHeight, BT_TREE_LEAVES_SNOW_1, colorgb)
					trunkHeight = trunkHeight*3 - 1
					for trunkHeight > 0 {
						chunks.voxel[V3_Vec3i16{int16(cx), int16(nInt + trunkHeight), int16(cz)}] = V3_VoxelData{
							V32_NO_FILL, 0xffffff, BT_TREE_TRUNK_GREY,
						}
						trunkHeight--
					}
				} else if Schneider(Vector2{float64(z + x), float64(x - z)}) > 0.975 {
					// Folliage snow generation //
					typeFolliageProba := Schneider(Vector2{float64(x + 10 + cx*2), float64(z + cz*2)})
					if typeFolliageProba < 0.30 {
						folliageType = BT_FOLLIAGE_DEAD_BUSH
					} else if typeFolliageProba < 0.6 {
						folliageType = BT_FOLLIAGE_GRASS_SMALL
					} else if typeFolliageProba < 0.85 {
						folliageType = BT_FOLLIAGE_GRASS_BIG
					} else {
						folliageType = BT_FOLLIAGE_SAPLING
					}
					chunks.voxel[V3_Vec3i16{int16(cx), int16(nInt + 1), int16(cz)}] = V3_VoxelData{
						V32_NO_FILL | V32_SPRITE, 0xffffff, folliageType,
					}
				}
			}

		} else if nDesertSnowBiome < 0.20 && nTerrain >= 0.25 {
			if Schneider(Vector2{float64(x), float64(z)}) > 0.9995 {
				// Tree desert generation //
				typeTreeProba := Schneider(Vector2{float64(x + 10 + (cx)*2), float64(z + (cz)*2)})
				if typeTreeProba > 0.33 {
					// Palm tree generation
					colorgb = color.RGBA{35, 117, 67, 0xff}
					var trunkHeight int = int(math.Round(Schneider(Vector2{float64(x + (cx)*2), float64(z + (cz)*2)})*4)) + 10
					// Trunk generation
					colorgb = color.RGBA{83, 53, 10, 0xff}
					probaDir := Schneider(Vector2{float64(x) + float64(cz), float64(z)})
					xDir := 0
					zDir := 0
					i := 0
					for i < trunkHeight {
						probaToMoveInDir := Schneider(Vector2{x + float64(i), float64(z + (cz)*2)})
						if probaToMoveInDir > 0.45 { // move trunk to a Direction
							if probaDir < 0.25 {
								xDir++
							} else if probaDir < 0.5 {
								xDir--
							} else if probaDir < 0.75 {
								zDir++
							} else {
								zDir--
							}
						}
						probaTreeTexture := Schneider(Vector2{x - float64(i), float64(z - (cz)*3)})
						if probaTreeTexture > 0.5 {
							bt = BT_TREE_TRUNK_PALM_1
						} else {
							bt = BT_TREE_TRUNK_PALM_2
						}
						CreateTreeVoxel(chunks,
							int16((cx + float64(xDir))),
							int16((nInt + i + 1)),
							int16((cz + float64(zDir))),
							V32_NO_FILL, colorgb, bt)
						i++
					}
					_generatePalmTreeLeaf(chunks, x, z, int16(cx), int16(cz), int16(nInt), int16(trunkHeight), int16(xDir), int16(zDir), int16(i), colorgb)
				} else {
					// Cactus tree generation
					var trunkHeight int = int(math.Round(Schneider(Vector2{float64(x + (cx)*2), float64(z + (cz)*2)})*5)) + 1
					colorgb = color.RGBA{36, 88, 33, 0xff}
					for trunkHeight > 0 {
						chunks.voxel[V3_Vec3i16{int16(cx), int16(nInt + trunkHeight), int16(cz)}] = V3_VoxelData{
							V32_NO_FILL, 0xffffff, BT_TREE_CACTUS,
						}
						trunkHeight--
					}
				}

			} else if Schneider(Vector2{float64(z + x), float64(x - z)}) > 0.995 {
				// Folliage desert generation //
				typeFolliageProba := Schneider(Vector2{float64(x + 10 + cx*2), float64(z + cz*2)})
				if typeFolliageProba < 0.75 {
					folliageType = BT_FOLLIAGE_BUSH
				} else if typeFolliageProba < 0.925 {
					folliageType = BT_FOLLIAGE_CACTUS_1
				} else {
					folliageType = BT_FOLLIAGE_CACTUS_2
					chunks.voxel[V3_Vec3i16{int16(cx), int16(nInt + 2), int16(cz)}] = V3_VoxelData{
						V32_NO_FILL | V32_SPRITE, 0xffffff, BT_FOLLIAGE_CACTUS_1,
					}
				}
				chunks.voxel[V3_Vec3i16{int16(cx), int16(nInt + 1), int16(cz)}] = V3_VoxelData{
					V32_NO_FILL | V32_SPRITE, 0xffffff, folliageType,
				}
			}
		} else {
			if nTerrain > 0.45 && nTerrain < 0.65 {
				// Tree autumn generation //
				if Schneider(Vector2{float64(x), float64(z)}) > 0.99 {
					var trunkHeight int = int(math.Round(Schneider(Vector2{float64(x + (cx)*2), float64(z + (cz)*2)}) * 4))
					typeTreeProba := Schneider(Vector2{float64(x + 10 + float64(cx)*2), float64(z + float64(cz)*2)})
					if typeTreeProba > 0.0 {
						_generateFirLeaf(chunks, int16(cx), int16(nInt), int16(cz), trunkHeight, BT_TREE_LEAVES_HIGH_DENSITY_1, colorgb)
						trunkHeight *= 3
						if trunkHeight > 3 {
							trunkHeight -= 2
						}
					} else {
						_generateMeinLeaf(chunks, int16(cx), int16(nInt), int16(cz), trunkHeight, V32_NO_FILL|V32_SPRITE, colorgb)
					}
					for trunkHeight > 0 {
						chunks.voxel[V3_Vec3i16{int16(cx), int16(nInt + trunkHeight), int16(cz)}] = V3_VoxelData{
							V32_NO_FILL, 0xffffff, BT_TREE_TRUNK_DARK,
						}
						trunkHeight--
					}
				} else if Schneider(Vector2{float64(z + x), float64(x - z)}) > 0.95 {
					// Folliage autumn generation //
					typeFolliageProba := Schneider(Vector2{float64(x + 10 + cx*2), float64(z + cz*2)})
					if typeFolliageProba < 0.3 {
						folliageType = BT_FOLLIAGE_AUTUMN_SMALL
					} else if typeFolliageProba < 0.6 {
						folliageType = BT_FOLLIAGE_AUTUMN_BIG_1
						chunks.voxel[V3_Vec3i16{int16(cx), int16(nInt + 2), int16(cz)}] = V3_VoxelData{
							V32_NO_FILL | V32_SPRITE, 0xffffff, BT_FOLLIAGE_AUTUMN_BIG_2,
						}
					} else if typeFolliageProba < 0.95 {
						folliageType = BT_FOLLIAGE_AUTUMN_FERN_1
						chunks.voxel[V3_Vec3i16{int16(cx), int16(nInt + 2), int16(cz)}] = V3_VoxelData{
							V32_NO_FILL | V32_SPRITE, 0xffffff, BT_FOLLIAGE_AUTUMN_FERN_2,
						}
					} else if typeFolliageProba < 0.955 {
						folliageType = BT_FOLLIAGE_MUSHROOM_BROWN
					} else {
						folliageType = BT_FOLLIAGE_MUSHROOM_FUNGUS
					}
					chunks.voxel[V3_Vec3i16{int16(cx), int16(nInt + 1), int16(cz)}] = V3_VoxelData{
						V32_NO_FILL | V32_SPRITE, 0xffffff, folliageType,
					}
				}

			} else if nTerrain < 0.85 {
				if Schneider(Vector2{float64(x), float64(z)}) > 0.9975 {
					// Tree basic generation //
					var trunkHeight int = int(math.Round(Schneider(Vector2{float64(x + (cx)*2), float64(z + (cz)*2)}) * 5))
					typeTreeProba := Schneider(Vector2{float64(x + 10 + cx*2), float64(z + cz*2)})
					if typeTreeProba < 0.33 {
						_generateBasicLeaf(chunks, chunkId, int16(cx), int16(nInt), int16(cz), trunkHeight, V32_NO_FILL|V32_SPRITE, colorgb)
					} else if typeTreeProba < 0.66 {
						_generateMeinLeaf(chunks, int16(cx), int16(nInt), int16(cz), trunkHeight, V32_NO_FILL|V32_SPRITE, colorgb)
					} else {
						_generateFirLeaf(chunks, int16(cx), int16(nInt), int16(cz), trunkHeight, BT_TREE_LEAVES_1, colorgb)
						trunkHeight *= 3
						if trunkHeight > 3 {
							trunkHeight -= 2
						}
					}
					for trunkHeight > 0 {
						chunks.voxel[V3_Vec3i16{int16(cx), int16(nInt + trunkHeight), int16(cz)}] = V3_VoxelData{
							V32_NO_FILL, 0xffffff, BT_TREE_TRUNK,
						}
						trunkHeight--
					}
				} else if Schneider(Vector2{float64(z + x), float64(x - z)}) > 0.875 {
					// Folliage basic generation //
					typeFolliageProba := Schneider(Vector2{float64(x + 10 + cx*2), float64(z + cz*2)})
					if typeFolliageProba < 0.3 {
						folliageType = BT_FOLLIAGE_GRASS_BIG
					} else if typeFolliageProba < 0.6 {
						folliageType = BT_FOLLIAGE_GRASS_SMALL
					} else if typeFolliageProba < 0.95 {
						folliageType = BT_FOLLIAGE_GRASS_FERN
					} else if typeFolliageProba < 0.955 {
						folliageType = BT_FOLLIAGE_MUSHROOM_RED
					} else {
						folliageType = BT_FOLLIAGE_DAISY
					}
					chunks.voxel[V3_Vec3i16{int16(cx), int16(nInt + 1), int16(cz)}] = V3_VoxelData{
						V32_NO_FILL | V32_SPRITE, 0xffffff, folliageType,
					}
				}
			}
		}
	}
}

func Biome(chunk *V3_ChunksList, x, z, cx, cz float64) int16 {
	nTerrain := Noise2dSimplex(x, (z), 0.5, 0.55, 0.0006, 8, 0)
	nTerrain = math.Pow(nTerrain*0.75, 2)
	nMountainBiome := Noise2dSimplex(x, (z), 0.0, 0.75, 0.00025, 1, 1)
	nDesertSnowBiome := Noise2dSimplex(x, (z), 0.0, 0.85, 0.00025, 5, 2)
	nRiverNoise := Noise2dSimplex(x, (z), 0.2, 1.0, 0.00025, 1, 3)

	nInt := int(nTerrain * chunkSizeYF * 0.5)
	blockType := uint8(BT_GRASS)

	if nMountainBiome > 0.75 {
		if nMountainBiome < 0.8 {
			// Interpolation between Mountain biome and others \\
			nMountainTerrain := Noise2dSimplex(x, float64(z), 0.0, 0.5, 0.00125, 3, 0)
			nMountainTerrain = math.Pow(nMountainTerrain, 1.25) * 0.75
			biomeRange := (0.8 - 0.75)
			terrainRange := (nMountainTerrain - nTerrain*0.5)
			interpolationValue := (((nMountainBiome - 0.75) * terrainRange) / biomeRange) + nTerrain*0.5
			if interpolationValue > 1 {
				interpolationValue = 1
			}
			nInt = int(interpolationValue * chunkSizeYF)
			if nInt > 140 {
				blockType = BT_MOUNTAIN_SNOW
			} else {
				proba := Schneider(Vector2{float64(x + (cx)*4), float64(z - (cz)*6)}) * 3
				if proba > 1 {
					blockType = BT_MOUNTAIN
				} else {
					blockType = BT_MOUNTAIN_MOSSY
				}
			}
			if nInt <= 32 {
				blockType = BT_WATER
			}
		} else {
			// Mountain Biome \\
			nMountainTerrain := Noise2dSimplex(x, float64(z), 0.0, 0.5, 0.00125, 3, 0)
			nMountainTerrain = math.Pow(nMountainTerrain, 1.25) * 0.75
			nInt = int(nMountainTerrain * chunkSizeYF)
			if nInt > 135 {
				blockType = BT_MOUNTAIN_SNOW
			} else {
				// Snow tree //
				proba := Schneider(Vector2{float64(x + (cx)*4), float64(z - (cz)*6)}) * 3
				if proba > 1 {
					blockType = BT_MOUNTAIN
				} else {
					blockType = BT_MOUNTAIN_MOSSY
				}
			}
			if nMountainTerrain > 1 {
				nMountainTerrain = 1
			}
			if nInt <= 32 {
				blockType = BT_WATER
			}
		}
	} else if nMountainBiome < 0.2 {
		if nMountainBiome > 0.19 {
			// Interpolation between Canyon biome and others \\
			nCanyonTerrain := Noise2dPerlin(x, float64(z), -0.2, 1.75, 0.0050, 8, 0)
			nIntCanyon := int(nCanyonTerrain * chunkSizeYF * 0.5)
			nIntCanyon = int(math.Max(float64(nIntCanyon), 34.0))
			nIntCanyon = int(math.Min(float64(nIntCanyon), 100.0*0.66666667))

			biomeRange := (0.19 - 0.2)
			terrainRange := float64(nIntCanyon - nInt)
			interpolationValue := (((nMountainBiome - 0.2) * terrainRange) / biomeRange) + float64(nInt)
			nInt = int(interpolationValue)
			if nInt <= 33 {
				nInt = 32
			}
			blockType = BT_CANYON

			if nInt > 42 && nInt < 46 {
				blockType = BT_CANYON_WHITE
			} else if nInt == 54 || nInt == 55 {
				blockType = BT_CANYON_WHITE
			} else if nInt < 54 && nInt > 45 {
				blockType = BT_CANYON_GRANITE
			} else if nInt > 45 {
				blockType = BT_CANYON_TOP
			}
			if nInt <= 32 {
				blockType = BT_WATER_CANYON
			}

		} else {
			// Canyon biome \\
			nCanyonTerrain := Noise2dPerlin(x, float64(z), -0.2, 1.75, 0.0050, 8, 0)
			nInt = int(nCanyonTerrain * chunkSizeYF * 0.5)
			nInt = int(math.Max(float64(nInt), 34.0))
			nInt = int(math.Min(float64(nInt), 100.0*0.66666667))
			blockType = BT_CANYON
			if nCanyonTerrain < 0.005 {
				blockType = BT_WATER_CANYON
				nInt -= 2
			} else if nCanyonTerrain < 0.175 {
				blockType = BT_CANYON_GROUND
				nInt -= 1
			} else if nInt < 43 {
				blockType = BT_CANYON
			} else if nInt == 43 || nInt == 44 || nInt == 45 {
				blockType = BT_CANYON_WHITE
			} else if nInt < 54 && nInt > 45 {
				blockType = BT_CANYON_GRANITE
			} else if nInt == 54 || nInt == 55 {
				blockType = BT_CANYON_WHITE
			} else if nCanyonTerrain < 0.85 {
				blockType = BT_CANYON_TOP
			} else if nCanyonTerrain < 0.95 {
				blockType = BT_CANYON_WHITE
				nInt += 1
			} else if nCanyonTerrain < 1.2 {
				blockType = BT_CANYON
				nInt += 2
			} else {
				blockType = BT_CANYON_WHITE
				nInt += 1
			}

			if nInt <= 32 {
				blockType = BT_WATER_CANYON
			}
		}
	} else if nDesertSnowBiome < 0.205 {
		if nDesertSnowBiome < 0.20 {
			// Desert Biome \\
			if nInt <= 32 {
				blockType = BT_WATER_DESERT
			} else {
				if _generateRiver(nRiverNoise, &nInt) {
					blockType = BT_WATER_DESERT
				} else {
					if Schneider(Vector2{float64(x), float64(z)}) > 0.9995 {
						// Desert tree //
						typeTreeProba := Schneider(Vector2{float64(x + 10 + (cx)*2), float64(z + (cz)*2)})
						if typeTreeProba > 0.33 {
							blockType = BT_TREE_TRUNK_PALM_1
						} else {
							blockType = BT_TREE_CACTUS
						}
					} else {
						blockType = BT_DESERT_SAND
					}
				}
			}
		} else {
			// Between Desert and Basic \\
			if nInt <= 32 {
				blockType = BT_MIX_WATER_BASIC_DESERT
			} else {
				if _generateRiver(nRiverNoise, &nInt) {
					blockType = BT_MIX_WATER_BASIC_DESERT
				} else {
					if nTerrain < 0.2625 {
						blockType = BT_MIX_BEACH_SAND
					} else {
						if nTerrain > 0.45 && nTerrain < 0.65 {
							blockType = BT_MIX_GRASS_SAND_AUTUMN
						} else {
							blockType = BT_MIX_GRASS_SAND
						}
					}
				}
			}
		}
	} else if nDesertSnowBiome > 0.7 {
		if nDesertSnowBiome < 0.705 {
			// Between Snow and Basic \\
			if nInt <= 32 {
				blockType = BT_MIX_WATER_BASIC_ICE
			} else {
				if _generateRiver(nRiverNoise, &nInt) {
					blockType = BT_MIX_WATER_BASIC_ICE
				} else {
					blockType = BT_MIX_GRASS_SNOW
				}
			}
		} else {
			// Snow Biome \\
			if nInt <= 32 {
				probaIceType := Schneider(Vector2{float64(x + (cx)*4), float64(z - (cz)*6)}) * 10
				if probaIceType < 8 {
					blockType = BT_ICE_1
				} else if probaIceType < 9.33 {
					blockType = BT_ICE_2
				} else {
					blockType = BT_ICE_3
				}
			} else {
				if _generateRiver(nRiverNoise, &nInt) {
					probaIceType := Schneider(Vector2{float64(x + (cx)*4), float64(z - (cz)*6)}) * 10
					if probaIceType < 8 {
						blockType = BT_ICE_1
					} else if probaIceType < 9.33 {
						blockType = BT_ICE_2
					} else {
						blockType = BT_ICE_3
					}
				} else {
					if nTerrain > 0.55 {
						// Snow tree //
						if Schneider(Vector2{float64(x), float64(z)}) > 0.99 {
							blockType = BT_TREE_TRUNK_DARK
						} else {
							blockType = BT_SNOW
						}
					} else {
						// Pink snow tree //
						if Schneider(Vector2{float64(x), float64(z)}) > 0.9975 {
							blockType = BT_TREE_TRUNK_GREY
						} else {
							blockType = BT_SNOW
						}
					}
				}
			}
		}
	} else {
		// Basic Biome \\
		if nInt <= 32 {
			blockType = BT_WATER
		} else if nTerrain < 0.2625 {
			if _generateRiver(nRiverNoise, &nInt) {
				blockType = BT_WATER
			} else {
				blockType = BT_BEACH
			}
		} else if nTerrain > 0.45 && nTerrain < 0.65 {
			if _generateRiver(nRiverNoise, &nInt) {
				blockType = BT_WATER
			} else {
				// Autumn tree //
				if Schneider(Vector2{float64(x), float64(z)}) > 0.99 {
					blockType = BT_TREE_TRUNK_DARK
				} else {
					blockType = BT_GRASS_AUTUMN
				}
			}
		} else if nTerrain < 0.85 {
			if _generateRiver(nRiverNoise, &nInt) {
				blockType = BT_WATER
			} else {
				// Basic tree //
				if Schneider(Vector2{float64(x), float64(z)}) > 0.9975 {
					blockType = BT_TREE_TRUNK
				} else {
					blockType = BT_GRASS
				}
			}
		} else {
			blockType = BT_MOUNTAIN_SNOW
		}
	}

	if chunk == nil {
		return int16(nInt)
	}
	// Generate clouds \\
	var nClouds float64 = 0
	if nMountainBiome < 0.19 {
		nClouds = 0
	} else if nDesertSnowBiome > 0.7 {
		nClouds = Noise2dSimplex(x, float64(z), 0.0, 1.0, 0.0075, 2, 4)
	} else if nDesertSnowBiome < 0.25 {
		nClouds = Noise2dSimplex(x, float64(z), 0.0, 0.65, 0.005, 2, 4)
	} else {
		nClouds = Noise2dSimplex(x, float64(z), 0.0, 1.0, 0.010, 2, 4)
	}
	nClouds = math.Pow(nClouds, 3)

	if nClouds > 0.5 {
		chunk.voxel[V3_Vec3i16{int16(cx), int16(200), int16(cz)}] = V3_VoxelData{
			V32_NO_FILL, 0xffffff, BT_CLOUD,
		}
		if nClouds > 0.9 {
			chunk.voxel[V3_Vec3i16{int16(cx), int16(200 + 2), int16(cz)}] = V3_VoxelData{
				V32_NO_FILL, 0xffffff, BT_CLOUD,
			}
			chunk.voxel[V3_Vec3i16{int16(cx), int16(200 - 2), int16(cz)}] = V3_VoxelData{
				V32_NO_FILL, 0xffffff, BT_CLOUD,
			}
		} else if nClouds > 0.7 {
			chunk.voxel[V3_Vec3i16{int16(cx), int16(200 + 1), int16(cz)}] = V3_VoxelData{
				V32_NO_FILL, 0xffffff, BT_CLOUD,
			}
			chunk.voxel[V3_Vec3i16{int16(cx), int16(200 - 1), int16(cz)}] = V3_VoxelData{
				V32_NO_FILL, 0xffffff, BT_CLOUD,
			}

		}
	}

	// Billow noise \\
	colorgb := color.RGBA{255, 255, 255, 1}
	colorgb = generateVarianteColorBillowNoise(colorgb, x, z)
	rgba := uint32(uint32(colorgb.R)<<16 | uint32(colorgb.G)<<8 | uint32(colorgb.B))

	if nInt <= 32 {
		chunk.voxel[V3_Vec3i16{int16(cx), int16(32), int16(cz)}] = V3_VoxelData{
			V32_NO_FILL, rgba, blockType,
		}
	}

	nHoleNoise := Noise2dSimplex(float64(x), float64(z), 0.0, 0.65, 0.0045, 2, 4)
	nHoleNoise = math.Pow(nHoleNoise, 2.0)
	if nHoleNoise <= .825 || nInt <= 33 || nMountainBiome < 0.2 || (nRiverNoise > 0.335 && nRiverNoise < 0.355) {
		// create final cube //
		chunk.voxel[V3_Vec3i16{int16(cx), int16(nInt), int16(cz)}] = V3_VoxelData{
			0x00, rgba, blockType,
		}
	} // else -> access to cavern

	return int16(nInt)
}
