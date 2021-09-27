package main

import "github.com/go-gl/gl/v4.1-core/gl"

func InitShadowmap(ogl *OpenGL) {
	SHADOW_WIDTH := int32(shadowMapWidth)
	SHADOW_HEIGHT := int32(shadowMapHeight)
	borderColor := []float32{1.0, 1.0, 1.0, 1.0}

	gl.GenFramebuffers(1, &ogl.fboShadowmap)
	gl.BindFramebuffer(gl.FRAMEBUFFER, ogl.fboShadowmap)
	gl.GenTextures(1, &ogl.texShadowmap)
	gl.BindTexture(gl.TEXTURE_2D, ogl.texShadowmap)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT, SHADOW_WIDTH, SHADOW_HEIGHT, 0, gl.DEPTH_COMPONENT, gl.FLOAT, gl.Ptr(nil))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &borderColor[0])
	gl.BindTexture(gl.TEXTURE_2D, ogl.texShadowmap)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, ogl.texShadowmap, 0)
	gl.DrawBuffer(gl.NONE)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func ShadowStage(ogl OpenGL) {
	gl.UseProgram(ogl.pShadowmap)
	gl.Viewport(0, 0, shadowMapWidth, shadowMapHeight)
	gl.BindFramebuffer(gl.FRAMEBUFFER, ogl.fboShadowmap)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}