package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	ID uint32
}

func (shader *Shader) Use() {
	gl.UseProgram(shader.ID)
}

func (shader *Shader) UniformSetMat4(uniformName string, mat4 *mgl32.Mat4) {
	uniform := gl.GetUniformLocation(shader.ID, GLString(uniformName))
	gl.UniformMatrix4fv(uniform, 1, false, &mat4[0])
}

func (shader *Shader) LoadFile(fileName string) {
	content, err := os.ReadFile(fileName + ".glsl_vert")
	if err != nil {
		log.Fatal(err)
	}
	vertexShader := string(content)

	content, err = os.ReadFile(fileName + ".glsl_frag")
	if err != nil {
		log.Fatal(err)
	}
	fragmentShader := string(content)

	// Configure the vertex and fragment shaders
	program, err := shader.CompileSource(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	shader.ID = program
}

func (shader *Shader) CompileSource(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := shader.CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := shader.CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func (shaderObj *Shader) CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
