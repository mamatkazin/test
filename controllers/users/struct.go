package users

// import (
// 	"nami/nami_ds/controllers/common"
// )

type S_Login struct {
	ID       int        `json:"id"`
	District S_District `json:"district"`
}

type S_District struct {
	ID int        `json:"id"`
	LL [2]float64 `json:"latlng"`
}
