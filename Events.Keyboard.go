package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func _keyboardTranslate(ctx *glfw.Window, vox *Vox, multiplier float32) {
	move := mgl32.Vec3{0, 0, 0}

	if ctx.GetKey(glfw.KeyW) == 1 {
		move = move.Add(mgl32.Vec3{0, 0, 0.10 * multiplier})
	}
	if ctx.GetKey(glfw.KeyS) == 1 {
		move = move.Add(mgl32.Vec3{0, 0, -0.10 * multiplier})
	}
	if ctx.GetKey(glfw.KeyA) == 1 {
		move = move.Add(mgl32.Vec3{0.10 * multiplier, 0, 0})

	}
	if ctx.GetKey(glfw.KeyD) == 1 {
		move = move.Add(mgl32.Vec3{-0.10 * multiplier, 0, 0})
	}
	if ctx.GetKey(glfw.KeySpace) == 1 {
		move = move.Add(mgl32.Vec3{0, -0.10 * multiplier, 0})
	}
	if ctx.GetKey(glfw.KeyLeftControl) == 1 {
		move = move.Add(mgl32.Vec3{0, 0.10 * multiplier, 0})
	}

	if move[0] != 0 || move[1] != 0 || move[2] != 0 {
		rotationMat := mgl32.HomogRotate3D(vox.rot.Y(), mgl32.Vec3{0, -1, 0})
		rotationMat = rotationMat.Mul4((mgl32.HomogRotate3D(vox.rot.X(), mgl32.Vec3{-1, 0, 0})))
		vox.pos = vox.pos.Add(rotationMat.Mul4x1(move.Vec4(1)).Vec3())
	}
}

func _keyboardRotation(ctx *glfw.Window, vox *Vox) {
	if ctx.GetKey(glfw.KeyLeft) == 1 {
		vox.rot = vox.rot.Sub(mgl32.Vec3{0, 0.05, 0})
	}
	if ctx.GetKey(glfw.KeyRight) == 1 {
		vox.rot = vox.rot.Add(mgl32.Vec3{0, 0.05, 0})
	}

	if ctx.GetKey(glfw.KeyUp) == 1 {
		vox.rot = vox.rot.Sub(mgl32.Vec3{0.05, 0, 0})
	}
	if ctx.GetKey(glfw.KeyDown) == 1 {
		vox.rot = vox.rot.Add(mgl32.Vec3{0.05, 0, 0})
	}
}

func EventsKeyboard(ctx *glfw.Window, vox *Vox) {
	var multiplier float32
	currentTime := time.Now().UnixNano() / 1000000

	if ctx.GetKey(glfw.KeyEscape) == 1 {
		os.Exit(1)
	}
	if ctx.GetKey(glfw.KeyLeftShift) == 1 {
		multiplier = 20.0
	} else if ctx.GetKey(glfw.KeyLeftAlt) == 1 {
		multiplier = 60.0
	} else {
		multiplier = 1
	}

	if ctx.GetKey(glfw.KeyGraveAccent) == 1 && currentTime-buttonTimeout > 250 {
		vox.showInformations = !vox.showInformations
		buttonTimeout = currentTime
	}
	if ctx.GetKey(glfw.Key3) == 1 && currentTime-buttonTimeout > 250 {
		vox.lightCameraMode = !vox.lightCameraMode
		buttonTimeout = currentTime
	}

	_keyboardTranslate(ctx, vox, multiplier)
	_keyboardRotation(ctx, vox)
}

func EventWorldReset(ctx *glfw.Window, store *V3_ChunkStore) {
	currentTime := time.Now().UnixNano() / 1000000

	store.mutexRemesh.Lock()
	currentRender := len(store.remesh)
	store.mutexRemesh.Unlock()

	if currentRender > 0 {
		return
	}

	if ctx.GetKey(glfw.Key1) == 1 && currentTime-buttonTimeout > 1500 {
		ResetWorld(store)
		buttonTimeout = currentTime
	}
	if ctx.GetKey(glfw.Key2) == 1 && currentTime-buttonTimeout > 1500 {
		currentSeed = float64(rand.Intn(15000-100) + 100)
		NoiseInitPermtables(currentSeed)
		time.Sleep(5000)
		ResetWorld(store)
		buttonTimeout = currentTime
	}
}
