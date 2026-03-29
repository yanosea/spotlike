package spotlike

import (
	"context"
	"strings"
	"time"

	"github.com/zmb3/spotify/v2"

	albumDomain "github.com/yanosea/spotlike/app/domain/spotify/album"
)

// getAlbumUseCase is a struct that contains the use case of getting for an album.
type getAlbumUseCase struct {
	albumRepo albumDomain.AlbumRepository
}

// NewGetAlbumUseCase returns a new instance of the GetAlbumUseCase struct.
func NewGetAlbumUseCase(albumRepo albumDomain.AlbumRepository) *getAlbumUseCase {
	return &getAlbumUseCase{
		albumRepo: albumRepo,
	}
}

// GetAlbumUseCaseOutputDto is a DTO struct that contains the output data of the getAlbumUseCase.
type GetAlbumUseCaseOutputDto struct {
	ID          string
	Name        string
	Artists     string
	ReleaseDate time.Time
}

// Run returns the get result of the album.
func (uc *getAlbumUseCase) Run(ctx context.Context, id string) (*GetAlbumUseCaseOutputDto, error) {
	album, err := uc.albumRepo.FindById(ctx, spotify.ID(id))
	if err != nil {
		return nil, err
	}

	artistNames := make([]string, len(album.Artists))
	for i, artist := range album.Artists {
		artistNames[i] = artist.Name
	}

	return &GetAlbumUseCaseOutputDto{
		ID:          album.ID.String(),
		Name:        album.Name,
		Artists:     strings.Join(artistNames, ", "),
		ReleaseDate: album.ReleaseDate,
	}, nil
}
