package storage

import (
	"path/filepath"
)

type Entity struct {
	DiskName      string
	PartitionName string
	Name          string
	EntityData
}

func (entity Entity) Location() string {
	return filepath.Join(entity.DiskName, entity.PartitionName, entity.Name)
}
