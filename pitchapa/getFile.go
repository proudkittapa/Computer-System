// func getFile(filename string) string {
// call_cache(filename)
// start := time.Now()
// f, err := os.Open(filename)
// f, err := os.Open("index.html")
// if err != nil {
// 	fmt.Println("File reading error", err)

// }
// defer func() {
// 	if err := f.Close(); err != nil {
// 		panic(err)
// 	}
// }()

// chunksize := 512
// reader := bufio.NewReader(f)
// part := make([]byte, chunksize)
// buffer := bytes.NewBuffer(make([]byte, 0))
// var bufferLen int
// for {
// 	count, err := reader.Read(part)
// 	if err != nil {
// 		break
// 	}
// 	bufferLen += count
// 	buffer.Write(part[:count])
// }
// fmt.Println("home")
// fmt.Println("Time get file: ", time.Since(start))
// return buffer.String()
// contentType = "text/html"
// headers = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n%s", bufferLen, contentType, buffer)

// }

package main

import (
	"fmt"
	"os"
	"sync"
)

type chunk struct {
	bufsize int
	offset  int64
}

func main() {
	const BufferSize = 100
	file, err := os.Open("text.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	filesize := int(fileinfo.Size())
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

	var wg sync.WaitGroup
	wg.Add(concurrency)
	store := make([]string, concurrency)

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
	var text string
	for i := 0; i < concurrency; i++ {
		text += store[i]
	}
	fmt.Println(text)
}
