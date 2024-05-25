package storage

import (
	"go.uber.org/zap"
	"testing"
)

func TestEvents_Add(t *testing.T) {
	type fields struct {
		Events []Event
	}
	type args struct {
		short       string
		url         string
		storagePath string
		baseUrl     string
		sugar       *zap.SugaredLogger
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &Events{
				Events: tt.fields.Events,
			}
			es.Add(tt.args.short, tt.args.url, tt.args.storagePath, tt.args.baseUrl, tt.args.sugar)
		})
	}
}
