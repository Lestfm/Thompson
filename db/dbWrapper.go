package db

import "github.com/syndtr/goleveldb/leveldb"

/*
	Враппер для levelDb, реализуцющий интерфейс tompson.Db
*/

type Db struct {
	db *leveldb.DB
}

func (db *Db) Put(key, val []byte) error {
	return db.db.Put(key, val, nil)
}

func (db *Db) Get(key []byte) ([]byte, error) {
	return db.db.Get(key, nil)
}

func (db *Db) Delete(key []byte) error {
	return db.db.Delete(key, nil)
}

func NewDb(fileName string) (*Db, error) {
	db, err := leveldb.OpenFile(fileName, nil)
	if err != nil{
		return nil, err
	}
	return &Db{
		db,
	}, nil
}
