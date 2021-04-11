package tencent_yun_client

type ExportDomainRecord struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	CodeDesc string `json:"codeDesc"`
	Data     struct {
		Domain struct {
			ID         string   `json:"id"`
			Name       string   `json:"name"`
			Punycode   string   `json:"punycode"`
			Grade      string   `json:"grade"`
			Owner      string   `json:"owner"`
			ExtStatus  string   `json:"ext_status"`
			TTL        int      `json:"ttl"`
			MinTTL     int      `json:"min_ttl"`
			DnspodNs   []string `json:"dnspod_ns"`
			Status     string   `json:"status"`
			QProjectID int      `json:"q_project_id"`
		} `json:"domain"`
		Info struct {
			SubDomains  string `json:"sub_domains"`
			RecordTotal string `json:"record_total"`
		} `json:"info"`
		Records []struct {
			ID         int    `json:"id"`
			TTL        int    `json:"ttl"`
			Value      string `json:"value"`
			Enabled    int    `json:"enabled"`
			Status     string `json:"status"`
			UpdatedOn  string `json:"updated_on"`
			QProjectID int    `json:"q_project_id"`
			Name       string `json:"name"`
			Line       string `json:"line"`
			LineID     string `json:"line_id"`
			Type       string `json:"type"`
			Remark     string `json:"remark"`
			Mx         int    `json:"mx"`
			Hold       string `json:"hold,omitempty"`
		} `json:"records"`
	} `json:"data"`
}
