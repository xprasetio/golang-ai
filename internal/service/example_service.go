package service

import (
	"context"
	"fmt"
	"golang-ai/internal/dto"
	"golang-ai/internal/repository"
)

type IExampleService interface {
	HelloWorld(ctx context.Context, req *dto.HelloWorldRequest) (*dto.HelloWorldResponse, error)
}

type exampleService struct {
	exampleRepository repository.IExampleRepository
}

func NewExampleService(exampleRepository repository.IExampleRepository) IExampleService {
	return &exampleService{
		exampleRepository: exampleRepository,
	}
}

func (c *exampleService) HelloWorld(ctx context.Context, req *dto.HelloWorldRequest) (*dto.HelloWorldResponse, error) {
	_, err := c.exampleRepository.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.HelloWorldResponse{
		Message: fmt.Sprintf(`Hello %s`, req.Name),
	}, nil
}
