package main

type score int
type converter func(string) score
type teamScore map[string]score
