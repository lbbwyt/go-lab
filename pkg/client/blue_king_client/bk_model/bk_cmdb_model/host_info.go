package bk_cmdb_model

type HostInfo struct {
	BkHostInnerIp string `json:"bk_host_innerip"`
	BkCloudID     int    `json:"bk_cloud_id"`
	ImportFrom    string `json:"import_from"`
}
