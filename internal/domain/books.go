package main

import (
	"time"
)

type book struct {
	ID       string
	Title    string
	Author   string
	createOn time.Time
}
