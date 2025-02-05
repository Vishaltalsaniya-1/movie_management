package consumer

import "log"

type ConsumerService struct {
	consumer ConsumerInterface
}

func NewConsumerService(consumer ConsumerInterface) *ConsumerService {
	return &ConsumerService{consumer: consumer}
}

func (cs *ConsumerService) Initialize() error {
	if cs.consumer == nil {
		log.Println("Consumer instance is nil")
		return nil
	}
	return cs.consumer.Initialize()
}

func (cs *ConsumerService) Consume() error {
	if cs.consumer == nil {
		log.Println("Consumer instance is nil")
		return nil
	}
	return cs.consumer.Consume()
}
