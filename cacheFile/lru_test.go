package cacheFile

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	// "time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Data struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

type Node struct {
	id    int
	value string
	prev  *node
	next  *node
}

type Lru_cache struct {
	limit int
	mp    map[int]*node
	head  *node
	last  *node
}

type JsonSave struct {
	ProductIDList []int `json:"productIDList"`
	Limit         int   `json:"limit"`
}

func Cache_cons(cap int) Lru_cache {
	return Lru_cache{limit: cap, mp: make(map[int]*Node, cap)}
}

func (list *Lru_cache) Cache(id int) (string, int){
	if node_val, ok := list.mp[id]; ok {
		fmt.Println("-----------HIT-----------")
		list.Move(node_val)
		// fmt.Printf("%T", node_val)
		fmt.Println(node_val.value)
		fmt.Printf("%T\n", node_val.value)
		return node_val.value
	} else {
		fmt.Println("-----------MISS-----------")
		if len(list.mp) >= list.limit {
			rm := list.Remove(list.last)
			delete(list.mp, rm)
		}
		json, _ := Db_query(id)
		node := Node{id: id, value: json}
		list.Add(&node)
		list.mp[id] = &node
		fmt.Println(node.value)
		fmt.Printf("%T\n", node.value)
		return node.value
		// return val.value
	}
}

func (list *Lru_cache) Move(node *Node) {
	if node == list.head {
		return
	}
	list.remove(node)
	list.add(node)
}

func (list *Lru_cache) Remove(node *Node) int {
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

func (list *Lru_cache) Add(node *Node) {
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

func Db_query(id int) (val string, quan int) {

	// fmt.Println("----------MISS----------")

	rows, _ := db.Query("SELECT name, quantity_in_stock, unit_price FROM products WHERE product_id = " + strconv.Itoa(id))

	for rows.Next() {
		var name string
		var quantity int
		var price int
		err := rows.Scan(&name, &quantity, &price)
		CheckErr(err)

		result := Data{Name: name, Quantity: quantity, Price: price}
		byteArray, err := json.Marshal(result)
		CheckErr(err)
		// fmt.Println(len(byteArray))

		val = string(byteArray)
		quan = result.Quantity
		// fmt.Println(val)

		// fmt.Println(byteArray)
		// fmt.Printf("%T", byteArray)
		// fmt.Printf("%T", result.Quantity)

	}

	return val, quan

}

func SaveFile(mp map[int]*Node, lru Lru_cache) {
	var prodIDList []int

	for prodID := 1; prodID <= len(mp); prodID++ {
		prodIDList = append(prodIDList, prodID)
	}

	tempList := jsonSave{ProductIDList: prodIDList, Limit: lru.limit}
	jsonIDList, _ := json.Marshal(tempList)
	_ = ioutil.WriteFile("cacheSave.json", jsonIDList, 0644)
}

// func saveFile_old(mp map[int]*node, lru lru_cache) {
// 	var cache_list []kv

// 	for productID := 1; productID < len(mp); productID++ {
// 		temp_kv := kv{Key: productID, Value: mp[productID].value}
// 		cache_list = append(cache_list, temp_kv)
// 	}

// 	tempCache := jsonCache{Cache: cache_list, Limit: lru.limit}
// 	// fmt.Println(lru.limit)

// 	jsonCacheList, _ := json.Marshal(tempCache)
// 	_ = ioutil.WriteFile("cacheSave.json", jsonCacheList, 0644)

// 	// fmt.Println(string(jsonCacheList))
// 	// fmt.Println(cache_list)
// 	// fmt.Println(tempCache)

// }

func ReadFile() {
	fromFile, err := ioutil.ReadFile("cacheSave.json")
	checkErr(err)

	var temp jsonSave
	err = json.Unmarshal(fromFile, &temp)

	c := cache_cons(temp.Limit)

	t := temp.ProductIDList
	// fmt.Println(t[0])
	for i := 0; i < len(t); i++ {
		fmt.Println(t[i])
		c.cache(t[i])
	}

}

// func readFile_old() lru_cache {
// 	fromFile, err := ioutil.ReadFile("cacheSave.json")
// 	checkErr(err)

// 	var tempStruct jsonSave
// 	err = json.Unmarshal(fromFile, &tempStruct)

// 	c := cache_cons(tempStruct.Limit)

// 	t := tempStruct.Cache
// 	for i := 0; i < len(t); i++ {
// 		for j := 1; j <= len(t); j++ {
// 			node := node{id: j, value: t[i].Value}
// 			c.add(&node)
// 			c.mp[j] = &node
// 			// fmt.Println(c)
// 		}
// 		// fmt.Println(t[i].Value)
// 		// fmt.Printf("%T\n", t[i].Value)
// 	}

// 	fmt.Println(t)
// 	fmt.Printf("%T\n", t)
// 	return c
// }

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

// func main() {
// 	db, _ = sql.Open("mysql", "root:62011212@tcp(127.0.0.1:3306)/prodj")

// 	// var prodIDList []int

// 	// c := cache_cons(1000)

// 	for i := 0; i < 10; i++ {
// 		// fmt.Println(i)
// 		for j := 0; j < 2; j++ {
// 			// start := time.Now()
// 			// c.cache(i)
// 			_, temp := db_query(i)
// 			fmt.Println(temp)
// 			// end := time.Since(start)
// 			// fmt.Printf("%v\n", end)
// 		}

// 	}

	// fmt.Println("last: ", c.last)
	// fmt.Println("head: ", c.head)

	// saveFile(c.mp, c)
	// fmt.Println(c.limit)
	// readFile()

	// fmt.Printf("%T\n", c.mp)

	// c.cache(1)
	// // c.Display()
	// // fmt.Println("last: ", c.last)
	// // fmt.Println("head: ", c.head)
	// c.cache(2)
	// // c.Display()
	// // fmt.Println("last: ", c.last)
	// // fmt.Println("head: ", c.head)
	// c.cache(1)
	// // c.Display()
	// // fmt.Println("last: ", c.last)
	// // fmt.Println("head: ", c.head)
	// c.cache(3)
	// c.cache(4)
	// c.cache(5)
	// c.cache(3)
	// c.cache(6)
	// c.cache(2)

	// saveFile_old(c.mp, c)
	// // c.Display()
	// fmt.Println("last: ", c.last)
	// fmt.Println("head: ", c.head)

// }
