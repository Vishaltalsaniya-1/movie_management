package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"movie_management/models"

	"github.com/sirupsen/logrus"
	"github.com/surendratiwari3/paota/config"
	"github.com/surendratiwari3/paota/schema"
	"github.com/surendratiwari3/paota/workerpool"
)

type Consumer struct {
	WorkerPool *workerpool.Pool
}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Initialize() error {
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

	workerPool, err := workerpool.NewWorkerPoolWithConfig(context.Background(), 10, "testmovie", cnf)
	if err != nil {
		logrus.Errorf("WorkerPool creation failed: %v", err)
		return err
	}

	c.WorkerPool = &workerPool

	if c.WorkerPool == nil {
		logrus.Error("WorkerPool is nil after initialization")
		return errors.New("failed to initialize worker pool")
	}

	regTasks := map[string]interface{}{
		"MovieCreatedTask": c.Print,
	}
	if err := workerPool.RegisterTasks(regTasks); err != nil {
		logrus.Errorf("Error registering tasks: %v", err)
		return err
	}

	if err := workerPool.Start(); err != nil {
		logrus.Errorf("Error starting worker: %v", err)
		return err
	}

	logrus.Info("WorkerPool consumer initialized and started successfully")
	return nil

}

func (c *Consumer) Consume(Data []byte, taskname string) error {
	if c.WorkerPool == nil {
		return errors.New("worker pool is not initialized")
	}

	logrus.Info("Starting to consume tasks...")

	task := &schema.Signature{
		Name: "MovieCreatedTask",
		Args: []schema.Arg{
			{
				Type:  "string",
				Value: string(Data),
			},
		},
		RetryCount:                  10,
		RoutingKey:                  "movie_add_binding_key",
		IgnoreWhenTaskNotRegistered: true,
	}

	state, err := (*c.WorkerPool).SendTaskWithContext(context.Background(), task)
	if err != nil {
		logrus.Errorf("Error consuming task: %v", err)
		return err
	}

	logrus.Infof("Task sent successfully. State: %+v", state)
	return nil
}

func (c *Consumer) Print(arg *schema.Signature) error {
	if len(arg.Args) == 0 {
		logrus.Info("No arguments found in the task")
		return nil
	}
	for _, argItem := range arg.Args {
		logrus.Infof("Task received - Arg Type: %s, Arg Value: %v", argItem.Type, argItem.Value)

		argStr, ok := argItem.Value.(string)
		if !ok {
			logrus.Errorf("Unexpected argument type. Expected string, got: %T", argItem.Value)
			return errors.New("invalid argument type")
		}

		var movie models.Movie
		if err := json.Unmarshal([]byte(argStr), &movie); err != nil {
			logrus.Errorf("Error unmarshalling movie data: %v", err)
			return err
		}
		logrus.Infof("Received movie: %+v", movie)
		    
	}

	return nil
}


	// for _, argItem := range arg.Args {
	// 	logrus.Infof("Task received - Arg Type: %s, Arg Value: %v", argItem.Type, argItem.Value)

	// 	var taskData map[string]interface{}
	// 	if err := json.Unmarshal([]byte(argItem.Value.(string)), &taskData); err != nil {
	// 		logrus.Errorf("Error unmarshalling task data: %v", err)
	// 		return err
	// 	}

	// 	logrus.Infof("Processed task with data: %+v", taskData)
	// }
	

// Insert into MySQL using GORM
// if err := db.DB.Create(&movie).Error; err != nil {
// 	logrus.Errorf("Error inserting movie into database: %v", err)
// 	return err
// }