

# 物模型定义

| 维度              | 说明                | 读写       | 举例         |
| --------------- | ----------------- | -------- | ---------- |
| **属性 Property** | 设备“状态量”：可读、可订阅 | RO  | 温度、开关、亮度   |
| **功能 Service**  | 设备“能力”：一条命令触发复杂逻辑 | IN / OUT | 固件升级、重启、扫码 |
| **事件 Event**    | 设备“通知”：主动上报、瞬时    | 只上报      | 高温告警、故障码   |


```
{
  "schema": "schema.json",
  "version": "1.0",
  "properties": [
    {
      "id": "temperature",
      "name": "当前温度",
      "accessMode": "r",
      "dataType": { "type": "float", "specs": { "min": -50, "max": 100, "unit": "℃" } }
    }
  ],
  "services": [
    {
      "id": "reboot",
      "name": "重启设备",
      "callType": "async",
      "input": [],
      "output": []
    }
  ],
  "events": [
    {
      "id": "overTemp",
      "name": "高温告警",
      "type": "alert",
      "output": [
        { "id": "temp", "dataType": { "type": "float" } }
      ]
    }
  ]
}
```