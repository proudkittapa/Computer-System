package main

import (
	"fmt"
	"pin2pre/Kittapa"
	"pin2pre/final1"

	// "pin2pre/transaction"
	// "pin2pre/cacheFile"
	"time"
)

var quan int = 0

// var cache cacheFile.Lru_cache
var t2 time.Duration
var t4 time.Duration
var t6 time.Duration
var t8 time.Duration

func main() {
	s := Kittapa.New()
	final1.InitDatabase()
	// final1.InitCache()
	// fmt.Println("head", C.head)
	// fmt.Println("last", C.last)
	// .Display()
	// cacheFile.C = cacheFile.Cache_cons(10)
	// db, _ = sql.Open("mysql", "root:62011139@tcp(localhost:3306)/prodj")
	// s.GET("/", abc)
	s.GET("/", getCacheFile)            //uye
	s.GET("/products", displayProducts) //all products
	s.GET("/products/:id", productID)
	// s.GET("/hitmiss", hitmiss)
	// s.GET("/hitmissFile", hitmissFile)
	// cache.ReCache(1)
	s.POST("/products/:id", postPreorder2)
	// s.POST("/products/:id", postPreorder)
	s.Start(":8081")
}

func productID() string {
	t := time.Now()
	fmt.Println("ID:", Kittapa.ID)
	a := final1.ReCache(Kittapa.ID)
	t6 = t6 + time.Since(t)
	fmt.Println("t6 productID(): ", t6)
	return a
}

func abc() string {
	return "abc"
}

func postPreorder2() string {
	// cacheFile.InitDatabase()
	t := time.Now()
	fmt.Println("ID", Kittapa.ID)
	a := final1.PostPreorder(Kittapa.ID, Kittapa.Result.Quantity)
	t8 = t8 + time.Since(t)
	fmt.Println("t8 postPreorder(): ", t8)
	return a
}

func getCacheFile() string {
	t := time.Now()
	a := final1.GetFile()
	t2 = t2 + time.Since(t)
	fmt.Println("t2 getFile(): ", t2)
	return a
}

func displayProducts() string {
	t := time.Now()
	// fmt.Println(Kittapa.LF)
	a := final1.Display_pro()
	t4 = t4 + time.Since(t)
	fmt.Println("t4 displayProducts(): ", t4)
	return a
}
