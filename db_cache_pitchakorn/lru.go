package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type data struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

type node struct {
	id    int
	value string
	prev  *node
	next  *node
}

type lru_cache struct {
	limit int
	mp    map[int]*node
	head  *node
	last  *node
}

type kv struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}

type jsonCache struct {
	Cache []kv `json:"cache"`
}

func cache_cons(cap int) lru_cache {
	return lru_cache{limit: cap, mp: make(map[int]*node, cap)}
}

func (list *lru_cache) cache(id int) string {
	if node_val, ok := list.mp[id]; ok {
		fmt.Println("-----------HIT-----------")
		list.move(node_val)
		// fmt.Println(val.value)
		return node_val.value
	} else {
		fmt.Println("-----------MISS-----------")
		if len(list.mp) >= list.limit {
			rm := list.remove(list.last)
			delete(list.mp, rm)
		}
		json := db_query(id)
		node := node{id: id, value: json}
		list.add(&node)
		list.mp[id] = &node
		// fmt.Println(node.value)
		return node.value
		// return val.value
	}
}

func (list *lru_cache) move(node *node) {
	if node == list.head {
		return
	}
	list.remove(node)
	list.add(node)
}

func (list *lru_cache) remove(node *node) int {
	if node == list.last {
		// fmt.Println("con 1")
		list.last = list.last.prev
	} else if node == list.head {
		// fmt.Println("con 2")
		list.head = list.head.next
	} else {
		// fmt.Println("con 3")
		node.prev.next = node.next
		node.next.prev = node.prev
	}
	return node.id
}

func (list *lru_cache) add(node *node) {
	if list.head != nil {
		list.head.prev = node
		node.next = list.head
		node.prev = nil
	}
	list.head = node
	if list.last == nil {
		list.last = node
	}
}

func db_query(id int) (val string) {

	// fmt.Println("----------MISS----------")

	rows, _ := db.Query("SELECT name, quantity_in_stock, unit_price FROM products WHERE product_id = " + strconv.Itoa(id))

	for rows.Next() {
		var name string
		var quantity int
		var price int
		err := rows.Scan(&name, &quantity, &price)
		checkErr(err)

		result := data{Name: name, Quantity: quantity, Price: price}
		byteArray, err := json.Marshal(result)
		checkErr(err)
		// fmt.Println(len(byteArray))

		val = string(byteArray)
		// fmt.Println(val)
	}

	return val
}

func saveFile(mp map[int]*node) {
	var cache_list []kv

	for productID := 1; productID < len(mp); productID++ {
		temp_kv := kv{Key: productID, Value: mp[productID].value}
		cache_list = append(cache_list, temp_kv)
	}

	tempCache := jsonCache{Cache: cache_list}

	jsonCacheList, _ := json.Marshal(tempCache)
	_ = ioutil.WriteFile("cacheSave.json", jsonCacheList, 0644)

	// fmt.Println(string(jsonCacheList))
	// fmt.Println(cache_list)
	// fmt.Println(tempCache)

}

func readFile() {
	fromFile, err := ioutil.ReadFile("cacheSave.json")
	checkErr(err)

	var tempStruct jsonCache
	err = json.Unmarshal(fromFile, &tempStruct)

	t := tempStruct.Cache
	// for i := 0; i < len(t); i++ {

	// }

	fmt.Println(t[0])
	fmt.Printf("%T\n", t[0])
}

// func (l *lru_cache) Display() {
// 	node := l.head
// 	for node != nil {
// 		fmt.Printf("%+v ->", node.id)
// 		node = node.next
// 	}
// 	fmt.Println()
// }

// func Display(node *node) {
// 	for node != nil {
// 		fmt.Printf("%v ->", node.id)
// 		node = node.next
// 	}
// 	fmt.Println()
// }

func main() {
	db, _ = sql.Open("mysql", "root:62011212@tcp(127.0.0.1:3306)/prodj")

	// defer profile.Start(profile.MemProfile).Stop()

	c := cache_cons(10)

	for i := 0; i < 10; i++ {
		for j := 0; j < 2; j++ {
			start := time.Now()
			c.cache(i)
			end := time.Since(start)
			fmt.Printf("%v\n", end)
			// fmt.Println(c.cache(i))
		}
	}

	// saveFile(c.mp)
	readFile()

	// fmt.Printf("%T\n", c.mp)

	// c.cache(1)
	// // c.Display()
	// fmt.Println("last: ", c.last)
	// fmt.Println("head: ", c.head)
	// c.cache(2)
	// // c.Display()
	// fmt.Println("last: ", c.last)
	// fmt.Println("head: ", c.head)
	// c.cache(1)
	// // c.Display()
	// fmt.Println("last: ", c.last)
	// fmt.Println("head: ", c.head)
	// c.cache(3)
	// c.cache(4)
	// c.cache(5)
	// c.cache(6)
	// fmt.Println("last: ", c.last)
	// fmt.Println("head: ", c.head)

}
