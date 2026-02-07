package spotlike

import (
	"context"
	"strings"
	"time"

	"github.com/zmb3/spotify/v2"

	trackDomain "github.com/yanosea/spotlike/app/domain/spotify/track"
)

// searchTrackUseCase is a struct that contains the use case of searching for an track.
type SearchTrackUseCase interface {
	Run(ctx context.Context, keywords []string, max int) ([]*SearchTrackUseCaseOutputDto, error)
}

// SearchTrackUseCaseStruct is a struct that contains the use case of searching for an track.
type SearchTrackUseCaseStruct struct {
	trackRepo trackDomain.TrackRepository
}

var (
	// NewSearchTrackUseCase is a function that returns a new instance of the SearchTrackUseCase struct.
	NewSearchTrackUseCase = newSearchTrackUseCase
)

// NewSearchTrackUseCaseStruct returns a new instance of the SearchTrackUseCase struct.
func newSearchTrackUseCase(trackRepo trackDomain.TrackRepository) *SearchTrackUseCaseStruct {
	return &SearchTrackUseCaseStruct{
		trackRepo: trackRepo,
	}
}

// SearchTrackUseCaseOutputDto is a DTO struct that contains the output data of the SearchTrackUseCase.
type SearchTrackUseCaseOutputDto struct {
	ID          string
	Artists     string
	Album       string
	Name        string
	TrackNumber spotify.Numeric
	ReleaseDate time.Time
}

// Run returns the search result of the track.
func (uc *SearchTrackUseCaseStruct) Run(ctx context.Context, keywords []string, max int) ([]*SearchTrackUseCaseOutputDto, error) {
	query := strings.Join(keywords, " ")

	tracks, err := uc.trackRepo.FindByNameLimit(ctx, query, max)
	if err != nil {
		return nil, err
	}

	var searchTrackResultDtos []*SearchTrackUseCaseOutputDto
	for _, track := range tracks {
		artistNames := make([]string, len(track.Artists))
		for i, artist := range track.Artists {
			artistNames[i] = artist.Name
		}

		searchTrackResultDtos = append(searchTrackResultDtos, &SearchTrackUseCaseOutputDto{
			ID:          track.ID.String(),
			Artists:     strings.Join(artistNames, ", "),
			Album:       track.Album.Name,
			Name:        track.Name,
			TrackNumber: track.TrackNumber,
			ReleaseDate: track.Album.ReleaseDateTime(),
		})
	}

	return searchTrackResultDtos, nil
}
