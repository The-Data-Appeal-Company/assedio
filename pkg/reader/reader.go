package reader

import (
	"bufio"
	"net/url"
	"os"
)

type StreamingReader interface {
	Read(fileName string, onConsumeFn func(url *url.URL), onCompleteFn func()) error
}

type FileStreamingReader struct {
}

func (f *FileStreamingReader) Read(fileName string, onConsumeFn func(url *url.URL), onCompleteFn func()) error {
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
