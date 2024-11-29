package models

type Alert struct {
	ID        string  `json:"id" bson:"_id,omitempty"`
	MetricID  string  `json:"metric_id"`
	Threshold float64 `json:"threshold"`
	Severity  string  `json:"severity"`
	Timestamp string  `json:"timestamp"`
}
