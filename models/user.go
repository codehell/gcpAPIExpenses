package models

import (
	"context"
	"google.golang.org/appengine/datastore"
	"time"
)

type User struct {
	Username string `json:"username,omitempty" datastore:"username,omitempty"`
	CreateAt int64  `json:"created_at"`
	UpdateAt int64  `json:"updated_at"`
}

func GetUsers(ctx context.Context) ([]User, error) {
	var users []User
	q := datastore.NewQuery("User")
	keys, err := q.GetAll(ctx, &users)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		users[i].Username = key.StringID()
	}
	return users, nil
}

func StoreUser(ctx context.Context) error {
	user := User{
		CreateAt: time.Now().Unix(),
		UpdateAt: time.Now().Unix(),
	}

	userKey := datastore.NewKey(ctx, "User", "codehell", 0, nil)

	if _, err := datastore.Put(ctx, userKey, &user); err != nil {
		return err
	}
	return nil
}
