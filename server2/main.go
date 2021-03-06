package main

import (
	"fmt"
	"pin2pre/Kittapa"
	"pin2pre/final2"
	// "pin2pre/transaction"
	// "pin2pre/cacheFile"
)

var quan int = 0

// var cache cacheFile.Lru_cache

func main() {
	s := Kittapa.New()
	final2.InitDatabase()
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
	fmt.Println("ID:", Kittapa.ID)
	a := final2.ReCache(Kittapa.ID)
	return a
}

func abc() string {
	return "abc"
}

func postPreorder2() string {
	// cacheFile.InitDatabase()
	fmt.Println("ID", Kittapa.ID)
	a := final2.PostPreorder(Kittapa.ID, Kittapa.Result.Quantity)
	return a
}

func getCacheFile() string {
	a := final2.Call_cache("index.html")
	return a
}

func displayProducts() string {
	// fmt.Println(Kittapa.LF)
	a := final2.Display_pro()
	return a
}
