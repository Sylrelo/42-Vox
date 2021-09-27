package main

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func InitSkybox(ogl *OpenGL) {
	gl.GenTextures(1, &ogl.texSkybox)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, ogl.texSkybox)
	path := []string{
		"Assets/Textures/base/Skybox/right.png",
		"Assets/Textures/base/Skybox/left.png",
		"Assets/Textures/base/Skybox/top.png",
		"Assets/Textures/base/Skybox/bottom.png",
		"Assets/Textures/base/Skybox/front.png",
		"Assets/Textures/base/Skybox/back.png",
	}
	for i, f := range path {
		img, w, h := LoadImage(f)
		gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), 0, gl.RGBA, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&img[0]))
	}
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	cubeSkybox := []float32{
		-1.0, 1.0, -1.0,
		-1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		-1.0, 1.0, -1.0,
		-1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, -1.0,
		-1.0, 1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, -1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		-1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, 1.0,
		-1.0, 1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, -1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		-1.0, -1.0, 1.0,
		1.0, -1.0, 1.0}
	gl.GenVertexArrays(1, &ogl.vaoSkybox)
	gl.BindVertexArray(ogl.vaoSkybox)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)

	gl.GenBuffers(1, &ogl.vboSkybox)
	gl.BindBuffer(gl.ARRAY_BUFFER, ogl.vboSkybox)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeSkybox)*4, gl.Ptr(cubeSkybox), gl.STATIC_DRAW)
}

func SkyboxStage(ogl OpenGL, camera, projection mgl32.Mat4) {
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, ogl.fboDeferred)
	gl.BlitFramebuffer(0, 0, windowWidth, windowHeight, 0, 0, windowWidth*HDPI_MULTIPLIER, windowHeight*HDPI_MULTIPLIER, gl.DEPTH_BUFFER_BIT, gl.NEAREST)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.UseProgram(ogl.pSkybox)
	gl.DepthRange(0.99, 1)
	gl.BindVertexArray(ogl.vaoSkybox)
	gl.BindBuffer(gl.ARRAY_BUFFER, ogl.vboSkybox)
	camera.SetCol(3, mgl32.Vec4{0, 0, 0, 1})
	view := projection.Mul4(camera)
	gl.UniformMatrix4fv(ogl.uSkybox.view, 1, false, &view[0])
	gl.Uniform1i(ogl.uSkybox.texSkybox, 0)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ogl.texSkybox)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
}
