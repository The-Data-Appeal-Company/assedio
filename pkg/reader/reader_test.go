package reader

import (
	"net/url"
	"testing"
)

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
	}{
		{
			name:    "should stream file",
			args:    args{
				fileName:     "test_data/targets",
				onConsumeFn:  nil,
				onCompleteFn: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileStreamingReader{}
			if err := f.Read(tt.args.fileName, tt.args.onConsumeFn, tt.args.onCompleteFn); (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
