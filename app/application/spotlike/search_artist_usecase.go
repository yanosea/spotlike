package spotlike

import (
	"context"
	"strings"

	artistDomain "github.com/yanosea/spotlike/app/domain/spotify/artist"
)

// SearchArtistUseCase is a struct that contains the use case of searching for an artist.
type SearchArtistUseCase interface {
	Run(ctx context.Context, keywords []string, max int) ([]*SearchArtistUseCaseOutputDto, error)
}

// SearchArtistUseCaseStruct is a struct that contains the use case of searching for an artist.
type SearchArtistUseCaseStruct struct {
	artistRepo artistDomain.ArtistRepository
}

var (
	// NewSearchArtistUseCase is a function that returns a new instance of the SearchArtistUseCaseStruct struct.
	NewSearchArtistUseCase = newSearchArtistUseCase
)

// newSearchArtistUseCase returns a new instance of the SearchArtistUseCase struct.
func newSearchArtistUseCase(artistRepo artistDomain.ArtistRepository) *SearchArtistUseCaseStruct {
	return &SearchArtistUseCaseStruct{
		artistRepo: artistRepo,
	}
}

// SearchArtistUseCaseOutputDto is a DTO struct that contains the output data of the SearchArtistUseCase.
type SearchArtistUseCaseOutputDto struct {
	ID   string
	Name string
}

// Run returns the search result of the artist.
func (uc *SearchArtistUseCaseStruct) Run(ctx context.Context, keywords []string, max int) ([]*SearchArtistUseCaseOutputDto, error) {
	query := strings.Join(keywords, " ")

	artists, err := uc.artistRepo.FindByNameLimit(ctx, query, max)
	if err != nil {
		return nil, err
	}

	var searchArtistResultDtos []*SearchArtistUseCaseOutputDto
	for _, artist := range artists {
		searchArtistResultDtos = append(searchArtistResultDtos, &SearchArtistUseCaseOutputDto{
			ID:   artist.ID.String(),
			Name: artist.Name,
		})
	}

	return searchArtistResultDtos, nil
}
