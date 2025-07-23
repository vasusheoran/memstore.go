-----

# **MemStore: A Simple Persisted In-Memory Key-Value Store**

memstore is a lightweight, concurrent-safe in-memory key-value store with built-in persistence. It's designed for applications that need fast data access and automatic saving of data to disk.

-----

## **Features**

* **In-Memory Speed:** Stores data directly in RAM for ultra-fast read and write operations.
* **Concurrency Safe:** Uses `sync.RWMutex` to ensure thread-safe access to data from multiple goroutines.
* **Automatic Persistence:** Periodically flushes data to a specified file on disk, preventing data loss on application restart.
* **Simple API:** Provides straightforward `Set`, `Get`, and `Delete` methods.
* **JSON Encoding:** Persists data to disk using standard JSON format, making it human-readable and easy to debug.

-----

## **Installation**

To use memstore, simply import the package into your Go project:

```bash
go get github.com/vasusheoran/memstore.go # (Replace with your actual module path)
```

-----

## **Usage**

Here's a quick example of how to use memstore:

```go
package main

import (
	"fmt"
	"time"

	memstore "github.com/vasusheoran/memstore.go"
)

func main() {
    // Initialize memstore with a flush path and a flush period of 5 seconds
    store := memstore.NewStorage("data.json", 5*time.Second)
    defer store.Close() // Ensure data is flushed on exit

    // Set some key-value pairs
    store.Set("name", "Alice")
    store.Set("age", 30)
    store.Set("city", "New York")

    // Get values
    if name, ok := store.Get("name"); ok {
        fmt.Println("Name:", name) // Output: Name: Alice
    }

    if age, ok := store.Get("age"); ok {
        fmt.Println("Age:", age)   // Output: Age: 30
    }

    // Update a value
    store.Set("age", 31)
    if age, ok := store.Get("age"); ok {
        fmt.Println("Updated Age:", age) // Output: Updated Age: 31
    }

    // Delete a value
    store.Delete("city")
    if _, ok := store.Get("city"); !ok {
        fmt.Println("City deleted successfully.") // Output: City deleted successfully.
    }

    // Data will be automatically flushed to data.json every 5 seconds.
    // You can also manually close the store to force a flush.
    fmt.Println("Data operations complete. Check data.json for persisted data.")
    time.Sleep(6 * time.Second) // Give it some time to flush
}

```

-----

## **API**

### `NewStorage(flushPath string, flushPeriod time.Duration) *Storage`

Initializes a new memstore instance.

* `flushPath`: The file path where the data will be persisted.
* `flushPeriod`: The interval at which data will be automatically flushed to disk.

### `func (s *Storage) Set(key string, value interface{})`

Stores a `value` associated with a given `key`. If the key already exists, its value will be updated.

### `func (s *Storage) Get(key string) (interface{}, bool)`

Retrieves the value associated with a given `key`. Returns the value and a boolean indicating whether the key was found.

### `func (s *Storage) Delete(key string)`

Removes the key-value pair associated with the given `key`.

### `func (s *Storage) Close() error`

Stops the periodic flushing and performs a final flush of the data to disk. It's crucial to call this method before exiting your application to ensure all data is saved.

-----

## **Error Handling**

Currently, errors during periodic flushing are ignored to prevent the background goroutine from crashing. However, errors from `Close()` or the initial `loadFromDisk()` will be returned.

-----

## **Contributing**

Feel free to open issues or pull requests if you have suggestions for improvements or find any bugs\!

-----