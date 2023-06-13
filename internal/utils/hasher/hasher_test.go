package hasher

import (
	"reflect"
	"testing"

	"github.com/popooq/collectimg-ma/internal/utils/encoder"
)

func TestMew(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want *Hash
	}{
		// TODO: Add test cases.
		{name: "nil",
			args: args{
				key: "",
			},
			want: &Hash{
				Key: []byte(""),
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mew() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHash_Hasher(t *testing.T) {
	var (
		counter int64 = 3
	)
	gauge := 3.16
	type fields struct {
		Key []byte
	}
	type args struct {
		metric *encoder.Encode
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "nil",
			fields: fields{
				Key: nil,
			},
			args: args{
				metric: &encoder.Encode{},
			},
			want: "",
		},
		{
			name: "not nil gauge",
			fields: fields{
				Key: []byte("secret"),
			},
			args: args{
				metric: &encoder.Encode{
					ID:    "alloc",
					MType: "gauge",
					Value: &gauge,
				},
			},
			want: "103dd1c679d2f535e610f0d7de5080a1721d98cc73af100adab251517c629833",
		},
		{
			name: "not nil counner",
			fields: fields{
				Key: []byte("secret"),
			},
			args: args{
				metric: &encoder.Encode{
					ID:    "counter",
					MType: "counter",
					Delta: &counter,
				},
			},
			want: "614779ded95129afb8de2578b932a7f77d9a3c5ed05ecb08b2b7a1beb685591b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hsh := &Hash{
				Key: tt.fields.Key,
			}
			if got := hsh.Hasher(tt.args.metric); got != tt.want {
				t.Errorf("Hash.Hasher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHash_HashChecker(t *testing.T) {
	var (
		counter int64 = 3
	)
	type fields struct {
		Key []byte
	}
	type args struct {
		hash   string
		metric encoder.Encode
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "positive",
			fields: fields{
				Key: []byte("secret"),
			},
			args: args{
				hash: "614779ded95129afb8de2578b932a7f77d9a3c5ed05ecb08b2b7a1beb685591b",
				metric: encoder.Encode{
					ID:    "counter",
					MType: "counter",
					Delta: &counter,
					Hash:  "614779ded95129afb8de2578b932a7f77d9a3c5ed05ecb08b2b7a1beb685591b",
				},
			},
			wantErr: false,
		},
		{
			name: "negative",
			fields: fields{
				Key: []byte("secret"),
			},
			args: args{
				hash: "614779ded95129afb8de2578b932a7f77d9a3c5ed05ecb08b2b7a1beb685591b",
				metric: encoder.Encode{
					ID:    "counter",
					MType: "counter",
					Delta: &counter,
					Hash:  "614779ded95129af8de2578b932a7f77d9a3c5ed05ecb08b2b7a1beb685591b",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hsh := &Hash{
				Key: tt.fields.Key,
			}
			if err := hsh.HashChecker(tt.args.hash, tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("Hash.HashChecker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
