package storage

import "github.com/RaghavSood/collectibles/types"

type Storage interface {
	GetCreators() ([]types.Creator, error)
}
