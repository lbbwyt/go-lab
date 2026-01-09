package gateway

//让网关在任何一次解码之前，就能从原始报文或连接上下文中拿到『产品Key + 设备名（或设备ID）』，然后用它作为后续所有环节的『主键』。”
//网关下发时， 需通过设备ID， 去注册中心找到网络会话句柄（MQTT Client、TCP Conn 等）， 由网络会话句柄下发消息给到设备端

// 做法 1：topic 解析（MQTT 最常见）
//设备端约定 Topic
///{productKey}/{deviceName}/up
//网关一收到 Publish 即可拆分：
//
//
//parts := strings.Split(msg.Topic(), "/")
//productKey, deviceName := parts[0], parts[1]
//
//decoded, _ := codec.Decode(msg.Payload())
//repo.SaveData(productKey, deviceName, decoded)

//做法 2：报文头带身份（TCP/UDP/私有协议）
//在二进制帧里固定 2~4 字节存设备 ID：
//| 2 byte DeviceID | N byte 数据 |
//Codec 解码前，网关先把头切出来：
//deviceID := binary.BigEndian.Uint16(frame[0:2])
//payload  := frame[2:]
//decoded, _ := codec.Decode(payload)
//productKey := repo.ProductOf(deviceID)
//repo.SaveData(productKey, strconv.Itoa(int(deviceID)), decoded)
