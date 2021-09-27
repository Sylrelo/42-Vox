package main

import (
	"fmt"
	"image"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unsafe"

	"image/draw"
	_ "image/jpeg"
	_ "image/png"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type BlockTexture struct {
	top    int
	bottom int
	sides  int
}

var Texture map[string]int
var blockTextures [BLOCK_COUNT]BlockTexture

func LoadImage(file string) ([]uint8, int32, int32) {
	imgFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer imgFile.Close()
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}
	rgba := image.NewRGBA(img.Bounds())

	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("Unsupported stride, only 32-bit colors supported")
	}

	width := int32(rgba.Rect.Size().X)
	height := int32(rgba.Rect.Size().Y)

	return rgba.Pix, width, height
}

func _GetTexturesInDirectory(path string) []string {
	var files []string

	fd, errOpen := os.Open("Assets/Textures/" + path)
	if errOpen != nil {
		fmt.Println(errOpen)
	}

	fileList, errReaddir := fd.Readdir(0)
	if errReaddir != nil {
		fmt.Println(errReaddir)
	}
	for _, v := range fileList {
		filePath := "Assets/Textures/" + path + v.Name()

		if v.IsDir() || filepath.Ext(filePath) != ".png" || filepath.Ext(filePath) == ".gif" || filepath.Ext(filePath) == ".jpg" {
			continue
		}
		files = append(files, filePath)
	}
	return files

}

func AssignTexture(blockType int, top, bottom, sides string) {
	blockTextures[blockType] = BlockTexture{
		Texture[top],
		Texture[bottom],
		Texture[sides],
	}
}

func InitTextures(win *glfw.Window, ogl *OpenGL) {

	fmt.Println("\x1b[94mLoading textures...\x1b[0m")
	files := _GetTexturesInDirectory("base/")
	Texture = make(map[string]int)

	var width, height, tmpWidth, tmpHeight int32
	var tmpImage, image []uint8

	for i, f := range files {
		tmpImage, tmpWidth, tmpHeight = LoadImage(f)
		if i > 0 && (tmpWidth != width || tmpHeight != height) {
			fmt.Println(f, tmpWidth, width, tmpWidth, height)
			panic(f + " is not the same size as other textures")
		}
		width = tmpWidth
		height = tmpHeight
		image = append(image, tmpImage...)

		filename := strings.TrimSuffix(filepath.Base(f), path.Ext(f))
		Texture[strings.ToUpper(filename)] = i
	}
	fmt.Println("\x1b[92mTextures successfully loaded.\x1b[0m")

	layerCount := int32(len(files))

	win.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.GenTextures(1, &ogl.blockTextures)
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, ogl.blockTextures)
	gl.TexStorage3D(gl.TEXTURE_2D_ARRAY, 1, gl.RGBA8, width, height, layerCount)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexSubImage3D(gl.TEXTURE_2D_ARRAY, 0, 0, 0, 0, width, height, layerCount, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&image[0]))

	for x := 0; x < BLOCK_COUNT; x++ {
		blockTextures[x] = BlockTexture{0, 0, 0}
	}

	gl.BindTexture(gl.TEXTURE_2D_ARRAY, ogl.blockTextures)
	gl.GenerateMipmap(gl.TEXTURE_2D_ARRAY)

}

