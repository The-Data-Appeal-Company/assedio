package reader

import (
	"bufio"
	"context"
	"net/url"
	"os"
)

type StreamingReader interface {
	Read(fileName string, ctx context.Context, onConsumeFn func(url *url.URL), onCompleteFn func()) error
}

type FileStreamingReader struct{}

func (f *FileStreamingReader) Read(fileName string, ctx context.Context, onConsumeFn func(url *url.URL), onCompleteFn func()) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer file.Close()
	defer onCompleteFn()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return scanner.Err()
		default:
			text := scanner.Text()
			if text == "" {
				continue
			}
			parsedUrl, err := url.Parse(text)

			if err != nil {
				return err
			}

			onConsumeFn(parsedUrl)
		}
	}

	return scanner.Err()
}
