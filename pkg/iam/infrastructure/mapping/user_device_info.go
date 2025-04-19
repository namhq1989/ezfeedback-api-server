package mapping

type UserDeviceInfo struct {
	OS        string `json:"os"`
	UserAgent string `json:"userAgent"`
	Ip        string `json:"ip"`
}
