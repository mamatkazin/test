package common

type S_ID struct {
	ID int `json:"id"`
}

type S_Data struct {
	Valid  bool        `json:"valid"`
	Errors []string    `json:"errors"`
	Items  interface{} `json:"items"`
}

type S_List struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CookieData struct {
	UserID int
}

type S_Bounds struct {
	SW S_Coordinate `json:"sw"`
	NE S_Coordinate `json:"ne"`
}

type S_Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type S_Geometry_Polygon struct {
	Coordinates [][][]float64 `json:"coordinates"`
	Type        string        `json:"type"`
}

type S_Geometry_MultiPolygon struct {
	Coordinates [][][][]float64 `json:"coordinates"`
	Type        string          `json:"type"`
}

type S_Geometry_LineString struct {
	Coordinates [][]float64 `json:"coordinates"`
	Type        string      `json:"type"`
}

type S_Geometry_MultiLineString struct {
	Coordinates [][][]float64 `json:"coordinates"`
	Type        string        `json:"type"`
}

type S_Geometry_Point struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

type S_Icon struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type S_IconExp struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Size  string `json:"size"`
}

// type S_Select struct {
// 	ID   int    `json:"value"`
// 	Name string `json:"label"`
// }

type S_Render struct {
	Position S_XYZ `json:"position"`
	Offset   S_XYZ `json:"offset"`
	Rotation S_XYZ `json:"rotation"`
	Scale    S_XYZ `json:"scale"`
}

type S_XYZ struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}
