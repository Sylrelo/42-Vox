package main

func grossePatate(v int32) int32 {
	return v + 100
}

func _shiftPosition(x, y, z, ox, oy, oz, col int16) int32 {
	return int32((int32(y+oy) << 23) | (grossePatate(int32(x+ox)) << 15) | grossePatate(int32(z+oz))<<7 | int32(col))
}

func _shiftTexCoords(x, y int16, billow, texture, faceType int) int32 {
	return int32((int32(x) << 26) | (int32(y) << 20) | (int32(billow&0xFF) << 12) | (int32(texture) << 5) | (int32(faceType) << 2))
}

// ############################################################

func GenerateLeftFace(x, y, z, cx, cy, cz int16, color int, blockType uint8) []int32 {
	v1 := _shiftPosition(x, y, z, 0, 0, 1+cz, 0)
	v2 := _shiftPosition(x, y, z, 0, 1+cy, 0, 0)
	v3 := _shiftPosition(x, y, z, 0, 0, 0, 0)
	v4 := _shiftPosition(x, y, z, 0, 1+cy, 1+cz, 0)

	texId := blockTextures[blockType].sides
	vt4 := _shiftTexCoords(0, 0, color, texId, V3_FACE_LEFT)
	vt1 := _shiftTexCoords(0, 1, color, texId, V3_FACE_LEFT)
	vt2 := _shiftTexCoords(1+cz, 0, color, texId, V3_FACE_LEFT)
	vt3 := _shiftTexCoords(1+cz, 1, color, texId, V3_FACE_LEFT)

	return []int32{v1, vt1, v2, vt2, v3, vt3, v1, vt1, v4, vt4, v2, vt2}
}

func GenerateRightFace(x, y, z, cx, cy, cz int16, color int, blockType uint8) []int32 {
	v1 := _shiftPosition(x, y, z, 1+cx, 0, 1+cz, 0)
	v2 := _shiftPosition(x, y, z, 1+cx, 1+cy, 0, 0)
	v3 := _shiftPosition(x, y, z, 1+cx, 0, 0, 0)
	v4 := _shiftPosition(x, y, z, 1+cx, 1+cy, 1+cz, 0)

	texId := blockTextures[blockType].sides
	vt4 := _shiftTexCoords(0, 0, color, texId, V3_FACE_RIGHT)
	vt1 := _shiftTexCoords(0, 1, color, texId, V3_FACE_RIGHT)
	vt2 := _shiftTexCoords(1+cz, 0, color, texId, V3_FACE_RIGHT)
	vt3 := _shiftTexCoords(1+cz, 1, color, texId, V3_FACE_RIGHT)

	return []int32{v1, vt1, v2, vt2, v3, vt3, v1, vt1, v4, vt4, v2, vt2}
}

func GenerateFrontFace(x, y, z, cx, cy, cz int16, color int, blockType uint8) []int32 {
	v1 := _shiftPosition(x, y, z, 0, 0, 1+cz, 0)
	v2 := _shiftPosition(x, y, z, 1+cx, 0, 1+cz, 0)
	v3 := _shiftPosition(x, y, z, 0, 1+cy, 1+cz, 0)
	v4 := _shiftPosition(x, y, z, 1+cx, 1+cy, 1+cz, 0)

	texId := blockTextures[blockType].sides
	vt1 := _shiftTexCoords(1+cx, 1, color, texId, V3_FACE_FRONT)
	vt3 := _shiftTexCoords(1+cx, 0, color, texId, V3_FACE_FRONT)
	vt4 := _shiftTexCoords(0, 0, color, texId, V3_FACE_FRONT)
	vt2 := _shiftTexCoords(0, 1, color, texId, V3_FACE_FRONT)

	return []int32{v1, vt1, v2, vt2, v3, vt3, v2, vt2, v4, vt4, v3, vt3}
}

