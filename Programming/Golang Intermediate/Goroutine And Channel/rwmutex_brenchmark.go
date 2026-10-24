package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

/*
labib0x9@TestGoCode $ go run main.go
Readers  RWMutex      Mutex
1        59.917µs     5.333µs
2        22.125µs     14.75µs
4        13.084µs     4.167µs

8        11.25µs      13.583µs

16       13.167µs     9.208µs
32       29.291µs     24.667µs
64       37.875µs     30.667µs

128      46.666µs     144µs
256      91.542µs     104.666µs
512      159.125µs    213.958µs

1024     344.541µs    2.1375ms

2048     645.458µs    840.833µs
4096     1.249792ms   1.4085ms
8192     2.3055ms     3.0915ms
16384    4.077292ms   5.806667ms
32768    6.924958ms   9.111333ms
65536    11.332875ms  14.175625ms
131072   18.032584ms  26.18125ms
262144   36.27475ms   51.253666ms
524288   72.74025ms   102.974792ms
*/

func main() {
	producer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			l.Lock()
			l.Unlock()
			time.Sleep(1)
		}
	}

	observer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		l.Lock()
		defer l.Unlock()
	}

	test := func(count int, mutex, rwmutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count + 1)
		beginTestTime := time.Now()
		go producer(&wg, mutex)
		for i := count; i > 0; i-- {
			go observer(&wg, rwmutex)
		}

		wg.Wait()
		return time.Since(beginTestTime)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()

	var m sync.RWMutex
	fmt.Fprintf(tw, "Readers\tRWMutex\tMutex\n")
	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(tw, "%d\t%v\t%v\n", count, test(count, &m, m.RLocker()), test(count, &m, &m))
	}
}
