package spotlike

import (
	"context"
	"strings"
	"time"

	albumDomain "github.com/yanosea/spotlike/app/domain/spotify/album"
)

// SearchAlbumUseCase is a struct that contains the use case of searching for an album.
type SearchAlbumUseCase interface {
	Run(ctx context.Context, keywords []string, max int) ([]*SearchAlbumUseCaseOutputDto, error)
}

// SearchAlbumUseCaseStruct is a struct that contains the use case of searching for an album.
type SearchAlbumUseCaseStruct struct {
	albumRepo albumDomain.AlbumRepository
}

var (
	// NewSearchAlbumUseCase is a function that returns a new instance of the SearchAlbumUseCase struct.
	NewSearchAlbumUseCase = newSearchAlbumUseCase
)

// newSearchAlbumUseCase returns a new instance of the SearchAlbumUseCase struct.
func newSearchAlbumUseCase(albumRepo albumDomain.AlbumRepository) *SearchAlbumUseCaseStruct {
	return &SearchAlbumUseCaseStruct{
		albumRepo: albumRepo,
	}
}

// SearchAlbumUseCaseOutputDto is a DTO struct that contains the output data of the SearchAlbumUseCase.
type SearchAlbumUseCaseOutputDto struct {
	ID          string
	Artists     string
	Name        string
	ReleaseDate time.Time
}

// Run returns the search result of the album.
func (uc *SearchAlbumUseCaseStruct) Run(ctx context.Context, keywords []string, max int) ([]*SearchAlbumUseCaseOutputDto, error) {
	query := strings.Join(keywords, " ")

	albums, err := uc.albumRepo.FindByNameLimit(ctx, query, max)
	if err != nil {
		return nil, err
	}

	var searchAlbumResultDtos []*SearchAlbumUseCaseOutputDto
	for _, album := range albums {
		artistNames := make([]string, len(album.Artists))
		for i, artist := range album.Artists {
			artistNames[i] = artist.Name
		}

		searchAlbumResultDtos = append(searchAlbumResultDtos, &SearchAlbumUseCaseOutputDto{
			ID:          album.ID.String(),
			Artists:     strings.Join(artistNames, ", "),
			Name:        album.Name,
			ReleaseDate: album.ReleaseDate,
		})
	}

	return searchAlbumResultDtos, nil
}
