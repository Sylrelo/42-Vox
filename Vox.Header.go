package main

import (
	"sync"
	"vox/glfont"

	"github.com/go-gl/mathgl/mgl32"
)

type V3_Vec3u8 [3]uint8
type V3_Vec3i8 [3]int8
type V3_Vec3i16 [3]int16
type Vec2i [2]int

var (
	currentSeed   = float64(42)
	buttonTimeout = int64(0)
)

// Chunk Flag
const (
	V3_C_NONE    = uint8(0x00)
	V3_C_WAITING = uint8(0x01)
	V3_C_BUILT   = uint8(0x02)
	V3_C_READY   = uint8(0x04)
	V3_C_MESHED  = uint8(0x08)
	V3_C_IN_VIEW = uint8(0x10)
	V3_C_GLOBY   = uint8(0x20)
	V3_C_MESHING = uint8(0x40)
	V3_C_REMESH  = uint8(0x80)
)

// Block Flag
const (
	V3_INACTIVE    = 0x00
	V32_FILLER     = 0x01
	V32_NO_FILL    = 0x02
	V32_SPRITE     = 0x04
	V32_NO_OCCLUDE = 0x08
	V32_NO_GREEDY  = 0x10
)

// Textured Block Type
const (
	BT_NONE                       = iota
	BT_GRASS                      = iota
	BT_GRASS_AUTUMN               = iota
	BT_CAVERN                     = iota
	BT_CAVERN_2                   = iota
	BT_CAVERN_GOLD_1              = iota
	BT_CAVERN_GOLD_2              = iota
	BT_CAVERN_LAPIS_1             = iota
	BT_CAVERN_LAPIS_2             = iota
	BT_TREE_TRUNK                 = iota
	BT_TREE_TRUNK_DARK            = iota
	BT_TREE_TRUNK_GREY            = iota
	BT_TREE_TRUNK_PALM_1          = iota
	BT_TREE_TRUNK_PALM_2          = iota
	BT_TREE_CACTUS                = iota
	BT_TREE_LEAVES_1              = iota
	BT_TREE_LEAVES_2              = iota
	BT_TREE_LEAVES_HIGH_DENSITY_1 = iota
	BT_TREE_LEAVES_HIGH_DENSITY_2 = iota
	BT_TREE_LEAVES_PINK_1         = iota
	BT_TREE_LEAVES_PINK_2         = iota
	BT_TREE_LEAVES_PALM_1         = iota
	BT_TREE_LEAVES_PALM_2         = iota
	BT_TREE_LEAVES_SNOW_1         = iota
	BT_TREE_LEAVES_SNOW_2         = iota
	BT_FOLLIAGE_HOUSTONIA         = iota
	BT_FOLLIAGE_GRASS_BIG         = iota
	BT_FOLLIAGE_GRASS_SMALL       = iota
	BT_FOLLIAGE_GRASS_FERN        = iota
	BT_FOLLIAGE_AUTUMN_SMALL      = iota
	BT_FOLLIAGE_AUTUMN_BIG_1      = iota
	BT_FOLLIAGE_AUTUMN_BIG_2      = iota
	BT_FOLLIAGE_AUTUMN_FERN_1     = iota
	BT_FOLLIAGE_AUTUMN_FERN_2     = iota
	BT_FOLLIAGE_TULIP             = iota
	BT_FOLLIAGE_LILY              = iota
	BT_FOLLIAGE_ORCHID            = iota
	BT_FOLLIAGE_MUSHROOM_RED      = iota
	BT_FOLLIAGE_MUSHROOM_BROWN    = iota
	BT_FOLLIAGE_MUSHROOM_FUNGUS   = iota
	BT_FOLLIAGE_DAISY             = iota
	BT_FOLLIAGE_DEAD_BUSH         = iota
	BT_FOLLIAGE_BUSH              = iota
	BT_FOLLIAGE_SAPLING           = iota
	BT_FOLLIAGE_CACTUS_1          = iota
	BT_FOLLIAGE_CACTUS_2          = iota
	BT_CANYON                     = iota
	BT_CANYON_GROUND              = iota
	BT_CANYON_GRANITE             = iota
	BT_CANYON_WHITE               = iota
	BT_CANYON_TOP                 = iota
	BT_WATER                      = iota
	BT_WATER_DESERT               = iota
	BT_WATER_CANYON               = iota
	BT_ICE_1                      = iota
	BT_ICE_2                      = iota
	BT_ICE_3                      = iota
	BT_DESERT_SAND                = iota
	BT_BEACH                      = iota
	BT_SNOW                       = iota
	BT_MOUNTAIN                   = iota
	BT_MOUNTAIN_MOSSY             = iota
	BT_MOUNTAIN_SNOW              = iota
	BT_MIX_WATER_BASIC_ICE        = iota
	BT_MIX_WATER_BASIC_DESERT     = iota
	BT_MIX_GRASS_SNOW             = iota
	BT_MIX_GRASS_SAND             = iota
	BT_MIX_BEACH_SAND             = iota
	BT_MIX_GRASS_SAND_AUTUMN      = iota
	BT_CLOUD                      = iota
	BT_DEBUG                      = iota
	BT_DEBUG_R                    = iota
	BT_DEBUG_G                    = iota
	BT_DEBUG_B                    = iota
	BLOCK_COUNT                   = iota
)

