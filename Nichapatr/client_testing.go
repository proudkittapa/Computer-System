package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Message struct {
	Name     string
	Quantity int
	Price    int
}
type Rate struct {
	Miss int `json:"miss"`
	Hit  int `json:"hit"`
}

// type PayInfo struct {
//  Name      string
//  ProductID int
//  Date      string
//  Time      string
//  imageName string
// }

var wg sync.WaitGroup

func send(conn net.Conn, host string, m string, p string, userid int, quan int) {
	// fmt.Println("sent")
	userid++
	if m == "GET" {
		// fmt.Println("sent GET")
		fmt.Fprintf(conn, createHeaderGET(p, userid))
		// } else if m == "POSE" && p == "/payment" {
		//  // fmt.Println("sent POST, img")
		//  fmt.Fprintf(conn, createHPimg(conn, userid))
		//  time.Sleep(1 * time.Millisecond)
		//  send_file(conn)
	} else {
		// fmt.Println("sent POST")
		fmt.Fprintf(conn, createHeaderPOST(userid, quan))
	}
}

var result Rate

func recv(conn net.Conn) {
	defer conn.Close()
	// fmt.Println("reading")
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		count_Fail++
		log.Println("failed to read contents", message)
		return
	} else {
		count_Res++
	}
	fmt.Print(message)
	result = getJson(message)
}

func client(m string, p string, quan int) {
	// t0 := time.Now()
	host := "178.128.94.63:8080"
	conn, err := net.Dial("tcp", host)
	if err != nil {
		count_Fail++
		log.Fatalln(err)
	}
	send(conn, host, m, p, userid, quan) //check parameter quan
	recv(conn)
	// fmt.Printf("Latency Time:   %v ", time.Since(t0))
	// wg.Done()
	// <-ch
}

var userid = 0
var count_Res = 0
var count_Fail = 0

// var n = flag.Int("n", 5, "Number of goroutines to create")
// var ch = make(chan byte)
func createHeaderGET(pathh string, u int) string {
	userID := u
	method := "GET"
	path := pathh
	host := "178.128.94.63:8080"
	contentLength := 0
	contentType := "text"
	headers := fmt.Sprintf("%s %s HTTP/1.1\r\nHost: %s\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n userID:%d",
		method, path, host, contentLength, contentType, userID)
	return headers
}

func createHeaderPOST(u int, quan int) string {

	userID := u
	method := "POST"
	path := "/products/" + strconv.Itoa(rand.Intn(10))
	host := "178.128.94.63:8080"

	contentType := "application/json"
	jsonStr := Message{Name: "mos", Quantity: quan}
	jsonData, err := json.Marshal(jsonStr)
	if err != nil {
		fmt.Println(err)
	}
	contentLength := len(string(jsonData))
	headers := fmt.Sprintf("%s %s HTTP/1.1\r\nHost: %s\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n%s userID:%d",
		method, path, host, contentLength, contentType, string(jsonData), userID)
	return headers
}

