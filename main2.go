package main

import (
	"encoding/json"
	"fmt"
	"pin2pre/Kittapa"
	"pin2pre/cacheFile"
	// "pin2pre/transaction"
	// "pin2pre/cacheFile"
)

var quan int = 0

// var cache cacheFile.Lru_cache

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
	s.GET("/", getCacheFile)            //uye
	s.GET("/products", displayProducts) //all products
	s.GET("/products/:id", productID)
	s.GET("/hitmiss", hitmiss)
	s.GET("/hitmissFile", hitmissFile)
	// cache.ReCache(1)
	s.POST("/products/:id", postPreorder2)
	// s.POST("/products/:id", postPreorder)
	s.Start(":8080")
}

func productID() string {
	fmt.Println("ID:", Kittapa.ID)
	a := cacheFile.C.ReCache(Kittapa.ID)
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
	// cacheFile.InitDatabase()
	fmt.Println("ID", Kittapa.ID)
	a := cacheFile.PostPreorder(Kittapa.ID, Kittapa.Result.Quantity)
	return a
}

func getCacheFile() string {
	a := cacheFile.Call_cache("index.html")
	return a
}

func hitmissFile() string {
	a, _ := json.Marshal(cacheFile.SendMissHitFile())
	// return "{miss: 1, hit: 2}"
	return string(a)
}

func displayProducts() string {
	// fmt.Println(Kittapa.LF)
	a := cacheFile.DisplayAllPro(Kittapa.LF.Limit, Kittapa.LF.Offset)
	return a
}
