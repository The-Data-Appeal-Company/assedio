package calculator

import (
	"assedio/pkg/model"
	"github.com/montanaflynn/stats"
)

type StatisticsCalculator interface {
	Calculate(records model.Slice) (model.Statistics, map[string]model.Statistics)
}

type AssedioStatisticsCalculator struct {
}

func (a *AssedioStatisticsCalculator) Calculate(records model.Slice) (model.Statistics, map[string]model.Statistics) {
	total := records.Len()
	errs := a.getErrors(records)

	latencies := make([]float64, 0)
	pathLatencies := make(map[string][]float64)

	for i := 0; i < records.Len(); i++ {
		rec := records.Get(i)
		if !rec.Error {
			pathLatencies[rec.Url.Path] = append(pathLatencies[rec.Url.Path], rec.Duration.Seconds())
			latencies = append(latencies, rec.Duration.Seconds())
		}
	}

	//calculate total stats
	data := stats.LoadRawData(latencies)
	median, _ := stats.Median(data)
	max, _ := stats.Max(data)
	min, _ := stats.Min(data)
	avg, _ := stats.Mean(data)

	//calculate stats per api

	pathStats := make(map[string]model.Statistics)

	for path, records := range pathLatencies {

		data := stats.LoadRawData(records)
		median, _ := stats.Median(data)
		max, _ := stats.Max(data)
		min, _ := stats.Min(data)
		avg, _ := stats.Mean(data)

		pathStats[path] = model.Statistics{
			LatencyStats: model.Latencies{
				AverageLatency: avg,
				MedianLatency:  median,
				MinLatency:     min,
				MaxLatency:     max,
			},
			Errors:       0,
			Total:        0,
			SuccessRatio: 0,
			ErrorRatio:   0,
		}
	}

	return model.Statistics{
		LatencyStats: model.Latencies{
			AverageLatency: avg,
			MedianLatency:  median,
			MinLatency:     min,
			MaxLatency:     max},
		Errors:       errs,
		Total:        total,
		SuccessRatio: (float64(total) - float64(errs)) / float64(total),
		ErrorRatio:   float64(errs) / float64(total),
	}, pathStats

}

func (a *AssedioStatisticsCalculator) getErrors(records model.Slice) int {
	errors := 0

	for i := 0; i < records.Len(); i++ {
		if records.Get(i).Error {
			errors++
		}
	}

	return errors

}
