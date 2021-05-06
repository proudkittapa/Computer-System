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
	// fmt.Println("send")
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
		fmt.Fprintf(conn, createHeaderPOST(userid, quan, p))
	}
	// fmt.Println("send done")
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
	// fmt.Println("client", userid)
	host := "178.128.94.63:8081"
	// host := "localhost:8081"
	conn, err := net.Dial("tcp", host)
	if err != nil {
		count_Fail++
		log.Fatalln(err)
	}
	send(conn, host, m, p, userid, quan) //check parameter quan
	a := receive2(conn)
	wg1.Done()
	// fmt.Println("client done")
	return a
	// fmt.Printf("Latency Time:   %v ", time.Since(t0))
	// <-ch
}

func clientNoGo(m string, p string, quan int) string {
	// t0 := time.Now()
	host := "178.128.94.63:8081"
	// host := "localhost:8081"

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

func createHeaderPOST(u int, quan int, p string) string {

	userID := u
	method := "POST"
	// path := "/products/" + strconv.Itoa(rand.Intn(10))
	path := p
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
			// client(&wg1, "GET", "/", 0)
			// client(&wg1, "GET", "/products", 0)
			client(&wg1, "GET", "/products/1", 0)
			// client(&wg1, "POST", "/products/1", 2)
		}()
	}
	wg1.Wait()

}
func tchecku(t01 float64, t02 float64, t03 float64) {
	if math.Abs(t03-t02) <= 1 {
		fmt.Println("Both are hit, so time is similar (success)") // t01 ad t02 are both hit, so time must be similar
	} else {
		fmt.Println("the different btw time of hit are not similar:")
		fmt.Println(t01 - t02)
	}
	if t03 < t01 && t02 < t01 {
		fmt.Println("t02 & t03 is faster, case3 and 2 Hit (success)") // t03 is time when it's hit; t01 is time when it's miss
	} else {
		fmt.Println("cache not make faster maybe not hit")
	}
}
func tcheckp(t01 float64, t02 float64, t03 float64) {
	if math.Abs(t01-t02) <= 1 {
		fmt.Println("Both are Miss, so time is similar (success)") // t01 ad t02 are both miss, so time must be similar
	} else {
		fmt.Println("the different btw time of misss are not similar:")
		fmt.Println(t01 - t02)
	}
	if t03 < t02 {
		fmt.Println("it is faster, case3 Hit (success)") // t03 is time when it's hit; t02 is time when it's miss
	} else {
		fmt.Println("cache not make faster maybe not hit")
	}
}
func test_time_check(wg1 sync.WaitGroup) { /*-------------------- time check --------------------*/
	/*
		t1 := time.Now()          //Uye
		clientNoGo("GET", "/", 0) //miss1
		t01 := float64(time.Since(t1)) / 1e6 / 5
		fmt.Printf("Latency Time:   %v ", t01)

		t2 := time.Now()
		for i := 6; i < 11; i++ {
			wg1.Add(1)
			go client(&wg1, "GET", "/", 0) //hit5
		}
		t02 := float64(time.Since(t2)) / 1e6 / 5
		fmt.Printf("Latency Time:   %v \n", t02)

		t3 := time.Now()
		for i := 6; i < 11; i++ {
			wg1.Add(1)
			go client(&wg1, "GET", "/", 0) //hit5
		}
		t03 := float64(time.Since(t3)) / 1e6 / 5
		fmt.Printf("Latency Time:   %v \n", t03)
		tchecku(t01, t02, t03)
	*/
	fmt.Println("-------------PUNE-----------")
	tp1 := time.Now() //Pune
	// var tp01 float64
	for i := 1; i < 6; i++ {
		wg1.Add(1)
		tp1 := time.Now()
		client(&wg1, "GET", "/products/"+strconv.Itoa(i), 0)
		tp01 := float64(time.Since(tp1)) / 1e6
		fmt.Printf("t01 Latency Time MISS:   %v \n", tp01)
	}
	tp01 := float64(time.Since(tp1)) / 1e6 / 5.0
	fmt.Printf("Latency Time MISS:   %v \n", tp01)
	/*
		tp2 := time.Now()
		// var tp02 float64
		for i := 6; i < 11; i++ {
			wg1.Add(1)
			fmt.Println(i)
			tp2 := time.Now()
			client(&wg1, "GET", "/products/"+strconv.Itoa(i), 0)
			tp02 := float64(time.Since(tp2)) / 1e6
			fmt.Printf("t02 Latency Time MISS:   %v \n", tp02)
		}
		tp02 := float64(time.Since(tp2)) / 1e6 / 5.0
		fmt.Printf("Latency Time MISS:   %v \n", tp02)

		tp3 := time.Now()
		// var tp03 float64
		for i := 6; i < 11; i++ {
			wg1.Add(1)
			fmt.Println(i)
			tp3 := time.Now()
			client(&wg1, "GET", "/products/"+strconv.Itoa(i), 0)
			tp03 := float64(time.Since(tp3)) / 1e6
			fmt.Printf("t03 Latency Time HIT:   %v \n", tp03)
		}
		tp03 := float64(time.Since(tp3)) / 1e6 / 5.0
		fmt.Printf("Latency Time HIT:   %v \n", tp03)
		tcheckp(tp01, tp02, tp03)

		/*--------------------time check (2)--------------------*/
	/*
		fmt.Println("-------------MIND-----------")
		t4 := time.Now() //Mind
		for i := 0; i < 5; i++ {
			wg1.Add(4)
			go client(&wg1, "POST", "/products/4", 2)   // stock must = 998
			go client(&wg1, "POST", "/products/4", 3)   // stock must = 995
			go client(&wg1, "POST", "/products/4", 5)   // stock must = 990
			go client(&wg1, "POST", "/products/4", 200) // stock must = 790
		}
		t04 := float64(time.Since(t4)) / 1e6
		fmt.Printf("Time:   %v ", t04)
		fmt.Printf("Latency Time:   %v ", (t04 / 20.0))

		// t5 := time.Now()
		// for i := 0; i < 2; i++ {
		// 	wg1.Add(1)
		// 	go client(&wg1, "POST", "/products/2", 10000) // stock must = 0
		// }
		// t05 := float64(time.Since(t5)) / 1e6 / 2
		// fmt.Printf("Latency Time:   %v ", t05)
		wg1.Wait()
	*/
}
func random(x int, y int) int {
	min := x
	max := y
	randomNum := min + rand.Intn(max-min+1)
	return randomNum
}
func quantity_check(wg1 sync.WaitGroup) { //Mind /*-------------------- quantity_check --------------------*/
	// 10 users && 1000 products in database (/product/3)
	// "The order is out of stock"

	fmt.Println("-----------------case1------------------------")
	for i := 0; i < 5; i++ {
		wg1.Add(1)
		go func() {
			a := client(&wg1, "POST", "/products/1", 200)
			mes1 := getJson(a)
			fmt.Println(qcheck(mes1.Mess, "transaction successful"))
		}()
	}
	wg1.Wait()
	wg1.Add(1)
	a := client(&wg1, "POST", "/products/1", 100)
	wg1.Wait()
	mes1 := getJson(a)
	fmt.Println(qcheck(mes1.Mess, "The order is out of stock"))
	fmt.Println("case 1 done")
	// 10 users && 10,000 products in database (/product/4) && random quantity in first Fifth orders, last order's quantity is more than stock quantity
	// "order more than stock quantity"

	fmt.Println("-----------------case2------------------------")

	for i := 0; i < 5; i++ {
		wg1.Add(1)
		go func() {
			a := client(&wg1, "POST", "/products/2", random(100, 150))
			mes1 := getJson(a)
			fmt.Println(qcheck(mes1.Mess, "transaction successful"))
		}()
	}
	wg1.Wait()
	fmt.Println("case 2 done")
	wg1.Add(1)
	a = client(&wg1, "POST", "/products/2", 500)

	mes1 = getJson(a)
	fmt.Println(qcheck(mes1.Mess, "order more than stock quantity"))
	// // unpredict result numer of "transaction successful"&"The order is out of stock"

	fmt.Println("-----------------case3------------------------")
	suc := 0
	for i := 0; i < 5; i++ {
		wg1.Add(2)
		go func() {
			a := client(&wg1, "POST", "/products/3", 100)
			mes1 := getJson(a)
			if qcheck2(mes1.Mess, "transaction successful") == "success" {
				suc++
			}
		}()
		go func() {
			a := client(&wg1, "POST", "/products/3", 200)
			mes1 := getJson(a)
			if qcheck2(mes1.Mess, "transaction successful") == "success" {
				suc++
			}
		}()
	}
	wg1.Wait()
	unpredictcheck(suc)

}

