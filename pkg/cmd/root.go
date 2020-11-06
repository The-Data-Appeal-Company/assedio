package cmd

import (
	"assedio/pkg/calculator"
	"assedio/pkg/model"
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var results = model.NewThreadSafeSlice()

var rootCmd = &cobra.Command{
	Use: "fight",
	Run: func(cmd *cobra.Command, args []string) {

		SetupCloseHandler()

		color.Red(fmt.Sprintf("Starting assedio with %d knights ...", flagConcurrency))

		errGroup, _ := errgroup.WithContext(cmd.Context())
		client := &http.Client{
			Timeout: 1 * time.Minute,
		}

		if flagFile == "" {
			log.Fatal("file not specified.")
		}

		urls := make(chan *url.URL)
		errGroup.Go(func() error {
			return GetUrl(flagFile, func(parsedUrl *url.URL) {
				urls <- parsedUrl
			}, func() {
				close(urls)
			})
		})

		for i := 0; i < flagConcurrency; i++ {
			errGroup.Go(func() error {
				for url := range urls {

					uri := url.String()

					startReqTime := time.Now()
					resp, err := client.Get(uri)
					requestTotTimeDt := time.Now().Sub(startReqTime)

					logUrl := fmt.Sprintf("%s", url.Path)

					if err != nil {
						color.Red(fmt.Sprintf("[ERR] -  Chiamata %s", logUrl))

						results.Append(model.Record{
							Status:   "ERROR",
							Duration: requestTotTimeDt,
							Url:      url,
							Error:    true,
						})

						continue
					}

					if resp.StatusCode == 200 {
						color.Cyan(fmt.Sprintf("[%s] %d ms Chiamata %s", resp.Status, requestTotTimeDt.Milliseconds(), logUrl))
						results.Append(model.Record{
							Status:   resp.Status,
							Duration: requestTotTimeDt,
							Url:      url,
							Error:    false,
						})

					} else {
						color.Magenta(fmt.Sprintf("[%s] %d ms Chiamata %s", resp.Status, requestTotTimeDt.Milliseconds(), logUrl))
						results.Append(model.Record{
							Status:   resp.Status,
							Duration: requestTotTimeDt,
							Url:      url,
							Error:    true,
						})
					}

				}
				return nil
			})
		}

		err := errGroup.Wait()

		showAssedioEsito()

		if err != nil {
			log.Fatal(err)
		}
	},
}

func showAssedioEsito() {
	calculator := calculator.AssedioStatisticsCalculator{}
	esito, pathsEsito := calculator.Calculate(results)

	fmt.Println(esito.String())

	for path, esitoPerPath := range pathsEsito {
		fmt.Println(path)
		fmt.Println(esitoPerPath.String())
	}
}

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Stopping the assedio, retreating all knights from battlefield.")
		showAssedioEsito()
		os.Exit(0)
	}()
}

var (
	flagConcurrency int
	flagFile        string
)

func init() {

	rootCmd.PersistentFlags().IntVarP(&flagConcurrency, "knights", "k", 1, "number of simultaneously knights making the assedio (default=1)")
	rootCmd.PersistentFlags().StringVarP(&flagFile, "file", "f", "", "assedio file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("sap")
	}
}

func GetUrl(fileName string, onConsumeFn func(url *url.URL), onCompleteFn func()) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer file.Close()
	defer onCompleteFn()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		parsedUrl, err := url.Parse(scanner.Text())

		if err != nil {
			return err
		}

		onConsumeFn(parsedUrl)
	}

	return scanner.Err()
}
