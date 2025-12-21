// Implements OpenGL shader program management.
// The Shader struct handles compilation, linking, and uniform management
// for GLSL vertex and fragment shaders.

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Shader represents an OpenGL shader program consisting of vertex and fragment shaders.
// It provides methods for loading, compiling, and setting uniform variables.
type Shader struct {
	ID uint32 // OpenGL shader program ID
}

// Use activates this shader program for subsequent rendering calls.
func (shader *Shader) Use() {
	gl.UseProgram(shader.ID)
}

// UniformSetMat4 sets a mat4 (4x4 matrix) uniform in the shader program.
// uniformName: Name of the uniform variable in the GLSL shader
// mat4: Pointer to the 4x4 matrix to upload
func (shader *Shader) UniformSetMat4(uniformName string, mat4 *mgl32.Mat4) {
	uniform := gl.GetUniformLocation(shader.ID, GLString(uniformName))
	gl.UniformMatrix4fv(uniform, 1, false, &mat4[0])
}

// UniformSetVec3 sets a vec3 (3-component vector) uniform in the shader program.
// uniformName: Name of the uniform variable in the GLSL shader
// vec3: Pointer to the 3-component vector to upload
func (shader *Shader) UniformSetVec3(uniformName string, vec3 *mgl32.Vec3) {
	uniform := gl.GetUniformLocation(shader.ID, GLString(uniformName))
	gl.Uniform3f(uniform, vec3[0], vec3[1], vec3[2])
}

// LoadFile loads vertex and fragment shaders from files and compiles them into a program.
// fileName: Base name of the shader files (without extension)
// Expected files: fileName.glsl_vert (vertex shader) and fileName.glsl_frag (fragment shader)
func (shader *Shader) LoadFile(fileName string) {
	// Read vertex shader source code
	content, err := os.ReadFile(fileName + ".glsl_vert")
	if err != nil {
		log.Fatal(err)
	}
	vertexShader := string(content)

	// Read fragment shader source code
	content, err = os.ReadFile(fileName + ".glsl_frag")
	if err != nil {
		log.Fatal(err)
	}
	fragmentShader := string(content)

	// Compile and link shaders into a program
	program, err := shader.CompileSource(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	shader.ID = program
}

// CompileSource compiles vertex and fragment shader source code and links them into a program.
// vertexShaderSource: GLSL source code for the vertex shader
// fragmentShaderSource: GLSL source code for the fragment shader
// Returns: OpenGL program ID or error if compilation/linking fails
func (shader *Shader) CompileSource(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	// Compile vertex shader
	vertexShader, err := shader.CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	// Compile fragment shader
	fragmentShader, err := shader.CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	// Create shader program
	program := gl.CreateProgram()

	// Attach shaders to program
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)

	// Link program (combines shaders into executable)
	gl.LinkProgram(program)

	// Check linking status
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		// Get error log
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	// Clean up individual shaders (they're now part of the program)
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

// CompileShader compiles a single shader from source code.
// source: GLSL source code for the shader
// shaderType: Type of shader (gl.VERTEX_SHADER, gl.FRAGMENT_SHADER, etc.)
// Returns: OpenGL shader ID or error if compilation fails
func (shaderObj *Shader) CompileShader(source string, shaderType uint32) (uint32, error) {
	// Create shader object
	shader := gl.CreateShader(shaderType)

	// Upload source code to GPU
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free() // Free C strings

	// Compile shader
	gl.CompileShader(shader)

	// Check compilation status
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		// Get compilation error log
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile shader (type %v): %v", shaderType, log)
	}

	return shader, nil
}
