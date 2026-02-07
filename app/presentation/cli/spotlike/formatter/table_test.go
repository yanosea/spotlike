package formatter

import (
	"errors"
	"reflect"
	"testing"
	"time"

	spotlikeApp "github.com/yanosea/spotlike/app/application/spotlike"

	"github.com/yanosea/spotlike/pkg/proxy"
	"github.com/yanosea/spotlike/pkg/utility"
	"go.uber.org/mock/gomock"
)

// setTableWriterUtil sets the package-level table writer util variable.
func setTableWriterUtil(tw utility.TableWriterUtil) {
	Tu = tw
}

// getDefaultTableWriterUtil returns the default table writer util.
func getDefaultTableWriterUtil() utility.TableWriterUtil {
	return utility.NewTableWriterUtil(proxy.NewTableWriter())
}

func TestNewTableFormatter(t *testing.T) {
	tests := []struct {
		name string
		want *TableFormatter
	}{
		{
			name: "positive testing",
			want: &TableFormatter{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTableFormatter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTableFormatter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableFormatter_Format(t *testing.T) {
	su := utility.NewStringsUtil()

	type args struct {
		result any
	}
	tests := []struct {
		name    string
		f       *TableFormatter
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "positive testing (result is SearchArtistUseCaseOutputDto)",
			f:    &TableFormatter{},
			args: args{
				result: []*spotlikeApp.SearchArtistUseCaseOutputDto{
					{
						ID:   "artist_id_1",
						Name: "artist_name_1",
					},
					{
						ID:   "artist_id_2",
						Name: "artist_name_2",
					},
				},
			},
			want:    "🆔ID🎤ARTISTartist_id_1artist_name_1artist_id_2artist_name_2TOTAL:2artists!",
			wantErr: false,
		},
		{
			name: "positive testing (result is GetArtistUseCaseOutputDto)",
			f:    &TableFormatter{},
			args: args{
				result: []*spotlikeApp.GetArtistUseCaseOutputDto{
					{
						ID:   "artist_id_1",
						Name: "artist_name_1",
					},
					{
						ID:   "artist_id_2",
						Name: "artist_name_2",
					},
				},
			},
			want:    "🆔ID🎤ARTISTartist_id_1artist_name_1artist_id_2artist_name_2TOTAL:2artists!",
			wantErr: false,
		},
		{
			name: "positive testing (result is SearchAlbumUseCaseOutputDto)",
			f:    &TableFormatter{},
			args: args{
				result: []*spotlikeApp.SearchAlbumUseCaseOutputDto{
					{
						ID:          "album_id_1",
						Artists:     "artist_name_1",
						Name:        "album_name_1",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:          "album_id_2",
						Artists:     "artist_name_2",
						Name:        "album_name_2",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want:    "🆔ID💿ALBUM🎤ARTISTS📅RELEASEDATEalbum_id_1album_name_1artist_name_12000-01-01album_id_2album_name_2artist_name_22000-01-01TOTAL:2albums!",
			wantErr: false,
		},
		{
			name: "positive testing (result is GetAlbumUseCaseOutputDto)",
			f:    &TableFormatter{},
			args: args{
				result: []*spotlikeApp.GetAlbumUseCaseOutputDto{
					{
						ID:          "album_id_1",
						Name:        "album_name_1",
						Artists:     "artist_name_1",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:          "album_id_2",
						Name:        "album_name_2",
						Artists:     "artist_name_2",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want:    "🆔ID💿ALBUM🎤ARTISTS📅RELEASEDATEalbum_id_1album_name_1artist_name_12000-01-01album_id_2album_name_2artist_name_22000-01-01TOTAL:2albums!",
			wantErr: false,
		},
		{
			name: "positive testing (result is SearchTrackUseCaseOutputDto)",
			f:    &TableFormatter{},
			args: args{
				result: []*spotlikeApp.SearchTrackUseCaseOutputDto{
					{
						ID:          "track_id_1",
						Artists:     "artist_name_1",
						Album:       "album_name_1",
						Name:        "track_name_1",
						TrackNumber: 1,
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:          "track_id_2",
						Artists:     "artist_name_2",
						Album:       "album_name_2",
						Name:        "track_name_2",
						TrackNumber: 2,
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want:    "🆔ID🔢NUMBER🎵TRACK💿ALBUM🎤ARTISTS📅RELEASEDATEtrack_id_11track_name_1album_name_1artist_name_12000-01-01track_id_22track_name_2album_name_2artist_name_22000-01-01TOTAL:2tracks!",
			wantErr: false,
		},
		{
			name: "positive testing (result is GetTrackUseCaseOutputDto)",
			f:    &TableFormatter{},
			args: args{
				result: []*spotlikeApp.GetTrackUseCaseOutputDto{
					{
						ID:          "track_id_1",
						Name:        "track_name_1",
						Artists:     "artist_name_1",
						Album:       "album_name_1",
						TrackNumber: 1,
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:          "track_id_2",
						Name:        "track_name_2",
						Artists:     "artist_name_2",
						Album:       "album_name_2",
						TrackNumber: 2,
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want:    "🆔ID🔢NUMBER🎵TRACK💿ALBUM🎤ARTISTS📅RELEASEDATEtrack_id_11track_name_1album_name_1artist_name_12000-01-01track_id_22track_name_2album_name_2artist_name_22000-01-01TOTAL:2tracks!",
			wantErr: false,
		},
		{
			name: "negative testing (result is invalid)",
			f:    &TableFormatter{},
			args: args{
				result: "invalid",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &TableFormatter{}
			got, err := f.Format(tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("TableFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStr := su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(got))); gotStr != tt.want {
				t.Errorf("TableFormatter.Format() = %v, want %v", gotStr, tt.want)
			}
		})
	}
}

func TestTableFormatter_formatSearchArtists(t *testing.T) {
	type args struct {
		items []*spotlikeApp.SearchArtistUseCaseOutputDto
	}
	tests := []struct {
		name string
		f    *TableFormatter
		args args
		want tableData
	}{
		{
			name: "positive testing (items is empty)",
			f:    &TableFormatter{},
			args: args{
				items: []*spotlikeApp.SearchArtistUseCaseOutputDto{},
			},
			want: tableData{
				header: []string{"🆔 ID", "🎤 Artist"},
				rows:   [][]string{},
			},
		},
		{
			name: "positive testing (items is not empty)",
			f:    &TableFormatter{},
			args: args{
				items: []*spotlikeApp.SearchArtistUseCaseOutputDto{
					{
						ID:   "artist_id_1",
						Name: "artist_name_1",
					},
					{
						ID:   "artist_id_2",
						Name: "artist_name_2",
					},
				},
			},
			want: tableData{
				header: []string{"🆔 ID", "🎤 Artist"},
				rows: [][]string{
					{"artist_id_1", "artist_name_1"},
					{"artist_id_2", "artist_name_2"},
					{"", ""},
					{"TOTAL : 2 artists!", ""},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &TableFormatter{}
			if got := f.formatSearchArtists(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TableFormatter.formatSearchArtists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableFormatter_formatGetArtists(t *testing.T) {
	type args struct {
		items []*spotlikeApp.GetArtistUseCaseOutputDto
	}
	tests := []struct {
		name string
		f    *TableFormatter
		args args
		want tableData
	}{
		{
			name: "positive testing (items is empty)",
			f:    &TableFormatter{},
			args: args{
				items: []*spotlikeApp.GetArtistUseCaseOutputDto{},
			},
			want: tableData{
				header: []string{"🆔 ID", "🎤 Artist"},
				rows:   [][]string{},
			},
		},
		{
			name: "positive testing (items is not empty)",
			f:    &TableFormatter{},
			args: args{
				items: []*spotlikeApp.GetArtistUseCaseOutputDto{
					{
						ID:   "artist_id_1",
						Name: "artist_name_1",
					},
					{
						ID:   "artist_id_2",
						Name: "artist_name_2",
					},
				},
			},
			want: tableData{
				header: []string{"🆔 ID", "🎤 Artist"},
				rows: [][]string{
					{"artist_id_1", "artist_name_1"},
					{"artist_id_2", "artist_name_2"},
					{"", ""},
					{"TOTAL : 2 artists!", ""},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &TableFormatter{}
			if got := f.formatGetArtists(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TableFormatter.formatGetArtists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableFormatter_formatSearchAlbums(t *testing.T) {
	type args struct {
		items []*spotlikeApp.SearchAlbumUseCaseOutputDto
	}
	tests := []struct {
		name string
		f    *TableFormatter
		args args
		want tableData
	}{
		{
			name: "positive testing (items is empty)",
			f:    &TableFormatter{},
			args: args{
				items: []*spotlikeApp.SearchAlbumUseCaseOutputDto{},
			},
			want: tableData{
				header: []string{"🆔 ID", "💿 Album", "🎤 Artists", "📅 Release Date"},
				rows:   [][]string{},
			},
		},
		{
			name: "positive testing (items is not empty)",
			f:    &TableFormatter{},
			args: args{
				items: []*spotlikeApp.SearchAlbumUseCaseOutputDto{
					{
						ID:          "album_id_1",
						Name:        "album_name_1",
						Artists:     "artist_name_1",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:          "album_id_2",
						Name:        "album_name_2",
						Artists:     "artist_name_2",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want: tableData{
				header: []string{"🆔 ID", "💿 Album", "🎤 Artists", "📅 Release Date"},
				rows: [][]string{
					{"album_id_1", "album_name_1", "artist_name_1", "2000-01-01"},
					{"album_id_2", "album_name_2", "artist_name_2", "2000-01-01"},
					{"", "", "", ""},
					{"TOTAL : 2 albums!", "", "", ""},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &TableFormatter{}
			if got := f.formatSearchAlbums(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TableFormatter.formatSearchAlbums() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableFormatter_formatGetAlbums(t *testing.T) {
	type args struct {
		items []*spotlikeApp.GetAlbumUseCaseOutputDto
	}
	tests := []struct {
		name string
		f    *TableFormatter
		args args
		want tableData
	}{
		{
			name: "positive testing (items is empty)",
			f:    &TableFormatter{},
			args: args{
				items: []*spotlikeApp.GetAlbumUseCaseOutputDto{},
			},
			want: tableData{
				header: []string{"🆔 ID", "💿 Album", "🎤 Artists", "📅 Release Date"},
				rows:   [][]string{},
			},
		},
		{
			name: "positive testing (items is not empty)",
			f:    &TableFormatter{},
			args: args{
				items: []*spotlikeApp.GetAlbumUseCaseOutputDto{
					{
						ID:          "album_id_1",
						Name:        "album_name_1",
						Artists:     "artist_name_1",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:          "album_id_2",
						Name:        "album_name_2",
						Artists:     "artist_name_2",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want: tableData{
				header: []string{"🆔 ID", "💿 Album", "🎤 Artists", "📅 Release Date"},
				rows: [][]string{
					{"album_id_1", "album_name_1", "artist_name_1", "2000-01-01"},
					{"album_id_2", "album_name_2", "artist_name_2", "2000-01-01"},
					{"", "", "", ""},
					{"TOTAL : 2 albums!", "", "", ""},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &TableFormatter{}
			if got := f.formatGetAlbums(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TableFormatter.formatGetAlbums() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableFormatter_formatSearchTracks(t *testing.T) {
	type args struct {
		items []*spotlikeApp.SearchTrackUseCaseOutputDto
	}
	tests := []struct {
		name string
		f    *TableFormatter
		args args
		want tableData
	}{
		{
			name: "positive testing (items is empty)",
			f:    &TableFormatter{},
			args: args{
				items: []*spotlikeApp.SearchTrackUseCaseOutputDto{},
			},
			want: tableData{
				header: []string{"🆔 ID", "🔢 Number", "🎵 Track", "💿 Album", "🎤 Artists", "📅 Release Date"},
				rows:   [][]string{},
			},
		},
		{
			name: "positive testing (items is not empty)",
			f:    &TableFormatter{},
			args: args{
				items: []*spotlikeApp.SearchTrackUseCaseOutputDto{
					{
						ID:          "track_id_1",
						Name:        "track_name_1",
						Album:       "album_name_1",
						Artists:     "artist_name_1",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						TrackNumber: 1,
					},
					{
						ID:          "track_id_2",
						Name:        "track_name_2",
						Album:       "album_name_2",
						Artists:     "artist_name_2",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						TrackNumber: 2,
					},
				},
			},
			want: tableData{
				header: []string{"🆔 ID", "🔢 Number", "🎵 Track", "💿 Album", "🎤 Artists", "📅 Release Date"},
				rows: [][]string{
					{"track_id_1", "1", "track_name_1", "album_name_1", "artist_name_1", "2000-01-01"},
					{"track_id_2", "2", "track_name_2", "album_name_2", "artist_name_2", "2000-01-01"},
					{"", "", "", "", "", ""},
					{"TOTAL : 2 tracks!", "", "", "", "", ""},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &TableFormatter{}
			if got := f.formatSearchTracks(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TableFormatter.formatSearchTracks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableFormatter_formatGetTracks(t *testing.T) {
	type args struct {
		items []*spotlikeApp.GetTrackUseCaseOutputDto
	}
	tests := []struct {
		name string
		f    *TableFormatter
		args args
		want tableData
	}{
		{
			name: "positive testing (items is empty)",
			f:    &TableFormatter{},
			args: args{
				items: []*spotlikeApp.GetTrackUseCaseOutputDto{},
			},
			want: tableData{
				header: []string{"🆔 ID", "🔢 Number", "🎵 Track", "💿 Album", "🎤 Artists", "📅 Release Date"},
				rows:   [][]string{},
			},
		},
		{
			name: "positive testing (items is not empty)",
			f:    &TableFormatter{},
			args: args{
				items: []*spotlikeApp.GetTrackUseCaseOutputDto{
					{
						ID:          "track_id_1",
						Name:        "track_name_1",
						Album:       "album_name_1",
						Artists:     "artist_name_1",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						TrackNumber: 1,
					},
					{
						ID:          "track_id_2",
						Name:        "track_name_2",
						Album:       "album_name_2",
						Artists:     "artist_name_2",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						TrackNumber: 2,
					},
				},
			},
			want: tableData{
				header: []string{"🆔 ID", "🔢 Number", "🎵 Track", "💿 Album", "🎤 Artists", "📅 Release Date"},
				rows: [][]string{
					{"track_id_1", "1", "track_name_1", "album_name_1", "artist_name_1", "2000-01-01"},
					{"track_id_2", "2", "track_name_2", "album_name_2", "artist_name_2", "2000-01-01"},
					{"", "", "", "", "", ""},
					{"TOTAL : 2 tracks!", "", "", "", "", ""},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &TableFormatter{}
			if got := f.formatGetTracks(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TableFormatter.formatGetTracks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableFormatter_addTotalRow(t *testing.T) {
	type args struct {
		rows        [][]string
		contentType string
	}
	tests := []struct {
		name string
		f    *TableFormatter
		args args
		want [][]string
	}{
		{
			name: "positive testing (rows is empty)",
			f:    &TableFormatter{},
			args: args{
				rows:        [][]string{},
				contentType: "artists",
			},
			want: [][]string{},
		},
		{
			name: "positive testing (rows is not empty)",
			f:    &TableFormatter{},
			args: args{
				rows: [][]string{
					{"id_1", "name_1"},
					{"id_2", "name_2"},
				},
				contentType: "artists",
			},
			want: [][]string{
				{"id_1", "name_1"},
				{"id_2", "name_2"},
				{"", ""},
				{"TOTAL : 2 artists!", ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &TableFormatter{}
			if got := f.addTotalRow(tt.args.rows, tt.args.contentType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TableFormatter.addTotalRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableFormatter_getTableString(t *testing.T) {
	su := utility.NewStringsUtil()

	type args struct {
		data tableData
	}
	tests := []struct {
		name    string
		f       *TableFormatter
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "positive testing (data is empty)",
			f:    &TableFormatter{},
			args: args{
				data: tableData{},
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "positive testing (data.Header is empty, data.Rows is empty)",
			f:    &TableFormatter{},
			args: args{
				data: tableData{
					header: []string{},
					rows:   [][]string{},
				},
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "positive testing (data.Header is not empty, data.Rows is empty)",
			f:    &TableFormatter{},
			args: args{
				data: tableData{
					header: []string{
						"header_1",
						"header_2",
					},
					rows: [][]string{},
				},
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "positive testing (data.Header is empty, data.Rows is not empty)",
			f:    &TableFormatter{},
			args: args{
				data: tableData{
					header: []string{},
					rows: [][]string{
						{"row_1_1", "row_1_2"},
					},
				},
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "positive testing (data.Header is not empty, data.Rows is not empty)",
			f:    &TableFormatter{},
			args: args{
				data: tableData{
					header: []string{
						"header_1",
						"header_2",
					},
					rows: [][]string{
						{"row_1_1", "row_1_2"},
					},
				},
			},
			want:    "HEADER1HEADER2row_1_1row_1_2",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &TableFormatter{}
			got, err := f.getTableString(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TableFormatter.getTableString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStr := su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(got))); gotStr != tt.want {
				t.Errorf("TableFormatter.getTableString() = %v, want %v", gotStr, tt.want)
			}
		})
	}
}

func TestTableFormatter_getTableString_bulkError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	origTu := getDefaultTableWriterUtil()
	defer setTableWriterUtil(origTu)

	mockTable := proxy.NewMockTable(mockCtrl)
	mockTable.EXPECT().Header(gomock.Any())
	mockTable.EXPECT().Bulk(gomock.Any()).Return(errors.New("bulk error"))

	mockTwu := utility.NewMockTableWriterUtil(mockCtrl)
	mockTwu.EXPECT().GetNewDefaultTable(gomock.Any()).Return(mockTable)

	setTableWriterUtil(mockTwu)

	f := &TableFormatter{}
	got, err := f.getTableString(tableData{
		header: []string{"id"},
		rows:   [][]string{{"1"}},
	})
	if err == nil {
		t.Errorf("TableFormatter.getTableString() error = nil, wantErr true")
	}
	if got != "" {
		t.Errorf("TableFormatter.getTableString() = %v, want empty", got)
	}
}

func TestTableFormatter_getTableString_renderError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	origTu := getDefaultTableWriterUtil()
	defer setTableWriterUtil(origTu)

	mockTable := proxy.NewMockTable(mockCtrl)
	mockTable.EXPECT().Header(gomock.Any())
	mockTable.EXPECT().Bulk(gomock.Any()).Return(nil)
	mockTable.EXPECT().Render().Return(errors.New("render error"))

	mockTwu := utility.NewMockTableWriterUtil(mockCtrl)
	mockTwu.EXPECT().GetNewDefaultTable(gomock.Any()).Return(mockTable)

	setTableWriterUtil(mockTwu)

	f := &TableFormatter{}
	got, err := f.getTableString(tableData{
		header: []string{"id"},
		rows:   [][]string{{"1"}},
	})
	if err == nil {
		t.Errorf("TableFormatter.getTableString() error = nil, wantErr true")
	}
	if got != "" {
		t.Errorf("TableFormatter.getTableString() = %v, want empty", got)
	}
}
