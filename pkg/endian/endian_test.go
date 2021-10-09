package endian

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
)

func TestName(t *testing.T) {
	BigEndianAndLittleEndianByLibrary()
}

func TestReadBool(t *testing.T) {
	var res bool
	var err error
	err = binary.Read(bytes.NewReader([]byte{0}), binary.BigEndian, &res)
	checkResult(t, "ReadBool", binary.BigEndian, err, res, false)

	res = false
	err = binary.Read(bytes.NewReader([]byte{1}), binary.BigEndian, &res)
	checkResult(t, "ReadBool", binary.BigEndian, err, res, true)
	res = false
	err = binary.Read(bytes.NewReader([]byte{2}), binary.BigEndian, &res)
	checkResult(t, "ReadBool", binary.BigEndian, err, res, true)
}

func TestReadBoolSlice(t *testing.T) {
	slice := make([]uint8, 4)
	err := binary.Read(bytes.NewReader([]byte{0, 1, 2, 255}), binary.BigEndian, slice)
	checkResult(t, "ReadBoolSlice", binary.BigEndian, err, slice, []uint8{0, 1, 2, 255})
}

func checkResult(t *testing.T, dir string, order binary.ByteOrder, err error, have, want interface{}) {
	if err != nil {
		t.Errorf("%v %v: %v", dir, order, err)
		return
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("%v %v:\n\thave %+v\n\twant %+v", dir, order, have, want)
	}
}
