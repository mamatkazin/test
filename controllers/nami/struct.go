package nami

//import (
//	"nami/nami_ds/controllers/common"
//)

type S_MultiPolygon struct {
	Features []S_Feature_MultiPolygon `json:"features"`
	Type     string                   `json:"type"`
}

type S_Feature_MultiPolygon struct {
	Geometry S_Geometry_MultiPolygon `json:"geometry"`
	//Properties interface{}             `json:"properties"`
	Type string `json:"type"`
}

type S_Geometry_MultiPolygon struct {
	Coordinates [][][][]float64 `json:"coordinates"`
	Type        string          `json:"type"`
}
