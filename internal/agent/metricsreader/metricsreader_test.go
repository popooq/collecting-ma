package metricsreader

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/popooq/collectimg-ma/internal/agent/config"
	"github.com/popooq/collectimg-ma/internal/agent/sender"
	"github.com/popooq/collectimg-ma/internal/utils/hasher"
)

func TestNew(t *testing.T) {
	// var hasher hasher.Hash
	type args struct {
		sndr         sender.Sender
		tickerpoll   time.Duration
		tickerreport time.Duration
		address      string
		rate         int
	}
	config := config.Config{
		Address:        "127.0.0.1:8080",
		ReportInterval: 10 * time.Second,
		PollInterval:   2 * time.Second,
		Rate:           100,
	}

	tests := []struct {
		name string
		args args
		want *Reader
	}{
		// TODO: Add test cases.
		{
			name: "Positive New",
			args: args{
				sndr:         sender.Sender{},
				tickerpoll:   config.PollInterval,
				tickerreport: config.ReportInterval,
				address:      config.Address,
				rate:         config.Rate,
			},
			want: &Reader{
				sndr:         sender.Sender{},
				tickerpoll:   2000000000,
				tickerreport: 10000000000,
				address:      "127.0.0.1:8080",
				rate:         100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.sndr, tt.args.tickerpoll, tt.args.tickerreport, tt.args.address, tt.args.rate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_worker_queueTask(t *testing.T) {
	type fields struct {
		workchan   chan metrics
		buffer     int
		wg         *sync.WaitGroup
		cancelFunc context.CancelFunc
		sndr       sender.Sender
	}
	type args struct {
		mem metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "positive test",
			fields: fields{
				workchan: make(chan metrics, 100),
				buffer:   100,
				wg:       new(sync.WaitGroup),
				sndr:     sender.New(hasher.Mew(""), "127.0.0.1:8080", ""),
			},
			args: args{
				mem: metrics{
					value: 3.16,
					name:  "Alloc",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &worker{
				workchan:   tt.fields.workchan,
				buffer:     tt.fields.buffer,
				wg:         tt.fields.wg,
				cancelFunc: tt.fields.cancelFunc,
				sndr:       tt.fields.sndr,
			}
			if err := w.queueTask(tt.args.mem); (err != nil) != tt.wantErr {
				t.Errorf("worker.queueTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
