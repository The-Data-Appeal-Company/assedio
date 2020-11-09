package cmd

import (
	"assedio/pkg/model"
	"assedio/pkg/reader"
	"assedio/pkg/render"
	"assedio/pkg/request"
	"context"
	"fmt"
	"net/url"
	"testing"
)

type NoOffensiveKnight struct{}

func (n *NoOffensiveKnight) Hit(urls chan *url.URL, results model.Slice) error {
	for u := range urls {
		fmt.Printf("hit %s\n", u)
	}
	return nil
}

type DontTellAnybodyBulletin struct{}

func (d *DontTellAnybodyBulletin) Render(results model.Slice) {}

func TestFightArmy_Fight(t *testing.T) {
	type fields struct {
		reader   reader.StreamingReader
		knight   request.Fighter
		bulletin render.WarBulletin
	}
	type args struct {
		ctx      context.Context
		fileName string
		nKnights int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "shouldErrorWhenNoFileGiven",
			fields: fields{
				reader:   &reader.FileStreamingReader{},
				knight:   &NoOffensiveKnight{},
				bulletin: &DontTellAnybodyBulletin{},
			},
			args: args{
				ctx:      context.TODO(),
				fileName: "",
				nKnights: 2,
			},
			wantErr: true,
		},
		{
			name: "shouldErrorWhenNoKnights",
			fields: fields{
				reader:   &reader.FileStreamingReader{},
				knight:   &NoOffensiveKnight{},
				bulletin: &DontTellAnybodyBulletin{},
			},
			args: args{
				ctx:      context.TODO(),
				fileName: "test_data/targets",
				nKnights: 0,
			},
			wantErr: true,
		},
		{
			name: "shouldFightHard",
			fields: fields{
				reader:   &reader.FileStreamingReader{},
				knight:   &NoOffensiveKnight{},
				bulletin: &DontTellAnybodyBulletin{},
			},
			args: args{
				ctx:      context.TODO(),
				fileName: "test_data/targets",
				nKnights: 10,
			},
			wantErr: false,
		},
		{
			name: "shouldErrorWhenNocSuchFile",
			fields: fields{
				reader:   &reader.FileStreamingReader{},
				knight:   &NoOffensiveKnight{},
				bulletin: &DontTellAnybodyBulletin{},
			},
			args: args{
				ctx:      context.TODO(),
				fileName: "test_data/non_esisto",
				nKnights: 10,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fa := NewFightArmy(
				tt.fields.reader,
				tt.fields.knight,
				tt.fields.bulletin,
			)
			if err := fa.Fight(tt.args.ctx, tt.args.fileName, tt.args.nKnights); (err != nil) != tt.wantErr {
				t.Errorf("Fight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
