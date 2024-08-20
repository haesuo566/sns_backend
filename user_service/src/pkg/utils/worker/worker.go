package worker

import (
	"os"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
)

type Job struct {
	CorrelationId string
	Key           string
	Value         []byte
	Task          func(string, string, []byte) error // error 뱉어야함
}

type WorkerPool struct {
	maxSize  int
	jobQueue chan Job
}

var initSync sync.Once
var pool *WorkerPool

func Run() {
	initSync.Do(func() {
		poolSize, err := strconv.Atoi(os.Getenv("WORKER_POOL_SIZE"))
		if err != nil {
			panic(err)
		}

		pool = &WorkerPool{
			maxSize:  poolSize,
			jobQueue: make(chan Job, poolSize),
		}

		execute()
	})
}

func AddJob(job Job) {
	// blocking
	for len(pool.jobQueue) >= pool.maxSize {
	}

	logrus.WithFields(logrus.Fields{
		"CorrelationId": job.CorrelationId,
		"Key":           job.Key,
	}).Info("goroutine start")
	pool.jobQueue <- job
}

func execute() {
	for i := 0; i < pool.maxSize; i++ {
		go func(id int) {
			for job := range pool.jobQueue {
				if err := job.Task(job.CorrelationId, job.Key, job.Value); err != nil {
					logrus.WithFields(logrus.Fields{
						"CorrelationId": job.CorrelationId,
						"Key":           job.Key,
					}).Error(err)
				} else {
					logrus.WithFields(logrus.Fields{
						"CorrelationId": job.CorrelationId,
						"Key":           job.Key,
					}).Info("goroutine end")
				}
			}
		}(i + 1)
	}
}
