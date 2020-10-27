package models

// Device структура входных данных с беспилотника
type SDevice struct {
	MAC       string  `json:"mac" binding:"required"`
	TS        int64   `json:"timestamp" binding:"required"`
	Lat       float64 `json:"lat" binding:"required"`
	Lng       float64 `json:"lng" binding:"required"`
	Speed     float64 `json:"speed" binding:"required"`
	Length    float64 `json:"length" binding:"required"`
	Hindrace  float64 `json:"distObject" binding:"required"`
	Direction float64 `json:"ori" binding:"required"`
}
