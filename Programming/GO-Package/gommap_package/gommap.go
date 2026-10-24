package main

import (
	"encoding/binary"
	"os"

	"github.com/tysonmote/gommap"
)

// Incomplete

func main() {

	// Where the data is stored
	file, err := os.OpenFile("data.index", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		//
	}
	defer file.Close()

	// you can also visualize like slice of []byte backed with file persistency
	const mapSize = 1024 // bytes -> 1kb
	if err = file.Truncate(mapSize); err != nil {
		//
	}

	// mmap maps the file into memory
	mmap, err := gommap.Map(
		file.Fd(),
		gommap.PROT_READ|gommap.PROT_WRITE|gommap.PROT_EXEC, // read, write mapped memory
		gommap.MAP_SHARED, // writes to memory also updates file
	)

	// PROT_READ: Allow reading
	// PROT_WRITE: Allow writing
	// PROT_EXEC: Allow execution (rarely used for data files)
	// PROT_NONE: No access

	// MAP_SHARED: Changes are visible to other processes and written to file
	// MAP_PRIVATE: Changes are private to this process (copy-on-write)

	defer mmap.UnsafeUnmap() // Unmaps the memory region

	copy(mmap, []byte("Hello"))

	// Sync changes to disk
	if err := mmap.Sync(gommap.MS_SYNC); err != nil {
		//
	}

	// mmap := gommap.MMap{}

	_ = mmap.Lock()                    // Lock memory to prevent swapping
	_ = mmap.Unlock()                  // Unlock memory
	_ = mmap.Protect(gommap.PROT_READ) // Change protection
	_, _ = mmap.IsResident()           // Check if pages are resident

	// map a file starting from a specific offset
	// gommap.MapAt()

	// store integer as bytes
	binary.BigEndian.PutUint64(mmap[0:4], uint64(123))

	// read integer from mmap
	_ = binary.BigEndian.Uint64(mmap[0:4])

	// MS_SYNC: Synchronous write (blocks until complete)
	// MS_ASYNC: Asynchronous write (returns immediately)
	// MS_INVALIDATE: Invalidate cached copies

	// Remap
	// This will not loss previous data
	mmap.Sync(gommap.MS_SYNC) // sync to persist, otherwise the data loss occurs
	mmap.UnsafeUnmap()        // umap clear the memory
	file.Truncate(2048)       // allocate more space
	mmap, err = gommap.Map(   // map again with new size
		file.Fd(),
		gommap.PROT_READ|gommap.PROT_WRITE|gommap.PROT_EXEC,
		gommap.MAP_SHARED,
	)

	// Need To Check this
	mmap, err = gommap.MapRegion( // Map the entire int(1024) cap
		file.Fd(),
		0,
		int64(1024),
		gommap.PROT_READ|gommap.PROT_WRITE|gommap.PROT_EXEC,
		gommap.MAP_SHARED,
	)
}

/**
// Regular Go slice
data := make([]byte, 1024)
data[0] = 'G'

// Memory-mapped file (mmap)
mmap, _ := gommap.Map(fd, PROT_READ|PROT_WRITE, MAP_SHARED)
mmap[0] = 'M'

**/
