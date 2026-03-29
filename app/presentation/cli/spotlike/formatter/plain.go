package formatter

import (
	"fmt"

	spotlikeApp "github.com/yanosea/spotlike/app/application/spotlike"
)

// PlainFormatter is a struct that formats the output of spotlike cli.
type PlainFormatter struct{}

// NewPlainFormatter returns a new instance of the PlainFormatter struct.
func NewPlainFormatter() *PlainFormatter {
	return &PlainFormatter{}
}

// Format formats the output of spotlike cli.
func (f *PlainFormatter) Format(result any) (string, error) {
	var formatted string
	switch v := result.(type) {
	case *spotlikeApp.GetVersionUseCaseOutputDto:
		formatted = fmt.Sprintf("spotlike version %s", v.Version)
	case []*spotlikeApp.SearchArtistUseCaseOutputDto:
		for i, item := range v {
			formatted += "[" + item.ID + "]" + " Artist : " + item.Name
			if i < len(v)-1 {
				formatted += "\n"
			}
		}
	case []*spotlikeApp.GetArtistUseCaseOutputDto:
		for i, item := range v {
			formatted += "[" + item.ID + "]" + " Artist : " + item.Name
			if i < len(v)-1 {
				formatted += "\n"
			}
		}
	case []*spotlikeApp.SearchAlbumUseCaseOutputDto:
		for i, item := range v {
			formatted += "[" + item.ID + "]" + " Album : " + item.Name + " released at " + item.ReleaseDate.Format("2006-01-02") + " by " + item.Artists
			if i < len(v)-1 {
				formatted += "\n"
			}
		}
	case []*spotlikeApp.GetAlbumUseCaseOutputDto:
		for i, item := range v {
			formatted += "[" + item.ID + "]" + " Album : " + item.Name + " released at " + item.ReleaseDate.Format("2006-01-02") + " by " + item.Artists
			if i < len(v)-1 {
				formatted += "\n"
			}
		}
	case []*spotlikeApp.SearchTrackUseCaseOutputDto:
		for i, item := range v {
			formatted += "[" + item.ID + "]" + " Track : #" + fmt.Sprint(item.TrackNumber) + " " + item.Name + " on " + item.Album + " released at " + item.ReleaseDate.Format("2006-01-02") + " by " + item.Artists
			if i < len(v)-1 {
				formatted += "\n"
			}
		}
	case []*spotlikeApp.GetTrackUseCaseOutputDto:
		for i, item := range v {
			formatted += "[" + item.ID + "]" + " Track : #" + fmt.Sprint(item.TrackNumber) + " " + item.Name + " on " + item.Album + " released at " + item.ReleaseDate.Format("2006-01-02") + " by " + item.Artists
			if i < len(v)-1 {
				formatted += "\n"
			}
		}
	default:
		formatted = ""
	}
	return formatted, nil
}
