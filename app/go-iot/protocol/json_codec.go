package protocol

import "encoding/json"

type JSONCodec struct{}

func (c JSONCodec) Name() string { return "json" }

func (c JSONCodec) Decode(b []byte) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	return m, err
}

func (c JSONCodec) Encode(m map[string]interface{}) ([]byte, error) {
	return json.Marshal(m)
}
