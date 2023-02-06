package opengl

import (
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Program struct {
	Id uint32
}

func (p *Program) Validate() {
	gl.ValidateProgram(p.Id)
	var params int32 = -1
	gl.GetProgramiv(p.Id, gl.VALIDATE_STATUS, &params)
	if params != gl.TRUE {
		log := GetInfoLog(p.Id, gl.GetProgramiv, gl.GetProgramInfoLog)
		fmt.Print(log)
		os.Exit(1)
	}
}
