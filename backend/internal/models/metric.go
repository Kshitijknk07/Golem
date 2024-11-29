package models

type Metric struct {
	ID    string  `json:"id" bson:"_id,omitempty"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Type  string  `json:"type"`
	Time  string  `json:"time"`
}
