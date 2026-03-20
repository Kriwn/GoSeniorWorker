package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/Kriwn/GoSeniorWorker/internal/infrastructures/ai"
	"github.com/Kriwn/GoSeniorWorker/internal/infrastructures/postgre"
	"github.com/Kriwn/GoSeniorWorker/internal/infrastructures/rabbitmq"
	"github.com/Kriwn/GoSeniorWorker/internal/repository"
	"github.com/Kriwn/GoSeniorWorker/internal/worker"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	workersStr := os.Getenv("WORKER_COUNT")

	workers, err := strconv.Atoi(workersStr)
	if err != nil {
		log.Fatalf("invalid WORKER_COUNT: %v", err)
	}
	ctx := context.Background()

	db, err := postgre.ConnectDB()
	if err != nil {
		log.Fatalf("DB init failed: %v", err)
	}

	predictor, err := ai.NewLstmPredictService()
	if err != nil {
		log.Fatalf("AI init failed: %v", err)
	}

	repo := repository.NewPredictionRepository(db)

	for i := 0; i < workers; i++ {

		go func(id int) {

			mq, err := rabbitmq.NewRabbitmqService()
			if err != nil {
				log.Fatalf("RabbitMQ init failed: %v", err)
			}
			log.Printf("worker %d started", id)

			err = worker.StartLstmWorker(ctx, repo, mq, predictor)
			if err != nil {
				log.Printf("worker %d crashed: %v", id, err)
			}

		}(i)
	}

	select {}
}
