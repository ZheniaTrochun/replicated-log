package repository

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

var db *DataStore

func InitDataStore() {
	store := DataStore{}

	store.items = make(map[int]Item)
	store.lastId = 0

	db = &store
}

func Insert(value string) Item {
	db.mu.Lock()
	defer db.mu.Unlock()

	id := db.lastId + 1
	item := Item{id, value, time.Now().UnixMilli()}
	db.lastId = id
	db.items[id] = item

	return item
}

func InsertById(item Item) bool {
	if db.lastId < item.Id {
		db.lastId = item.Id
	}

	_, ok := db.items[item.Id]

	db.items[item.Id] = item

	return ok
}

func GetAll() []Item {
	var acc = make([]Item, 0, db.lastId+1)

	for i := 1; i <= db.lastId; i++ {
		item, ok := db.items[i]

		// early return if missing some elements
		if !ok {
			return acc
		}

		acc = append(acc, item)
	}

	return acc
}
