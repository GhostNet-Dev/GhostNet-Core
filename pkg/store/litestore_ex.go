package store

import (
	"fmt"

	"github.com/boltdb/bolt"
)

// 대량으로 선택한다.
func (liteStore *LiteStore) LoadEntryEx(parentTable, table []byte) (key, result [][]byte, err error) {
	err = liteStore.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(parentTable).Bucket(table)
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

// 대량으로 선택한다.
func (liteStore *LiteStore) LoadEntryDescEx(parentTable, table []byte) (key, result [][]byte, err error) {
	err = liteStore.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(parentTable).Bucket(table)
		c := b.Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			key = append(key, k)
			result = append(result, v)
		}
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not set up buckets, %v", err)
	}
	return key, result, nil
}

// 대량으로 선택한다.
func (liteStore *LiteStore) LoadEntryDescLimitEx(parentTable, table []byte, count int) (key, result [][]byte, err error) {
	err = liteStore.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(parentTable).Bucket(table)
		c := b.Cursor()
		i := 0
		for k, v := c.Last(); k != nil && i < count; k, v = c.Prev() {
			key = append(key, k)
			result = append(result, v)
			i++
		}
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not set up buckets, %v", err)
	}
	return key, result, nil
}

// 하나만 선택한다.
func (liteStore *LiteStore) SelectEntryEx(parentTable, table []byte, key []byte) (result []byte, err error) {
	err = liteStore.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(parentTable).Bucket(table)
		result = b.Get(key)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not set up buckets, %v", err)
	}
	return result, nil
}

func (liteStore *LiteStore) SaveEntryEx(parentTable, table, k, v []byte) error {
	err := liteStore.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket(parentTable).Bucket(table).Put(k, v)
		return err
	})
	return err
}
