package storage

import (
	"reflect"
	"testing"

	"github.com/popooq/collectimg-ma/internal/utils/encoder"
)

type keeper struct {
}

func (k keeper) SaveMetric(metric *encoder.Encode) error {
	return nil
}
func (k keeper) SaveAllMetrics(metric encoder.Encode) error {
	return nil
}
func (k keeper) LoadMetrics() ([]encoder.Encode, error) {
	return nil, nil
}
func (k keeper) KeeperCheck() error {
	return nil
}

func TestMetricsStorage_InsertMetric(t *testing.T) {
	var keeper keeper
	type fields struct {
		Keeper         Keeper
		MetricsGauge   map[string]float64
		MetricsCounter map[string]int64
	}
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "positive",
			fields: fields{
				Keeper:         keeper,
				MetricsGauge:   make(map[string]float64),
				MetricsCounter: make(map[string]int64),
			},
			args: args{
				name:  "alloc",
				value: 3.16,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MetricsStorage{
				Keeper:         tt.fields.Keeper,
				MetricsGauge:   tt.fields.MetricsGauge,
				MetricsCounter: tt.fields.MetricsCounter,
			}
			ms.InsertMetric(tt.args.name, tt.args.value)
		})
	}
}

func TestMetricsStorage_CountCounterMetric(t *testing.T) {
	var keeper keeper
	type fields struct {
		Keeper         Keeper
		MetricsGauge   map[string]float64
		MetricsCounter map[string]int64
	}
	type args struct {
		name  string
		value int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "positive",
			fields: fields{
				Keeper:         keeper,
				MetricsGauge:   make(map[string]float64),
				MetricsCounter: make(map[string]int64),
			},
			args: args{
				name:  "counter",
				value: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MetricsStorage{
				Keeper:         tt.fields.Keeper,
				MetricsGauge:   tt.fields.MetricsGauge,
				MetricsCounter: tt.fields.MetricsCounter,
			}
			ms.CountCounterMetric(tt.args.name, tt.args.value)
		})
	}
}

