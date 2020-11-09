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
}

func NewFightArmy(reader reader.StreamingReader, knight request.Fighter, bulletin render.WarBulletin) *FightArmy {
	return &FightArmy{reader: reader, knight: knight, bulletin: bulletin}
}

func (fa *FightArmy) Fight(parentCtx context.Context, fileName string, nKnights int) error {
	if fileName == "" {
		return fmt.Errorf("file not specified")
	}

	if nKnights < 1 {
		return fmt.Errorf("you can't fight without knights")
	}

	results := model.NewThreadSafeSlice()
	ctx, cancel := context.WithCancel(parentCtx)
	color.Red(fmt.Sprintf("Fighting with %d knights ...", nKnights))
	errGroup, ctx := errgroup.WithContext(ctx)
	fa.SetupCloseHandler(cancel)

	urls := make(chan *url.URL)
	errGroup.Go(func() error {
		return fa.reader.Read(fileName, ctx, func(parsedUrl *url.URL) {
			urls <- parsedUrl
		}, func() {
			close(urls)
		})
	})

	for i := 0; i < nKnights; i++ {
		errGroup.Go(func() error {
			return fa.knight.Hit(urls, results)
		})
	}

	err := errGroup.Wait()

	fa.bulletin.Render(results)
	return err
}

func (fa *FightArmy) SetupCloseHandler(cancel context.CancelFunc) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Stopping the Fight, retreating all knights from battlefield.")
		cancel()
	}()
}
