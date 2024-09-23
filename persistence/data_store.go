package persistence

import (
	"sync"
	"time"
)

type Item struct {
	Id        int
	Value     string
	Timestamp int64
}

type DataStore struct {
	mu     sync.Mutex // maybe concurrent map instead of mutex
	items  map[int]Item
	lastId int
}

func NewStore() *DataStore {
	dataStore := DataStore{}

	dataStore.items = make(map[int]Item)
	dataStore.lastId = -1

	return &dataStore
}

func (s *DataStore) Insert(value string) Item {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.lastId + 1
	item := Item{id, value, time.Now().UnixMilli()}
	s.lastId = id
	s.items[id] = item

	return item
}

func (s *DataStore) InsertById(item Item) bool {
	if s.lastId < item.Id {
		s.lastId = item.Id
	}

	_, ok := s.items[item.Id]

	s.items[item.Id] = item

	return ok
}

func (s DataStore) GetAll() []Item {
	var acc = make([]Item, 0, s.lastId+1)

	for i := 0; i <= s.lastId; i++ {
		item, ok := s.items[i]

		// early return if missing some elements
		if !ok {
			return acc
		}

		acc = append(acc, item)
	}

	return acc
}
