package pubsub

type EventType string

const (
	EventCreate EventType = "create"

	EventUpdate EventType = "update"

	EventDelete EventType = "delete"

	EventError EventType = "error"
)

type DataSource string

const (
	// DataSource 命名规范：a.b.c
	DataSourceAppTest DataSource = "app.test"
)

type EncodingType string

const (
	EncodingTypeJSON    EncodingType = "json"
	EncodingTypeMSGPACK EncodingType = "msgpackage"
)
