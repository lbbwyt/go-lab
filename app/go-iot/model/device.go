package model

import "time"

type DevDevice struct {
	Id           uint       `json:"id"             description:""`
	Key          string     `json:"key"            description:"设备标识"`
	Name         string     `json:"name"           description:"设备名称"`
	ProductKey   string     `json:"productKey"     description:"所属产品KEY"`
	Status       int        `json:"status"         description:"状态：0=未启用,1=离线,2=在线"`
	RegistryTime *time.Time `json:"registryTime"   description:"激活时间"`
}
