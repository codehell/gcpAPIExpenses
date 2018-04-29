package models

import (
	"google.golang.org/appengine/datastore"
	"context"
)

type Movement struct {
	ID          int64    `json:"id,omitempty" datastore:"-"`
	Username    string   `json:"username"`
	Amount      int32    `json:"amount"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	CreateAt    int64    `json:"created_at"`
	UpdateAt    int64    `json:"updated_at"`
}

func GetMovements(ctx context.Context) ([]Movement, error) {
	var expenses []Movement
	q := datastore.NewQuery("Movement")
	keys, err := q.GetAll(ctx, &expenses)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		expenses[i].ID = key.IntID()
	}

	return expenses, nil
}

func StoreMovement(ctx context.Context, expense *Movement) error {
	key := datastore.NewIncompleteKey(ctx, "Movement", nil)
	_, err := datastore.Put(ctx, key, expense)
	if err != nil {
		return err
	}
	return nil
}
