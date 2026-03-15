package ai

import (
	"context"
	"fmt"
	ort "github.com/yalue/onnxruntime_go"
	"os"
)

func NewArimaPredictService() (*ArimaPredictService, error) {

	model := os.Getenv("MODEL_NAME")
	if model == "" {
		return nil, fmt.Errorf("MODEL_NAME environment variable is required")
	}

	var wegihtModelPath, heightModelPath string

	wegihtModelPath = "internal/infrastructures/ai/arima/forcast12/arima_weight_model.onnx"
	heightModelPath = "internal/infrastructures/ai/arima/forcast12/arima_height_model.onnx"

	ort.SetSharedLibraryPath("/opt/homebrew/lib/libonnxruntime.dylib")
	var ortInitialized bool

	if !ortInitialized {
		ort.SetSharedLibraryPath("/opt/homebrew/lib/libonnxruntime.dylib")
		err := ort.InitializeEnvironment()
		if err != nil {
			return nil, err
		}
		ortInitialized = true
	}

	weightHistoryShape := ort.NewShape(1, 3)
	weightHistoryTensor, err := ort.NewEmptyTensor[float32](weightHistoryShape)
	if err != nil {
		return nil, fmt.Errorf("failed to create weight history tensor: %w", err)
	}

	weightErrorShape := ort.NewShape(1, 1)
	weightErrorTensor, err := ort.NewEmptyTensor[float32](weightErrorShape)
	if err != nil {
		return nil, fmt.Errorf("failed to create error tensor: %w", err)
	}

	weightForecastShape := ort.NewShape(1, 12)
	weightForecastTensor, err := ort.NewEmptyTensor[float32](weightForecastShape)
	if err != nil {
		return nil, fmt.Errorf("failed to create forecast tensor: %w", err)
	}

	weightSession, err := ort.NewAdvancedSession(
		wegihtModelPath,
		[]string{"history", "errors"},
		[]string{"forecast"},
		[]ort.Value{weightHistoryTensor, weightErrorTensor},
		[]ort.Value{weightForecastTensor},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create weight session: %w", err)
	}

	fmt.Println("Weight session created:", weightSession)

	heightHistoryShape := ort.NewShape(1, 3)
	heightHistoryTensor, err := ort.NewEmptyTensor[float32](heightHistoryShape)
	if err != nil {
		return nil, fmt.Errorf("failed to create height history tensor: %w", err)
	}

	heightErrorShape := ort.NewShape(1, 1)
	heightErrorTensor, err := ort.NewEmptyTensor[float32](heightErrorShape)
	if err != nil {
		return nil, fmt.Errorf("failed to create error tensor: %w", err)
	}

	heightForecastShape := ort.NewShape(1, 12)
	heightForecastTensor, err := ort.NewEmptyTensor[float32](heightForecastShape)
	if err != nil {
		return nil, fmt.Errorf("failed to create height forecast tensor: %w", err)
	}

	heightSession, err := ort.NewAdvancedSession(
		heightModelPath,
		[]string{"history", "errors"},
		[]string{"forecast"},
		[]ort.Value{heightHistoryTensor, heightErrorTensor},
		[]ort.Value{heightForecastTensor},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create height session: %w", err)
	}

	return &ArimaPredictService{
		weightSession:        weightSession,
		weightHistoryTensor:  weightHistoryTensor,
		weightErrorTensor:    weightErrorTensor,
		weightForecastTensor: weightForecastTensor,
		heightSession:        heightSession,
		heightHistoryTensor:  heightHistoryTensor,
		heightErrorTensor:    heightErrorTensor,
		heightForecastTensor: heightForecastTensor,
	}, nil
}

func (p *ArimaPredictService) predictWeight(ctx context.Context, data []float32) ([]float32, error) {

	copy(p.weightHistoryTensor.GetData(), data)

	errs := make([]float32, len(data))
	copy(p.weightErrorTensor.GetData(), errs)

	err := p.weightSession.Run()
	if err != nil {
		return nil, err
	}

	return p.weightForecastTensor.GetData(), nil
}

func (p *ArimaPredictService) predictHeight(ctx context.Context, data []float32) ([]float32, error) {

	copy(p.heightHistoryTensor.GetData(), data)

	errs := make([]float32, len(data))
	copy(p.heightErrorTensor.GetData(), errs)
	err := p.heightSession.Run()

	if err != nil {
		return nil, err
	}

	return p.heightForecastTensor.GetData(), nil
}

func (p *ArimaPredictService) PredictArima(ctx context.Context, input ArimaPredictionInput) (ArimaPredictionOutput, error) {

	weightPred, err := p.predictWeight(ctx, input.Weight)
	if err != nil {
		return ArimaPredictionOutput{}, err
	}

	heightPred, err := p.predictHeight(ctx, input.Height)
	if err != nil {
		return ArimaPredictionOutput{}, err
	}

	return ArimaPredictionOutput{
		Weight: weightPred,
		Height: heightPred,
	}, nil
}
