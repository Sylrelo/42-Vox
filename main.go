package main

import (
	"flag"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"runtime"
	"time"
	"vox/glfont"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	windowWidth     = int32(2560)
	windowHeight    = int32(1440)
	shadowMapWidth  = windowWidth * 2
	shadowMapHeight = windowHeight * 2
	HDPI_MULTIPLIER = int32(2)
	MAX_CUBE_VISION = float32(240)
)

const (
	chunkSizeYF  = 255
	chunkSize    = 26
	chunkSizeF   = 26.0
	chunkSizeInv = 1.0 / chunkSizeF
)

func init() {
	runtime.LockOSThread()
}

func GetCameraMatrix(vox Vox) mgl32.Mat4 {
	camera := mgl32.HomogRotate3D(vox.rot.X(), mgl32.Vec3{1, 0, 0})
	camera = camera.Mul4(mgl32.HomogRotate3D(vox.rot.Y(), mgl32.Vec3{0, 1, 0}))
	camera = camera.Mul4(mgl32.Translate3D(vox.pos.X(), vox.pos.Y(), vox.pos.Z()))

	return camera
}

func InitUniforms(ogl *OpenGL) {
	ogl.uShadowmap.light = gl.GetUniformLocation(ogl.pShadowmap, gl.Str("light\x00"))

	ogl.uDeferred.view = gl.GetUniformLocation(ogl.pDeferred, gl.Str("view\x00"))
	ogl.uDeferred.blockTextures = gl.GetUniformLocation(ogl.pDeferred, gl.Str("textures\x00"))
	ogl.uDeferred.texShadowmap = gl.GetUniformLocation(ogl.pDeferred, gl.Str("texShadowmap\x00"))
	ogl.uDeferred.light = gl.GetUniformLocation(ogl.pDeferred, gl.Str("light\x00"))
	ogl.uDeferred.matWorld = gl.GetUniformLocation(ogl.pDeferred, gl.Str("matWorld\x00"))

	ogl.uQuads.texColor = gl.GetUniformLocation(ogl.pQuads, gl.Str("texColor\x00"))
	ogl.uQuads.texNormal = gl.GetUniformLocation(ogl.pQuads, gl.Str("texNormal\x00"))
	ogl.uQuads.texPosition = gl.GetUniformLocation(ogl.pQuads, gl.Str("texPosition\x00"))
	ogl.uQuads.texShadowmap = gl.GetUniformLocation(ogl.pQuads, gl.Str("texShadowmap\x00"))
	ogl.uQuads.texSsao = gl.GetUniformLocation(ogl.pQuads, gl.Str("texSsao\x00"))

	ogl.uSsao.projection = gl.GetUniformLocation(ogl.pSsao, gl.Str("projection\x00"))
	ogl.uSsao.texPosition = gl.GetUniformLocation(ogl.pSsao, gl.Str("texPosition\x00"))
	ogl.uSsao.texNormal = gl.GetUniformLocation(ogl.pSsao, gl.Str("texNormal\x00"))

	ogl.uSkybox.texSkybox = gl.GetUniformLocation(ogl.pSkybox, gl.Str("texSkybox\x00"))
	ogl.uSkybox.view = gl.GetUniformLocation(ogl.pSkybox, gl.Str("view\x00"))
}

func LoadShaders(ogl *OpenGL) {
	var err error

	ogl.pShadowmap, err = CreateProgram("Shaders/sm.vertex.glsl", "Shaders/sm.fragment.glsl")
	if err != nil {
		panic(err)
	}

	ogl.pQuads, err = CreateProgram("Shaders/qd.vertex.glsl", "Shaders/qd.fragment.glsl")
	if err != nil {
		panic(err)
	}

	ogl.pDeferred, err = CreateProgram("Shaders/deferred.vertex.glsl", "Shaders/deferred.fragment.glsl")
	if err != nil {
		panic(err)
	}

	ogl.pSkybox, err = CreateProgram("Shaders/skybox.vertex.glsl", "Shaders/skybox.fragment.glsl")
	if err != nil {
		panic(err)
	}

	ogl.pSsao, err = CreateProgram("Shaders/ssao.vertex.glsl", "Shaders/ssao.fragment.glsl")
	if err != nil {
		panic(err)
	}
}

