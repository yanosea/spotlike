package api

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/zmb3/spotify/v2"
)

func TestGetClient(t *testing.T) {
	tests := []struct {
		name    string
		want    *spotify.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setAuthInfo(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setAuthInfo(); (err != nil) != tt.wantErr {
				t.Errorf("setAuthInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_authenticate(t *testing.T) {
	tests := []struct {
		name    string
		want    *spotify.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := authenticate()
			if (err != nil) != tt.wantErr {
				t.Errorf("authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authenticate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_completeAuthenticate(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			completeAuthenticate(tt.args.w, tt.args.r)
		})
	}
}
