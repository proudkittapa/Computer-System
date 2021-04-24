package main

import (
	"fmt"
	"pin2pre/Kittapa"
	"pin2pre/cacheFile"
	"pin2pre/transaction"
	// "pin2pre/cacheFile"
)

var quan int = 0

// var cache cacheFile.Lru_cache

func main() {
	s := Kittapa.New()
	cacheFile.InitDatabase()
	cacheFile.C = cacheFile.Cache_cons(10)
	// fmt.Println("head", C.head)
	// fmt.Println("last", C.last)
	// .Display()
	// cacheFile.C = cacheFile.Cache_cons(10)
	// db, _ = sql.Open("mysql", "root:62011139@tcp(localhost:3306)/prodj")
	s.GET("/", abc)
	s.GET("/products/:id", productID)
	s.GET("/hitmiss", hitmiss)
	// cache.ReCache(1)
	s.POST("/products/:id", postPreorder)
	s.Start(":8080")
}

func productID() string {
	fmt.Println("ID:", Kittapa.ID)
	a := cacheFile.C.ReCache(Kittapa.ID)
	return a
}

func postPreorder() string {
	a := transaction.PostPreorder(Kittapa.ID, Kittapa.Result.Quantity)
	return a
}

func abc() string {
	return "abc"
}

func hitmiss() string {
	// a, _ := json.Marshal(cacheFile.SendHitMiss())
	return "{miss: 1, hit: 2}"
}
