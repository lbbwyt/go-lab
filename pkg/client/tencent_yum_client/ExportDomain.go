package tencent_yun_client

//type ExportDomain struct {
//	Code     int             `json:"code"`
//	Message  string          `json:"message"`
//	CodeDesc string          `json:"codeDesc"`
//	Data     *DomainInfoData `json:"data"`
//}

type DomainInfoData struct {
	Info    DomainInfo `json:"info"`
	Domains []*Domain  `json:"domains"`
}

type DomainInfo struct {
	DomainTotal int `json:"domain_total"`
}
type Domain struct {
	ID               int    `json:"id"`
	Status           string `json:"status"`
	GroupID          string `json:"group_id"`
	SearchEnginePush string `json:"searchengine_push"`
	IsMark           string `json:"is_mark"`
	TTL              string `json:"ttl"`
	CnameSpeedup     string `json:"cname_speedup"`
	Remark           string `json:"remark"`
	CreatedOn        string `json:"created_on"`
	UpdatedOn        string `json:"updated_on"`
	QProjectID       int    `json:"q_project_id"`
	Punycode         string `json:"punycode"`
	ExtStatus        string `json:"ext_status"`
	SrcFlag          string `json:"src_flag"`
	Name             string `json:"name"`
	Grade            string `json:"grade"`
	GradeTitle       string `json:"grade_title"`
	IsVip            string `json:"is_vip"`
	Owner            string `json:"owner"`
	Records          string `json:"records"`
	MinTTL           int    `json:"min_ttl"`
}

type ExportDomain struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	CodeDesc string `json:"codeDesc"`
	Data     struct {
		Info struct {
			DomainTotal int `json:"domain_total"`
		} `json:"info"`
		Domains []struct {
			ID               int    `json:"id"`
			Status           string `json:"status"`
			GroupID          string `json:"group_id"`
			SearchenginePush string `json:"searchengine_push"`
			IsMark           string `json:"is_mark"`
			TTL              string `json:"ttl"`
			CnameSpeedup     string `json:"cname_speedup"`
			Remark           string `json:"remark"`
			CreatedOn        string `json:"created_on"`
			UpdatedOn        string `json:"updated_on"`
			QProjectID       int    `json:"q_project_id"`
			Punycode         string `json:"punycode"`
			ExtStatus        string `json:"ext_status"`
			SrcFlag          string `json:"src_flag"`
			Name             string `json:"name"`
			Grade            string `json:"grade"`
			GradeTitle       string `json:"grade_title"`
			IsVip            string `json:"is_vip"`
			Owner            string `json:"owner"`
			Records          string `json:"records"`
			MinTTL           int    `json:"min_ttl"`
		} `json:"domains"`
	} `json:"data"`
}
