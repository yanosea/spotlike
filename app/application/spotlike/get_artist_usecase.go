package spotlike

import (
	"context"

	"github.com/zmb3/spotify/v2"

	artistDomain "github.com/yanosea/spotlike/app/domain/spotify/artist"
)

// GetArtistUseCase is an interface that defines the use case of getting for an artist.
type GetArtistUseCase interface {
	Run(ctx context.Context, id string) (*GetArtistUseCaseOutputDto, error)
}

// GetArtistUseCaseStruct is a struct that implements the GetArtistUseCase interface.
type GetArtistUseCaseStruct struct {
	artistRepo artistDomain.ArtistRepository
}

var (
	// NewGetArtistUseCase is a function that returns a new instance of the GetArtistUseCaseStruct struct.
	NewGetArtistUseCase = newGetArtistUseCase
)

// NewGetArtistUseCase returns a new instance of the GetArtistUseCase struct.
func newGetArtistUseCase(artistRepo artistDomain.ArtistRepository) *GetArtistUseCaseStruct {
	return &GetArtistUseCaseStruct{
		artistRepo: artistRepo,
	}
}

// GetArtistUseCaseOutputDto is a DTO struct that contains the output data of the getArtistUseCase.
type GetArtistUseCaseOutputDto struct {
	ID   string
	Name string
}

// Run returns the get result of the artist.
func (uc *GetArtistUseCaseStruct) Run(ctx context.Context, id string) (*GetArtistUseCaseOutputDto, error) {
	artist, err := uc.artistRepo.FindById(ctx, spotify.ID(id))
	if err != nil {
		return nil, err
	}

	return &GetArtistUseCaseOutputDto{
		ID:   artist.ID.String(),
		Name: artist.Name,
	}, nil
}
