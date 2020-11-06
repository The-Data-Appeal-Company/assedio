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
	table := newTable([]string{"Path", "Average", "Median", "Min", "Max", "Total", "Errors", "Success Ratio", "Error Ratio"})
	table.Render()

	i := 1
	for path, group := range groupedBulletin {
		row := append([]string{path}, group.Strings()...)
		if i%2 == 1 {
			table.Rich(row, []tablewriter.Colors{{tablewriter.FgCyanColor}})
		} else {
			table.Rich(row, []tablewriter.Colors{{tablewriter.FgWhiteColor}})
		}
		i++
	}
	table.SetFooter(append([]string{"Total"}, globalBulletin.Strings()...))
	table.Render()
}

func newTable(header []string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetBorder(false)
	return table
}
