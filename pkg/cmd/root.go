package cmd

import (
	"assedio/pkg/reader"
	"assedio/pkg/render"
	"assedio/pkg/request"
	"github.com/spf13/cobra"
	"log"
)

var (
	flagConcurrency int
	flagFile        string
)

var rootCmd = &cobra.Command{
	Use: "Fight",
	Run: func(cmd *cobra.Command, args []string) {
		army := NewFightArmy(
			&reader.FileStreamingReader{},
			&request.Knight{},
			render.NewAsciiWarBulletin(),
		)
		err := army.Fight(cmd.Context(), flagFile, flagConcurrency)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&flagConcurrency, "knights", "k", 1, "number of simultaneously knights making the Fight (default=1)")
	rootCmd.PersistentFlags().StringVarP(&flagFile, "file", "f", "", "Fight file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("we cannot fight my Lord", err)
	}
}
