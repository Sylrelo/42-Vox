package main

import (
	"math"
	"sort"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func (store *V3_ChunkStore) CreateChunkRoutine() {
	for {
		store.mutexChunks.RLock()
		var chunks []Vec2i
		chunkFlags := make(map[Vec2i]uint8, len(store.chunks))

		for c := range store.chunks {
			chunks = append(chunks, c)
			chunkFlags[c] = store.chunks[c].flag
		}
		store.mutexChunks.RUnlock()

		store.mutexPosition.RLock()
		sort.Slice(chunks, func(i, j int) bool {
			return getDistance(store.position, chunks[i]) < getDistance(store.position, chunks[j])
		})
		store.mutexPosition.RUnlock()

		var pendingCount int
		var wg sync.WaitGroup

		for c := range chunks {
			chunkId := chunks[c]
			chunkFlag := chunkFlags[chunkId]
			store.CreateNewChunks(chunkId)
			if hasBit(chunkFlag, V3_C_WAITING) && pendingCount < 2 {
				wg.Add(1)
				go store.GenerateChunk(chunkId, &wg)
				pendingCount++
			}
			if hasBit(chunkFlag, V3_C_READY) && !hasBit(chunkFlag, V3_C_MESHED) && !hasBit(chunkFlag, V3_C_MESHING) {
				store.mutexRemesh.Lock()
				store.remesh[chunks[c]] = V3_Remesh{}
				store.mutexRemesh.Unlock()
			}
		}

		if pendingCount > 0 {
			wg.Wait()
		} else {
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (store *V3_ChunkStore) StartMeshing(chunkId Vec2i, wg *sync.WaitGroup) {
	defer wg.Done()
	var meshing V3_Remesh

	store.mutexChunks.Lock()
	tmp := copyChunkData(store.chunks[chunkId])
	setBit(&tmp.flag, V3_C_MESHING)
	store.chunks[chunkId] = copyChunkData(tmp)
	store.mutexChunks.Unlock()

	store.GreedyFaces(&meshing, chunkId, &tmp)

	store.mutexRemesh.Lock()
	store.remesh[chunkId] = copyRemeshData(meshing)
	store.mutexRemesh.Unlock()

	remBit(&tmp.flag, V3_C_MESHING)
	setBit(&tmp.flag, V3_C_MESHED)
	store.mutexChunks.Lock()
	store.chunks[chunkId] = copyChunkData(tmp)
	store.mutexChunks.Unlock()

	tmp.voxel = make(map[V3_Vec3i16]V3_VoxelData)
}

func (store *V3_ChunkStore) HandleMeshingQueue() {
	var lastCount int
	var lastTime int64

	for {
		store.mutexChunks.RLock()
		meshingFlags := make(map[Vec2i]uint8, len(store.chunks))
		for c := range store.chunks {
			meshingFlags[c] = store.chunks[c].flag
		}
		store.mutexChunks.RUnlock()

		store.mutexRemesh.RLock()
		var meshingId []Vec2i
		meshingCpy := make(map[Vec2i]V3_Remesh)

		remeshCount := len(store.remesh)
		currentTime := time.Now().UnixNano() / 1000000

		if remeshCount == 0 {
			lastCount = 0
			lastTime = currentTime
		} else if remeshCount != 0 && lastCount != remeshCount {
			lastCount = remeshCount
			lastTime = currentTime
		}

		if remeshCount != 0 && lastCount == remeshCount && currentTime-lastTime >= 1000 {
			store.mutexChunks.Lock()
			for c := range store.remesh {
				delete(store.remesh, c)
				delete(store.chunks, c)
			}
			lastCount = 0
			lastTime = currentTime
			store.mutexChunks.Unlock()
		}

		for c := range store.remesh {

			if store.remesh[c].done || hasBit(meshingFlags[c], V3_C_MESHING) {
				continue
			}
			meshingId = append(meshingId, c)
			meshingCpy[c] = store.remesh[c]
		}
		store.mutexRemesh.RUnlock()

		store.mutexPosition.RLock()
		sort.Slice(meshingId, func(i, j int) bool {
			return getDistance(store.position, meshingId[i]) < getDistance(store.position, meshingId[j])
		})
		store.mutexPosition.RUnlock()

		var pendingCount int
		var wg sync.WaitGroup

		if len(meshingCpy) > 0 {
			for c := range meshingId {
				if pendingCount >= 2 {
					break
				}
				wg.Add(1)
				pendingCount++
				go store.StartMeshing(meshingId[c], &wg)
			}
			wg.Wait()
		} else {
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (store *V3_ChunkStore) HandleBuffers() {
	mesh := make(map[Vec2i]V3_Remesh)

	store.mutexRemesh.RLock()
	for c := range store.remesh {
		if !store.remesh[c].done || store.remesh[c].glBufferLen == 0 {
			continue
		}
		mesh[c] = store.remesh[c]
	}
	store.mutexRemesh.RUnlock()

	if len(mesh) > 0 {
		for c := range mesh {
			currentVbo, vboExists := store.vbo[c]
			if vboExists {
				currentVbo.glBufferLen = mesh[c].glBufferLen
				gl.BindBuffer(gl.ARRAY_BUFFER, currentVbo.glVbo)
				gl.BufferData(gl.ARRAY_BUFFER, currentVbo.glBufferLen*4, gl.Ptr(mesh[c].glBuffer), gl.STATIC_DRAW)
				store.vbo[c] = currentVbo
				store.mutexRemesh.Lock()
				delete(store.remesh, c)
				store.mutexRemesh.Unlock()
			} else {
				var vbo V3_Vbo
				vbo.glBufferLen = mesh[c].glBufferLen
				gl.GenBuffers(1, &vbo.glVbo)
				gl.BindBuffer(gl.ARRAY_BUFFER, vbo.glVbo)
				gl.BufferData(gl.ARRAY_BUFFER, vbo.glBufferLen*4, gl.Ptr(mesh[c].glBuffer), gl.STATIC_DRAW)
				store.vbo[c] = vbo
				store.mutexRemesh.Lock()
				delete(store.remesh, c)
				store.mutexRemesh.Unlock()
			}
		}
	}
}

func (store *V3_ChunkStore) RenderVbo(ogl OpenGL, projection, camera mgl32.Mat4, vox Vox, isShadow bool) (int32, int32) {
	var sortedList []Vec2i
	triangleCount := int32(0)

	store.mutexPosition.RLock()
	frustrum := store.frustrum
	position := store.position
	store.mutexPosition.RUnlock()

	maxCircleVision := getTheoricalMaxCircleVision(vox.pos, vox.rot)

	for c := range store.vbo {
		px := c[0] * chunkSize
		pz := c[1] * chunkSize
		inFrustrum := InsideFrustrum(frustrum, -px, -px+chunkSize, -pz, -pz+chunkSize)

		if isShadow && (getDistance(position, c) > (maxCircleVision * .40)) {
			continue
		}
		if !isShadow && (!inFrustrum || getDistance(position, c) > maxCircleVision) {
			continue
		}
		sortedList = append(sortedList, c)
	}

	sort.Slice(sortedList, func(i, j int) bool {
		return getDistance(position, sortedList[i]) < getDistance(position, sortedList[j])
	})

	camera.SetCol(3, mgl32.Vec4{0, 0, 0, 1})

	orthoProjection := mgl32.Ortho(-280, 280, -280, 280, -280, 280)
	lightProjection := orthoProjection.Mul4(mgl32.LookAt(1, 3, 2, 0, 0, 0, 0, 1, 0))

	for i := range sortedList {
		c := sortedList[i]

		gl.BindVertexArray(ogl.vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, store.vbo[c].glVbo)

		matTranslation := mgl32.Translate3D(vox.pos[0]-(float32(c[0]*chunkSize)), vox.pos[1], vox.pos[2]-float32(c[1]*chunkSize))

		view := projection.Mul4(camera)
		model := view.Mul4(matTranslation)
		lightCamera := lightProjection.Mul4(matTranslation)
		matWithoutProjection := camera.Mul4(matTranslation)

		if isShadow {
			gl.UniformMatrix4fv(ogl.uShadowmap.light, 1, false, &lightCamera[0])
		} else {
			if vox.lightCameraMode {
				model = lightCamera
			}
			gl.UniformMatrix4fv(ogl.uDeferred.view, 1, false, &model[0])
			gl.UniformMatrix4fv(ogl.uDeferred.light, 1, false, &lightCamera[0])
			gl.UniformMatrix4fv(ogl.uDeferred.matWorld, 1, false, &matWithoutProjection[0])
			_ = model
		}

		gl.VertexAttribPointerWithOffset(0, 1, gl.FLOAT, false, 8, 0)
		gl.VertexAttribPointerWithOffset(1, 1, gl.FLOAT, false, 8, 4)
		gl.EnableVertexAttribArray(0)
		gl.EnableVertexAttribArray(1)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(store.vbo[c].glBufferLen)/2)
		gl.DisableVertexAttribArray(0)
		gl.DisableVertexAttribArray(1)
		triangleCount += int32(store.vbo[c].glBufferLen) / 2 / 3
	}
	return triangleCount, int32(len(sortedList))
}

func (store *V3_ChunkStore) Display(ogl OpenGL, camera, projection mgl32.Mat4, vox Vox) {
	timeStart := time.Now()

	gl.ClearColor(0, 0, 0, 1)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthRange(0, 1)

	ShadowStage(ogl)
	store.RenderVbo(ogl, projection, camera, vox, true)

	GeometryStage(ogl)
	triangleCount, chunkCount := store.RenderVbo(ogl, projection, camera, vox, false)

	SsaoStage(ogl, projection)

	RenderStage(ogl)
	SkyboxStage(ogl, camera, projection)

	gl.Finish()

	if vox.showInformations {
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		vox.font.Printf(6, 20, 1.0, "Displayed Triangles  : %s", humanize.Comma(int64(triangleCount)))
		vox.font.Printf(6, 40, 1.0, "Displayed Chunks     : %d", chunkCount)

		store.mutexChunks.RLock()
		totalWaiting := 0
		totalReady := 0
		totalDone := 0

		for c := range store.chunks {
			flag := store.chunks[c].flag
			if hasBit(flag, V3_C_WAITING) {
				totalWaiting++
			}
			if hasBit(flag, V3_C_READY) {
				totalReady++
			}
			if hasBit(flag, V3_C_MESHED) {
				totalDone++
			}
		}
		vox.font.Printf(6, 80, 1.0, "Chunks Waiting   : %d", totalWaiting)
		vox.font.Printf(6, 100, 1.0, "Chunks Ready     : %d", totalReady)
		vox.font.Printf(6, 120, 1.0, "Chunks Done      : %d", totalDone)
		vox.font.Printf(6, 140, 1.0, "Chunks Total     : %d", len(store.chunks))
		store.mutexChunks.RUnlock()

		store.mutexRemesh.RLock()
		vox.font.Printf(6, 160, 1.0, "Meshing queue    : %d", len(store.remesh))
		store.mutexRemesh.RUnlock()

		vox.font.Printf(6, 240, 1.0, "Total Render Time  : %d ms", time.Since(timeStart).Milliseconds())
		if time.Since(timeStart).Milliseconds() > 0 {
			vox.font.Printf(6, 260, 1.0, "Total FPS          : %d", 1000/time.Since(timeStart).Milliseconds())
		}

		vox.font.Printf(6, 300, 1.0, "Position    : (%.2f, %.2f, %.2f)", vox.pos[0], -vox.pos[1], vox.pos[2])
		vox.font.Printf(6, 320, 1.0, "Angle       : (%.2f, %.2f, %.2f)", vox.rot[0], vox.rot[1], vox.rot[2])
		vox.font.Printf(6, 340, 1.0, "Seed        : %.f", currentSeed)

		vox.font.Printf(6, 355, 0.7, "(fps are lower with text on screen)")
	}
}

func (store *V3_ChunkStore) UnloadUnseenChunk() {
	var chunks []Vec2i

	store.mutexChunks.RLock()
	for c := range store.chunks {
		px := c[0] * chunkSize
		pz := c[1] * chunkSize
		Visible := InsideFrustrum(store.frustrum, -px, -px+chunkSize, -pz, -pz+chunkSize)

		if hasBit(store.chunks[c].flag, V3_C_WAITING) && !Visible && getDistance(store.position, c) > 2 {
			chunks = append(chunks, c)
		}
		if getDistance(store.position, c) > 10 && !Visible {
			chunks = append(chunks, c)
		}
	}
	store.mutexChunks.RUnlock()

	if len(chunks) > 0 {
		store.mutexChunks.Lock()
		store.mutexRemesh.Lock()
		for c := range chunks {
			delete(store.chunks, chunks[c])
			if chunkVbo, exists := store.vbo[chunks[c]]; exists {
				gl.DeleteBuffers(1, &chunkVbo.glVbo)
				delete(store.vbo, chunks[c])
			}
			delete(store.remesh, chunks[c])
		}
		store.mutexRemesh.Unlock()
		store.mutexChunks.Unlock()
	}
}

func (store *V3_ChunkStore) CreateEmptyChunk() {
	fx := math.Floor(float64(store.position.X() / chunkSizeF))
	fz := math.Floor(float64(store.position.Z() / chunkSizeF))

	chk := Vec2i{int(fx + 1), int(fz + 1)}

	store.mutexChunks.RLock()
	_, existsInQueue := store.chunks[chk]
	store.mutexChunks.RUnlock()
	if existsInQueue {
		return
	}
	store.mutexChunks.Lock()
	store.chunks[chk] = V3_ChunksList{flag: V3_C_WAITING}
	store.mutexChunks.Unlock()
}
