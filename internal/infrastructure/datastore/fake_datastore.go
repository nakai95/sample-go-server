package datastore

import "sample-go-server/internal/adapter/repository"

type fakeDataStore struct {
}

func NewFakeDataStore() repository.DataStore {
	return &fakeDataStore{}
}
