package database

import (
	"github.com/solher/arangolite"
	"github.com/solher/snakepit"
)

type Manager struct {
	db *snakepit.ArangoDBManager
}

func NewManager(db *snakepit.ArangoDBManager) *Manager {
	return &Manager{db: db}
}

func (d *Manager) Create(rootName, rootPassword string) error {
	if err := d.db.Create(rootName, rootPassword); err != nil {
		return err
	}

	if _, err := d.db.Run(&arangolite.SetCacheProperties{Mode: "demand"}); err != nil {
		return err
	}

	return nil
}

func (d *Manager) Migrate() error {
	if err := d.db.Migrate(); err != nil {
		return err
	}

	return nil
}

func (d *Manager) Drop(rootName, rootPassword string) error {
	if err := d.db.Drop(rootName, rootPassword); err != nil {
		return err
	}

	return nil
}

func (d *Manager) SyncSeeds() error {
	if err := d.db.SyncSeeds(); err != nil {
		return err
	}

	return nil
}
