package main

import (
	"fmt"
	"os"
	"time"
)

// wrap -> add content to error

func OpenFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("[error] [%v] OpenFile <---> %w", time.Now().Format("2006-01-02 15:04:05"), err)
	}
	defer file.Close()
	
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("[error]  Stat() function fails <---> %w", err)
	}
	_ = info
	return nil
}

func main() {

	fmt.Println(OpenFile("test.txt"))

}