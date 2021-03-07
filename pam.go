package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	// "time"
	// "strconv"
)

func client() {
	// time.Sleep(100 * time.Millisecond)
	/*GET requests are for retrieving information,
	POST requests are for creating data,
	PUT requests are for updating existing records.*/
	// count_sent++
	// l0 := time.Now().UnixNano()
	// t0 := time.Now()
	con, err := net.Dial("tcp", "0.0.0.0:8081")
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()

	// clientReader := bufio.NewReader(os.Stdin)
	// serverReader := bufio.NewReader(con)
	var request = make([]byte, 100)

	for {
		// log.Printf("Type something:")
		// clientRequest, err := clientReader.ReadString('\n')

		// switch err {
		// case nil:
		// 	clientRequest := strings.TrimSpace(clientRequest)
		// 	if _, err = con.Write([]byte(clientRequest + "\n")); err != nil {
		// 		log.Printf("failed to send the client request: %v\n", err)
		// 	}
		// default:
		// 	log.Printf("client error: %v\n", err)
		// 	return
		// }
		// serverResponse, err := serverReader.ReadString('\n')

		// switch err {
		// case nil:
		// 	log.Println(strings.TrimSpace(serverResponse))
		// default:
		// 	log.Printf("server error: %v\n", err)
		// 	return
		// }
		_, err = con.Read(request)

		if err != nil {
			log.Println("failed to read request contents")
			return
		}
		// fmt.Printf(" Latency Time:   %v ", time.Since(t0))
		log.Println(&con, string(request))
		// count_res++
		request = make([]byte, 100)

	}
	// fmt.Printf(" Latency Time:   %v ", time.Since(t0))

	// l1 := time.Now().UnixNano()
	// fmt.Printf("  Time:   %f µs\n", float64(l1-l0)/float64(count_sent)/1e3)
	
}

// var count_sent = 0
// var count_res = 0

//https://gist.github.com/AntoineAugusti/80e99edfe205baf7a094
func main() {
	

	//var start int = time.Now()
	Maxroutine := flag.Int("maxNbConcurrentGoroutines", 10, "the number of goroutines that are allowed to run concurrently")
	nbclients := flag.Int("nbJobs", 10, "the number of jobs that we need to do")
	flag.Parse()
	//concurrentGoroutines
	ch := make(chan struct{}, *Maxroutine)

	for i := 0; i < *Maxroutine; i++ {
		ch <- struct{}{}
	}
	done := make(chan bool)
	waitForAllclients := make(chan bool)

	// Collect all the jobs, and since the job is finished, we can
	// release another spot for a goroutine.
	go func() {
		for i := 0; i < *nbclients; i++ {
			<-done
			// Say that another goroutine can now start.
			ch <- struct{}{}
		}
		// We have collected all the jobs, the program
		// can now terminate
		waitForAllclients <- true
	}()
	// var total int
	// Try to start nbclients jobs
	for i := 1; i <= *nbclients; i++ {
		//time start
		// start := time.Now()
		fmt.Printf("ID: %v: waiting to launch!\n", i)
		// Try to receive from the concurrentGoroutines channel. When we have something,
		// it means we can start a new goroutine because another one finished.
		// Otherwise, it will block the execution until an execution spot is available.
		<-ch
		//fmt.Printf("ID: %v: it's my turn!\n", i)
		go func(id int) {
			client()
			fmt.Printf("ID: %v: all done!\n", id)
			done <- true
		}(i)
		// new := start.UnixNano() / int64(time.Millisecond)
		// end := time.Now().UnixNano() / int64(time.Millisecond)
		// fmt.Println("newwww", end-new)
		// total += int(end - new)
	}
	// fmt.Printf(" Total Timeeeeeeeeeeeeeeeeee:   %d ", total)
	// con, err := net.Dial("tcp", "0.0.0.0:8081")
	// if _, err = con.Write([]byte(10)); err != nil{
	// 	log.Printf("failed to respond to client: %v\n", err)
	// }
	fmt.Printf(" done")
	<-waitForAllclients // Wait for all clients to finish
	// con, err := net.Dial("tcp", "0.0.0.0:8081")
	// if err != nil{
	// 	fmt.Println(err)
	// 	return
	// }


	//time sine (total from all requests)

	// fmt.Printf("Number Sent: %d\n", count_sent)
	// fmt.Printf("Number Response: %d\n", count_res)

}
