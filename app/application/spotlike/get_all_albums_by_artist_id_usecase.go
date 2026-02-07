package spotlike

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/zmb3/spotify/v2"

	albumDomain "github.com/yanosea/spotlike/app/domain/spotify/album"
)

// getAllAlbumsByArtistIdUseCase is a struct that contains the use case of getting for an artist.
type getAllAlbumsByArtistIdUseCase struct {
	albumRepo albumDomain.AlbumRepository
}

// NewGetAllAlbumsByArtistIdUseCase returns a new instance of the GetAllAlbumsByArtistIdUseCase struct.
func NewGetAllAlbumsByArtistIdUseCase(albumRepo albumDomain.AlbumRepository) *getAllAlbumsByArtistIdUseCase {
	return &getAllAlbumsByArtistIdUseCase{
		albumRepo: albumRepo,
	}
}

// GetAllAlbumsByArtistIdUseCaseOutputDto is a DTO struct that contains the output data of the getAllAlbumsUseCase.
type GetAllAlbumsByArtistIdUseCaseOutputDto struct {
	ID          string
	Artists     string
	Name        string
	ReleaseDate time.Time
}

// Run returns the get result of the albums.
func (uc *getAllAlbumsByArtistIdUseCase) Run(ctx context.Context, id string) ([]*GetAllAlbumsByArtistIdUseCaseOutputDto, error) {
	albums, err := uc.albumRepo.FindByArtistId(ctx, spotify.ID(id))
	if err != nil {
		return nil, err
	}
	sort.Slice(albums, func(i, j int) bool {
		return albums[i].ReleaseDate.Before(albums[j].ReleaseDate)
	})

	var getAllAlbumsByArtistIdUseCaseOutputDtos []*GetAllAlbumsByArtistIdUseCaseOutputDto
	for _, album := range albums {
		artistNames := make([]string, len(album.Artists))
		for i, artist := range album.Artists {
			artistNames[i] = artist.Name
		}

		getAllAlbumsUseCaseOutputDto := &GetAllAlbumsByArtistIdUseCaseOutputDto{
			ID:          album.ID.String(),
			Artists:     strings.Join(artistNames, ", "),
			Name:        album.Name,
			ReleaseDate: album.ReleaseDate,
		}
		getAllAlbumsByArtistIdUseCaseOutputDtos = append(getAllAlbumsByArtistIdUseCaseOutputDtos, getAllAlbumsUseCaseOutputDto)
	}

	return getAllAlbumsByArtistIdUseCaseOutputDtos, nil
}
