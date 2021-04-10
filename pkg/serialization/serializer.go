package serialization

const (
	Default = "default"
	Json    = "json"
	MsgP    = "msgpackage"
)

type Serializer interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
	Name() string
}

var (
	defaultSerializer Serializer
	jsoner            Serializer
	msgper            Serializer
)

func init() {
	defaultSerializer = new(JsonSerializer)
	jsoner = new(JsonSerializer)
	msgper = new(MsgpSerializer)
}

func SetSerializer(t string) {
	switch t {
	case Json:
		defaultSerializer = new(JsonSerializer)
	case MsgP:
		defaultSerializer = new(MsgpSerializer)
	default:
		defaultSerializer = new(JsonSerializer)
	}
}

func JSON() Serializer {
	return jsoner
}

func MSGP() Serializer {
	return msgper
}

func Marshal(v interface{}) ([]byte, error) {
	return defaultSerializer.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return defaultSerializer.Unmarshal(data, v)
}
