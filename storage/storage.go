package storage

import "github.com/RaghavSood/collectibles/types"

type Storage interface {
	GetCreators() ([]types.Creator, error)
	GetCreator(slug string) (*types.Creator, error)

	GetSeries() ([]types.Series, error)
}
