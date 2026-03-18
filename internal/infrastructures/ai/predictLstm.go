package ai

import (
	"context"
	"fmt"
	"os"

	ort "github.com/yalue/onnxruntime_go"
)

func NewLstmPredictService() (*LstmPredictService, error) {

	libPath := os.Getenv("LIB_PATH")
	if libPath == "" {
		return nil, fmt.Errorf("LIB_PATH environment variable is required")
	}


	modelPath := "internal/infrastructures/ai/lstm/forcast6/growth_lstm_model.onnx"

	ort.SetSharedLibraryPath(libPath)

	err := ort.InitializeEnvironment()
	if err != nil {
		return nil, err
	}

	// history shape = (1,3,3)
	historyShape := ort.NewShape(1, 3, 3)
	historyTensor, err := ort.NewEmptyTensor[float32](historyShape)
	if err != nil {
		return nil, err
	}

	// forecast shape = (1,6,2)
	forecastShape := ort.NewShape(1, 6, 2)
	forecastTensor, err := ort.NewEmptyTensor[float32](forecastShape)
	if err != nil {
		return nil, err
	}

	session, err := ort.NewAdvancedSession(
		modelPath,
		[]string{"history"},
		[]string{"forecast"},
		[]ort.Value{historyTensor},
		[]ort.Value{forecastTensor},
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &LstmPredictService{
		Session:        session,
		HistoryTensor:  historyTensor,
		ForecastTensor: forecastTensor,
	}, nil
}

func (p *LstmPredictService) lstm(ctx context.Context, data []float32) (LstmPredictionOutput, error) {

	copy(p.HistoryTensor.GetData(), data)

	err := p.Session.Run()
	if err != nil {
		return LstmPredictionOutput{}, err
	}

	temp := p.ForecastTensor.GetData()

	horizon := len(temp) / 2
	weights := make([]float32, horizon)
	heights := make([]float32, horizon)

	for i := 0; i < horizon; i++ {
    	weights[i] = temp[i*2]
    	heights[i] = temp[i*2+1]
	}

	return LstmPredictionOutput{
		Weight: weights,
		Height: heights,
	}, nil
}

func (p *LstmPredictService) PredictLstm(ctx context.Context, input LstmPredictionInput) (LstmPredictionOutput, error) {

	seqLen := len(input.Weight)

	scaledInput := converToScalar(input)

	history := make([]float32, seqLen*3)

	for i := 0; i < seqLen; i++ {
		history[i*3+0] = scaledInput.Weight[i]
		history[i*3+1] = scaledInput.Height[i]
		history[i*3+2] = scaledInput.Sex
	}

	pred, err := p.lstm(ctx, history)
	if err != nil {
		return LstmPredictionOutput{}, err
	}

	len := len(pred.Weight)
	scaledOutput := inverseScale(pred)
	weights := make([]float32, len)
	heights := make([]float32, len)

	for i := 0; i < len; i++ {
		weights[i] = scaledOutput.Weight[i]
		heights[i] = scaledOutput.Height[i]
	}

	return LstmPredictionOutput{
		Weight: weights,
		Height: heights,
	}, nil
}

func converToScalar(input LstmPredictionInput) LstmPredictionInput {
	weightMean := float32(10.43825597)
	weightStd := float32(3.2316772)

	heightMean := float32(77.58884794)
	heightStd := float32(12.7406077)

	SexMean := float32(1.48714286)
	SexStd := float32(0.49983467)

	len := len(input.Weight)
	weight := make([]float32, len)
	height := make([]float32, len)
	sex	:= (float32(input.Sex) - SexMean) / SexStd
	for i := 0; i < len; i++ {

		scaledWeight := (input.Weight[i] - weightMean) / weightStd
		scaledHeight := (input.Height[i] - heightMean) / heightStd
		weight[i] = scaledWeight
		height[i] = scaledHeight
	}

	return LstmPredictionInput{
		ChildId: input.ChildId,
		Weight:  weight,
		Height:  height,
		Sex: sex,
	}
}

func inverseScale(pred LstmPredictionOutput) LstmPredictionOutput {
	weightMean := float32(10.43825597)
	weightStd := float32(3.2316772)

	heightMean := float32(77.58884794)
	heightStd := float32(12.7406077)

	height := make([]float32, len(pred.Height))
	weight := make([]float32, len(pred.Weight))

	for i := 0; i < len(pred.Weight); i++ {
		weight[i] = pred.Weight[i]*weightStd + weightMean
		height[i] = pred.Height[i]*heightStd + heightMean
	}

	return LstmPredictionOutput{
		Weight: weight,
		Height: height,
	}

}
