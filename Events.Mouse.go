package main

import (
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var _oldMousePosX, _oldMousePosY float64

func EventMouseRay(ctx *glfw.Window, store *V3_ChunkStore, vox Vox, camera, projection mgl32.Mat4) {
	currentTime := time.Now().UnixNano() / 1000000

	if ctx.GetMouseButton(glfw.MouseButton1) == 1 && currentTime-buttonTimeout >= 250 {
		store.SendDeleteRay(vox, camera, projection)
		buttonTimeout = currentTime
	}
}

func _mouseMovements(ctx *glfw.Window, vox *Vox) {
	posX, posY := ctx.GetCursorPos()

	if _oldMousePosX == 0 {
		_oldMousePosX = posX
	}
	if _oldMousePosY == 0 {
		_oldMousePosY = posY
	}

	vox.rot = vox.rot.Add(mgl32.Vec3{
		-float32((_oldMousePosY - posY) * 0.001),
		-float32((_oldMousePosX - posX) * 0.001),
		0,
	})
	if vox.rot[0] > 1.0 {
		vox.rot[0] = 1.0
	}
	if vox.rot[0] < -1.0 {
		vox.rot[0] = -1.0
	}
	_oldMousePosX = posX
	_oldMousePosY = posY
}

func EventsMouse(ctx *glfw.Window, vox *Vox) {
	_mouseMovements(ctx, vox)
}
