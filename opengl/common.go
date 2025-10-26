package opengl

import (
	"strings"

	"github.com/go-gl/gl/all-core/gl"
)

const (
	Float32ByteSize = 4
)

func GetInfoLog(
	id uint32,
	getIVFn func(uint32, uint32, *int32),
	getInfoLogFn func(uint32, int32, *int32, *uint8),
) string {
	var logLength int32
	getIVFn(id, gl.INFO_LOG_LENGTH, &logLength)
	log := strings.Repeat("\x00", int(logLength+1))
	getInfoLogFn(id, logLength, nil, gl.Str(log))
	return log
}
