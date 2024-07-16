package main

import "context"

type orderStore struct {
}

func NewOrderStore() *orderStore {
	return &orderStore{}
}

func (s *orderStore) Create(context.Context) error {
	return nil
}
