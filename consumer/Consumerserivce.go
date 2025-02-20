package consumer

import (
	"log"
)

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

// func (cs *ConsumerService) Consume(Data []byte, taskname string) error {
// 	if cs.consumer == nil {
// 		log.Println("Consumer instance is nil")
// 		return nil
// 	}
// 	return cs.consumer.Consume(Data, taskname)
// }

// func (cs *ConsumerService) Print(arg *schema.Signature) error {
// 	if cs.consumer == nil {
// 		log.Println("Consumer instance is nil")
// 		return nil
// 	}
// 	return cs.consumer.Print(arg)
// }
