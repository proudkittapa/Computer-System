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

var quan int = 0

// var cache cacheFile.Lru_cache
var t2 time.Duration
var t4 time.Duration
var t6 time.Duration
var t8 time.Duration

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
	s.GET("/resetTime", resetTime)
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

func resetTime() string {
	fmt.Println("t2 getFile(): ", t2)
	fmt.Println("t4 displayProducts(): ", t4)
	fmt.Println("t6 productID(): ", t6)
	fmt.Println("t8 postPreorder(): ", t8)
	t2 = 0
	t4 = 0
	t6 = 0
	t8 = 0
	return ""
}

func productID() string {
	t := time.Now()
	fmt.Println("ID:", Kittapa.ID)
	a := cacheFile.C.ReCache(Kittapa.ID)
	t6 = t6 + time.Since(t)
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
	t := time.Now()
	a := cacheFile.PostPreorder(Kittapa.ID, Kittapa.Result.Quantity)
	// fmt.Println(Kittapa.Result.Quantity == 200)
	t8 = t8 + time.Since(t)
	return a
}

func getCacheFile() string {
	t := time.Now()
	a := cacheFile.Call_cache("index.html")
	t2 = t2 + time.Since(t)
	return a
}

func hitmissFile() string {
	a, _ := json.Marshal(cacheFile.SendMissHitFile())
	// return "{miss: 1, hit: 2}"
	return string(a)
}

func displayProducts() string {
	t := time.Now()
	// fmt.Println(Kittapa.LF)
	a := cacheFile.DisplayAllPro(Kittapa.LF.Limit, Kittapa.LF.Offset)
	t4 = t4 + time.Since(t)
	fmt.Println("t4", time.Since(t))
	return a
}