func TestMetricsStorage_GetMetricGauge(t *testing.T) {
	var keeper keeper
	type fields struct {
		Keeper         Keeper
		MetricsGauge   map[string]float64
		MetricsCounter map[string]int64
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "positive",
			fields: fields{
				Keeper:         keeper,
				MetricsGauge:   map[string]float64{"alloc": 3.16},
				MetricsCounter: map[string]int64{"counter": 3},
			},
			args: args{
				name: "alloc",
			},
			want:    3.16,
			wantErr: false,
		},
		{
			name: "negative",
			fields: fields{
				Keeper:         keeper,
				MetricsGauge:   map[string]float64{"alloc": 3.16},
				MetricsCounter: map[string]int64{"counter": 3},
			},
			args: args{
				name: "aloc",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MetricsStorage{
				Keeper:         tt.fields.Keeper,
				MetricsGauge:   tt.fields.MetricsGauge,
				MetricsCounter: tt.fields.MetricsCounter,
			}
			got, err := ms.GetMetricGauge(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricsStorage.GetMetricGauge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MetricsStorage.GetMetricGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricsStorage_GetMetricJSONGauge(t *testing.T) {
	var (
		keeper keeper
	)
	gauge := 3.16
	type fields struct {
		Keeper         Keeper
		MetricsGauge   map[string]float64
		MetricsCounter map[string]int64
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *float64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "positive",
			fields: fields{
				Keeper:         keeper,
				MetricsGauge:   map[string]float64{"alloc": 3.16},
				MetricsCounter: map[string]int64{"counter": 3},
			},
			args: args{
				name: "alloc",
			},
			want:    &gauge,
			wantErr: false,
		},
		{
			name: "negative",
			fields: fields{
				Keeper:         keeper,
				MetricsGauge:   map[string]float64{"alloc": 3.16},
				MetricsCounter: map[string]int64{"counter": 3},
			},
			args: args{
				name: "aloc",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MetricsStorage{
				Keeper:         tt.fields.Keeper,
				MetricsGauge:   tt.fields.MetricsGauge,
				MetricsCounter: tt.fields.MetricsCounter,
			}
			got, err := ms.GetMetricJSONGauge(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricsStorage.GetMetricJSONGauge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if *got != *tt.want {
					t.Errorf("MetricsStorage.GetMetricJSONGauge() = %f, want %f", *got, *tt.want)
				}
			}
			if tt.want == nil {
				if got != tt.want {
					t.Errorf("MetricsStorage.GetMetricJSONCounter() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestMetricsStorage_GetMetricCounter(t *testing.T) {
	var (
		keeper  keeper
		counter int64 = 3
	)
	type fields struct {
		Keeper         Keeper
		MetricsGauge   map[string]float64
		MetricsCounter map[string]int64
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "positive",
			fields: fields{
				Keeper:         keeper,
				MetricsGauge:   map[string]float64{"alloc": 3.16},
				MetricsCounter: map[string]int64{"counter": 3},
			},
			args: args{
				name: "counter",
			},
			want:    counter,
			wantErr: false,
		},
		{
			name: "negative",
			fields: fields{
				Keeper:         keeper,
				MetricsGauge:   map[string]float64{"alloc": 3.16},
				MetricsCounter: map[string]int64{"counter": 3},
			},
			args: args{
				name: "aloc",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MetricsStorage{
				Keeper:         tt.fields.Keeper,
				MetricsGauge:   tt.fields.MetricsGauge,
				MetricsCounter: tt.fields.MetricsCounter,
			}
			got, err := ms.GetMetricCounter(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricsStorage.GetMetricCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MetricsStorage.GetMetricCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricsStorage_GetMetricJSONCounter(t *testing.T) {
	var (
		keeper  keeper
		counter int64 = 3
	)
	type fields struct {
		Keeper         Keeper
		MetricsGauge   map[string]float64
		MetricsCounter map[string]int64
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "positive",
			fields: fields{
				Keeper:         keeper,
				MetricsGauge:   map[string]float64{"alloc": 3.16},
				MetricsCounter: map[string]int64{"counter": 3},
			},
			args: args{
				name: "counter",
			},
			want:    &counter,
			wantErr: false,
		},
		{
			name: "negative",
			fields: fields{
				Keeper:         keeper,
				MetricsGauge:   map[string]float64{"alloc": 3.16},
				MetricsCounter: map[string]int64{"counter": 3},
			},
			args: args{
				name: "aloc",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MetricsStorage{
				Keeper:         tt.fields.Keeper,
				MetricsGauge:   tt.fields.MetricsGauge,
				MetricsCounter: tt.fields.MetricsCounter,
			}
			got, err := ms.GetMetricJSONCounter(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricsStorage.GetMetricJSONCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if *got != *tt.want {
					t.Errorf("MetricsStorage.GetMetricJSONCounter() = %v, want %v", got, tt.want)
				}
			}
			if tt.want == nil {
				if got != tt.want {
					t.Errorf("MetricsStorage.GetMetricJSONCounter() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestMetricsStorage_InsertMetrics(t *testing.T) {
	var (
		keeper  keeper
		counter int64 = 3
	)
	gauge := 3.16
	type fields struct {
		Keeper         Keeper
		MetricsGauge   map[string]float64
		MetricsCounter map[string]int64
	}
	type args struct {
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
			name: "positive gauge",
			fields: fields{
				Keeper:         keeper,
				MetricsGauge:   map[string]float64{"alloc": 3.16},
				MetricsCounter: map[string]int64{"counter": 3},
			},
			args: args{
				encoder.Encode{
					MType: "gauge",
					ID:    "alloc",
					Value: &gauge,
				},
			},
			wantErr: false,
		},
		{
			name: "positive counter",
			fields: fields{
				Keeper:         keeper,
				MetricsGauge:   map[string]float64{"alloc": 3.16},
				MetricsCounter: map[string]int64{"counter": 3},
			},
			args: args{
				encoder.Encode{
					MType: "counter",
					ID:    "counter",
					Delta: &counter,
				},
			},
			wantErr: false,
		},
		{
			name: "negative",
			fields: fields{
				Keeper:         keeper,
				MetricsGauge:   map[string]float64{"alloc": 3.16},
				MetricsCounter: map[string]int64{"counter": 3},
			},
			args: args{
				encoder.Encode{
					MType: "counddter",
					ID:    "counter",
					Delta: &counter,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MetricsStorage{
				Keeper:         tt.fields.Keeper,
				MetricsGauge:   tt.fields.MetricsGauge,
				MetricsCounter: tt.fields.MetricsCounter,
			}
			if err := ms.InsertMetrics(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("MetricsStorage.InsertMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetricsStorage_GetAllMetrics(t *testing.T) {
	var (
		keeper  keeper
		counter int64 = 3
	)
	gauge := 3.16
	type fields struct {
		Keeper         Keeper
		MetricsGauge   map[string]float64
		MetricsCounter map[string]int64
	}
	tests := []struct {
		name   string
		fields fields
		want   []encoder.Encode
	}{
		// TODO: Add test cases.
		{
			name: "positive",
			fields: fields{
				Keeper:         keeper,
				MetricsGauge:   map[string]float64{"alloc": 3.16},
				MetricsCounter: map[string]int64{"counter": 3},
			},
			want: []encoder.Encode{
				{ID: "alloc", MType: "gauge", Delta: nil, Value: &gauge, Hash: ""},
				{ID: "counter", MType: "counter", Delta: &counter, Value: nil, Hash: ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MetricsStorage{
				Keeper:         tt.fields.Keeper,
				MetricsGauge:   tt.fields.MetricsGauge,
				MetricsCounter: tt.fields.MetricsCounter,
			}
			if got := ms.GetAllMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricsStorage.GetAllMetrics() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
