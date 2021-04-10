package serialization

import "encoding/json"

type JsonSerializer struct{}

func (s *JsonSerializer) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (s *JsonSerializer) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (s *JsonSerializer) Name() string {
	return "json"
}
