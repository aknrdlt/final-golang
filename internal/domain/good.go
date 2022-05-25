package main

import (
	"time"
)

type Good struct {
	ID       string
	Name     string
	Type     string
	createOn time.Time
}


func NewGood(ID string, Name string, Type string) *Good {
	g := new(Good)
	g.ID = ID
	g.Name = Name
	g.Type = Type
	g.createOn = time.Now()
	return g
}
