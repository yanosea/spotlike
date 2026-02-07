package config

import (
	"reflect"
	"testing"

	"github.com/yanosea/spotlike/pkg/proxy"
)

func TestNewConfigurator(t *testing.T) {
	envconfig := proxy.NewEnvconfig()

	type args struct {
		envconfigProxy proxy.Envconfig
	}
	tests := []struct {
		name string
		args args
		want *BaseConfigurator
	}{
		{
			name: "positive testing",
			args: args{
				envconfigProxy: envconfig,
			},
			want: &BaseConfigurator{
				Envconfig: envconfig,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfigurator(tt.args.envconfigProxy); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigurator() = %v, want %v", got, tt.want)
			}
		})
	}
}
