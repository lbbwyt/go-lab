package tencent_yun_client

// 证书列表接口返回， 证书公私钥等其他信息， 都在证书的详情接口中返回
type ExportCertificate struct {
	Response struct {
		Certificates []struct {
			From             string   `json:"From"`           //证书来源。注意：此字段可能返回 null，表示取不到有效值。
			SubjectAltName   []string `json:"SubjectAltName"` //证书包含的多个域名（包含主域名）
			BoundResource    []string `json:"BoundResource"`  //关联的云资源，暂不可用
			CertificateExtra struct { //证书扩展信息
				OriginCertificateID string `json:"OriginCertificateId"` //原始证书 ID。
				ReplacedBy          string `json:"ReplacedBy"`          //重颁发证书原始 ID。
				ReplacedFor         string `json:"ReplacedFor"`         //重颁发证书新 ID。
				DomainNumber        string `json:"DomainNumber"`        //证书可配置域名数量。
				RenewOrder          string `json:"RenewOrder"`          //新订单证书 ID。
			} `json:"CertificateExtra"`
			StatusName      string   `json:"StatusName"`      //状态名称。
			RenewAble       bool     `json:"RenewAble"`       //是否可重颁发证书。
			Status          int      `json:"Status"`          //状态值 0：审核中，1：已通过，2：审核失败，3：已过期，4：已添加 DNS 解析记录，5：OV/EV 证书，待提交资料，6：订单取消中，7：已取消，8：已提交资料， 待上传确认函。
			IsDv            bool     `json:"IsDv"`            //是否为 DV 版证书。
			CertBeginTime   string   `json:"CertBeginTime"`   //证书生效时间。
			IsVulnerability bool     `json:"IsVulnerability"` //是否启用了漏洞扫描功能。
			VerifyType      string   `json:"VerifyType"`      //验证类型：DNS_AUTO = 自动DNS验证，DNS = 手动DNS验证，FILE = 文件验证，EMAIL = 邮件验证。
			StatusMsg       string   `json:"StatusMsg"`       //状态信息。
			ProjectID       string   `json:"ProjectId"`       //项目 ID。 注意：此字段可能返回 null，表示取不到有效值。
			OwnerUin        string   `json:"OwnerUin"`        //用户 UIN。 注意：此字段可能返回 null，表示取不到有效值
			ProjectInfo     struct { //项目信息。
				ProjectCreatorUin int    `json:"ProjectCreatorUin"` //项目创建用户 UIN。
				ProjectCreateTime string `json:"ProjectCreateTime"` //项目创建时间。
				ProjectID         string `json:"ProjectId"`         //项目 ID。
				OwnerUin          int    `json:"OwnerUin"`          //用户 UIN。
				ProjectResume     string `json:"ProjectResume"`     //项目信息简述。
				ProjectName       string `json:"ProjectName"`       //项目名称。
			} `json:"ProjectInfo"`
			ProductZhName       string     `json:"ProductZhName"`   //颁发者。
			CertEndTime         string     `json:"CertEndTime"`     //证书过期时间。
			PackageType         string     `json:"PackageType"`     //证书类型名称。
			InsertTime          string     `json:"InsertTime"`      //创建时间。
			CertificateType     string     `json:"CertificateType"` //证书类型：CA = 客户端证书，SVR = 服务器证书。
			IsVip               bool       `json:"IsVip"`           //是否为 VIP 客户。
			ValidityPeriod      string     `json:"ValidityPeriod"`  //证书有效期，单位（月）。
			Domain              string     `json:"Domain"`          //主域名。
			CertificateID       string     `json:"CertificateId"`   //证书 ID。
			Alias               string     `json:"Alias"`
			IsWildcard          bool       `json:"IsWildcard"`          //是否为泛域名证书。
			PackageTypeName     string     `json:"PackageTypeName"`     //证书套餐类型
			VulnerabilityStatus string     `json:"VulnerabilityStatus"` //漏洞扫描状态：INACTIVE = 未开启，ACTIVE = 已开启
			Deployable          bool       `json:"Deployable"`          //是否可部署。
			Tags                []struct { //标签列表
				TagKey   string `json:"TagKey"`
				TagValue string `json:"TagValue"`
			} `json:"Tags"`
		} `json:"Certificates"`
		TotalCount int    `json:"TotalCount"`
		RequestID  string `json:"RequestId"`
	} `json:"Response"`
}
