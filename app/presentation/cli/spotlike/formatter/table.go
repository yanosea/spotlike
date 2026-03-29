package formatter

import (
	"fmt"
	"strings"

	spotlikeApp "github.com/yanosea/spotlike/app/application/spotlike"

	"github.com/yanosea/spotlike/pkg/proxy"
	"github.com/yanosea/spotlike/pkg/utility"
)

// TableFormatter is a struct that formats the output of spotlike cli.
type TableFormatter struct{}

// NewTableFormatter returns a new instance of the TableFormatter struct.
func NewTableFormatter() *TableFormatter {
	return &TableFormatter{}
}

var (
	// Tu is a variable to store the table writer with the default values for injecting the dependencies in testing.
	Tu = utility.NewTableWriterUtil(proxy.NewTableWriter())
)

// tableData is a struct that holds the data of a table.
type tableData struct {
	header []string
	rows   [][]string
}

// Format formats the output of spotlike cli.
func (f *TableFormatter) Format(result any) (string, error) {
	var data tableData
	switch v := result.(type) {
	case []*spotlikeApp.SearchArtistUseCaseOutputDto:
		data = f.formatSearchArtists(v)
	case []*spotlikeApp.GetArtistUseCaseOutputDto:
		data = f.formatGetArtists(v)
	case []*spotlikeApp.SearchAlbumUseCaseOutputDto:
		data = f.formatSearchAlbums(v)
	case []*spotlikeApp.GetAlbumUseCaseOutputDto:
		data = f.formatGetAlbums(v)
	case []*spotlikeApp.SearchTrackUseCaseOutputDto:
		data = f.formatSearchTracks(v)
	case []*spotlikeApp.GetTrackUseCaseOutputDto:
		data = f.formatGetTracks(v)
	default:
		return "", nil
	}

	return f.getTableString(data)
}

// formatSearchArtists formats the output of the search artists use case.
func (f *TableFormatter) formatSearchArtists(items []*spotlikeApp.SearchArtistUseCaseOutputDto) tableData {
	header := []string{"🆔 ID", "🎤 Artist"}
	var rows [][]string
	for _, artist := range items {
		rows = append(rows, []string{
			artist.ID,
			artist.Name,
		})
	}
	rows = f.addTotalRow(rows, "artists")

	return tableData{header: header, rows: rows}
}

// formatGetArtists formats the output of the get artists use case.
func (f *TableFormatter) formatGetArtists(items []*spotlikeApp.GetArtistUseCaseOutputDto) tableData {
	var header []string
	var rows [][]string
	header = []string{"🆔 ID", "🎤 Artist"}
	for _, item := range items {
		rows = append(rows, []string{
			item.ID,
			item.Name,
		})
	}
	rows = f.addTotalRow(rows, "artists")

	return tableData{header: header, rows: rows}
}

// formatSearchAlbums formats the output of the search albums use case.
func (f *TableFormatter) formatSearchAlbums(items []*spotlikeApp.SearchAlbumUseCaseOutputDto) tableData {
	header := []string{"🆔 ID", "💿 Album", "🎤 Artists", "📅 Release Date"}
	var rows [][]string
	for _, album := range items {
		rows = append(rows, []string{
			album.ID,
			album.Name,
			album.Artists,
			album.ReleaseDate.Format("2006-01-02"),
		})
	}
	rows = f.addTotalRow(rows, "albums")

	return tableData{header: header, rows: rows}
}

// formatGetAlbums formats the output of the get albums use case.
func (f *TableFormatter) formatGetAlbums(items []*spotlikeApp.GetAlbumUseCaseOutputDto) tableData {
	var header []string
	var rows [][]string
	header = []string{"🆔 ID", "💿 Album", "🎤 Artists", "📅 Release Date"}
	for _, item := range items {
		rows = append(rows, []string{
			item.ID,
			item.Name,
			item.Artists,
			item.ReleaseDate.Format("2006-01-02"),
		})
	}
	rows = f.addTotalRow(rows, "albums")

	return tableData{header: header, rows: rows}
}

// formatSearchTracks formats the output of the search tracks use case.
func (f *TableFormatter) formatSearchTracks(items []*spotlikeApp.SearchTrackUseCaseOutputDto) tableData {
	header := []string{"🆔 ID", "🔢 Number", "🎵 Track", "💿 Album", "🎤 Artists", "📅 Release Date"}
	var rows [][]string
	for _, track := range items {
		rows = append(rows, []string{
			track.ID,
			fmt.Sprint(track.TrackNumber),
			track.Name,
			track.Album,
			track.Artists,
			track.ReleaseDate.Format("2006-01-02"),
		})
	}
	rows = f.addTotalRow(rows, "tracks")

	return tableData{header: header, rows: rows}
}

// formatGetTracks formats the output of the get tracks use case.
func (f *TableFormatter) formatGetTracks(items []*spotlikeApp.GetTrackUseCaseOutputDto) tableData {
	var header []string
	var rows [][]string
	header = []string{"🆔 ID", "🔢 Number", "🎵 Track", "💿 Album", "🎤 Artists", "📅 Release Date"}
	for _, item := range items {
		rows = append(rows, []string{
			item.ID,
			fmt.Sprint(item.TrackNumber),
			item.Name,
			item.Album,
			item.Artists,
			item.ReleaseDate.Format("2006-01-02"),
		})
	}
	rows = f.addTotalRow(rows, "tracks")

	return tableData{header: header, rows: rows}
}

// addTotalRow adds a total row to the table.
func (f *TableFormatter) addTotalRow(rows [][]string, contentType string) [][]string {
	if len(rows) == 0 {
		return [][]string{}
	}

	emptyRow := make([]string, len(rows[0]))
	for i := range emptyRow {
		emptyRow[i] = ""
	}

	totalRow := make([]string, len(rows[0]))
	totalRow[0] = fmt.Sprintf("TOTAL : %d %s!", len(rows), contentType)

	rows = append(rows, emptyRow)
	rows = append(rows, totalRow)

	return rows
}

// getTableString returns a string representation of a table.
func (f *TableFormatter) getTableString(data tableData) (string, error) {
	if len(data.header) == 0 || len(data.rows) == 0 {
		return "", nil
	}

	tableString := &strings.Builder{}
	table := Tu.GetNewDefaultTable(tableString)
	table.Header(data.header)
	if err := table.Bulk(data.rows); err != nil {
		return "", err
	}
	if err := table.Render(); err != nil {
		return "", err
	}
	return strings.TrimSuffix(tableString.String(), "\n"), nil
}
