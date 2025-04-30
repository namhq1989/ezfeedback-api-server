package domain

import "github.com/namhq1989/go-utilities/appcontext"

type CachingRepository interface {
	GetIpLocationData(ctx *appcontext.AppContext, ip string) (*IpLocationData, error)
	SetIpLocationData(ctx *appcontext.AppContext, ip string, data IpLocationData) error
}
