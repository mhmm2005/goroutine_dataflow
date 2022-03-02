package main

import (
	"fmt"
	"time"
)

type Worker struct {
	Data    chan int
	Quit    chan bool
	Stopped bool
}

func (w *Worker) Stop() {
	w.Quit <- true
	w.Stopped = true
}

func main() {
	var forEver = make(chan bool)

	w := Worker{
		Data:    make(chan int),
		Quit:    make(chan bool),
		Stopped: false,
	}

	go w.thread1()
	go w.thread3()

	go func() {
		time.Sleep(10 * time.Second)
		forEver <- true
	}()
	for {
		select {
		case <-forEver:
			return
		}
	}
}

func (w *Worker) thread1() {
	i := 0
	for {
		if w.Stopped {
			return
		}
		time.Sleep(time.Second)
		i = i + 1
		fmt.Println("value sent to thread3 from thread1")
		w.Data <- i

		if i == 5 {
			w.Stop()
			return
		}
	}
}

func (w *Worker) thread2() {
	for {
		if w.Stopped {
			return
		}
		select {
		case d := <-w.Data:
			fmt.Println("value received in thread2 from thread3 ", d)
			fmt.Println("===============")
			if w.Stopped {
				return
			}
		case <-w.Quit:
			return
		}
	}

}

func (w *Worker) thread3() {
	w2 := Worker{
		Data:    make(chan int),
		Quit:    make(chan bool),
		Stopped: false,
	}
	go w2.thread2()

	for {
		if w.Stopped {
			return
		}
		select {
		case d := <-w.Data:
			fmt.Println("value received in thread3 from thread1 ", d)
			w2.Data <- d
			fmt.Println("value sent from thread3 to thread2 ", d)

			if w.Stopped {
				return
			}
		case <-w.Quit:
			return
		}
	}

}
