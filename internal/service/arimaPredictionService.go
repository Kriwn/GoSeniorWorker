package service

import (
	"context"
	"encoding/json"
	"fmt"
	// "os"

	"github.com/Kriwn/GoSeniorWorker/internal/infrastructures/ai"
	"github.com/Kriwn/GoSeniorWorker/internal/infrastructures/models"
	"github.com/Kriwn/GoSeniorWorker/internal/infrastructures/rabbitmq"
	"github.com/Kriwn/GoSeniorWorker/internal/repository"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ArimaPredictionService struct {
	repo      *repository.PredictionRepository
	mq        *rabbitmq.RabbitmqService
	predictor *ai.ArimaPredictService
}

func NewArimaPredictionService(repo *repository.PredictionRepository, mq *rabbitmq.RabbitmqService, ai *ai.ArimaPredictService) *ArimaPredictionService {
	return &ArimaPredictionService{
		repo:      repo,
		mq:        mq,
		predictor: ai,
	}
}

func (s *ArimaPredictionService) Consume(ctx context.Context) error {
	if s == nil || s.mq == nil {
		return fmt.Errorf("arima prediction service is not initialized")
	}

	return s.mq.Consume(ctx, s.handleMessage)
}

func (s *ArimaPredictionService) handleMessage(ctx context.Context, msg amqp.Delivery) error {
	if s == nil || s.repo == nil {
		return fmt.Errorf("prediction repository is not initialized")
	}

	var req models.CreatePredictionRequest
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return fmt.Errorf("invalid request payload: %w", err)
	}

	if req.ChildId <= 0 {
		return fmt.Errorf("childId must be greater than 0")
	}

	if len(req.Weight) == 0 || len(req.Height) == 0 {
		return fmt.Errorf("weight and height cannot be empty")
	}

	PredictionOutput, err := s.predictor.PredictArima(ctx, ai.ArimaPredictionInput{
		ChildId: req.ChildId,
		Weight:  req.Weight,
		Height:  req.Height,
		Sex:     req.Sex,
	})

	if err != nil {
		return fmt.Errorf("predict failed: %w", err)
	}

	fmt.Printf("Prediction result: %+v\n", PredictionOutput)

	// id, err := s.repo.CreateAiPrediction(req.ChildId,model)
	// if err != nil {
	// 	return fmt.Errorf("create ai prediction failed: %w", err)
	// }

	// res := models.CreatePredictionResponse{
	// 	ID:      id,
	// 	ChildId: req.ChildId,
	// 	Model:   model,
	// }

	// body, err := json.Marshal(res)
	// if err != nil {
	// 	return fmt.Errorf("marshal response failed: %w", err)
	// }

	// if err := s.mq.CreatePublish(ctx, body); err != nil {
	// 	return fmt.Errorf("publish response failed: %w", err)
	// }

	return nil
}
