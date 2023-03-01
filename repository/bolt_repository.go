package repository

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"secret_keeper/pb"

	"github.com/boltdb/bolt"
)

const (
	usersBucket = "users"
)

var (
	ErrExists    = errors.New("already exists")
	ErrNotExists = errors.New("not found")
)

type BoltRepository struct {
	db *bolt.DB
}

type User struct {
	*pb.User
	Password string `json:"password"`
}

func (br *BoltRepository) CreateUser(user *User) error {
	return br.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(usersBucket))

		res := bucket.Get([]byte(user.Email))
		if len(res) != 0 {
			return ErrExists
		}

		encoded, err := json.Marshal(user)
		if err != nil {
			return fmt.Errorf("encoding error %s", err)
		}

		return bucket.Put([]byte(user.Email), encoded)
	})
}

func (br *BoltRepository) GetUser(email string) (*User, error) {
	user := &User{}

	err := br.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(usersBucket))

		res := bucket.Get([]byte(email))
		if len(res) == 0 {
			return ErrNotExists
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

func (br *BoltRepository) CreateSecret(secret *pb.Secret, email string) error {
	return br.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(email))

		if err != nil {
			return fmt.Errorf("error creating bucket %s", err)
		}

		secretId, _ := bucket.NextSequence()
		secret.Id = int64(secretId)

		encoded, err := json.Marshal(secret)
		if err != nil {
			return fmt.Errorf("encoding error %s", err)
		}

		return bucket.Put(itob(secretId), encoded)
	})
}

func (br *BoltRepository) SecretsList(email string) (secrets []*pb.Secret, err error) {
	err = br.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(email))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			secret := &pb.Secret{}
			err := json.Unmarshal(v, secret)
			if err != nil {
				return err
			}
			secrets = append(secrets, secret)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return secrets, nil
}

func (br *BoltRepository) GetSecret(id uint64, email string) (*pb.Secret, error) {
	secret := &pb.Secret{}

	err := br.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(email))
		if bucket == nil {
			return ErrNotExists
		}

		res := bucket.Get(itob(id))
		if len(res) == 0 {
			return ErrNotExists
		}

		return json.Unmarshal(res, secret)
	})

	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (br *BoltRepository) DeleteSecret(id uint64, email string) error {
	return br.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(email))

		bytesId := itob(id)
		secret := bucket.Get(itob(id))

		if len(secret) == 0 {
			return nil
		}

		return bucket.Delete(bytesId)
	})
}

func (br *BoltRepository) DeleteAllSecrets(email string) error {
	return br.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(email))
	})
}

func NewBoltRepository(db *bolt.DB) *BoltRepository {
	return &BoltRepository{
		db: db,
	}
}

func SetupDb(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(usersBucket))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	})
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func btoi(v []byte) uint64 {
	return binary.BigEndian.Uint64(v)
}
