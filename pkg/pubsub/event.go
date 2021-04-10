package pubsub

type EventBaseInfo struct {
	DataSource DataSource `json:"data_source"`
	EventType  EventType  `json:"event_type"`
	TimeStamp  int64      `json:"time_stamp"`
}
