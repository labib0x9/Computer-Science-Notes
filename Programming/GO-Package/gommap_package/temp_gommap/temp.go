package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/tysonmote/gommap"
)

type Mapper struct {
	File *os.File
	mmap gommap.MMap
}

func NewMapper(filename string) *Mapper {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	mmap, err := gommap.Map(
		file.Fd(),
		gommap.PROT_READ|gommap.PROT_WRITE,
		gommap.MAP_SHARED,
	)
	if err != nil {
		file.Close()
		panic(err)
	}

	return &Mapper{
		File: file,
		mmap: mmap,
	}
}

func (m *Mapper) Resize(maxSize int64) error {
	if err := m.File.Truncate(maxSize); err != nil {
		return err
	}

	info, err := m.File.Stat()
	if err != nil {
		return err
	}

	debug("[resize] : ", info.Size())

	mmap, err := gommap.Map(
		m.File.Fd(),
		gommap.PROT_READ|gommap.PROT_WRITE,
		gommap.MAP_SHARED,
	)
	if err != nil {
		m.File.Close()
		panic(err)
	}

	m.mmap = mmap

	return nil
}

func (m *Mapper) Close() {
	m.File.Close()
}

func debug(arr ...any) {
	fmt.Println(arr...)
}

func main() {

	filename := "store.bin"
	mapper := NewMapper(filename)
	defer mapper.Close()

	info, err := mapper.File.Stat()
	if err != nil {
		panic(err)
	}

	offset := uint64(info.Size())
	debug("[initial] : ", offset)

	encoder := binary.BigEndian

	// Panics because slices out of bound.
	// encoder.PutUint64(mapper.mmap[offset:offset + 8], uint64(12))

	mapper.Resize(int64(offset) + int64(8))
	encoder.PutUint64(mapper.mmap[offset:offset+8], uint64(12))
	offset += uint64(8)

	mapper.Resize(int64(offset) + int64(8))
	encoder.PutUint64(mapper.mmap[offset:offset+8], uint64(24))
	offset += uint64(8)

	mapper.Resize(int64(offset) + int64(8))
	encoder.PutUint64(mapper.mmap[offset:offset+8], uint64(36))
	offset += uint64(8)

	info, _ = mapper.File.Stat()
	debug(info.Size() == int64(offset))

	startPosition := uint64(0)
	for i := 0; i < 3; i++ {
		debug(encoder.Uint64(mapper.mmap[startPosition : startPosition+8]))
		startPosition += 8
	}
}
