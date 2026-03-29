package spotlike

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/zmb3/spotify/v2"

	trackDomain "github.com/yanosea/spotlike/app/domain/spotify/track"

	"go.uber.org/mock/gomock"
)

func TestNewUnlikeTrackUseCase(t *testing.T) {
	type args struct {
		trackDomain trackDomain.TrackRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *unlikeTrackUseCase
		setup func(mockCtrl *gomock.Controller, tt *args) *unlikeTrackUseCase
	}{
		{
			name: "positive testing",
			args: args{
				trackDomain: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *unlikeTrackUseCase {
				mockTrackRepo := trackDomain.NewMockTrackRepository(mockCtrl)
				tt.trackDomain = mockTrackRepo
				return &unlikeTrackUseCase{
					trackDomain: mockTrackRepo,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.want = tt.setup(mockCtrl, &tt.args)
			}
			if got := NewUnlikeTrackUseCase(tt.args.trackDomain); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUnlikeTrackUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unlikeTrackUseCase_Run(t *testing.T) {
	type fields struct {
		trackDomain trackDomain.TrackRepository
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				trackDomain: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  "test_track_id",
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockTrackRepo := trackDomain.NewMockTrackRepository(mockCtrl)
				mockTrackRepo.EXPECT().Unlike(gomock.Any(), spotify.ID("test_track_id")).Return(nil)
				tt.trackDomain = mockTrackRepo
			},
		},
		{
			name: "negative testing (uc.trackDomain.Unlike() failed)",
			fields: fields{
				trackDomain: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  "test_track_id",
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockTrackRepo := trackDomain.NewMockTrackRepository(mockCtrl)
				mockTrackRepo.EXPECT().Unlike(gomock.Any(), spotify.ID("test_track_id")).Return(errors.New("failed to unlike track"))
				tt.trackDomain = mockTrackRepo
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
			uc := &unlikeTrackUseCase{
				trackDomain: tt.fields.trackDomain,
			}
			if err := uc.Run(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("unlikeTrackUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