func onerun() {
	// for i := 0; i < 200; i++ {
	// client("GET", "/", 0)
	// client("GET", "/products", 0)
	client("GET", "/products/1", 0)
	// client("POST", "/products/1", 2)
	// client6("POST", "/payment", 0)
	// }
}
func test_check() {
	/*--------------------Cache check (2)--------------------*/
	// t5 := time.Now()
	// for i := 0; i < 1000; i++ {
	//  client6("POST", "/products/1", 2) // stock must = 0
	// }
	// t05 = float64(time.Since(t5)) / 1e6 / 5
	// fmt.Printf("Latency Time:   %v ", t05)
	// fmt.Printf("Number Response: %d\nIf number of Responses = 1000, is it success or not since it out of stock at Order500?", count_Res)
	/*--------------------Cache check (1)--------------------*/
	t1 := time.Now()
	for i := 1; i < 6; i++ {
		client("GET", "/products/"+strconv.Itoa(i), 0)
	}
	t01 := float64(time.Since(t1)) / 1e6 / 5
	fmt.Printf("Latency Time:   %v ", t01)

	t2 := time.Now()
	for i := 6; i < 11; i++ {
		client("GET", "/products/"+strconv.Itoa(i), 0)
	}
	t02 := float64(time.Since(t2)) / 1e6 / 5
	fmt.Printf("Latency Time:   %v \n", t02)

	t3 := time.Now()
	for i := 6; i < 11; i++ {
		client("GET", "/products/"+strconv.Itoa(i), 0)
	}
	t03 := float64(time.Since(t3)) / 1e6 / 5
	fmt.Printf("Latency Time:   %v \n", t03)
	if math.Abs(t01-t02) <= 1 {
		fmt.Println("miss?")
	} else {
		fmt.Println("something is not right(1) :")
		fmt.Println(t01 - t02)
	}
	if t03 <= t02 {
		fmt.Println("faster")
	} else {
		fmt.Println("cache not make faster maybe not hit")
	}
	/*--------------------Cache check (2)--------------------*/
	t4 := time.Now()
	for i := 0; i < 2; i++ {
		client("POST", "/products/1", 2)    // stock must = 998
		client("POST", "/products/1", 3)    // stock must = 995
		client("POST", "/products/1", 5)    // stock must = 990
		client("POST", "/products/1", 1000) // stock must = 0
	}
	t04 := float64(time.Since(t4)) / 1e6 / 4
	fmt.Printf("Latency Time:   %v ", t04)

	t5 := time.Now()
	for i := 0; i < 2; i++ {
		client("POST", "/products/2", 10000) // stock must = 0
	}
	t05 := float64(time.Since(t5)) / 1e6 / 2
	fmt.Printf("Latency Time:   %v ", t05)
}

var num_user float64 = 100

func user_model() {
	go func() {
		for i := 0.0; i < (num_user * 0.60); i++ {
			go func() {
				client("GET", "/", 0)
				client("GET", "/products", 0)
			}()
		}
	}()
	go func() {
		for i := 0.0; i < (num_user * 0.25); i++ {
			go func() {
				client("GET", "/", 0)
				client("GET", "/products", 0)
				client("GET", "/products/"+strconv.Itoa(rand.Intn(967)), 0)
			}()
		}
	}()

	go func() {
		for i := 0.0; i < (num_user * 0.15); i++ {
			go func() {
				client("GET", "/", 0)
				client("GET", "/products", 0)
				client("GET", "/products/"+strconv.Itoa(rand.Intn(967)), 0)
				client("POST", "/products/"+strconv.Itoa(rand.Intn(967)), 2)
			}()
		}
	}()

}
func check() {
	//declare variables pid
	check1 := []string{"miss", "miss", "miss", "miss", "miss"}
	// check2 := []string{"miss", "miss", "miss", "miss", "miss"}
	// check3 := []string{"hit", "hit", "hit", "hit", "hit"}

	for i := 1; i < 6; i++ {
		client("GET", "/products/"+strconv.Itoa(i), 0)
	}
	//check
	for i := range check1 {
		if check1[1] != "00" {
			fmt.Printf("fail at %d", i)
		} else {
			return
		}
	}
	fmt.Printf("success")

	for i := 6; i < 11; i++ {
		client("GET", "/products/"+strconv.Itoa(i), 0)
	}

	for i := 6; i < 11; i++ {
		client("GET", "/products/"+strconv.Itoa(i), 0)
	}

	// check4 := []string{"miss", "hit", "hit", "hit", "hit"}
}

func main() {
	// flag.Parse()
	start := time.Now()
	// test_check()
	// user_model()
	onerun()
	// wg.Wait()
	// time.Sleep(100 * time.Millisecond)
	t := time.Since(start)
	fmt.Printf("\n \nTotal TIME: %v\n", t)
	fmt.Printf("Number Response: %d\n", count_Res)
	fmt.Printf("Number fail: %d\n", count_Fail)
	tt := float64(t) / 1e6
	rate := float64(count_Res) / (tt / 1000)
	fmt.Printf("Rate per Sec: %f", rate)
	client("GET", "/hitmiss", 0)
	fmt.Println("HIT:", result.Hit)
	fmt.Println("Miss:", result.Miss)
}

func getJson(message string) Rate {
	var result Rate
	if strings.ContainsAny(string(message), "}") {

		r, _ := regexp.Compile("{([^)]+)}")
		match := r.FindString(message)
		// fmt.Println(match)
		fmt.Printf("%T\n", match)
		json.Unmarshal([]byte(match), &result)
		// fmt.Println("data", result)
	}
	return result
}
