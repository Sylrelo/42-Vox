package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

func RenderStage(ogl OpenGL) {
	gl.Viewport(0, 0, windowWidth*HDPI_MULTIPLIER, windowHeight*HDPI_MULTIPLIER)
	gl.UseProgram(ogl.pQuads)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.DepthRange(0, 1)
	gl.Uniform1i(ogl.uQuads.texColor, 0)
	gl.Uniform1i(ogl.uQuads.texSsao, 1)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, ogl.texColor)
	gl.ActiveTexture(gl.TEXTURE0 + 1)
	gl.BindTexture(gl.TEXTURE_2D, ogl.texSsao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)

}
