package producer

import "log"

type ProducerService struct {
	producer ProducerInterface
}

func NewProducerService(producer ProducerInterface) *ProducerService {
	return &ProducerService{producer: producer}
}

func (ps *ProducerService) Initialize() error {
	if ps.producer == nil {
		log.Println("Producer instance is nil")
		return nil
	}
	return ps.producer.Initialize()
}

func (ps *ProducerService) Publish(message []byte, taskName string) error {
	if ps.producer == nil {
		log.Println("Producer instance is nil")
		return nil
	}
	return ps.producer.Publish(message, taskName)
}
