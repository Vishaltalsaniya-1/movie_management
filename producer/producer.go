package producer

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/surendratiwari3/paota/config"
	"github.com/surendratiwari3/paota/schema"
	"github.com/surendratiwari3/paota/workerpool"
)

type RMP struct {
	WorkerPool *workerpool.Pool
}

func NewProducer() ProducerInterface {
	return &RMP{}
}

func (rmp *RMP) Initialize() error {
	cnf := config.Config{
		Broker:        "amqp",
		TaskQueueName: "movie_add",
		AMQP: &config.AMQPConfig{
			Url:                "amqp://guest:guest@localhost:5672/",
			Exchange:           "movie_add_exchange",
			ExchangeType:       "direct",
			BindingKey:         "movie_add_binding_key",
			PrefetchCount:      100,
			ConnectionPoolSize: 10,
			DelayedQueue:       "movie_add_delay_test",
		},
	}

	workerPool, err := workerpool.NewWorkerPoolWithConfig(context.Background(), 10, "testWorker", cnf)
	if err != nil {
		logrus.Fatalf("WorkerPool creation failed: %v", err)
		return err
	}

	rmp.WorkerPool = &workerPool

	if rmp.WorkerPool == nil {
		logrus.Fatal("WorkerPool is nil after initialization")
		return errors.New("failed to initialize worker pool")
	}

	logrus.Info("WorkerPool initialized successfully")
	return nil
}

func (rmp *RMP) Publish(Data []byte, taskname string) error {
	if rmp.WorkerPool == nil {
		return errors.New("worker pool is not initialized")
	}

	task := &schema.Signature{
		Name: taskname, 
		Args: []schema.Arg{
			{
				Type:  "string",
				Value: string(Data),
			},
		},
		RetryCount:                  10,
		IgnoreWhenTaskNotRegistered: true,
	}

	logrus.Infof("Created task: %+v", task)

	state, err := (*rmp.WorkerPool).SendTaskWithContext(context.Background(), task)
	if err != nil {
		logrus.Error("Failed to publish task:", err)
		return err
	}

	logrus.Infof("Task State: %+v", state)
	logrus.Info("Task published successfully.")
	return nil
}






























// package producer


// import (
// 	"context"
// 	"errors"

// 	"github.com/sirupsen/logrus"
// 	"github.com/surendratiwari3/paota/config"
// 	"github.com/surendratiwari3/paota/schema"
// 	"github.com/surendratiwari3/paota/workerpool"
// )

// type RMP struct {
// 	WorkerPool *workerpool.Pool
// }

// func (rmp *RMP) Initialize() error {
// 	cnf := config.Config{
// 		Broker:        "amqp",
// 		TaskQueueName: "movie_add",
// 		AMQP: &config.AMQPConfig{
// 			Url:                "amqp://guest:guest@localhost:5672/",
// 			Exchange:           "movie_add_exchange",
// 			ExchangeType:       "direct",
// 			BindingKey:         "movie_add_binding_key",
// 			PrefetchCount:      100,
// 			ConnectionPoolSize: 10,
// 			DelayedQueue:       "movie_add_delay_test",
// 		},
// 	}

// 	workerPool, err := workerpool.NewWorkerPoolWithConfig(context.Background(), 10, "testWorker", cnf)
// 	if err != nil {
// 		logrus.Fatalf("WorkerPool creation failed: %v", err)
// 		return err
// 	}

// 	rmp.WorkerPool = &workerPool

// 	if rmp.WorkerPool == nil {
// 		logrus.Fatal("WorkerPool is nil after initialization")
// 		return errors.New("failed to initialize worker pool")
// 	}

// 	logrus.Info("WorkerPool initialized successfully")
// 	return nil
// }

// func (rmp *RMP) Publish(Data []byte) error {
// 	if rmp.WorkerPool == nil {
// 		return errors.New("worker pool is not initialized")
// 	}

// 	task := &schema.Signature{
// 		Name: "MovieCreatedTask",
// 		Args: []schema.Arg{
// 			{
// 				Type:  "string",
// 				Value: string(Data),
// 			},
// 		},
// 		RetryCount:                  10,
// 		IgnoreWhenTaskNotRegistered: true,
// 	}

// 	logrus.Infof("Created task: %+v", task)

// 	state, err := (*rmp.WorkerPool).SendTaskWithContext(context.Background(), task)
// 	if err != nil {
// 		logrus.Error("Failed to publish task:", err)
// 		return err
// 	}

// 	logrus.Infof("Task State: %+v", state)
// 	logrus.Info("Task published successfully.")
// 	return nil
// }
