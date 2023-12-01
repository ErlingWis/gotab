package storage

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type Partition struct {
	Name     string
	DiskName string
	Logger   *slog.Logger
}

func (p Partition) Location() string {
	return filepath.Join(p.DiskName, p.Name)
}

func (p Partition) Delete() error {
	p.Logger.Info("Deleting partition", "partition", p)
	return os.RemoveAll(p.Location())
}

func (p Partition) ListEntities() ([]string, error) {
	p.Logger.Info("Listing entities", "partition", p)
	entries, err := os.ReadDir(p.Location())

	if err != nil {
		return nil, err
	}

	entities := make([]string, 0)

	for _, entry := range entries {
		entities = append(entities, entry.Name())
	}

	return entities, nil
}

func (p Partition) EnsureCreated() error {
	if _, err := os.Stat(p.Location()); err != nil {

		if !os.IsNotExist(err) {
			return err
		}

		if err := os.Mkdir(p.Location(), 0755); err != nil {
			return err
		}

	}
	return nil
}

func (p Partition) UpdateEntity(entity Entity) error {
	p.Logger.Info("Updating entity", "partition", p.Name, "entity", entity.Name)

	if _, err := os.Stat(entity.Location()); os.IsNotExist(err) {
		return p.CreateEntity(entity)
	} else if err != nil {
		return err
	}

	existing, err := p.GetEntity(entity.Name)

	if err != nil {
		return err
	}
	p.Logger.Info("Asserting etags", "existing", existing.ETag, "entity", entity.ETag)
	if existing.EntityData.ETag != entity.EntityData.ETag {
		return errors.New("etag mismatch. Pull new verison of entity before updating")
	}

	p.CreateEntity(entity)

	return nil
}

func (p Partition) CreateEntity(entity Entity) error {
	p.Logger.Info("Creating entity", "partition", p.Name, "entity", entity.Name)

	if err := p.EnsureCreated(); err != nil {
		return err
	}

	entity.EntityData.LastModified = time.Now()
	entity.EntityData.ETag = time.Now().UnixNano()
	encoding, err := entity.encode()

	if err != nil {
		return err
	}

	return os.WriteFile(entity.Location(), encoding, 0755)
}

func (p Partition) DeleteEntity(e Entity) error {
	p.Logger.Info("Deleting entity", "partition", p, "entity", e)
	return os.Remove(e.Location())
}

func (p Partition) GetEntity(name string) (Entity, error) {
	p.Logger.Info("Getting entity", "partition", p.Name, "entity", name)

	location := filepath.Join(p.Location(), name)

	file, err := os.Open(location)

	if err != nil {
		return Entity{}, err
	}

	defer file.Close()

	data, err := decodeData(file)

	if err != nil {
		return Entity{}, err
	}

	return Entity{
		DiskName:      p.DiskName,
		PartitionName: p.Name,
		Name:          name,
		EntityData:    data,
	}, nil
}
