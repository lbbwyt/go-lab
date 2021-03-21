package event

//定义一个 事件channel 类型（DataEvent）
type DataChannel chan DataEvent

type DataChannelSlice []DataChannel
