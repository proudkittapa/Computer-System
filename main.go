package main

import (
	"encoding/json"
	"fmt"
	"pin2pre/Kittapa"
	"pin2pre/cacheFile"
	"time"
	// "pin2pre/transaction"
	// "pin2pre/cacheFile"
)

var user = 100 + 1
var quan int = 0

// var cache cacheFile.Lru_cache
var t2 time.Time
var t4 time.Time
var t6 time.Time
var t8 time.Time

var tt2 time.Duration
var tt4 time.Duration
var tt6 time.Duration
var tt8 time.Duration

var Counter2 = 0
var Counter4 = 0
var Counter6 = 0
var Counter8 = 0

var l []time.Duration

func main() {
	s := Kittapa.New()
	cacheFile.InitDatabase()
	cacheFile.InitCache()
	// fmt.Println("head", C.head)
	// fmt.Println("last", C.last)
	// .Display()
	// cacheFile.C = cacheFile.Cache_cons(10)
	// db, _ = sql.Open("mysql", "root:62011139@tcp(localhost:3306)/prodj")
	// s.GET("/", abc)
	s.GET("/timeFunction", timeFunction)
	s.GET("/", getCacheFile)            //uye
	s.GET("/products", displayProducts) //all products
	s.GET("/products/:id", productID)
	s.GET("/hitmiss", hitmiss)
	s.GET("/hitmissFile", hitmissFile)
	// cache.ReCache(1)
	s.POST("/products/:id", postPreorder2)
	// s.POST("/products/:id", postPreorder)
	s.Start(":8081")
}

func timeFunction() string {
	// fmt.Println(l)
	fmt.Println("t2 getFile(): ", tt2)
	fmt.Println("t4 displayProducts(): ", tt4)
	fmt.Println("t6 productID(): ", tt6)
	fmt.Println("t8 postPreorder(): ", tt8)
	fmt.Println("counter cache:", cacheFile.Count)
	fmt.Println("L1:", cacheFile.L1)
	fmt.Println("L2:", cacheFile.L2)
	fmt.Println("L3:", cacheFile.L3)
	// fmt.Println("l:", l)
	// fmt.Println("T:", Kittapa.T)
	return ""
}

func productID() string {
	// t := time.Now()
	if Counter6 == 0 {
		t6 = time.Now()
	}
	fmt.Println("ID:", Kittapa.ID)
	a := cacheFile.C.ReCache(Kittapa.ID)
	// t6 = t6 + time.Since(t)
	Counter6++
	if Counter6 == user {
		fmt.Println("------------------------------------")
		fmt.Println("productID():", time.Since(t6))
		fmt.Println("------------------------------------")
		tt6 = time.Since(t6)
		Counter6 = 0
	}

	return a
}

func abc() string {
	return "abc"
}

func hitmiss() string {
	a, _ := json.Marshal(cacheFile.SendMissHit())
	// return "{miss: 1, hit: 2}"
	return string(a)
}

func postPreorder2() string {
	// t := time.Now()
	if Counter8 == 0 {
		t8 = time.Now()
	}
	a := cacheFile.PostPreorder(Kittapa.ID, Kittapa.Result.Quantity)
	// fmt.Println(Kittapa.Result.Quantity == 200)
	// t8 = t8 + time.Since(t)
	Counter8++
	if Counter8 == user {
		fmt.Println("------------------------------------")
		fmt.Println("postPreorder():", time.Since(t8))
		fmt.Println("------------------------------------")
		tt8 = time.Since(t8)
		Counter8 = 0
	}

	return a
}

func getCacheFile() string {
	// t := time.Now()

	if Counter2 == 1 {
		t2 = time.Now()
	}
	Counter2++
	a := cacheFile.Call_cache("index.html")
	// l = append(l, time.Since(t))
	if Counter2 == user {
		fmt.Println("------------------------------------")
		fmt.Println("GetFile():", time.Since(t2))
		fmt.Println("------------------------------------")
		tt2 = time.Since(t2)
		Counter2 = 0
	}

	// t2 = t2 + time.Since(t)
	return a
}

func hitmissFile() string {
	a, _ := json.Marshal(cacheFile.SendMissHitFile())
	// return "{miss: 1, hit: 2}"
	return string(a)
}

func displayProducts() string {
	if Counter4 == 1 {
		t4 = time.Now()
	}
	// fmt.Println(Kittapa.LF)

	a := cacheFile.DisplayAllPro(Kittapa.LF.Limit, Kittapa.LF.Offset)
	Counter4++
	if Counter4 == user {
		fmt.Println("------------------------------------")
		fmt.Println("displayProducts():", time.Since(t4))
		fmt.Println("------------------------------------")
		tt4 = time.Since(t4)
		Counter4 = 0
	}

	return a
}
