package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

func InitDeferred(ogl *OpenGL) {

	width := windowWidth
	height := windowHeight

	gl.GenFramebuffers(1, &ogl.fboDeferred)
	gl.BindFramebuffer(gl.FRAMEBUFFER, ogl.fboDeferred)

	gl.GenTextures(1, &ogl.texPosition)
	gl.BindTexture(gl.TEXTURE_2D, ogl.texPosition)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB16F, width, height, 0, gl.RGB, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, ogl.texPosition, 0)

	gl.GenTextures(1, &ogl.texNormal)
	gl.BindTexture(gl.TEXTURE_2D, ogl.texNormal)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA16F, width, height, 0, gl.RGBA, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+1, gl.TEXTURE_2D, ogl.texNormal, 0)

	gl.GenTextures(1, &ogl.texColor)
	gl.BindTexture(gl.TEXTURE_2D, ogl.texColor)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+2, gl.TEXTURE_2D, ogl.texColor, 0)

	gl.GenTextures(1, &ogl.texDepth)
	gl.BindTexture(gl.TEXTURE_2D, ogl.texDepth)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT32F, width, height, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, ogl.texDepth, 0)

	attachments := [3]uint32{gl.COLOR_ATTACHMENT0, gl.COLOR_ATTACHMENT1, gl.COLOR_ATTACHMENT2}
	gl.DrawBuffers(3, &attachments[0])
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
}

func GeometryStage(ogl OpenGL) {
	gl.Viewport(0, 0, windowWidth, windowHeight)
	gl.BindFramebuffer(gl.FRAMEBUFFER, ogl.fboDeferred)
	gl.DepthRange(0, 0.99)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(ogl.pDeferred)
	gl.Uniform1i(ogl.uDeferred.blockTextures, 0)
	gl.Uniform1i(ogl.uDeferred.texShadowmap, 1)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ogl.blockTextures)
	gl.ActiveTexture(gl.TEXTURE0 + 1)
	gl.BindTexture(gl.TEXTURE_2D, ogl.texShadowmap)
}
