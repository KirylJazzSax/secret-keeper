package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/KirylJazzSax/secret-keeper/internal/common/errors"
	"github.com/KirylJazzSax/secret-keeper/internal/user/domain"
	"github.com/boltdb/bolt"
)

const (
	UsersBucket = "users"
)

type BoltUserRepository struct {
	db *bolt.DB
}

func (r *BoltUserRepository) CreateUser(ctx context.Context, u *domain.User) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(UsersBucket))

		res := bucket.Get([]byte(u.Email))
		if len(res) != 0 {
			return errors.ErrExists
		}

		encoded, err := json.Marshal(u)
		if err != nil {
			return fmt.Errorf("encoding error %s", err)
		}

		return bucket.Put([]byte(u.Email), encoded)
	})
}

func (r *BoltUserRepository) GetUser(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}

	err := r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(UsersBucket))

		res := bucket.Get([]byte(email))
		if len(res) == 0 {
			return errors.ErrNotExists
		}

		err := json.Unmarshal(res, user)
		if err != nil {
			return fmt.Errorf("decoding error %s", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewBoltUserRepository(db *bolt.DB) *BoltUserRepository {
	return &BoltUserRepository{
		db: db,
	}
}
