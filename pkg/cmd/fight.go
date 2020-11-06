package cmd

import (
	"assedio/pkg/model"
	"assedio/pkg/reader"
	"assedio/pkg/render"
	"assedio/pkg/request"
	"context"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/sync/errgroup"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

type FightArmy struct {
	reader   reader.StreamingReader
	knight   request.Fighter
	bulletin render.WarBulletin
	results  model.Slice
}

func NewFightArmy(reader reader.StreamingReader, knight request.Fighter, bulletin render.WarBulletin) *FightArmy {
	return &FightArmy{reader: reader, knight: knight, bulletin: bulletin, results: model.NewThreadSafeSlice()}
}

func (fa *FightArmy) Fight(ctx context.Context, fileName string, nKnights int) error {
	if fileName == "" {
		return fmt.Errorf("file not specified")
	}

	if nKnights < 1 {
		return fmt.Errorf("you can't fight without knights")
	}

	fa.SetupCloseHandler()
	color.Red(fmt.Sprintf("Fighting with %d knights ...", nKnights))

	errGroup, _ := errgroup.WithContext(ctx)

	urls := make(chan *url.URL)
	errGroup.Go(func() error {
		return fa.reader.Read(fileName, func(parsedUrl *url.URL) {
			urls <- parsedUrl
		}, func() {
			close(urls)
		})
	})

	for i := 0; i < nKnights; i++ {
		errGroup.Go(func() error {
			return fa.knight.Hit(urls, fa.results)
		})
	}

	err := errGroup.Wait()

	fa.bulletin.Render(fa.results)
	return err
}

func (fa *FightArmy) SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Stopping the Fight, retreating all knights from battlefield.")
		fa.bulletin.Render(fa.results)
		os.Exit(0)
	}()
}
