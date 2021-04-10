package serialization

import (
	msgpack "github.com/vmihailenco/msgpack/v4"
)

type MsgpSerializer struct {
}

func (s *MsgpSerializer) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (s *MsgpSerializer) Unmarshal(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}

func (s *MsgpSerializer) Name() string {
	return "messagepack"
}
