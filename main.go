package main

import (
	"log"
	cnf "movie_management/config"
	"movie_management/consumer"
	"movie_management/controller"
	"movie_management/db"
	"movie_management/middlewares"
	"movie_management/producer"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartConsumerService() error {
	consumerInstance := consumer.NewConsumer()
	if err := consumerInstance.Initialize(); err != nil {
		log.Fatalf("Failed to start consumer service: %v", err)
		return err
	}
	log.Println("Consumer service started successfully")
	return nil
}

func StartProducerService() error {
	producerInstance := producer.NewProducer()
	if err := producerInstance.Initialize(); err != nil {
		log.Fatalf("Failed to start producer service: %v", err)
		return nil
	}
	log.Println("Producer service started successfully")
	return nil
}

func main() {
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	cnf.Loadcosumer()
	if cnf.Consumerconfig.RunConsumer {
		go StartConsumerService()
	}
	if cnf.Consumerconfig.RunProducer {
		go StartProducerService()
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/movies", controller.CreateMovie, middlewares.AuthMiddleware)
	e.PUT("/movies/:id", controller.UpdateMovie)
	e.DELETE("/movies/:id", controller.DeleteMovie)
	e.GET("/movies", controller.ListMovies)
	e.GET("/movies/analytics", controller.GetMovieAnalytics)
	e.GET("/movies/:id", controller.GetMoviesById)

	e.POST("/register", controller.Register)
	e.POST("/login", controller.Login)

	// api := e.Group("/api")
	// api.Use(middleware.AuthMiddleware)

	// api.GET("/profile", func(c echo.Context) error {
	// 	return c.JSON(200, map[string]string{"message": "This is a protected profile route"})
	// })

	// e.Logger.Fatal(e.Start(":8081"))
	go func() {
		log.Println("Server started on port 8081")
		if err := e.Start(":8081"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan
	log.Println(" Shutting down application gracefully...")

}

// func Startconsumerserivce() error {
// 	// consumerConfig := cnf.Consumerconfig

// 	AMQP := cnf.ConsumerConfig{
// 		Url:                cnf.Consumerconfig.Url,
// 		Exchange:           cnf.Consumerconfig.Exchange,
// 		ExchangeType:       "direct",
// 		BindingKeyName:     cnf.Consumerconfig.BindingKeyName,
// 		PrefetchCount:      cnf.Consumerconfig.PrefetchCount,
// 		ConnectionPoolSize: cnf.Consumerconfig.ConnectionPoolSize,
// 		DelayedQueueName:   cnf.Consumerconfig.DelayedQueueName,
// 	}

// 	_, err := workerpool.NewWorkerPoolWithConfig(context.Background(), 10, "testmovie", AMQP)
// 	if err != nil {
// 		logrus.Fatalf("WorkerPool creation failed: %v", err)
// 		return err
// 	}

// 	consumerInstance := consumer.NewConsumer()
// 	consumerService := consumer.NewConsumerService(consumerInstance)

// 	if err := consumerService.Initialize(); err != nil {
// 		log.Fatalf("Failed to initialize consumer: %v", err)
// 		return err
// 	}

// 	log.Println("Consumer and worker pool initialized successfully!")
// 	return nil
// }

// func (c *Consumer) Initialize() error {
// 	consumerConfig := cnf.Consumerconfig

// 	paotaConfig := config.Config{
// 		Broker:        "amqp",
// 		TaskQueueName: consumerConfig.QueueTaskName,
// 		AMQP: &config.AMQPConfig{
// 			Url:                consumerConfig.Url,
// 			Exchange:           consumerConfig.Exchange,
// 			ExchangeType:       "direct",
// 			BindingKey:         consumerConfig.BindingKeyName,
// 			PrefetchCount:      consumerConfig.PrefetchCount,
// 			ConnectionPoolSize: consumerConfig.ConnectionPoolSize,
// 			DelayedQueue:       consumerConfig.DelayedQueueName,
// 		},
// 	}

// 	workerPool, err := workerpool.NewWorkerPoolWithConfig(context.Background(), 10, "testmovie", paotaConfig)
// 	if err != nil {
// 		logrus.Errorf("WorkerPool creation failed: %v", err)
// 		return err
// 	}

// 	c.WorkerPool = &workerPool
// 	if c.WorkerPool == nil {
// 		logrus.Error("WorkerPool is nil after initialization")
// 		return errors.New("failed to initialize worker pool")
// 	}

// 	regTasks := map[string]interface{}{
// 		consumerConfig.QueueTaskName: c.Print,
// 	}
// 	if err := workerPool.RegisterTasks(regTasks); err != nil {
// 		logrus.Errorf("Error registering tasks: %v", err)
// 		return err
// 	}

// 	if err := workerPool.Start(); err != nil {
// 		logrus.Errorf("Error starting worker: %v", err)
// 		return err
// 	}

// 	logrus.Info("WorkerPool consumer initialized and started successfully")
// 	return nil
// }
