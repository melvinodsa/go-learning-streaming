package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var readingDelay = 100 * time.Microsecond
var writingDelay = 1000 * time.Microsecond

const readingBufferSize = 100

func main() {
	queue, err := readFile("input.txt")
	if err != nil {
		panic(err)
	}
	cl := writeToFile("output.txt", queue)
	<-cl
}

// Queue represents a dynamically growing queue of byte slices
type Queue struct {
	mu    sync.Mutex
	queue [][]byte
	close bool
}

// NewQueue creates a new instance of Queue
func NewQueue() *Queue {
	return &Queue{
		queue: make([][]byte, 0), // Start with an empty queue
	}
}

// Enqueue adds a new byte slice to the queue
func (q *Queue) Enqueue(data []byte) {

	d := make([]byte, len(data))
	copy(d, data)
	q.mu.Lock()
	defer q.mu.Unlock()

	// Append data to the queue
	q.queue = append(q.queue, d)
}

// Dequeue removes and returns the first byte slice from the queue
func (q *Queue) Dequeue() ([]byte, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.queue) == 0 {
		return nil, false
	}

	// Get the first element in the queue
	data := q.queue[0]

	// Remove the first element from the queue
	q.queue = q.queue[1:]

	return data, true
}

func (q *Queue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.close = true
}

// Size returns the current number of items in the queue
func (q *Queue) IsClosed() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.close
}

func (q *Queue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue) == 0
}

func readFile(path string) (*Queue, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	queue := NewQueue()
	go keepReading(f, queue)
	return queue, nil
}

func keepReading(f *os.File, queue *Queue) {
	now := time.Now()
	i := 0
	b := make([]byte, readingBufferSize)
	for {
		i++
		n, err := f.Read(b)
		if err != nil {
			break
		}
		queue.Enqueue(b[:n])
		time.Sleep(readingDelay)
	}
	queue.Close()
	t := time.Since(now)
	fmt.Println("read", i, "times")
	fmt.Println("finished in", t)
}

func writeToFile(path string, queue *Queue) chan struct{} {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	cl := make(chan struct{})
	go keepWriting(f, queue, cl)
	return cl
}

func keepWriting(f io.Writer, queue *Queue, close chan struct{}) {
	var now time.Time
	i := 0
	for !queue.IsClosed() || !queue.IsEmpty() {
		if now.IsZero() {
			now = time.Now()
		}
		time.Sleep(writingDelay)
		b, ok := queue.Dequeue()
		if !ok {
			continue
		}
		i++
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
