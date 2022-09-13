package kvs

import (
	"context"
	"errors"
)

var ErrEmptyObject = errors.New("empty object")
var ErrEmptyID = errors.New("empty ID")

type KvSID string

type KvsStorage interface {
	GetAll() map[string]interface{}
	Get(ctx context.Context, key KvSID) (interface{}, error)
	Save(ctx context.Context, key KvSID, newObject interface{}) error
}

type Kvs struct {
	connStr string
	objects map[KvSID]interface{}
}

func NewKvs(connStr string) *Kvs {
	return &Kvs{
		connStr: connStr,
		objects: map[KvSID]interface{}{},
	}
}

func (kvs *Kvs) GetAll() (map[KvSID]interface{}, error) {
	if kvs.connStr != "dummy" {
		return nil, errors.New("network required")
	}
	return kvs.objects, nil
}

func (kvs *Kvs) Get(ctx context.Context, key KvSID) (interface{}, error) {
	if kvs.connStr != "dummy" {
		return nil, errors.New("network required")
	}

	if key == "" {
		return nil, ErrEmptyID
	}

	if val, ok := kvs.objects[key]; ok {
		return val, nil
	}

	return nil, nil
}

func (kvs *Kvs) Save(ctx context.Context, key KvSID, newObject interface{}) error {
	if kvs.connStr != "dummy" {
		return errors.New("network required")
	}
	if newObject == nil {
		return ErrEmptyObject
	}

	if key == "" {
		return ErrEmptyID
	}

	kvs.objects[key] = newObject

	return nil
}
