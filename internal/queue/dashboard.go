package queue

import (
	"github.com/hibiken/asynqmon"
)

const DashboardPath = "/q"

func EnableDashboard(redisURL string) *asynqmon.HTTPHandler {
	redisConn := getRedisConnFromURL(redisURL)

	return asynqmon.New(asynqmon.Options{
		RootPath:     DashboardPath,
		RedisConnOpt: redisConn,
	})
}
