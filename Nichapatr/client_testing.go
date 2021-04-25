package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
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
type Mess struct {
	Mess string `json:"mess"`
}

// var wg1 sync.WaitGroup

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

func receive2(conn net.Conn) string {
	defer conn.Close()
	// fmt.Println("reading")
	message := ""
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		// fmt.Println(string(buffer[:n]))
		if !strings.Contains(string(buffer[:n]), "HTTP") {
			if _, err := conn.Write([]byte("Recieved\n")); err != nil {
				log.Printf("failed to respond to client: %v\n", err)
			}
			break
		}
		message = string(buffer[:n])
		count_Res++
		// fmt.Println("before out of loop")
		break
	}
	// fmt.Println("out of loop")

	return message
}

func client(wg1 *sync.WaitGroup, m string, p string, quan int) string {
	// t0 := time.Now()
	fmt.Println("client", userid)
	host := "178.128.94.63:8080"
	conn, err := net.Dial("tcp", host)
	if err != nil {
		count_Fail++
		log.Fatalln(err)
	}
	send(conn, host, m, p, userid, quan) //check parameter quan
	a := receive2(conn)
	wg1.Done()
	fmt.Println("client done")
	return a
	// fmt.Printf("Latency Time:   %v ", time.Since(t0))
	// <-ch
}

func clientNoGo(m string, p string, quan int) string {
	// t0 := time.Now()
	host := "178.128.94.63:8080"
	conn, err := net.Dial("tcp", host)
	if err != nil {
		count_Fail++
		log.Fatalln(err)
	}
	send(conn, host, m, p, userid, quan) //check parameter quan
	a := receive2(conn)
	// wg1.Done()
	return a
	// fmt.Printf("Latency Time:   %v ", time.Since(t0))
	// <-ch
}

var userid = 0
var count_Res = 0
var count_Fail = 0

func createHeaderGET(pathh string, u int) string {
	// fmt.Println("headerGET")
	userID := u
	method := "GET"
	path := pathh
	host := "178.128.94.63:8080"
	contentLength := len("userID:" + string(userID))
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

func onerun2(wg1 sync.WaitGroup) {

	for i := 0; i < 1000; i++ {
		wg1.Add(1)
		go func() {
			client(&wg1, "GET", "/products/1", 0)
		}()
	}
	wg1.Wait()
	// client(&wg, "GET", "/", 0)
	// client(&wg, "GET", "/products", 0)
	// client(&wg1, "GET", "/products/1", 0)
	// client(&wg, "POST", "/products/1", 2)
}
func test_time_check(wg1 sync.WaitGroup) {
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
		wg1.Add(1)
		go client(&wg1, "GET", "/"+strconv.Itoa(i), 0)
	}
	t01 := float64(time.Since(t1)) / 1e6 / 5
	fmt.Printf("Latency Time:   %v ", t01)

	t2 := time.Now()
	for i := 6; i < 11; i++ {
		wg1.Add(1)
		go client(&wg1, "GET", "/products/"+strconv.Itoa(i), 0)
	}
	t02 := float64(time.Since(t2)) / 1e6 / 5
	fmt.Printf("Latency Time:   %v \n", t02)

	t3 := time.Now()
	for i := 6; i < 11; i++ {
		wg1.Add(1)
		go client(&wg1, "GET", "/products/"+strconv.Itoa(i), 0)
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
		wg1.Add(4)
		go client(&wg1, "POST", "/products/1", 2)    // stock must = 998
		go client(&wg1, "POST", "/products/1", 3)    // stock must = 995
		go client(&wg1, "POST", "/products/1", 5)    // stock must = 990
		go client(&wg1, "POST", "/products/1", 1000) // stock must = 0
	}
	t04 := float64(time.Since(t4)) / 1e6 / 4
	fmt.Printf("Latency Time:   %v ", t04)

	t5 := time.Now()
	for i := 0; i < 2; i++ {
		wg1.Add(1)
		go client(&wg1, "POST", "/products/2", 10000) // stock must = 0
	}
	t05 := float64(time.Since(t5)) / 1e6 / 2
	fmt.Printf("Latency Time:   %v ", t05)
	wg1.Wait()
}

var num_user float64 = 100

