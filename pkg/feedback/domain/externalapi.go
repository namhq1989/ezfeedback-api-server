package domain

import "github.com/namhq1989/go-utilities/appcontext"

type ExternalAPIRepository interface {
	GetIpLocationData(ctx *appcontext.AppContext, ip string) (*IpLocationData, error)
}

type IpLocationData struct {
	Country string
}
