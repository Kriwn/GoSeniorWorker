package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Kriwn/GoSeniorWorker/internal/infrastructures/ai"
	"github.com/Kriwn/GoSeniorWorker/internal/infrastructures/models"
	"github.com/Kriwn/GoSeniorWorker/internal/infrastructures/rabbitmq"
	"github.com/Kriwn/GoSeniorWorker/internal/repository"
	amqp "github.com/rabbitmq/amqp091-go"
)

type LstmPredictionService struct {
	repo      *repository.PredictionRepository
	mq        *rabbitmq.RabbitmqService
	predictor *ai.LstmPredictService
}

func NewLstmPredictionService(repo *repository.PredictionRepository, mq *rabbitmq.RabbitmqService, ai *ai.LstmPredictService) *LstmPredictionService {
	return &LstmPredictionService{
		repo:      repo,
		mq:        mq,
		predictor: ai,
	}
}

func (s *LstmPredictionService) Consume(ctx context.Context) error {
	if s == nil || s.mq == nil {
		return fmt.Errorf("lstm prediction service is not initialized")
	}

	return s.mq.Consume(ctx, s.handleMessage)
}

func (s *LstmPredictionService) handleMessage(ctx context.Context, msg amqp.Delivery) error {

	fmt.Println("DO consume Job")

	fmt.Println(string(msg.Body));

	var req models.CreatePredictionRequest
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return fmt.Errorf("invalid request payload: %w", err)
	}

	fmt.Println("Creat Prediction Request:", req)

	PredictionOutput, err := s.predictor.PredictLstm(ctx, ai.LstmPredictionInput{
		ChildId: req.ChildId,
		Weight:  req.Weight,
		Height:  req.Height,
		Sex:     float32(req.Sex),
	})
	if err != nil {
		return err
	}
	fmt.Println(PredictionOutput);
	monthUsed := len(req.Weight)
	id, err := s.repo.CreateAiPrediction(req.ChildId, "lstm", PredictionOutput, monthUsed)
	if err != nil {
		return err
	}

	res := models.CreatePredictionResponse{
		Id: id,
	}

	body, err := json.Marshal(res)
	if err != nil {
		return err
	}

	if err := s.mq.CreatePublish(ctx, body); err != nil {
		return err
	}

	return nil
}
