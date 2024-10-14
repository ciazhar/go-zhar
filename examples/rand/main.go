package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

func main() {
	//1.83f7400b5ac0bb4.17def347a78e61f0
	//1.fc3d5dcf5686c20c.17def35244d55680
	println(RandomMessageID())
}

// IdCounter struct to hold the counter
type IdCounter struct {
	counter int
	mu      sync.Mutex
}

// Next returns the next counter value
func (c *IdCounter) Next() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counter++
	return c.counter
}

var idCounterInstance *IdCounter
var once sync.Once

// GetIdCounterInstance returns a singleton instance of IdCounter
func GetIdCounterInstance() *IdCounter {
	once.Do(func() {
		idCounterInstance = &IdCounter{}
	})
	return idCounterInstance
}

func RandomMessageID() string {
	counter := GetIdCounterInstance().Next()
	randomULong, err := rand.Int(rand.Reader, big.NewInt(1<<63-1))
	if err != nil {
		panic(err)
	}
	currentTime := time.Now().UnixNano()

	return fmt.Sprintf("%x.%x.%x", counter, randomULong.Uint64(), currentTime)
}

func RandomMessageID() string {
	counter := GetIdCounterInstance().Next()
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomULong := rand.Uint64()
	currentTime := time.Now().UnixNano()

	return fmt.Sprintf("%x.%x.%x", counter, randomULong, currentTime)
}
