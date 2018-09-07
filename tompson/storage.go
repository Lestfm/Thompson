package tompson

/*
	Реализация инициализации, хранения игровых комнат.
	Определяет интерфейс необходимый для сохранения комнат в БД
*/

import (
	"sync"
	"encoding/json"
	"fmt"
)

type Db interface {
	Put([]byte, []byte) error
	Get([]byte) ([]byte, error)
	Delete([]byte) error
}

type Storage struct {
	db Db
	identity sync.Map
}

func NewStorage(db Db) *Storage{
	return &Storage{
		db:db,
	}
}

func (s *Storage)Create(id string, outCount int, in []InputVec) (*Room, error) {
	var val *Room
	var err error
	if room, ok := s.identity.Load(id); !ok {
		val, err = InitRoom(id, outCount, in)
		if err != nil {
			return nil, err
		}
		s.identity.Store(id, val)
	} else {
		val = room.(*Room)
	}
	return val, nil
}

func (s *Storage) Get(id string) *Room {
	if room, ok := s.identity.Load(id); ok {
		return room.(*Room)
	} else {
		return nil
	}
}

func (s *Storage) Load(id string) *Room{
	if b,err := s.db.Get([]byte(id)); err == nil{
		room := &Room{}
		json.Unmarshal(b, room)
		return room
	}
	return nil
}

func (s *Storage) Dump(id string) error{
	if room, ok := s.identity.Load(id); !ok {
		b,_ := json.Marshal(room.(*Room))
		s.identity.Delete(id)
		if err := s.db.Put([]byte(id), b); err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("комната отсутствует")
	}
}