func BaseRessourcePack() {
	AssignTexture(BT_GRASS, "GRASS_TOP", "CAVERN_BOTTOM", "GRASS_SIDES")
	AssignTexture(BT_GRASS_AUTUMN, "AUTUMN_TOP", "CAVERN_BOTTOM", "AUTUMN_SIDES")
	AssignTexture(BT_CAVERN, "CAVERN", "CAVERN", "CAVERN")
	AssignTexture(BT_CAVERN_2, "CAVERN_2", "CAVERN_2", "CAVERN_2")
	AssignTexture(BT_CAVERN_GOLD_1, "CAVERN_GOLD_1", "CAVERN_GOLD_1", "CAVERN_GOLD_2")
	AssignTexture(BT_CAVERN_GOLD_2, "CAVERN_GOLD_2", "CAVERN_GOLD_2", "CAVERN_GOLD_1")
	AssignTexture(BT_CAVERN_LAPIS_1, "CAVERN_LAPIS_1", "CAVERN_LAPIS_1", "CAVERN_LAPIS_2")
	AssignTexture(BT_CAVERN_LAPIS_2, "CAVERN_LAPIS_2", "CAVERN_LAPIS_2", "CAVERN_LAPIS_1")

	AssignTexture(BT_TREE_TRUNK, "TRUNK_TOP", "CAVERN_BOTTOM", "TRUNK_SIDE")
	AssignTexture(BT_TREE_TRUNK_DARK, "TRUNK_DARK_TOP", "CAVERN_BOTTOM", "TRUNK_DARK_SIDE")
	AssignTexture(BT_TREE_TRUNK_GREY, "SNOW_TOP", "CAVERN_BOTTOM", "TRUNK_GREY_SIDE")
	AssignTexture(BT_TREE_TRUNK_PALM_1, "PALM_TOP", "CAVERN_BOTTOM", "PALM_SIDE_1")
	AssignTexture(BT_TREE_TRUNK_PALM_2, "PALM_TOP", "CAVERN_BOTTOM", "PALM_SIDE_2")
	AssignTexture(BT_TREE_CACTUS, "CACTUS_TOP", "CAVERN_BOTTOM", "CACTUS_SIDE")

	AssignTexture(BT_TREE_LEAVES_1, "LEAF_1", "LEAF_1", "LEAF_1")
	AssignTexture(BT_TREE_LEAVES_2, "LEAF_2", "LEAF_2", "LEAF_2")
	AssignTexture(BT_TREE_LEAVES_HIGH_DENSITY_1, "LEAF_AUTUMN_1", "LEAF_AUTUMN_1", "LEAF_AUTUMN_1")
	AssignTexture(BT_TREE_LEAVES_HIGH_DENSITY_2, "LEAF_AUTUMN_2", "LEAF_AUTUMN_2", "LEAF_AUTUMN_2")
	AssignTexture(BT_TREE_LEAVES_PINK_1, "SNOW_TOP", "LEAF_PINK_1", "LEAF_PINK_1")
	AssignTexture(BT_TREE_LEAVES_PINK_2, "SNOW_TOP", "LEAF_PINK_2", "LEAF_PINK_2")
	AssignTexture(BT_TREE_LEAVES_PALM_1, "LEAF_PALM_1", "LEAF_PALM_1", "LEAF_PALM_1")
	AssignTexture(BT_TREE_LEAVES_PALM_2, "LEAF_PALM_2", "LEAF_PALM_2", "LEAF_PALM_2")

	AssignTexture(BT_TREE_LEAVES_SNOW_1, "SNOW_TOP", "LEAF_SNOW_1", "LEAF_SNOW_1")
	AssignTexture(BT_TREE_LEAVES_SNOW_2, "SNOW_TOP", "LEAF_SNOW_2", "LEAF_SNOW_2")

	AssignTexture(BT_FOLLIAGE_HOUSTONIA, "FOL_HOUSTONIA", "FOL_HOUSTONIA", "FOL_HOUSTONIA")
	AssignTexture(BT_FOLLIAGE_GRASS_SMALL, "FOL_GRASS_SMALL", "FOL_GRASS_SMALL", "FOL_GRASS_SMALL")
	AssignTexture(BT_FOLLIAGE_GRASS_BIG, "FOL_GRASS_BIG", "FOL_GRASS_BIG", "FOL_GRASS_BIG")
	AssignTexture(BT_FOLLIAGE_GRASS_FERN, "FOL_GRASS_FERN", "FOL_GRASS_FERN", "FOL_GRASS_FERN")
	AssignTexture(BT_FOLLIAGE_AUTUMN_SMALL, "FOL_AUTUMN_SMALL", "FOL_AUTUMN_SMALL", "FOL_AUTUMN_SMALL")
	AssignTexture(BT_FOLLIAGE_AUTUMN_BIG_1, "FOL_AUTUMN_BIG_1", "FOL_AUTUMN_BIG_1", "FOL_AUTUMN_BIG_1")
	AssignTexture(BT_FOLLIAGE_AUTUMN_BIG_2, "FOL_AUTUMN_BIG_2", "FOL_AUTUMN_BIG_2", "FOL_AUTUMN_BIG_2")
	AssignTexture(BT_FOLLIAGE_AUTUMN_FERN_1, "FOL_AUTUMN_FERN_1", "FOL_AUTUMN_FERN_1", "FOL_AUTUMN_FERN_1")
	AssignTexture(BT_FOLLIAGE_AUTUMN_FERN_2, "FOL_AUTUMN_FERN_2", "FOL_AUTUMN_FERN_2", "FOL_AUTUMN_FERN_2")
	AssignTexture(BT_FOLLIAGE_TULIP, "FOL_ORANGE_TULIP", "FOL_ORANGE_TULIP", "FOL_ORANGE_TULIP")
	AssignTexture(BT_FOLLIAGE_DAISY, "FOL_DAISY", "FOL_DAISY", "FOL_DAISY")
	AssignTexture(BT_FOLLIAGE_ORCHID, "FOL_ORCHID", "FOL_ORCHID", "FOL_ORCHID")
	AssignTexture(BT_FOLLIAGE_LILY, "FOL_LILY", "FOL_LILY", "FOL_LILY")
	AssignTexture(BT_FOLLIAGE_DEAD_BUSH, "FOL_DEAD_BUSH", "FOL_DEAD_BUSH", "FOL_DEAD_BUSH")
	AssignTexture(BT_FOLLIAGE_BUSH, "FOL_BUSH", "FOL_BUSH", "FOL_BUSH")
	AssignTexture(BT_FOLLIAGE_SAPLING, "FOL_SAPLING", "FOL_SAPLING", "FOL_SAPLING")
	AssignTexture(BT_FOLLIAGE_CACTUS_1, "FOL_CACTUS_1", "FOL_CACTUS_1", "FOL_CACTUS_1")
	AssignTexture(BT_FOLLIAGE_CACTUS_2, "FOL_CACTUS_2", "FOL_CACTUS_2", "FOL_CACTUS_2")
	AssignTexture(BT_FOLLIAGE_MUSHROOM_RED, "FOL_MUSHROOM_RED", "FOL_MUSHROOM_RED", "FOL_MUSHROOM_RED")
	AssignTexture(BT_FOLLIAGE_MUSHROOM_BROWN, "FOL_MUSHROOM_BROWN", "FOL_MUSHROOM_BROWN", "FOL_MUSHROOM_BROWN")
	AssignTexture(BT_FOLLIAGE_MUSHROOM_FUNGUS, "FOL_MUSHROOM_FUNGUS", "FOL_MUSHROOM_FUNGUS", "FOL_MUSHROOM_FUNGUS")

	AssignTexture(BT_CANYON_TOP, "CANYON_TOP", "CAVERN_BOTTOM", "CANYON_TOP")
	AssignTexture(BT_CANYON_WHITE, "CANYON_LINE", "CAVERN_BOTTOM", "CANYON_LINE")
	AssignTexture(BT_CANYON, "CANYON_BROWN", "CAVERN_BOTTOM", "CANYON_BROWN")
	AssignTexture(BT_CANYON_GROUND, "CANYON_GROUND", "CAVERN_BOTTOM", "CANYON_GROUND")
	AssignTexture(BT_CANYON_GRANITE, "CANYON_GRANITE", "CAVERN_BOTTOM", "CANYON_GRANITE")
	AssignTexture(BT_DESERT_SAND, "SAND", "CAVERN_BOTTOM", "SAND")
	AssignTexture(BT_BEACH, "BEACH", "CAVERN_BOTTOM", "BEACH")

	AssignTexture(BT_WATER, "WATER_BASIC", "CAVERN_BOTTOM", "WATER_BASIC")
	AssignTexture(BT_WATER_DESERT, "WATER_DESERT", "CAVERN_BOTTOM", "WATER_DESERT")
	AssignTexture(BT_WATER_CANYON, "WATER_CANYON", "CAVERN_BOTTOM", "WATER_CANYON")

	AssignTexture(BT_ICE_1, "ICE_1", "CAVERN_BOTTOM", "ICE_1")
	AssignTexture(BT_ICE_2, "ICE_2", "CAVERN_BOTTOM", "ICE_2")
	AssignTexture(BT_ICE_3, "ICE_3", "CAVERN_BOTTOM", "ICE_3")

	AssignTexture(BT_SNOW, "SNOW_TOP", "CAVERN_BOTTOM", "SNOW_SIDES")

	AssignTexture(BT_MOUNTAIN, "MOUNTAIN", "CAVERN_BOTTOM", "MOUNTAIN")
	AssignTexture(BT_MOUNTAIN_MOSSY, "MOUNTAIN_MOSSY", "CAVERN_BOTTOM", "MOUNTAIN_MOSSY")
	AssignTexture(BT_MOUNTAIN_SNOW, "SNOW_TOP", "CAVERN_BOTTOM", "MOUNTAIN_SNOW")

	AssignTexture(BT_MIX_WATER_BASIC_ICE, "MIX_WATER_BASIC_ICE", "MIX_WATER_BASIC_ICE", "MIX_WATER_BASIC_ICE")
	AssignTexture(BT_MIX_WATER_BASIC_DESERT, "MIX_WATER_BASIC_DESERT", "MIX_WATER_BASIC_DESERT", "MIX_WATER_BASIC_DESERT")

	AssignTexture(BT_MIX_GRASS_SNOW, "MIX_GRASS_SNOW_TOP", "CAVERN_BOTTOM", "MIX_GRASS_SNOW_SIDES")
	AssignTexture(BT_MIX_GRASS_SAND, "MIX_GRASS_SAND_TOP", "CAVERN_BOTTOM", "MIX_GRASS_SAND_SIDES")
	AssignTexture(BT_MIX_BEACH_SAND, "MIX_BEACH_SAND", "CAVERN_BOTTOM", "MIX_BEACH_SAND")
	AssignTexture(BT_MIX_GRASS_SAND_AUTUMN, "MIX_GRASS_SAND_AUTUMN_TOP", "CAVERN_BOTTOM", "MIX_GRASS_SAND_AUTUMN_SIDES")

	AssignTexture(BT_CLOUD, "SNOW_TOP", "SNOW_TOP", "SNOW_TOP")
	AssignTexture(BT_DEBUG, "DEBUG", "DEBUG", "DEBUG")
	AssignTexture(BT_DEBUG_R, "DEBUG_R", "DEBUG_R", "DEBUG_R")
	AssignTexture(BT_DEBUG_G, "DEBUG_G", "DEBUG_G", "DEBUG_G")
	AssignTexture(BT_DEBUG_B, "DEBUG_B", "DEBUG_B", "DEBUG_B")
}
