package teams

//import (
//	"nami/nami_ds/controllers/common"
//)

type S_Point struct {
	Features []S_Feature_Point `json:"features"`
	Type     string            `json:"type"`
}

type S_Feature_Point struct {
	Geometry S_Geometry_Point `json:"geometry"`
	Type     string           `json:"type"`
}

type S_Geometry_Point struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}