func user_model(wg1 sync.WaitGroup) {
	for i := 0.0; i < (num_user * 0.60); i++ {
		wg1.Add(2)
		go func() {
			client(&wg1, "GET", "/", 0)
			client(&wg1, "GET", "/products", 0)
		}()
	}
	fmt.Println("here")
	for i := 0.0; i < (num_user * 0.25); i++ {
		wg1.Add(3)
		go func() {
			client(&wg1, "GET", "/", 0)
			client(&wg1, "GET", "/products", 0)
			client(&wg1, "GET", "/products/"+strconv.Itoa(rand.Intn(967)), 0)
		}()
	}
	for i := 0.0; i < (num_user * 0.15); i++ {
		wg1.Add(4)
		go func() {
			client(&wg1, "GET", "/", 0)
			client(&wg1, "GET", "/products", 0)
			client(&wg1, "GET", "/products/"+strconv.Itoa(rand.Intn(967)), 0)
			client(&wg1, "POST", "/products/"+strconv.Itoa(rand.Intn(967)), 2)
		}()
	}
	wg1.Wait()
	fmt.Println("after wait group")
}

func check(expect Rate, get Rate) {
	if get != expect {
		fmt.Printf("smt wrong!") //("expected v ==== %v \n, get v ==== %v \n", expect, get)
	} else {
		fmt.Printf("success : v ==== %v \n", get)
	}
	fmt.Println("expect:", expect)
	fmt.Println("get:", get)
}

func misshit_check() {
	//declare variables pid
	// check1 := []string{"miss", "miss", "miss", "miss", "miss"}
	// check2 := []string{"miss", "miss", "miss", "miss", "miss"}
	// check3 := []string{"hit", "hit", "hit", "hit", "hit"}

	checkU1 := Rate{Miss: 1, Hit: 4}
	for i := 1; i < 6; i++ {
		clientNoGo("GET", "/", 0)
	}
	fmt.Println("before hitmissFile")
	m := clientNoGo("GET", "/hitmissFile", 0)
	j1 := getJson(m)
	fmt.Println("j1:", j1)
	k1 := getJson2(j1.Mess)
	fmt.Println("k1:", k1)
	check(checkU1, k1) //check miss, hit
	fmt.Println("Hit for /", k1.Hit)
	fmt.Println("Miss for /", k1.Miss)

	/*-------------check(2)-------------*/

	checkP1 := Rate{Miss: 5, Hit: 0}
	for i := 1; i < 6; i++ {
		clientNoGo("GET", "/products/"+strconv.Itoa(i), 0)
	}
	m1 := clientNoGo("GET", "/hitmiss", 0)
	l1 := getJson(m1)
	n1 := getJson2(l1.Mess)
	check(checkP1, n1) //check miss, hit
	fmt.Println("Hit for /products/:id", n1.Hit)
	fmt.Println("Miss for /products/:id", n1.Miss)

	checkP2 := Rate{Miss: 10, Hit: 0}
	for i := 6; i < 11; i++ {
		clientNoGo("GET", "/products/"+strconv.Itoa(i), 0)
	}
	m2 := clientNoGo("GET", "/hitmiss", 0)
	l2 := getJson(m2)
	n2 := getJson2(l2.Mess)
	check(checkP2, n2)
	fmt.Println("Hit for /products/:id", n1.Hit)
	fmt.Println("Miss for /products/:id", n1.Miss)

	checkP3 := Rate{Miss: 10, Hit: 5}
	for i := 6; i < 11; i++ {
		clientNoGo("GET", "/products/"+strconv.Itoa(i), 0)
	}
	m3 := clientNoGo("GET", "/hitmiss", 0)
	l3 := getJson(m3)
	n3 := getJson2(l3.Mess)
	check(checkP3, n3)
	fmt.Println("Hit for /products/:id", n1.Hit)
	fmt.Println("Miss for /products/:id", n1.Miss)
}

func main() {
	// flag.Parse()
	var wg1 sync.WaitGroup
	start := time.Now()
	// misshit_check()
	// test_time_check(wg1)
	// onerun2(wg1)
	// start := time.Now()
	user_model(wg1)
	fmt.Println("after usermodel")
	// for i := 0; i < 1000; i++ {
	// 	wg1.Add(1)
	// 	go client(&wg1, "GET", "/products/1", 0)
	// // wg1.Add(1000)
	// // onerun2(wg1)
	// }
	// wg1.Wait()
	// time.Sleep(100 * time.Millisecond)
	t := time.Since(start)
	fmt.Printf("\n \nTotal TIME: %v\n", t)
	fmt.Printf("Number Response: %d\n", count_Res)
	fmt.Printf("Number fail: %d\n", count_Fail)
	tt := float64(t) / 1e6
	rate := float64(count_Res) / (tt / 1000)
	fmt.Printf("Rate per Sec: %f", rate)
	// client("GET", "/hitmiss", 0)
	// fmt.Println("HIT:", result.Hit)
	// fmt.Println("Miss:", result.Miss)
}

func getJson(message string) Mess {
	var result Mess
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
func getJson2(message string) Rate {
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
