package storage

import (
	"errors"
	"log/slog"
	"os"
)

type Disk struct {
	Name   string
	Logger *slog.Logger
}

func (s Disk) Location() string {
	return s.Name
}

func (s Disk) Validate() error {

	s.Logger.Info("Validating disk", "disk", s)

	if stat, err := os.Stat(s.Name); os.IsNotExist(err) {
		return err
	} else if !stat.IsDir() {
		return errors.New(s.Name + " is not a directory")
	}
	s.Logger.Info("Disk is valid", "disk", s)
	return nil
}

func (s Disk) Delete() error {
	s.Logger.Info("Deleting disk", "disk", s)
	return os.RemoveAll(s.Name)
}

func (s Disk) GetPartition(name string) Partition {
	s.Logger.Info("Getting partition", "disk", s, "partition", name)
	return Partition{DiskName: s.Location(), Name: name, Logger: s.Logger}
}
func (s Disk) ListPartitionNames() (partitions []string, err error) {
	s.Logger.Info("Listing partition names", "disk", s)

	entires, err := os.ReadDir(s.Name)

	if err != nil {
		return nil, err
	}

	for _, entry := range entires {
		if !entry.IsDir() {
			continue
		}

		partitions = append(partitions, entry.Name())
	}

	return partitions, nil
}
func (s Disk) ListPartitions() ([]Partition, error) {
	s.Logger.Info("Listing partitions", "disk", s)

	entires, err := s.ListPartitionNames()

	if err != nil {
		return nil, err
	}

	partitions := make([]Partition, 0)

	for _, entry := range entires {

		partitions = append(partitions, s.GetPartition(entry))
	}

	return partitions, nil
}
