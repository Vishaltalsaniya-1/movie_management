package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	cnf "movie_management/config"
	"movie_management/models"

	"github.com/beego/beego/v2/client/orm"
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
	consumerConfig := cnf.Consumerconfig

	cnf := config.Config{
		Broker:        "amqp",
		TaskQueueName: consumerConfig.QueueTaskName,
		AMQP: &config.AMQPConfig{
			Url:                consumerConfig.Url,
			Exchange:           consumerConfig.Exchange,
			ExchangeType:       "direct",
			BindingKey:         consumerConfig.BindingKeyName,
			PrefetchCount:      consumerConfig.PrefetchCount,
			ConnectionPoolSize: consumerConfig.ConnectionPoolSize,
			DelayedQueue:       consumerConfig.DelayedQueueName,
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
		consumerConfig.QueueTaskName: c.ProcessTask,
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

func (c *Consumer) ProcessTask(arg *schema.Signature) error {

	log.Println("Received a new task to process...")

	if len(arg.Args) == 0 {
		logrus.Warn("No arguments received in the task")
		return nil
	}

	for _, argItem := range arg.Args {
		logrus.Infof("Processing task - Arg Type: %s, Arg Value: %v", argItem.Type, argItem.Value)

		argStr, ok := argItem.Value.(string)
		if !ok {
			logrus.Error("Received argument is not a string")
			return errors.New("invalid argument type")
		}

		// var user models.Movie
		// logrus.Infof("Raw task argument: %s", argStr)
		// if err := json.Unmarshal([]byte(argStr), &user); err != nil {
		// 	logrus.Errorf("Error decoding user data: %v", err)
		// 	return err
		// }
		// logrus.Infof("Decoded User: %+v", user)
		var movie models.Movie
		if err := json.Unmarshal([]byte(argStr), &movie); err != nil {
			logrus.Errorf("Error decoding movie data: %v", err)
			return err
		}

		o := orm.NewOrm()
		if o == nil {
			logrus.Error("Database connection is nil")
			return errors.New("database connection not established")
		}

		if _, err := o.Insert(&movie); err != nil {
			logrus.Errorf("Database insertion failed: %v", err)
			return err
		}

		logrus.Info("User data successfully inserted into the database.")
	}

	return nil
}

// func (c *Consumer) Consume(data []byte, taskName string) error {
// 	if c.WorkerPool == nil {
// 		logrus.Error("Worker pool is not initialized")
// 		return errors.New("worker pool is not initialized")
// 	}

// 	task := &schema.Signature{
// 		Name: taskName,
// 		Args: []schema.Arg{
// 			{
// 				Name:  taskName,
// 				Type:  "string",
// 				Value: string(data),
// 			},
// 		},
// 		RetryCount:                  10,
// 		RoutingKey:                  cnf.Consumerconfig.BindingKeyName,
// 		IgnoreWhenTaskNotRegistered: true,
// 	}

// 	state, err := (*c.WorkerPool).SendTaskWithContext(context.Background(), task)
// 	if err != nil {
// 		logrus.Errorf("Error consuming task: %v", err)
// 		return err
// 	}

// 	logrus.Infof("Task sent successfully. State: %+v", state)
// 	return nil
// }

// func (c *Consumer) Print(arg *schema.Signature) error {
// 	if len(arg.Args) == 0 {
// 		logrus.Warn("No arguments received in the task")
// 		return nil
// 	}

// 	for _, argItem := range arg.Args {
// 		logrus.Infof("Processing task - Arg Type: %s, Arg Value: %v", argItem.Type, argItem.Value)

// 		argStr, ok := argItem.Value.(string)
// 		if !ok {
// 			logrus.Error("Received argument is not a string")
// 			return errors.New("invalid argument type")
// 		}

// 		var movie models.Movie
// 		if err := json.Unmarshal([]byte(argStr), &movie); err != nil {
// 			logrus.Errorf("Error decoding movie data: %v", err)
// 			return err
// 		}

// 		o := orm.NewOrm()
// 		if o == nil {
// 			logrus.Error("Database connection is nil")
// 			return errors.New("database connection not established")
// 		}

// 		if _, err := o.Insert(&movie); err != nil {
// 			logrus.Errorf("Database insertion failed: %v", err)
// 			return err
// 		}

// 		logrus.Info("Movie successfully inserted into the database.")
// 	}

// 	return nil
// }
