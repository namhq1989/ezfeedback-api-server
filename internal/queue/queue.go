package queue

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

type Operations interface {
	GetServer() *asynq.ServeMux
	GetScheduler() *asynq.Scheduler

	GenerateTypename(name string) string
	RunTask(queue string, payload interface{}, retryTimes int) (*asynq.TaskInfo, error)
	ScheduleTask(typename string, payload interface{}, cronSpec string, retryTimes int) (string, error)
	RemoveScheduler(id string) error
}

type Queue struct {
	Client    *asynq.Client
	Server    *asynq.ServeMux
	Scheduler *asynq.Scheduler
}

func Init(redisURL string, concurrency int) *Queue {
	// redis connection
	redisConn := getRedisConnFromURL(redisURL)

	// init instance
	q := &Queue{}
	q.Server = initServer(redisConn, concurrency)
	q.Scheduler = initScheduler(redisConn)
	q.Client = initClient(redisConn)

	fmt.Printf("⚡️ [queue]: initialized \n")

	return q
}

func initServer(redisConn asynq.RedisClientOpt, concurrency int) *asynq.ServeMux {
	// set unrecognized for concurrency
	if concurrency == 0 {
		concurrency = 100
	}

	// unrecognized retry delay in 10s
	retryDelayFunc := func(n int, e error, t *asynq.Task) time.Duration {
		return 10 * time.Second
	}

	// init server
	server := asynq.NewServer(redisConn, asynq.Config{
		Concurrency:    concurrency,
		RetryDelayFunc: retryDelayFunc,
		Queues: map[string]int{
			queueCronjob: 10,
			queueDefault: 3,
		},
	})

	// init mux server
	mux := asynq.NewServeMux()

	// run server
	go func() {
		if err := server.Run(mux); err != nil {
			msg := fmt.Sprintf("error when initializing queue SERVER: %s", err.Error())
			panic(msg)
		}
	}()

	return mux
}

func initScheduler(redisConn asynq.RedisClientOpt) *asynq.Scheduler {
	// init scheduler
	scheduler := asynq.NewScheduler(redisConn, nil)

	// run scheduler
	go func() {
		if err := scheduler.Run(); err != nil {
			msg := fmt.Sprintf("error when initializing queue SCHEDULER: %s", err.Error())
			panic(msg)
		}
	}()

	return scheduler
}

func initClient(redisConn asynq.RedisClientOpt) *asynq.Client {
	client := asynq.NewClient(redisConn)
	if client == nil {
		panic("error when initializing queue CLIENT")
	}
	return client
}

func (q Queue) GetServer() *asynq.ServeMux {
	return q.Server
}

func (q Queue) GetScheduler() *asynq.Scheduler {
	return q.Scheduler
}
