
package cacheFile
/*
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

	aboutfile, err := file.stat()
	if err != nil {
		fmt.Println("error to get info of file")
		return
	}

	filesize := aboutfile.Size()
	endding := filesize - 1
	sizelast := 0

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
*/
