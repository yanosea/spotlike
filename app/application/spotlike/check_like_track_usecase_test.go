package spotlike

import (
	"context"
	"reflect"
	"testing"

	trackDomain "github.com/yanosea/spotlike/app/domain/spotify/track"

	"go.uber.org/mock/gomock"
)

func TestNewCheckLikeTrackUseCase(t *testing.T) {
	type args struct {
		trackRepo trackDomain.TrackRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *checkLikeTrackUseCase
		setup func(mockCtrl *gomock.Controller, tt *args) *checkLikeTrackUseCase
	}{
		{
			name: "positive testing",
			args: args{
				trackRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *checkLikeTrackUseCase {
				mockTrackRepo := trackDomain.NewMockTrackRepository(mockCtrl)
				tt.trackRepo = mockTrackRepo
				return &checkLikeTrackUseCase{
					trackRepo: mockTrackRepo,
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
			if got := NewCheckLikeTrackUseCase(tt.args.trackRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCheckLikeTrackUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkLikeTrackUseCase_Run(t *testing.T) {
	type fields struct {
		trackRepo trackDomain.TrackRepository
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				trackRepo: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  "test_track_id",
			},
			want:    true,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockTrackRepo := trackDomain.NewMockTrackRepository(mockCtrl)
				mockTrackRepo.EXPECT().IsLiked(gomock.Any(), gomock.Any()).Return(true, nil)
				tt.trackRepo = mockTrackRepo
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
			uc := &checkLikeTrackUseCase{
				trackRepo: tt.fields.trackRepo,
			}
			got, err := uc.Run(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkLikeTrackUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkLikeTrackUseCase.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
