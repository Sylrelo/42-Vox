package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func SsaoStage(ogl OpenGL, projection mgl32.Mat4) {
	gl.UseProgram(ogl.pSsao)
	gl.BindFramebuffer(gl.FRAMEBUFFER, ogl.fboSsao)

	gl.Uniform1i(ogl.uSsao.texPosition, 0)
	gl.Uniform1i(ogl.uSsao.texNormal, 1)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ogl.texPosition)
	gl.ActiveTexture(gl.TEXTURE0 + 1)
	gl.BindTexture(gl.TEXTURE_2D, ogl.texNormal)

	gl.UniformMatrix4fv(ogl.uSsao.projection, 1, false, &projection[0])
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func InitSSAO(ogl *OpenGL) {
	width := windowWidth
	height := windowHeight

	gl.GenFramebuffers(1, &ogl.fboSsao)
	gl.BindFramebuffer(gl.FRAMEBUFFER, ogl.fboSsao)

	gl.GenTextures(1, &ogl.texSsao)
	gl.BindTexture(gl.TEXTURE_2D, ogl.texSsao)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, width, height, 0, gl.RED, gl.FLOAT, gl.Ptr(nil))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, ogl.texSsao, 0)

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
