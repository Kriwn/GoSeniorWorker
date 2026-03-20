package ai

import (
	"context"
	ort "github.com/yalue/onnxruntime_go"
)

type Predictor interface {
	Predict(ctx context.Context, input any) (any, error)
}

type LstmPredictionInput struct {
	ChildId int
	Weight  []float32
	Height  []float32
	Sex     float32
}

type LstmPredictionOutput struct {
	Weight []float32
	Height []float32
}

type LstmPredictService struct {
	Session        *ort.AdvancedSession
	HistoryTensor  *ort.Tensor[float32]
	ForecastTensor *ort.Tensor[float32]
}
