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

	pathRecords := make(map[string][]model.Record)

	for _, rec := range records.ToSlice() {
		pathRecords[rec.Url.Path] = append(pathRecords[rec.Url.Path], rec)
	}

	totalStats := a.getStatistic(records.ToSlice())

	pathStats := make(map[string]model.Statistics)

	for path, records := range pathRecords {
		pathStats[path] = a.getStatistic(records)
	}

	return totalStats, pathStats

}

func (a *AssedioStatisticsCalculator) getStatistic(records []model.Record) model.Statistics {

	total := len(records)
	errs := a.getErrors(records)

	latencies := make([]float64, 0)

	for _, rec := range records {
		if !rec.Error {
			latencies = append(latencies, rec.Duration.Seconds())
		}
	}

	//calculate total stats
	data := stats.LoadRawData(latencies)
	median, _ := stats.Median(data)
	max, _ := stats.Max(data)
	min, _ := stats.Min(data)
	avg, _ := stats.Mean(data)

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
	}

}

func (a *AssedioStatisticsCalculator) getErrors(records []model.Record) int {
	errors := 0
	for _, record := range records {
		if record.Error {
			errors++
		}
	}

	return errors

}
