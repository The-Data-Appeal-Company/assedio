package reader

import (
	"assedio/pkg/test"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

var consumed []*url.URL
var completedCalled bool

func TestFileStreamingReader_Read(t *testing.T) {
	type args struct {
		fileName     string
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
				onConsumeFn: func(url *url.URL) {
					consumed = append(consumed, url)
				},
				onCompleteFn: func() {
					completedCalled = true
				},
			},
			wantErr: false,
			want: []*url.URL{
				test.ParseUrlOrDie("http://trippa.io"),
				test.ParseUrlOrDie("http://lampre.dotto"),
				test.ParseUrlOrDie("http://antani.it"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumed = make([]*url.URL, 0)
			completedCalled = false
			f := &FileStreamingReader{}
			if err := f.Read(tt.args.fileName, tt.args.onConsumeFn, tt.args.onCompleteFn); (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.True(t, completedCalled)
			require.Equal(t, tt.want, consumed)
		})
	}
}
