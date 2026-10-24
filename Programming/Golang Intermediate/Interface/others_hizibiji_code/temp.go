package main

import (
	"io"
	"time"
)

type Artifact interface {
	Title() string
	Creators() []string
	CreatedAt() time.Time
}

type Text interface {
	Pages() int
	Words() int
	PageSize() int
}

type Audio interface {
	Stream() (io.ReadCloser, error)
	RunningTime() time.Duration
	Format() string
}

type Video interface {
	Stream() (io.ReadCloser, error)
	RunningTime() time.Duration
	Format() string
	Resulation() (x, y int)
}