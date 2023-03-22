package services

type producerService struct {
}

func NewProducerService() *producerService {
	return &producerService{}
}

func (s *producerService) readFile(filePath string) error {
	return nil
}