func GenerateBackFace(x, y, z, cx, cy, cz int16, color int, blockType uint8) []int32 {
	v1 := _shiftPosition(x, y, z, 0, 0, 0, 0)
	v2 := _shiftPosition(x, y, z, 1+cx, 0, 0, 0)
	v3 := _shiftPosition(x, y, z, 0, 1+cy, 0, 0)
	v4 := _shiftPosition(x, y, z, 1+cx, 1+cy, 0, 0)

	texId := blockTextures[blockType].sides
	vt1 := _shiftTexCoords(1+cx, 1, color, texId, V3_FACE_BACK)
	vt3 := _shiftTexCoords(1+cx, 0, color, texId, V3_FACE_BACK)
	vt4 := _shiftTexCoords(0, 0, color, texId, V3_FACE_BACK)
	vt2 := _shiftTexCoords(0, 1, color, texId, V3_FACE_BACK)

	return []int32{v1, vt1, v2, vt2, v3, vt3, v2, vt2, v4, vt4, v3, vt3}
}

func GenerateTopFace(x, y, z, cx, cy, cz int16, color int, blockType uint8) []int32 {
	v1 := _shiftPosition(x, y, z, 0, 1+cy, 0, 0)
	v2 := _shiftPosition(x, y, z, 0, 1+cy, 1+cz, 0)
	v3 := _shiftPosition(x, y, z, 1+cx, 1+cy, 0, 0)
	v4 := _shiftPosition(x, y, z, 1+cx, 1+cy, 1+cz, 0)

	texId := blockTextures[blockType].top
	vt1 := _shiftTexCoords(0, 0, color, texId, V3_FACE_TOP)
	vt2 := _shiftTexCoords(0, 1, color, texId, V3_FACE_TOP)
	vt3 := _shiftTexCoords(1+cx, 0, color, texId, V3_FACE_TOP)
	vt4 := _shiftTexCoords(1+cx, 1, color, texId, V3_FACE_TOP)

	return []int32{v1, vt1, v2, vt2, v3, vt3, v3, vt3, v2, vt2, v4, vt4}
}

func GenerateBottomFace(x, y, z, cx, cy, cz int16, color int, blockType uint8) []int32 {
	a := _shiftPosition(x, y, z, 0, 0, 0, 0)
	b := _shiftPosition(x, y, z, 0, 0, 1+cz, 0)
	c := _shiftPosition(x, y, z, 1+cx, 0, 1+cz, 0)
	d := _shiftPosition(x, y, z, 1+cx, 0, 0, 0)

	texId := blockTextures[blockType].bottom
	tc := _shiftTexCoords(0, 0, color, texId, V3_FACE_BOTTOM)
	td := _shiftTexCoords(0, 1, color, texId, V3_FACE_BOTTOM)
	ta := _shiftTexCoords(1+cx, 1, color, texId, V3_FACE_BOTTOM)
	tb := _shiftTexCoords(1+cx, 0, color, texId, V3_FACE_BOTTOM)

	return []int32{a, ta, b, tb, d, td, b, tb, d, td, c, tc}
}

func GenerateSpriteBlock(x, y, z int16, color int, blockType uint8) []int32 {
	a := _shiftPosition(x, y, z, 0, 1, 0, 0)
	b := _shiftPosition(x, y, z, 1, 1, 1, 0)
	c := _shiftPosition(x, y, z, 1, 0, 1, 0)
	d := _shiftPosition(x, y, z, 0, 0, 0, 0)

	a2 := _shiftPosition(x, y, z, 0, 1, 1, 0)
	b2 := _shiftPosition(x, y, z, 1, 1, 0, 0)
	c2 := _shiftPosition(x, y, z, 1, 0, 0, 0)
	d2 := _shiftPosition(x, y, z, 0, 0, 1, 0)

	texId := blockTextures[blockType].sides
	ta := _shiftTexCoords(0, 0, color, texId, V3_SPRITE_LEFT)
	tb := _shiftTexCoords(1, 0, color, texId, V3_SPRITE_LEFT)
	tc := _shiftTexCoords(1, 1, color, texId, V3_SPRITE_LEFT)
	td := _shiftTexCoords(0, 1, color, texId, V3_SPRITE_LEFT)
	ta2 := _shiftTexCoords(0, 0, color, texId, V3_SPRITE_RIGHT)
	tb2 := _shiftTexCoords(1, 0, color, texId, V3_SPRITE_RIGHT)
	tc2 := _shiftTexCoords(1, 1, color, texId, V3_SPRITE_RIGHT)
	td2 := _shiftTexCoords(0, 1, color, texId, V3_SPRITE_RIGHT)

	return []int32{
		a, ta, b, tb, d, td, b, tb, c, tc, d, td,
		a2, ta2, b2, tb2, d2, td2, b2, tb2, c2, tc2, d2, td2,
	}
}
