package main

import (
	"fmt"
	"os"
	"pin2pre/cacheFile"
	"sync"
	"time"
)

type chunk struct {
	bufsize int
	offset  int64
}

var cacheObject cacheFile.Cache = cacheFile.NewCache()

func main() {
	// fmt.Println(cacheFile.NewCache())
	// const BufferSize = 500
	// start := time.Now()
	// file, err := os.Open("index.html")
	// call_cache("index.html")
	// if err != nil {
	// 	fmt.Println("File reading error", err)
	// 	return
	// }
	// defer func() {
	// 	if err := file.Close(); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// fileinfo, err := file.Stat()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// filesize := int(fileinfo.Size())
	// fmt.Println(filesize)
	// // Number of go routines we need to spawn.
	// concurrency := filesize / BufferSize
	// // buffer sizes that each of the go routine below should use. ReadAt
	// // returns an error if the buffer size is larger than the bytes returned
	// // from the file.
	// chunksizes := make([]chunk, concurrency)

	// // All buffer sizes are the same in the normal case. Offsets depend on the
	// // index. Second go routine should start at 100, for example, given our
	// // buffer size of 100.
	// for i := 0; i < concurrency; i++ {
	// 	chunksizes[i].bufsize = BufferSize
	// 	chunksizes[i].offset = int64(BufferSize * i)
	// }

	// // check for any left over bytes. Add the residual number of bytes as the
	// // the last chunk size.
	// if remainder := filesize % BufferSize; remainder != 0 {
	// 	c := chunk{bufsize: remainder, offset: int64(concurrency * BufferSize)}
	// 	concurrency++
	// 	chunksizes = append(chunksizes, c)
	// }

	// var wg sync.WaitGroup
	// wg.Add(concurrency)
	// store := make([]string, concurrency)
	// start2 := time.Now()
	// for i := 0; i < concurrency; i++ {
	// 	go func(chunksizes []chunk, i int) {
	// 		defer wg.Done()

	// 		chunk := chunksizes[i]
	// 		buffer := make([]byte, chunk.bufsize)
	// 		_, err := file.ReadAt(buffer, chunk.offset)

	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return
	// 		}
	// 		store[i] = string(buffer)
	// 		// fmt.Println("bytes read, string(bytestream): ", bytesread)
	// 		// fmt.Println("bytestream to string: ", string(buffer))
	// 	}(chunksizes, i)
	// }

	// wg.Wait()
	// fmt.Printf("time: %v\n", time.Since(start))
	//fmt.Printf("hello")
	// fmt.Printf("time2: %v\n", time.Since(start2))

	// var text string
	// for i := 0; i < concurrency; i++ {
	// 	text += store[i]
	// }
	// fmt.Println(text)
	fmt.Println(call_cache("index.html"))

}

func call_cache(filename string) string {
	start := time.Now()
	d, err := cacheObject.Check(filename)
	if err != nil {
		fmt.Println(err)
		a := getFile("index.html")
		cacheObject.Add(filename, a)
		d, _ = cacheObject.Check(filename)
		cacheObject.Display()

		fmt.Println("Time calling cache miss: ", time.Since(start))
		return d
	} else {
		cacheObject.Display()

		fmt.Println("Time calling cache hit: ", time.Since(start))
		return d
	}

}

func getFile(filename string) string {
	// call_cache("index.html")
	const BufferSize = 1000
	start := time.Now()
	file, err := os.Open(filename)
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
	// Number of go routines we need to spawn.
	concurrency := filesize / BufferSize
	// buffer sizes that each of the go routine below should use. ReadAt
	// returns an error if the buffer size is larger than the bytes returned
	// from the file.
	chunksizes := make([]chunk, concurrency)

	// All buffer sizes are the same in the normal case. Offsets depend on the
	// index. Second go routine should start at 100, for example, given our
	// buffer size of 100.
	for i := 0; i < concurrency; i++ {
		chunksizes[i].bufsize = BufferSize
		chunksizes[i].offset = int64(BufferSize * i)
	}

	// check for any left over bytes. Add the residual number of bytes as the
	// the last chunk size.
	if remainder := filesize % BufferSize; remainder != 0 {
		c := chunk{bufsize: remainder, offset: int64(concurrency * BufferSize)}
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
		go func(chunksizes []chunk, i int) {
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
