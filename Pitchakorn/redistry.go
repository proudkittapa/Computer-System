package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

var rdb *redis.Client

type data struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var ctx = context.Background()

func cache(id string) {
	val, err := rdb.Get(ctx, id).Result()
	// switch {
	// case err == redis.Nil:
	// 	fmt.Println("----------MISS----------")
	// 	db_query(id)
	// case err != nil:
	// 	fmt.Println("failed", err)
	// case val == "":
	// 	fmt.Println("value is empty")
	// }

	if err == redis.Nil {
		fmt.Println("----------MISS----------")
		db_query(id)
	} else if val == "" {
		fmt.Println("value is empty")
	} else {
		fmt.Println("----------HIT----------")
		fmt.Printf(val)
	}

}

func db_query(id string) {
	db, err := sql.Open("mysql", "root:62011212@tcp(127.0.0.1:3306)/prodj")
	checkError(err)
	var byteArray []byte
	rows, err := db.Query("SELECT name, quantity_in_stock, unit_price FROM products WHERE product_id = " + id)
	checkError(err)

	for rows.Next() {
		var name string
		var quantity int
		var price int
		err = rows.Scan(&name, &quantity, &price)

		result := data{Name: name, Quantity: quantity, Price: price}
		byteArray, err = json.Marshal(result)
		checkError(err)
		// fmt.Println(string(byteArray))
	}
	add_redis(id, byteArray)
}

func add_redis(id string, jso []byte) {

	err := rdb.Set(ctx, id, string(jso), 0).Err()

	val, err := rdb.Get(ctx, id).Result()
	checkError(err)

	fmt.Println(val)

}

// func main() {
// 	rdb = redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})

// 	cache("6")

// }

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var intSlice [10]int
	sum := 0

	for i := 0; i < 10; i++ {
		for j := 0; j < 2; j++ {
			start := time.Now()
			num := i + 1
			inc := strconv.Itoa(num)
			cache(inc)

			end := time.Since(start)
			fmt.Printf("\n%v\n", end)
			intSlice[i] = int(end)
		}
	}

	for k := 0; k < 10; k++ {
		sum += (intSlice[k])
	}

	avg := (float64(sum)) / (float64(10))
	fmt.Println("Average = ", avg, "µs")

}
