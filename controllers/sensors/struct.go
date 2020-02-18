package sensors

type S_Sensor struct {
	MAC       string `json:"mac"`
	Timestamp int64  `json:"timestamp"`
	Value     bool   `json:"value"`
}