func main() {
	var soundEngine SoundEngine
	var ogl OpenGL
	var vox Vox
	var window *glfw.Window
	var err error

	customSeed := flag.Float64("seed", 42, "Seed used for procedural generation")
	customPosX := flag.Float64("px", 49272.25, "Starting X position")
	customPosY := flag.Float64("py", 90, "Starting Y position")
	customPosZ := flag.Float64("pz", 82937.15, "Starting Z position")
	isFullscreen := flag.Bool("fs", true, "Set fullscreen mode")
	maxVision := flag.Float64("maxVision", 240, "Set max cube vision")

	flag.Parse()

	vox.pos = mgl32.Vec3{float32(*customPosX), -float32(*customPosY), float32(*customPosZ)}
	currentSeed = *customSeed

	if *maxVision >= 500 {
		MAX_CUBE_VISION = 500
	} else if *maxVision < 0 {
		MAX_CUBE_VISION = 0
	} else {
		MAX_CUBE_VISION = float32(*maxVision)
	}

	NoiseInitPermtables(currentSeed)

	// InitMusics(&soundEngine)

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
		os.Exit(0)
	}
	defer glfw.Terminate()

	vidMode := glfw.GetPrimaryMonitor().GetVideoMode()

	windowWidth = int32(vidMode.Width)
	windowHeight = int32(vidMode.Height)

	windowWidth = 2560
	windowHeight = 1440

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	if *isFullscreen {
		HDPI_MULTIPLIER = 1
		windowWidth = 2560
		windowHeight = 1440
		window, err = glfw.CreateWindow(int(windowWidth), int(windowHeight), "Vox", glfw.GetPrimaryMonitor(), nil)
	} else {
		window, err = glfw.CreateWindow(int(windowWidth), int(windowHeight), "Vox", nil, nil)
	}

	if err != nil {
		panic(err)
	}
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	err = gl.Init()

	if err != nil {
		panic(err)
	}

	InitShadowmap(&ogl)
	InitTextures(window, &ogl)
	BaseRessourcePack()
	InitSkybox(&ogl)
	InitSSAO(&ogl)
	InitDeferred(&ogl)

	LoadShaders(&ogl)
	InitUniforms(&ogl)

	gl.Disable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.GenVertexArrays(1, &ogl.vao)

	vox.font, err = glfont.LoadFont("Assets/Fonts/SourceCodePro-Regular.ttf", int32(14), int(windowWidth), int(windowHeight))
	if err != nil {
		log.Panicf("LoadFont: %v", err)
	}

	store := &V3_ChunkStore{
		remesh: make(map[Vec2i]V3_Remesh),
		chunks: make(map[Vec2i]V3_ChunksList),
		vbo:    make(map[Vec2i]V3_Vbo),
	}

	go store.HandleMeshingQueue()
	go store.CreateChunkRoutine()

	for !window.ShouldClose() {
		EventsKeyboard(window, &vox)
		EventsMouse(window, &vox)

		camera := GetCameraMatrix(vox)
		projection := mgl32.Perspective(mgl32.DegToRad(80.0), float32(windowWidth)/float32(windowHeight), 0.1, MAX_CUBE_VISION)

		EventMouseRay(window, store, vox, camera, projection)
		EventWorldReset(window, store)
		timeStart := time.Now()

		_ = soundEngine
		// handleMusic(vox, soundEngine)

		store.mutexPosition.Lock()
		store.frustrum = ExtractViewFrustrumPlanes(projection, camera)
		store.rotation = vox.rot
		store.position = vox.pos
		store.CreateEmptyChunk()
		store.UnloadUnseenChunk()
		store.mutexPosition.Unlock()

		store.HandleBuffers()
		store.Display(ogl, camera, projection, vox)

		window.SwapBuffers()
		glfw.PollEvents()

		window.SetTitle(
			fmt.Sprintf("ft_vox | %4.1f ms %3d fps",
				float32(time.Since(timeStart).Microseconds()/100.0)*0.1,
				1000/(time.Since(timeStart).Milliseconds()+1),
			))
	}

}
