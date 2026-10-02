package main

import (
	"fmt"
	"sync"
)

type Database struct {
	connectionString string
}

var (
	instance *Database
	once     sync.Once
)

// GetInstance returns the single shared Database instance.
func GetInstance() *Database {
	once.Do(func() {
		fmt.Println("Creating database instance...")
		instance = &Database{connectionString: "postgres://localhost:5432/mydb"}
	})
	return instance
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			db := GetInstance()
			fmt.Printf("Goroutine %d got instance at %p\n", id, db)
		}(i)
	}

	wg.Wait()
}
