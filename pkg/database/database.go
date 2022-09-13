package database

import (
	"context"
	"errors"

	"github.com/mercadolibre/fury_cx-example/pkg/utils"
)

var ErrEmptyObject = errors.New("empty object")
var ErrEmptyID = errors.New("empty ID")

type DatabaseObject interface {
	GetID() string
	SetID(string)
}

type Database struct {
	objects []DatabaseObject
}

func (db *Database) GetAll() []DatabaseObject {
	return db.objects
}

func (db *Database) GetObject(ctx context.Context, ID string) (DatabaseObject, error) {
	if ID == "" {
		return nil, ErrEmptyID
	}

	for _, v := range db.objects {
		if v.GetID() == ID {
			return v, nil
		}
	}

	return nil, nil
}

func (db *Database) CreateObject(ctx context.Context, newObject DatabaseObject) error {
	if newObject == nil {
		return ErrEmptyObject
	}

	id := utils.RandStringBytes(10)

	newObject.SetID(id)

	db.objects = append(db.objects, newObject)

	return nil
}