// Face Identifier
const (
	V3_FACE_LEFT    = 0
	V3_FACE_RIGHT   = 1
	V3_FACE_TOP     = 2
	V3_FACE_BOTTOM  = 3
	V3_FACE_FRONT   = 4
	V3_FACE_BACK    = 5
	V3_SPRITE_RIGHT = 6
	V3_SPRITE_LEFT  = 7
)

type V3_VoxelData struct {
	flags     uint8
	color     uint32
	blockType uint8
}

type V3_ChunksList struct {
	flag  uint8
	maxY  [chunkSize][chunkSize]uint8
	voxel map[V3_Vec3i16]V3_VoxelData
}

type V3_Remesh struct {
	done        bool
	glBufferLen int
	glBuffer    []int32
}

type V3_Vbo struct {
	glVbo       uint32
	glBufferLen int
}

type V3_ChunkStore struct {
	chunks        map[Vec2i]V3_ChunksList
	remesh        map[Vec2i]V3_Remesh
	vbo           map[Vec2i]V3_Vbo
	position      mgl32.Vec3
	rotation      mgl32.Vec3
	frustrum      FrustrumPlanes
	mutexRemesh   sync.RWMutex
	mutexChunks   sync.RWMutex
	mutexPosition sync.RWMutex
	mutexVbo      sync.RWMutex
}

type Uniforms struct {
	light      int32
	view       int32
	projection int32
	matWorld   int32

	texShadowmap  int32
	blockTextures int32

	texPosition int32
	texNormal   int32
	texColor    int32
	texSkybox   int32
	texSsao     int32
}

type OpenGL struct {
	vao    uint32
	pQuads uint32
	uQuads Uniforms

	uShadowmap   Uniforms
	pShadowmap   uint32
	fboShadowmap uint32
	texShadowmap uint32

	blockTextures uint32

	pDeferred   uint32
	fboDeferred uint32
	texPosition uint32
	texNormal   uint32
	texColor    uint32
	texDepth    uint32
	uDeferred   Uniforms

	pSkybox   uint32
	uSkybox   Uniforms
	texSkybox uint32
	vboSkybox uint32
	vaoSkybox uint32

	fboSsao uint32
	texSsao uint32
	pSsao   uint32
	uSsao   Uniforms
}

type Vox struct {
	pos              mgl32.Vec3
	rot              mgl32.Vec3
	showInformations bool
	lightCameraMode  bool
	font             *glfont.Font
}
