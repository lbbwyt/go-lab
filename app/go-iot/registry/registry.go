package registry

type SessionMeta struct {
	DeviceId string
	Network  string // mqtt / tcp / udp
	ClientID string // mqtt 需要
	Remote   string // TCP/UDP 远端地址
	Extra    string `json:"extra,omitempty"` // 任意扩展字段
}

// Session 内部使用，不可序列化
type Session struct {
	Meta   SessionMeta // 嵌套元数据
	Client interface{} // 真实句柄
	// 具体类型断言：
	//   *mqtt.Client  (network == "mqtt")
	//   *net.Conn     (network == "tcp")
	//   *net.UDPConn  (network == "udp")
}

// Registry 注册中心仅用于保存设备的会话信息，方便记录上线下时间, 需要持久化(内存， Redis等方案)，  真正的网络会话句柄需要另外单独在内存中存储， 无法持久化
type Registry interface {
	Register(deviceID string, s *Session) error
	Unregister(deviceID string) error
	Lookup(deviceID string) (*Session, bool)
	Heartbeat(deviceID string) error
}
