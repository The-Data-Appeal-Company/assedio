package render

import (
	"assedio/pkg/calculator"
	"assedio/pkg/model"
	"fmt"
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
	table := newTable()
	globalBulletin, groupedBulletin := a.calculator.Calculate(results)
	table.Append(globalBulletin.Strings())
	table.Render()
	for path, group := range groupedBulletin {
		fmt.Println(path)
		table = newTable()
		table.Append(group.Strings())
		table.Render()
	}
}

func newTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Average", "Median", "Min", "Max", "Total", "Errors", "Success Ratio", "Error Ratio"})
	table.SetBorder(false)
	return table
}
