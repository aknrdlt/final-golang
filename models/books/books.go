package main

import (
	"fmt"
	"time"
)

type good struct {
	ID       string
	Name     string
	Type     string
	createOn time.Time
}

func NewGood(ID string, Name string, Type string) *good {
	g := new(good)
	g.ID = ID
	g.Name = Name
	g.Type = Type
	g.createOn = time.Now()
	return g
}

// key-value store
var store = make(map[string]good)

func putElement(k string, b good) bool {
	if k == "" {
		return false
	}
	if getElement(k) == nil {
		store[k] = b
		return true
	}
	return false
}

func getElement(k string) *good {
	_, ok := store[k]
	if ok {
		b := store[k]
		return &b
	}
	return nil
}

func deleteElement(k string) bool {
	if getElement(k) != nil {
		delete(store, k)
		return true
	}
	return false
}

func updateElement(k string, b good) bool {
	store[k] = b
	return true
}

func getAll() {
	for k, d := range store {
		fmt.Printf("key: %s value: %v\n", k, d)
	}
}
