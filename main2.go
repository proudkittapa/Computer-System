package main

import (
	"fmt"
	"pin2pre/Kittapa"
	"pin2pre/cacheFile"
	// "pin2pre/cacheFile"
)

var quan int = 0

var cache cacheFile.Lru_cache

func main() {
	s := Kittapa.New()
	cacheFile.InitCache()
	// cacheFile.C = cacheFile.Cache_cons(10)
	// db, _ = sql.Open("mysql", "root:62011139@tcp(localhost:3306)/prodj")
	// s.GET("/products/:id", abc)
	s.GET("/products/:id", productID)
	// cache.ReCache(1)
	s.Start(":8080")
}

func productID() string {
	fmt.Println("ID:", Kittapa.ID)
	a := cache.ReCache(Kittapa.ID)
	return a
}
