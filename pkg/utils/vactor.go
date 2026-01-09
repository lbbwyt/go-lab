package utils

import (
	"encoding/binary"
	"math"
)

func FloatsToBytes(fs []float32) []byte {
	buf := make([]byte, len(fs)*4)

	for i, f := range fs {
		u := math.Float32bits(f)
		binary.NativeEndian.PutUint32(buf[i*4:], u)
	}

	return buf
}

func ToFloat64(in []float32) []float64 {
	out := make([]float64, len(in))
	for i, v := range in {
		out[i] = float64(v)
	}
	return out
}
func ToFloat32(in []float64) []float32 {
	out := make([]float32, len(in))
	for i, v := range in {
		out[i] = float32(v)
	}
	return out
}

func Vector2Bytes(vector []float64) []byte {
	float32Arr := make([]float32, len(vector))
	for i, v := range vector {
		float32Arr[i] = float32(v)
	}
	bytes := make([]byte, len(float32Arr)*4)
	for i, v := range float32Arr {
		binary.LittleEndian.PutUint32(bytes[i*4:], math.Float32bits(v))
	}
	return bytes
}
