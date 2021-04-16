package main

import (
	"github.com/labstack/echo"
	"github.com/meyskens/helmet/template"
	"reflect"
	"testing"
)

func Test_buildManifests(t *testing.T) {
	type args struct {
		c     echo.Context
		data  putData
		name  string
		chart *template.Chart
	}
	tests := []struct {
		name    string
		args    args
		want    [][]byte
		wantErr bool
	}{
		{
			name: "test simple chart",
			args: args{
				c: nil,
				data: putData{
					Values: map[string]interface{}{"test": "ok"},
				},
				name:  "test",
				chart: template.New(&template.ChartInfo{Name: "test"}, map[interface{}]interface{}{"test": "nok"}, map[string][]byte{"test": []byte("{{.Values.test}}")}),
			},
			want:    [][]byte{[]byte("ok")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chart = tt.args.chart
			got, err := buildManifests(tt.args.c, tt.args.data, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildManifests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildManifests() got = %v, want %v", got, tt.want)
			}
		})
	}
}
