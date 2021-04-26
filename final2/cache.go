//ref https://github.com/Lebonesco/go_lru_cache/blob/master/main.go
package final2

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"time"
)

const SIZE = 5 // size of cache

// maps data to node in Queue
type Hash (map[string]string)

// type hash map[int]byte

type Cache struct {
	Hash Hash
}
type HM struct {
	Miss int `json:"miss"`
	Hit  int `json:"hit"`
}

var cacheObject Cache = NewCache()

var miss_num int

var hit_num int

func NewCache() Cache {
	return Cache{Hash: Hash{}}
}

func (c *Cache) Check(str string) (string, error) {
	if _, ok := c.Hash[str]; ok {
		return c.Hash[str], nil
	} else {
		return "", errors.New("key doesn't exists")
	}
}

func (c *Cache) Remove(key string) {
	fmt.Printf("remove key: %s\n", key)
	delete(c.Hash, key)
}

func (c *Cache) Add(key string, value string) {
	c.Hash[key] = value
}

func (c *Cache) Display() {
	for key, _ := range c.Hash {
		fmt.Printf("{%s}\n", key)
	}
}

func getFile(filename string) string {
	// call_cache(filename)
	start := time.Now()
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("File reading error", err)

	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	chunksize := 512
	reader := bufio.NewReader(f)
	part := make([]byte, chunksize)
	buffer := bytes.NewBuffer(make([]byte, 0))
	var bufferLen int
	for {
		count, err := reader.Read(part)
		if err != nil {
			break
		}
		bufferLen += count
		buffer.Write(part[:count])
	}
	// fmt.Println("home")
	fmt.Println("Time get file: ", time.Since(start))
	return buffer.String()
	// contentType = "text/html"
	// headers = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n%s", bufferLen, contentType, buffer)

}

func Call_cache(filename string) string {
	start := time.Now()

	d, err := cacheObject.Check(filename)
	if err != nil {
		fmt.Println(err)
		a := getFile("/root/go/src/Computer-System/pre-order/" + filename)
		cacheObject.Add(filename, a)
		d, _ = cacheObject.Check(filename)
		cacheObject.Display()
		miss_num += 1
		fmt.Println("Cache miss: ", miss_num)
		fmt.Println("Time calling cache miss: ", time.Since(start))
		return d
	} else {
		cacheObject.Display()
		hit_num += 1

		fmt.Println("Cache hit: ", hit_num)
		fmt.Println("Time calling cache hit: ", (time.Since(start)))
		return d
	}

}
