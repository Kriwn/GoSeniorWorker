package repository

import (
	"github.com/Kriwn/GoSeniorWorker/internal/infrastructures/ai"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PredictionRepository struct {
	db *gorm.DB
}

func NewPredictionRepository(db *gorm.DB) *PredictionRepository {
	return &PredictionRepository{db: db}
}

func (r *PredictionRepository) CreateAiPrediction(
	childID int,
	model string,
	output ai.LstmPredictionOutput,
	monthUsed int,

) (int, error) {

	var id int

	err := r.db.Raw(`
		INSERT INTO "aiPrediction" (
			"childId",
			"modelUsed",
			height,
			weight,
			"dataMonthsUsed"
		)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id
	`,
		childID,
		model,
		pq.Array(output.Height),
		pq.Array(output.Weight),
		monthUsed,
	).Scan(&id).Error

	if err != nil {
		return 0, err
	}

	return id, nil
}
