package reader

import (
	"assedio/pkg/test"
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
	"time"
)

var consumed []*url.URL
var completedCalled bool

func TestFileStreamingReader_Read(t *testing.T) {
	type args struct {
		fileName     string
		ctx          context.Context
		onConsumeFn  func(url *url.URL)
		onCompleteFn func()
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []*url.URL
	}{
		{
			name: "should stream file skipping empty lines",
			args: args{
				fileName: "test_data/targets",
				ctx:      context.TODO(),
				onConsumeFn: func(url *url.URL) {
					consumed = append(consumed, url)
				},
				onCompleteFn: func() {
					completedCalled = true
				},
			},
			wantErr: false,
			want: []*url.URL{
				test.ParseUrlOrDie("https://trippa.io:8000"),
				test.ParseUrlOrDie("lampre.dotto"),
				test.ParseUrlOrDie("http://antani.it?clacsonava=true&supercazzola=sinistra"),
			},
		},
		{
			name: "should error when no such file",
			args: args{
				fileName: "test_data/nulla",
				ctx:      context.TODO(),
				onConsumeFn: func(url *url.URL) {
					consumed = append(consumed, url)
				},
				onCompleteFn: func() {
					completedCalled = true
				},
			},
			wantErr: true,
		},
		{
			name: "should error when invalid url",
			args: args{
				fileName: "test_data/targets_invalid_urls",
				ctx:      context.TODO(),
				onConsumeFn: func(url *url.URL) {
					consumed = append(consumed, url)
				},
				onCompleteFn: func() {
					completedCalled = true
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumed = make([]*url.URL, 0)
			completedCalled = false
			f := &FileStreamingReader{}
			if err := f.Read(tt.args.fileName, tt.args.ctx, tt.args.onConsumeFn, tt.args.onCompleteFn); (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				require.True(t, completedCalled)
				require.Equal(t, tt.want, consumed)
			}
		})
	}
}

func TestShouldStopReadingWhenContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	f := &FileStreamingReader{}
	go f.Read("test_data/targets_lots_of", ctx, func(url *url.URL) {
		consumed = append(consumed, url)
	}, func() {
		completedCalled = true
	})
	time.Sleep(25 * time.Microsecond)
	cancel()
	time.Sleep(1 * time.Millisecond)
	require.Less(t, len(consumed), 200)
	require.True(t, completedCalled)
}
