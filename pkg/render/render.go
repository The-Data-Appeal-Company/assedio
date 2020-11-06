package render

import (
	"assedio/pkg/calculator"
	"assedio/pkg/model"
	"fmt"
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
	esito, pathsEsito := a.calculator.Calculate(results)

	fmt.Println(esito.String())

	for path, esitoPerPath := range pathsEsito {
		fmt.Println(path)
		fmt.Println(esitoPerPath.String())
	}
}
