package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Serialization format
// Incomplete

type A struct {
	a int
	b int
}

type B struct {
	X int
	Y int
}

func main() {

	// // What is a byte order ?
	// // sequence in which bytes are arranged in memory or in binary data streams
	// byteOrder := binary.NativeEndian
	lsbOrder := binary.LittleEndian // stores the least significant byte first
	msbOrder := binary.BigEndian    // stores the most significant byte first

	var buf bytes.Buffer
	binary.Write(&buf, msbOrder, uint32(124)) // write binary-encoded data into writer (buf)
	binary.Write(&buf, lsbOrder, uint32(124))

	var value uint32
	binary.Read(&buf, msbOrder, &value) // read binary-encoded data from reader (buf) and store to value

	// variable-length encoding of integers

	buffer := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buffer, 300)
	fmt.Println(buffer[:n]) // Compact byte representation

	x, n := binary.Uvarint(buffer)
	fmt.Println(x) // 300

	// binary.PutVarint()
	// binary.Varint()

	binary.Size(12) // number of bytes needed to encode a value

	// convert uint64 into byte slice and store the slice in buffer named slice
	buffer = make([]byte, 8)
	msbOrder.PutUint64(buffer, uint64(123456789))

	// converts byte slice into uint64
	_ = msbOrder.Uint64(buffer)

	// Size of struct
	binary.Size(A{2, 3}) // private
	binary.Size(B{2, 3}) // public

}

// Use Cases:

//     Custom network protocols

//     Embedded systems

//     Writing binary file formats

//     Parsing file headers, e.g., PNG, BMP, WAV

/**
 binary.BigEndian

	Network Protocols

    TCP/IP, UDP, HTTP headers, DNS, etc.

    Referred to as "network byte order".

	File Formats

    Some formats store data in Big-Endian for standardization:

        JPEG (image format)

        TIFF (Tagged Image File Format) – sometimes uses Big-Endian

        MP4 video files

        PNG (Portable Network Graphics)

		. Cryptography and Hashing Algorithms

    Many cryptographic algorithms define data in Big-Endian byte order.

    Example: SHA-1, MD5 (when padding data blocks)
**/

/**
 binary.LittleEndian
	Microsoft Binary Formats

    BMP (Bitmap images) – uses Little-Endian

    WAV (Waveform Audio File Format) – audio file format using Little-Endian

    AVI (Audio Video Interleave) – Little-Endian video format

	USB and Device Drivers

    Many firmware and driver-level binary data are in Little-Endian because they're designed for x86 or ARM architectures.
 **/
