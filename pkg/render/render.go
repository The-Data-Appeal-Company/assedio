package render

import (
	"assedio/pkg/calculator"
	"assedio/pkg/model"
	"github.com/olekukonko/tablewriter"
	"os"
)

type WarBulletin interface {
	Render(results model.Slice)
}

type AsciiWarBulletin struct {
	calculator calculator.StatisticsCalculator
}

func NewAsciiWarBulletin() *AsciiWarBulletin {
	return &AsciiWarBulletin{calculator: &calculator.AssedioStatisticsCalculator{}}
}

func (a *AsciiWarBulletin) Render(results model.Slice) {
	globalBulletin, groupedBulletin := a.calculator.Calculate(results)
	headers := []string{"Path", "Average", "Median", "Min", "Max", "Total", "Errors", "Success Ratio", "Error Ratio"}
	table := newTable(headers)
	oddColors := getColor(tablewriter.Colors{tablewriter.FgHiBlackColor}, len(headers))
	evenColors := getColor(tablewriter.Colors{tablewriter.FgWhiteColor}, len(headers))

	i := 1
	for path, group := range groupedBulletin {
		row := append([]string{path}, group.Strings()...)
		if i%2 == 1 {
			table.Rich(row, oddColors)
		} else {
			table.Rich(row, evenColors)
		}
		i++
	}
	table.SetFooter(append([]string{"Total"}, globalBulletin.Strings()...))
	table.Render()
}

func getColor(color tablewriter.Colors, len int) []tablewriter.Colors {
	colors := make([]tablewriter.Colors, len)
	for i := 0; i < len; i++ {
		colors[i] = color
	}
	return colors
}

func newTable(header []string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetBorder(false)
	return table
}
