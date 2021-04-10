package pubsub

import (
	"fmt"
	"go-lab/pkg/serialization"
)

type Publication struct {
	Data         []byte       `json:"data,omitempty"`
	Topic        string       `json:"topic,omitempty"`
	EncodingType EncodingType `json:"encodingType,omitempty"` //序列化\反序列化类型
}

func NewPublication(topic string, encodingType EncodingType, obj interface{}) (*Publication, error) {
	publication := &Publication{
		Topic:        topic,
		EncodingType: encodingType,
	}

	if err := publication.Encode(obj); err != nil {
		return nil, err
	}
	return publication, nil
}

func (p *Publication) Encode(obj interface{}) error {
	if obj == nil {
		return fmt.Errorf("encode received a nil object")
	}

	var serializer serialization.Serializer

	switch p.EncodingType {
	case EncodingTypeMSGPACK:
		serializer = serialization.MSGP()
	default:
		serializer = serialization.JSON()
	}

	data, err := serializer.Marshal(obj)
	if err != nil {
		return err
	}

	p.Data = data

	return nil
}

func (p *Publication) Decode(dest interface{}) error {
	var serializer serialization.Serializer

	switch p.EncodingType {
	case EncodingTypeMSGPACK:
		serializer = serialization.MSGP()
	default:
		serializer = serialization.JSON()
	}

	return serializer.Unmarshal(p.Data, dest)
}
