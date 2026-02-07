package config

import (
	"errors"
	"reflect"
	"testing"

	baseConfig "github.com/yanosea/spotlike/app/config"

	"github.com/yanosea/spotlike/pkg/proxy"

	"go.uber.org/mock/gomock"
)

func TestNewSpotlikeCliConfigurator(t *testing.T) {
	envconfig := proxy.NewEnvconfig()

	type args struct {
		envconfigProxy proxy.Envconfig
	}
	tests := []struct {
		name string
		args args
		want SpotlikeCliConfigurator
	}{
		{
			name: "positive testing",
			args: args{
				envconfigProxy: envconfig,
			},
			want: &cliConfigurator{
				BaseConfigurator: baseConfig.NewConfigurator(
					envconfig,
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSpotlikeCliConfigurator(tt.args.envconfigProxy); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSpotlikeCliConfigurator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cliConfigurator_GetConfig(t *testing.T) {
	type fields struct {
		BaseConfigurator *baseConfig.BaseConfigurator
	}
	tests := []struct {
		name    string
		fields  fields
		want    *SpotlikeCliConfig
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				BaseConfigurator: &baseConfig.BaseConfigurator{
					Envconfig: nil,
				}},
			want: &SpotlikeCliConfig{
				SpotlikeConfig: baseConfig.SpotlikeConfig{
					SpotifyID:           "test_id",
					SpotifySecret:       "test_secret",
					SpotifyRedirectUri:  "test_redirect_uri",
					SpotifyRefreshToken: "test_refresh_token",
				},
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockEnvconfig := proxy.NewMockEnvconfig(mockCtrl)
				mockEnvconfig.EXPECT().Process("", gomock.Any()).DoAndReturn(
					func(_ string, cfg *envConfig) error {
						cfg.SpotifyID = "test_id"
						cfg.SpotifySecret = "test_secret"
						cfg.SpotifyRedirectUri = "test_redirect_uri"
						cfg.SpotifyRefreshToken = "test_refresh_token"
						return nil
					},
				)
				tt.BaseConfigurator.Envconfig = mockEnvconfig
			},
		},
		{
			name: "negative testing (c.Envconfig.Process(\"\", &config) failed)",
			fields: fields{
				BaseConfigurator: &baseConfig.BaseConfigurator{
					Envconfig: nil,
				}},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockEnvconfig := proxy.NewMockEnvconfig(mockCtrl)
				mockEnvconfig.EXPECT().Process("", gomock.Any()).Return(errors.New("Envconfig.Process() failed"))
				tt.BaseConfigurator.Envconfig = mockEnvconfig
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			c := &cliConfigurator{
				BaseConfigurator: tt.fields.BaseConfigurator,
			}
			got, err := c.GetConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("cliConfigurator.GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cliConfigurator.GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
