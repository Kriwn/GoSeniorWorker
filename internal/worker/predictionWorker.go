package worker

import (
	"context"

	"github.com/Kriwn/GoSeniorWorker/internal/infrastructures/ai"
	"github.com/Kriwn/GoSeniorWorker/internal/infrastructures/rabbitmq"
	"github.com/Kriwn/GoSeniorWorker/internal/repository"
	"github.com/Kriwn/GoSeniorWorker/internal/service"
)

func StartLstmWorker(
	ctx context.Context,
	repo *repository.PredictionRepository,
	mq *rabbitmq.RabbitmqService,
	predictor *ai.LstmPredictService,
) error {

	svc :=  service.NewLstmPredictionService(repo, mq, predictor)

	return svc.Consume(ctx)
}
