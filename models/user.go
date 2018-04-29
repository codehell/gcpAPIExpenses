package models

import (
	"context"
	"google.golang.org/appengine/datastore"
	"time"
)

type User struct {
	Username string `json:"username,omitempty" datastore:"-"`
	Password string `json:"password,omitempty"`
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

func (user *User) StoreUser(ctx context.Context) error {
	user.CreateAt = time.Now().Unix()
	user.UpdateAt = time.Now().Unix()

	userKey := datastore.NewKey(ctx, "User", user.Username, 0, nil)
	if _, err := datastore.Put(ctx, userKey, user); err != nil {
		return err
	}
	return nil
}


func (user *User) GetUser(ctx context.Context) error {
	q := datastore.NewQuery("User").Filter("__key__=", user.Username)
	t := q.Run(ctx)
	if _, err := t.Next(user); err != nil {
		return err
	}
	return nil
}
