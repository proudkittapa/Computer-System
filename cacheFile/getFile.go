//https://kgrz.io/reading-files-in-go-an-overview.html#reading-file-chunks-concurrently
package cacheFile

import (
	"fmt"
	"os"

	// "pin2pre/milestone2/cacheFile"

	// "pin2pre/cacheFile"
	"sync"
	"time"
)

type chunk struct {
	bufsize int
	offset  int64
}

type HM struct {
	Miss int `json:"miss"`
	Hit  int `json:"hit"`
}

var cacheObject Cache = NewCache()

var miss_num int

var hit_num int

func Call_cache(filename string) string {
	start := time.Now()

	d, err := cacheObject.Check(filename)
	if err != nil {
		fmt.Println(err)
		a := getFile("index.html")
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

func getFile(filename string) string {
	// call_cache("index.html")
	const BufferSize = 300
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
	// Num of go routines
	concurrency := filesize / BufferSize

	chunksizes := make([]chunk, concurrency)

	for i := 0; i < concurrency; i++ {
		chunksizes[i].bufsize = BufferSize
		chunksizes[i].offset = int64(BufferSize * i)
	}

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
	// start2 := time.Now()
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

	// fmt.Printf("hello")
	// fmt.Printf("time2: %v\n", time.Since(start2))

	var text string
	for i := 0; i < concurrency; i++ {
		text += store[i]
	}
	// fmt.Println(text)
	fmt.Printf("time: %v\n", time.Since(start))
	return text
}

func SendMissHitFile() HM {
	result := HM{Miss: miss_num, Hit: hit_num}

	// byteArray, err := json.Marshal(result)
	// CheckErr(err)

	// tmp := string(byteArray)

	return result
}
