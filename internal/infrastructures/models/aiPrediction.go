package models

import "time"

type AIPrediction struct {
	ID                  int       `json:"id"`
	ChildID             int       `json:"childId"`
	HeightDevelopmentID int       `json:"heightDevelopmentId"`
	WeightDevelopmentID int       `json:"weightDevelopmentId"`
	DataMonthsUsed      int16     `json:"dataMonthsUsed"`
	ModelUsed           string    `json:"modelUsed"`
	Height              []float32 `json:"height"`
	Weight              []float32 `json:"weight"`
	CreatedAt           time.Time `json:"createdAt"`

	// Relation IDs are stored above. Add relation structs here when those types exist.
	// Child             Child       `json:"child"`
	// HeightDevelopment Development `json:"heightDevelopment"`
	// WeightDevelopment Development `json:"weightDevelopment"`
}
