package protocol

import (
	"encoding/binary"
	"fmt"
	"github.com/howeyc/crc16"
)

//规则：
//1 byte 头 0xAA + 2 byte 温度(uint16, 小端) + 1 byte 开关(bool 0/1) + 2 byte 校验(CRC16)

type BinaryCodec struct{}

func (c BinaryCodec) Name() string { return "binary_test1" }

func (c BinaryCodec) Decode(b []byte) (map[string]interface{}, error) {
	if len(b) != 6 || b[0] != 0xAA {
		return nil, fmt.Errorf("invalid frame")
	}
	// CRC 校验
	if crc16.ChecksumCCITTFalse(b[:4]) != binary.LittleEndian.Uint16(b[4:6]) {
		return nil, fmt.Errorf("crc error")
	}
	return map[string]interface{}{
		"temp":   float32(binary.LittleEndian.Uint16(b[1:3])) / 100, // 精确到 0.01
		"switch": b[3] != 0,
	}, nil
}
