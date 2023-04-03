package encoder

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Encode
	}{
		// TODO: Add test cases.
		{
			name: "Positive New",
			want: &Encode{},
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

func TestEncode_Decode(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
		Hash  string
	}
	type args struct {
		body io.Reader
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:   "Positive Decode",
			fields: fields{},
			args: args{
				body: bytes.NewBuffer([]byte(`{"id": "PollCount", "delta": 2345234211616163, "type": "counter"}`)),
			},
			wantErr: false,
		},
		{
			name:   "Negative Decode",
			fields: fields{},
			args: args{
				body: bytes.NewBuffer([]byte(`{id": "PollCount", "delta": 2345234211616163, "type": "counter"}`)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Encode{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
				Hash:  tt.fields.Hash,
			}
			if err := m.Decode(tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("Encode.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncode_Encode(t *testing.T) {
	b := 1.1
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
		Hash  string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Positive Encode",
			fields: fields{
				ID:    "Malloc",
				MType: "gauge",
				Value: &b,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Encode{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
				Hash:  tt.fields.Hash,
			}
			body := &bytes.Buffer{}
			if err := m.Encode(body); (err != nil) != tt.wantErr {
				t.Errorf("Encode.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
