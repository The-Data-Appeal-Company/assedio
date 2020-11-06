package calculator

import (
	"assedio/pkg/model"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
	"time"
)

func TestAssedioStatisticsCalculator_Calculate(t *testing.T) {

	poiListUrl, _ := url.Parse("http://localhost:8081/rates/poi/list?from=2020-11-04&to=2021-02-04&country=italy&state=veneto&poi_category=hotel")
	overallGroupedUrl, _ := url.Parse("http://localhost:8080/overall/grouped?from=2020-04-01&to=2020-10-31&group_type=poi_type&country=italy&state=veneto&county=provincia+di+verona")
	overallGroupedUrl1, _ := url.Parse("http://localhost:8080/overall/grouped?from=2020-04-01&to=2021-10-31&group_type=poi_type&country=italy&state=veneto&county=provincia+di+verona")

	records := model.NewThreadSafeSlice()

	records.Append(model.Record{
		Status:   "200 OK",
		Duration: 1 * time.Second,
		Url:      poiListUrl,
		Error:    false,
	})

	records.Append(model.Record{
		Status:   "200 OK",
		Duration: 2 * time.Second,
		Url:      overallGroupedUrl,
		Error:    false,
	})

	records.Append(model.Record{
		Status:   "200 OK",
		Duration: 1 * time.Second,
		Url:      overallGroupedUrl1,
		Error:    false,
	})

	type args struct {
		records model.Slice
	}
	tests := []struct {
		name          string
		args          args
		wantStats     model.Statistics
		wantPathStats map[string]model.Latencies
	}{
		{
			name: "shouldCalculate",
			args: args{
				records: records,
			},
			wantStats: model.Statistics{
				LatencyStats: model.Latencies{
					AverageLatency: 1.3333333333333333,
					MedianLatency:  1,
					MinLatency:     1,
					MaxLatency:     2,
				},
				Errors:       0,
				Total:        3,
				SuccessRatio: 1.0,
				ErrorRatio:   0.0,
			},
			wantPathStats: map[string]model.Latencies{
				"/rates/poi/list": {
					AverageLatency: float64(1),
					MedianLatency:  float64(1),
					MinLatency:     float64(1),
					MaxLatency:     float64(1),
				},
				"/overall/grouped": {
					AverageLatency: 1.5,
					MedianLatency:  1.5,
					MinLatency:     float64(1),
					MaxLatency:     float64(2),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AssedioStatisticsCalculator{}
			got, got1 := a.Calculate(tt.args.records)
			require.Equal(t, tt.wantStats, got)
			require.Equal(t, tt.wantPathStats, got1)
		})
	}
}
