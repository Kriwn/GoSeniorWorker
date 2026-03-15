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


type ArimaPredictionInput struct {
	ChildId int
	Weight  []float32
	Height  []float32
	Sex     int
}

type ArimaPredictionOutput struct {
	Weight []float32
	Height []float32
}

type ArimaPredictService struct {
	weightSession        *ort.AdvancedSession
	weightHistoryTensor  *ort.Tensor[float32]
	weightErrorTensor    *ort.Tensor[float32]
	weightForecastTensor *ort.Tensor[float32]

	heightSession        *ort.AdvancedSession
	heightHistoryTensor  *ort.Tensor[float32]
	heightErrorTensor    *ort.Tensor[float32]
	heightForecastTensor *ort.Tensor[float32]
}
