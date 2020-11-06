package cmd

import (
	"assedio/pkg/model"
	"assedio/pkg/reader"
	"assedio/pkg/render"
	"assedio/pkg/request"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

var (
	flagConcurrency int
	flagFile        string
	results         = model.NewThreadSafeSlice()
)

var rootCmd = &cobra.Command{
	Use: "Fight",
	Run: func(cmd *cobra.Command, args []string) {
		army := FightArmy{
			reader:   &reader.FileStreamingReader{},
			knight:   &request.Knight{},
			bulletin: render.NewAsciiWarBulletin(),
		}
		err := army.Fight(cmd.Context())
		if err != nil {
			log.Fatal(err)
		}
	},
}

type FightArmy struct {
	reader   reader.StreamingReader
	knight   request.Fighter
	bulletin render.WarBulletin
}

func (fa *FightArmy) Fight(ctx context.Context) error {
	if flagFile == "" {
		return fmt.Errorf("file not specified")
	}

	fa.SetupCloseHandler()
	color.Red(fmt.Sprintf("Starting Fight with %d knights ...", flagConcurrency))

	errGroup, _ := errgroup.WithContext(ctx)

	urls := make(chan *url.URL)
	errGroup.Go(func() error {
		return fa.reader.Read(flagFile, func(parsedUrl *url.URL) {
			urls <- parsedUrl
		}, func() {
			close(urls)
		})
	})

	for i := 0; i < flagConcurrency; i++ {
		errGroup.Go(func() error {
			return fa.knight.Hit(urls, results)
		})
	}

	err := errGroup.Wait()

	fa.bulletin.Render(results)
	return err
}

func (fa *FightArmy) SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Stopping the Fight, retreating all knights from battlefield.")
		fa.bulletin.Render(results)
		os.Exit(0)
	}()
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&flagConcurrency, "knights", "k", 1, "number of simultaneously knights making the Fight (default=1)")
	rootCmd.PersistentFlags().StringVarP(&flagFile, "file", "f", "", "Fight file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("sap")
	}
}
