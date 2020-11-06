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
	totTable := newTable([]string{"Average", "Median", "Min", "Max", "Total", "Errors", "Success Ratio", "Error Ratio"})
	totTable.Append(globalBulletin.Strings())
	totTable.Render()

	groupTable := newTable([]string{"Path", "Average", "Median", "Min", "Max", "Total", "Errors", "Success Ratio", "Error Ratio"})
	for path, group := range groupedBulletin {
		groupTable.Append(append([]string{path}, group.Strings()...))
	}
	groupTable.Render()
}

func newTable(header []string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetBorder(false)
	return table
}
