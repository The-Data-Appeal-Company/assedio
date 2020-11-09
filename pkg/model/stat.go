package model

import (
	"fmt"
	"net/url"
	"time"
)

type Record struct {
	Status   string
	Duration time.Duration
	Url      *url.URL
	Error    bool
}

type Latencies struct {
	AverageLatency float64
	MedianLatency  float64
	MinLatency     float64
	MaxLatency     float64
}

func (s *Latencies) String() string {
	return fmt.Sprintf(`

	AverageLatency %f
	MedianLatency  %f
	MinLatency     %f
	MaxLatency     %f

`, s.AverageLatency, s.MedianLatency, s.MinLatency, s.MaxLatency)
}

type Statistics struct {
	LatencyStats Latencies
	Errors       int
	Total        int
	SuccessRatio float64
	ErrorRatio   float64
}

func (s *Statistics) String() string {
	return fmt.Sprintf(`

	Errors         %d
	Total          %d
	AverageLatency %f
	MedianLatency  %f
	MinLatency     %f
	MaxLatency     %f
	SuccessRatio   %f
	ErrorRatio     %f

`, s.Errors, s.Total, s.LatencyStats.AverageLatency, s.LatencyStats.MedianLatency, s.LatencyStats.MinLatency, s.LatencyStats.MaxLatency, s.SuccessRatio, s.ErrorRatio)
}

func (s *Statistics) Strings() []string {
	return []string{
		fmt.Sprintf("%f", s.LatencyStats.AverageLatency),
		fmt.Sprintf("%f", s.LatencyStats.MedianLatency),
		fmt.Sprintf("%f", s.LatencyStats.MinLatency),
		fmt.Sprintf("%f", s.LatencyStats.MaxLatency),
		fmt.Sprintf("%d", s.Total),
		fmt.Sprintf("%d", s.Errors),
		fmt.Sprintf("%f", s.SuccessRatio),
		fmt.Sprintf("%f", s.ErrorRatio),
	}
}
