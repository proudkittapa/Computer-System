package cacheFile

import (
	"fmt"
	"os"

	// "pin2pre/cacheFile"
	"sync"
	"time"
)

type Chunk struct {
	bufsize int
	offset  int64
}

var CacheObject Cache = NewCache()

func Call_cache(filename string) string {
	start := time.Now()
	d, err := CacheObject.Check(filename)
	if err != nil {
		fmt.Println(err)
		a := GetFile("index.html")
		CacheObject.Add(filename, a)
		d, _ = CacheObject.Check(filename)
		CacheObject.Display()

		fmt.Println("Time calling cache miss: ", time.Since(start))
		return d
	} else {
		CacheObject.Display()

		fmt.Println("Time calling cache hit: ", time.Since(start))
		return d
	}

}

func GetFile(filename string) string {
	// call_cache("index.html")
	const BufferSize = 500
	start := time.Now()
	file, err := os.Open("../pre-order/" + filename)
	if err != nil {
		fmt.Println("File reading error", err)
		return ""
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	filesize := int(fileinfo.Size())
	fmt.Println(filesize)
	// Num of go routines
	concurrency := filesize / BufferSize

	chunksizes := make([]Chunk, concurrency)

	for i := 0; i < concurrency; i++ {
		chunksizes[i].bufsize = BufferSize
		chunksizes[i].offset = int64(BufferSize * i)
	}

	if remainder := filesize % BufferSize; remainder != 0 {
		c := Chunk{bufsize: remainder, offset: int64(concurrency * BufferSize)}
		concurrency++
		chunksizes = append(chunksizes, c)
	}

	// var wg sync.WaitGroup
	// wg.Add(concurrency)
	// store := make([]string, concurrency)
	var wg sync.WaitGroup
	wg.Add(concurrency)
	store := make([]string, concurrency)
	start2 := time.Now()
	for i := 0; i < concurrency; i++ {
		go func(chunksizes []Chunk, i int) {
			defer wg.Done()

			chunk := chunksizes[i]
			buffer := make([]byte, chunk.bufsize)
			_, err := file.ReadAt(buffer, chunk.offset)

			if err != nil {
				fmt.Println(err)
				return
			}
			store[i] = string(buffer)
			// fmt.Println("bytes read, string(bytestream): ", bytesread)
			// fmt.Println("bytestream to string: ", string(buffer))
		}(chunksizes, i)
	}

	wg.Wait()
	fmt.Printf("time: %v\n", time.Since(start))
	fmt.Printf("hello")
	fmt.Printf("time2: %v\n", time.Since(start2))

	var text string
	for i := 0; i < concurrency; i++ {
		text += store[i]
	}
	fmt.Println(text)

	return text
}
