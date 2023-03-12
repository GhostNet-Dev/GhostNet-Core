// reference https://github.com/zupzup/boltdb-example/blob/main/main.go
package store

import (
	"fmt"

	"github.com/boltdb/bolt"
)

type LiteStore struct {
	db     *bolt.DB
	tables []string
	dbPath string
}

func NewLiteStore(dbPath string, tables []string) *LiteStore {
	return &LiteStore{
		dbPath: dbPath,
		tables: tables,
	}
}

func (liteStore *LiteStore) OpenStore() error {
	db, err := bolt.Open(liteStore.dbPath+"litestore.db", 0600, nil)
	if err != nil {
		return fmt.Errorf("could not open bolt db = %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("litestore"))
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}
		for _, table := range liteStore.tables {
			_, err = root.CreateBucketIfNotExists([]byte(table))
			if err != nil {
				return fmt.Errorf("could not create %v bucket: %v", table, err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("could not set up buckets, %v", err)
	}
	liteStore.db = db
	return nil
}

func (liteStore *LiteStore) Close() error {
	return liteStore.db.Close()
}

// 대량으로 선택한다.
func (liteStore *LiteStore) LoadEntry(table string) (key, result [][]byte, err error) {
	err = liteStore.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("litestore")).Bucket([]byte(table))
		b.ForEach(func(k, v []byte) error {
			key = append(key, k)
			result = append(result, v)
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not set up buckets, %v", err)
	}
	return key, result, nil
}

// 하나만 선택한다.
func (liteStore *LiteStore) SelectEntry(table string, key []byte) (result []byte, err error) {
	err = liteStore.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("litestore")).Bucket([]byte(table))
		result = b.Get(key)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not set up buckets, %v", err)
	}
	return result, nil
}

func (liteStore *LiteStore) SaveEntry(table string, k, v []byte) error {
	err := liteStore.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("litestore")).Bucket([]byte(table)).Put(k, v)
		return err
	})
	return err
}
