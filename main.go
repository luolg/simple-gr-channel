package main

import (
	"fmt"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

type Demo struct {
	Number int
	lock   sync.Mutex
}

func newDemo(number int) *Demo {
	return &Demo{
		Number: number,
	}
}

func (d *Demo) addNumber() {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.Number += 1000
}

func f1(ch1 chan *Demo) {
	for i := 0; i < 999999; i++ {
		line := i

		wg.Add(1)
		go func() {
			defer wg.Done()
			d := newDemo(line)
			fmt.Println("function f1 start new goroutine to handle : ", d)
			ch1 <- d
		}()

	}
	wg.Wait()
	close(ch1)
}

func f2(ch1, ch2 chan *Demo) {
	for data := range ch1 {
		wg.Add(1)
		tmp := data
		go func() {
			defer wg.Done()
			tmp.addNumber()
			fmt.Println("function f2 start new goroutine to handle : ", tmp)

			ch2 <- tmp
		}()
	}
	wg.Wait()
	close(ch2)
}

func main() {
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)

	ch1 := make(chan *Demo,2)
	ch2 := make(chan *Demo)

	go f1(ch1)
	go f2(ch1, ch2)

	for data := range ch2 {
		fmt.Println("finally we got data", data)
	}
}
