package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func _getShaderSource(file string) string {

	data, error := ioutil.ReadFile(file)
	if error != nil {
		fmt.Println("Could not open", file)
		os.Exit(1)
	}
	return string(data) + "\x00"
}

func _compileShader(source string, shaderType uint32) (uint32, error) {
	var status int32
	var logLength int32

	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)

	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func _newProgram(vertex uint32, fragment uint32) (uint32, error) {
	var status int32
	var logLength int32

	program := gl.CreateProgram()

	gl.AttachShader(program, vertex)
	gl.AttachShader(program, fragment)
	gl.LinkProgram(program)

	gl.GetProgramiv(program, gl.LINK_STATUS, &status)

	if status == gl.FALSE {
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertex)
	gl.DeleteShader(fragment)

	return program, nil
}

func CreateProgram(vertex, fragment string) (uint32, error) {
	var err error
	var cVertex, cFragment, program uint32

	vertx := _getShaderSource(vertex)
	fragx := _getShaderSource(fragment)

	cVertex, err = _compileShader(vertx, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	cFragment, err = _compileShader(fragx, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	fmt.Printf("%-30s %p\n", vertex, &cVertex)
	fmt.Printf("%-30s %p\n", fragment, &cFragment)

	program, err = _newProgram(cVertex, cFragment)
	if err != nil {
		return 0, err
	}
	return program, nil
}
