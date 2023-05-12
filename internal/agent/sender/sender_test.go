package sender

import (
	"reflect"
	"testing"

	"github.com/popooq/collectimg-ma/internal/storage"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
)

var endpoint string = "127.0.0.1:8080"

func TestNew(t *testing.T) {
	type args struct {
		hasher   *hasher.Hash
		endpoint string
	}
	tests := []struct {
		name string
		args args
		want Sender
	}{
		// TODO: Add test cases.
		{
			name: "testNew",
			args: args{
				hasher:   nil,
				endpoint: endpoint,
			},
			want: Sender{
				endpoint: endpoint,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.hasher, tt.args.endpoint, tt.want.encryptor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestSender_Go(t *testing.T) {
	type fields struct {
		hasher   *hasher.Hash
		endpoint string
	}
	type args struct {
		value any
		name  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "positive test",
			fields: fields{
				hasher:   hasher.Mew(""),
				endpoint: endpoint,
			},
			args: args{
				value: 3.16,
				name:  "Alloc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sender{
				hasher:   tt.fields.hasher,
				endpoint: tt.fields.endpoint,
			}
			s.Go(tt.args.value, tt.args.name)
		})
	}
}

func TestSender_bodyBuild(t *testing.T) {
	var counter storage.Counter = 3
	type fields struct {
		hasher   *hasher.Hash
		endpoint string
	}
	type args struct {
		value any
		name  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		// TODO: Add test cases.
		{
			name: "positive test gauge",
			fields: fields{
				hasher:   hasher.Mew(""),
				endpoint: endpoint,
			},
			args: args{
				value: 3.16,
				name:  "Alloc",
			},
			want: []byte(`{"id":"Alloc","type":"gauge","value":3.16,"hash":"a6eac70297c6622a6ac0ddbe0e8dc82a18a2063ae12776c2bdccbd49561ae530"}`),
		},
		{
			name: "positive test gauge",
			fields: fields{
				hasher:   hasher.Mew(""),
				endpoint: endpoint,
			},
			args: args{
				value: counter,
				name:  "Counter",
			},
			want: []byte(`{"id":"Counter","type":"counter","delta":3,"hash":"623fbf1b4ac4be1b3dcc64214386b55375b99a5c448568f53570f7b7a49706ff"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sender{
				hasher:   tt.fields.hasher,
				endpoint: tt.fields.endpoint,
			}
			if got := s.bodyBuild(tt.args.value, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sender.bodyBuild() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
