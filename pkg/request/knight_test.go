package request

import (
	"assedio/pkg/model"
	"assedio/pkg/test"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type Castle struct{}

func (d *Castle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/error" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

var urls = make(chan *url.URL)

func TestKnight_Hit(t *testing.T) {
	type args struct {
		urls    []string
		results model.Slice
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []model.Record
	}{
		{
			name: "should hit strong",
			args: args{
				urls: []string{
					"/a",
					"/b",
					"/c",
				},
				results: model.NewThreadSafeSlice(),
			},
			wantErr: false,
			want: []model.Record{
				{
					Status: "200 OK",
					Url:    test.ParseUrlOrDie("http://localhost:8080/a"),
					Error:  false,
				},
				{
					Status: "200 OK",
					Url:    test.ParseUrlOrDie("http://localhost:8080/b"),
					Error:  false,
				},
				{
					Status: "200 OK",
					Url:    test.ParseUrlOrDie("http://localhost:8080/c"),
					Error:  false,
				},
			},
		},
		{
			name: "should handle errors",
			args: args{
				urls: []string{
					"/error",
					"/a",
				},
				results: model.NewThreadSafeSlice(),
			},
			wantErr: false,
			want: []model.Record{
				{
					Status: "500 Internal Server Error",
					Url:    test.ParseUrlOrDie("http://localhost:8080/error"),
					Error:  true,
				},
				{
					Status: "200 OK",
					Url:    test.ParseUrlOrDie("http://localhost:8080/a"),
					Error:  false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(&Castle{})
			defer server.Close()

			k := &Knight{}
			c := make(chan *url.URL, len(tt.args.urls))
			for _, s := range tt.args.urls {
				c <- test.ParseUrlOrDie(server.URL + s)
			}
			close(c)
			if err := k.Hit(c, tt.args.results); (err != nil) != tt.wantErr {
				t.Errorf("Hit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, len(tt.want), tt.args.results.Len())
			for i, expected := range tt.want {
				actual := tt.args.results.Get(i)
				require.Equal(t, expected.Status, actual.Status)
				require.Equal(t, expected.Error, actual.Error)
				require.Equal(t, expected.Url.Path, actual.Url.Path)
				require.Equal(t, expected.Url.Query(), actual.Url.Query())
			}
		})
	}
}
