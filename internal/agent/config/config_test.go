package config

import (
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		// TODO: Add test cases.
		{
			name: "Positive congif",
			want: &Config{
				Address:        "127.0.0.1:8080",
				ReportInterval: 10 * time.Second,
				PollInterval:   2 * time.Second,
				Key:            "",
				Rate:           100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
