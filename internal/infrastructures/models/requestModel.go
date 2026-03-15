package models

type CreatePredictionRequest struct {
	ChildId int       `json:"childId"`
	Weight  []float32 `json:"weight_kg"`
	Height  []float32 `json:"height_cm"`
	Sex     int       `json:"sex"`
}
