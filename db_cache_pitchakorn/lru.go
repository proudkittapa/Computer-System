package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
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
	end   *node
}

func cache_cons(cap int) lru_cache {
	return lru_cache{limit: cap, mp: make(map[int]*node, cap)}
}

func (list *lru_cache) cache(id int) string {
	if val, ok := list.mp[id]; ok {
		fmt.Println("-----------HIT-----------")
		list.move_front(val)
		// fmt.Println(val.value)
		return val.value
	} else {
		fmt.Println("-----------MISS-----------")
		if len(list.mp) >= list.limit {
			rm := list.remove(list.head)
			delete(list.mp, rm)
		}
		json := db_query(id)
		node := node{id: id, value: json}
		list.add(&node)
		list.mp[id] = &node
		return ""
		// return val.value
	}
}

func (l *lru_cache) move_front(node *node) {
	if node == l.end {
		return
	}
	l.remove(node)
	l.add(node)
}

func (l *lru_cache) remove(node *node) int {
	if node == l.end {
		l.end = l.end.prev
	} else if node == l.head {
		l.head = l.head.next
	} else {
		node.prev.next = node.next
		node.next.prev = node.prev
	}
	return node.id
}

func (l *lru_cache) add(node *node) {
	if l.end != nil {
		l.end.next = node
		node.prev = l.end
		node.next = nil
	}
	l.end = node
	if l.head == nil {
		l.head = node
	}
}

func db_query(id int) (val string) {
	db, err := sql.Open("mysql", "root:62011212@tcp(127.0.0.1:3306)/prodj")
	checkErr(err)

	// fmt.Println("----------MISS----------")

	rows, err := db.Query("SELECT name, quantity_in_stock, unit_price FROM products WHERE product_id = " + strconv.Itoa(id))
	checkErr(err)

	for rows.Next() {
		var name string
		var quantity int
		var price int
		err = rows.Scan(&name, &quantity, &price)

		result := data{Name: name, Quantity: quantity, Price: price}
		byteArray, err := json.Marshal(result)
		checkErr(err)
		// fmt.Println(len(byteArray))

		val = string(byteArray)
		// fmt.Println(val)
	}
	return val
}

func main() {

	// var c lru_cache
	// c := lru_cache{limit: 5, mp: make(map[int]*node, 5)}
	// c.cache_cons(5)
	c := cache_cons(5)
	fmt.Println(c.limit)
	c.cache(1)
	c.cache(2)
	c.cache(3)
	c.cache(4)
	c.cache(1)
	c.cache(1)
}
