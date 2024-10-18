package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

var readingDelay = 10 * time.Microsecond
var writingDelay = 10 * time.Microsecond

const bufferSize = 20
const readingBufferSize = 100

func main() {
	ch, err := readFile("input.txt")
	if err != nil {
		panic(err)
	}
	outChan := buffer(ch)
	cl := writeToFile("output.txt", outChan)
	<-cl
}

func readFile(path string) (chan []byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	ch := make(chan []byte)
	go keepReading(f, ch)
	return ch, nil
}

func keepReading(f *os.File, ch chan []byte) {
	now := time.Now()
	i := 0
	b := make([]byte, readingBufferSize)
	for {
		i++
		n, err := f.Read(b)
		if err != nil {
			break
		}
		d := make([]byte, n)
		copy(d, b[:n])
		ch <- d
		time.Sleep(readingDelay)
	}
	close(ch)
	t := time.Since(now)
	fmt.Println("read", i, "times")
	fmt.Println("finished in", t)
}

func buffer(ch chan []byte) chan []byte {
	bf := bytes.NewBuffer([]byte{})
	outChan := make(chan []byte)
	go func() {
		for b := range ch {
			bf.Write(b)
		}
		p := [bufferSize]byte{}
		for {
			n, err := bf.Read(p[:])
			if err != nil {
				break
			}

			d := make([]byte, n)
			copy(d, p[:n])
			outChan <- d
		}
		close(outChan)
	}()
	return outChan
}

func writeToFile(path string, ch chan []byte) chan struct{} {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	cl := make(chan struct{})
	go keepWriting(f, ch, cl)
	return cl
}

func keepWriting(f io.Writer, ch chan []byte, close chan struct{}) {
	var now time.Time
	i := 0
	for b := range ch {
		i++
		if now.IsZero() {
			now = time.Now()
		}
		time.Sleep(writingDelay)
		_, err := f.Write(b)
		if err != nil {
			break
		}
	}
	t := time.Since(now)
	fmt.Println("wrote", i, "times")
	fmt.Println("finished in", t)
	close <- struct{}{}
}
