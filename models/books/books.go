package main

import (
	"fmt"
	"time"
)

type book struct {
	ID     string
	Title  string
	Author string
	createOn time.Time;
}


func NewBook(ID string,  Title string, Author string, createOn time.Time) *book {
    b := new(book)
    b.ID = ID
    b.Title = Title
    b.Author = Author
	b.createOn = time.Now();
    return b
}

// key-value store
var store = make(map[string]book)


func add(k string, b book) bool {
	if k == "" {
		return false
	}
	if lookUp(k) == nil {
		store[k] = b
		return true
	}
	return false
}

func lookUp(k string) *book {
	_, ok := store[k]
	if ok {
		b := store[k]
		return &b
	}
	return nil
}

func deleteElement(k string) bool {
	if lookUp(k) != nil {
		delete(store, k)
		return true
	}
	return false
}

func change(k string, b book) bool {
	store[k] = b
	return true
}


func printAll() {
	for k, d := range store {
		fmt.Printf("key: %s value: %v\n", k, d)
	}
}