func qcheck(message string, expect string) string {
	if message == "" {
		fmt.Println("No message")
	} else if message == expect {
		fmt.Printf("-------success------ expect: %s, \nget: %s\n", expect, message)
		return "success"
	} else {
		fmt.Printf("-------Fail------ expect: %s, \nget: %s\n", expect, message)
		return "fail"
	}
	return "fail"
}
func qcheck2(message string, expect string) string {
	if message == "" {
		fmt.Println("No message")
	} else if message == expect {
		fmt.Println("-success-")
		return "success"
	} else {
		fmt.Println("-Fail-")
		return "fail"
	}
	return "fail"
}

func unpredictcheck(success int) {
	if success == 7 || success == 5 || success == 6 {
		fmt.Printf("-------------------success get %d success in this senario------------------ \n", success)
	} else {
		fmt.Println("--------------------fail-------------------")
	}
}

var num_user float64 = 3000

func user_model(wg1 sync.WaitGroup) { /*-------------------- user_model --------------------*/

	t7 := time.Now()
	for i := 0.0; i < (num_user * 0.15); i++ {
		wg1.Add(1)
		go client(&wg1, "POST", "/products/"+strconv.Itoa(rand.Intn(967)), 2)
	}
	wg1.Wait()
	fmt.Printf("\n------> TIME t7: %v\n", time.Since(t7))

	t5 := time.Now()
	for i := 0.0; i < (num_user * 0.4); i++ {
		wg1.Add(1)
		go client(&wg1, "GET", "/products/"+strconv.Itoa(rand.Intn(967)), 0)
	}
	wg1.Wait()
	fmt.Printf("\n------> TIME t5: %v\n", time.Since(t5))

	t3 := time.Now()
	for i := 0.0; i < (num_user * 1.00); i++ {
		wg1.Add(1)
		go client(&wg1, "GET", "/products?limit=5&offset=0", 0)
	}
	wg1.Wait()
	fmt.Printf("\n------> TIME t3: %v\n", time.Since(t3))

	t1 := time.Now()
	for i := 0.0; i < (num_user * 1.00); i++ {
		wg1.Add(1)
		go client(&wg1, "GET", "/", 0)
	}
	wg1.Wait()
	fmt.Printf("\n------> TIME t1: %v\n", time.Since(t1))

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

func misshit_check() { /*------------------------------------------------ miss/hit_check -------------------------------*/
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
	fmt.Println("Hit for /products/:id", n2.Hit)
	fmt.Println("Miss for /products/:id", n2.Miss)

	checkP3 := Rate{Miss: 10, Hit: 5}
	for i := 6; i < 11; i++ {
		clientNoGo("GET", "/products/"+strconv.Itoa(i), 0)
	}
	m3 := clientNoGo("GET", "/hitmiss", 0)
	l3 := getJson(m3)
	n3 := getJson2(l3.Mess)
	check(checkP3, n3)
	fmt.Println("Hit for /products/:id", n3.Hit)
	fmt.Println("Miss for /products/:id", n3.Miss)
}

func completed_flow1() { /*-------------------------------------- baseline -------------------------------*/
	t1 := time.Now()
	clientNoGo("GET", "/", 0)
	fmt.Printf("\n------> TIME t1: %v\n", time.Since(t1))
	t3 := time.Now()
	clientNoGo("GET", "/products", 0)
	fmt.Printf("\n------> TIME t3: %v\n", time.Since(t3))
	t5 := time.Now()
	clientNoGo("GET", "/products/10", 0)
	fmt.Printf("\n------> TIME t5: %v\n", time.Since(t5))
	t7 := time.Now()
	clientNoGo("POST", "/products/10", 2)
	fmt.Printf("\n------> TIME t7: %v\n", time.Since(t7))
}
func completed_flowN() { /*-------------------------------------- baseline No Go -------------------------------*/
	n := 10000
	t1 := time.Now()
	for i := 0; i < n; i++ {
		clientNoGo("GET", "/", 0)
	}
	fmt.Printf("\n------> TIME t1: %v\n", time.Since(t1))

	t3 := time.Now()
	for i := 0; i < n; i++ {
		clientNoGo("GET", "/products", 0)
	}
	fmt.Printf("\n------> TIME t3: %v\n", time.Since(t3))
	t5 := time.Now()
	for i := 0; i < n; i++ {
		clientNoGo("GET", "/products/10", 0)
	}
	fmt.Printf("\n------> TIME t5: %v\n", time.Since(t5))
	t7 := time.Now()
	for i := 0; i < n; i++ {
		clientNoGo("POST", "/products/10", 2)
	}
	fmt.Printf("\n------> TIME t7: %v\n", time.Since(t7))
	clientNoGo("GET", "/resetTime", 0)
}

func completed_flow(wg sync.WaitGroup, n int) { /*-------------------------------------- baseline with go -------------------------------*/
	t1 := time.Now()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go client(&wg, "GET", "/", 0)
	}
	wg.Wait()
	fmt.Printf("\n------> TIME t1: %v\n", time.Since(t1))

	t3 := time.Now()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go client(&wg, "GET", "/products?limit=5&offset=0", 0)
	}
	wg.Wait()
	fmt.Printf("\n------> TIME t3: %v\n", time.Since(t3))

	t5 := time.Now()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go client(&wg, "GET", "/products/10", 0)
	}
	wg.Wait()
	fmt.Printf("\n------> TIME t5: %v\n", time.Since(t5))

	t7 := time.Now()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go client(&wg, "POST", "/products/"+strconv.Itoa(rand.Intn(967)), 2)
	}
	wg.Wait()
	fmt.Printf("\n------> TIME t7: %v\n", time.Since(t7))

	// clientNoGo("GET", "/timeFunction", 0)
}

func main() {
	var wg1 sync.WaitGroup
	start := time.Now()

	// fmt.Println("---------------miss hit check---------------")
	// misshit_check()
	// fmt.Println("---------------quantity_check---------------")
	// quantity_check(wg1)
	// fmt.Println("-----------------time_check-----------------")
	// test_time_check(wg1)
	fmt.Println("-----------------RUN-----------------")
	// completed_flow1()
	// completed_flowN()
	completed_flow(wg1, 1000)
	// onerun2(wg1)
	// user_model(wg1)
	fmt.Println("-----------------END-----------------")
	t := time.Since(start)
	fmt.Printf("\n \nTotal TIME: %v\n", t)
	fmt.Printf("Number Response: %d\n", count_Res)
	fmt.Printf("Number fail: %d\n", count_Fail)
	tt := float64(t) / 1e6
	rate := float64(count_Res) / (tt / 1000)
	fmt.Printf("Rate per Sec: %f\n", rate)
	clientNoGo("GET", "/timeFunction", 0)
	clientNoGo("GET", "/hitmissFile", 0)
	// clientNoGo("GET", "/resetTime", 0)
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
