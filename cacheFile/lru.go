package cacheFile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
)

var (
//db *sql.DB
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
	prev  *Node
	next  *Node
}

type Lru_cache struct {
	limit int
	mp    map[int]*Node
	head  *Node
	last  *Node
}

type JsonSave struct {
	ProductIDList []int `json:"productIDList"`
	Limit         int   `json:"limit"`
}

func Cache_cons(cap int) Lru_cache {
	// db, _ = sql.Open("mysql", "root:62011139@tcp(localhost:3306)/prodj")
	// // db.SetMaxIdleConns(200000)
	// db.SetMaxOpenConns(200000)
	return Lru_cache{limit: cap, mp: make(map[int]*Node, cap)}
}

func (list *Lru_cache) GetCache(id int) string {
	if node_val, ok := list.mp[id]; ok {
		fmt.Println("-----------HIT-----------")
		list.Move(node_val)
		// fmt.Println(val.value)
		return node_val.value
	} else {
		fmt.Println("-----------MISS-----------")
		return ""
	}

}

func (list *Lru_cache) Set(id int, val Data) {

	byteArray, err := json.Marshal(val)
	CheckErr(err)
	// fmt.Println(len(byteArray))

	temp := string(byteArray)

	if prod, ok := list.mp[id]; ok || len(list.mp) >= list.limit {
		fmt.Println("if 1")
		if len(list.mp) >= list.limit {
			fmt.Println("if 2")
			rm := list.Remove(list.last)
			delete(list.mp, rm)

		} else if _, ok := list.mp[id]; ok {
			fmt.Println("else")
			// list.Move(prod)
			rm := list.Remove(prod)
			delete(list.mp, rm)

		}
	}

	node := Node{id: id, value: temp}
	list.AddNode(&node)
	list.mp[id] = &node
	// fmt.Println(node.value)
}

// mind -> cache MISS
// mind -> Query
// mind -> set query
// func set

func (list *Lru_cache) Move(node *Node) {
	if node == list.head {
		return
	}
	list.Remove(node)
	list.AddNode(node)
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

func (list *Lru_cache) AddNode(node *Node) {
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

// ref https://www.tutorialfor.com/blog-259822.htm
//     https://medium.com/@fazlulkabir94/lru-cache-golang-implementation-92b7bafb76f0

// func Db_query(id int) (val string) {

// 	// fmt.Println("----------MISS----------")
// 	fmt.Println(id)
// 	rows := db.QueryRow("SELECT name, quantity_in_stock, unit_price FROM products WHERE product_id = " + strconv.Itoa(id))

// 	var name string
// 	var quantity int
// 	var price int
// 	err := rows.Scan(&name, &quantity, &price)
// 	CheckErr(err)

// 	// fmt.Println(val)

// 	return val
// }

func SaveFile(mp map[int]*Node, lru Lru_cache) {
	var prodList []int
	t := lru.mp

	keys := make([]int, 0, len(t))
	for k := range t {
		keys = append(keys, k)
	}

	for i := 0; i < len(t); i++ {
		prodList = append(prodList, keys[i])
	}

	fmt.Println(prodList)

	tempList := JsonSave{ProductIDList: prodList, Limit: lru.limit}
	jsonIDList, _ := json.Marshal(tempList)
	_ = ioutil.WriteFile("cacheSave.json", jsonIDList, 0644)
}

// ref https://stackoverflow.com/questions/47898327/properly-create-a-json-file-and-read-from-it

// func ReadFile() Lru_cache {

// 	fromFile, err := ioutil.ReadFile("cacheSave.json")
// 	CheckErr(err)

// 	var temp JsonSave
// 	err = json.Unmarshal(fromFile, &temp)

// 	c := Cache_cons(temp.Limit)

// 	t := temp.ProductIDList
// 	// fmt.Println(t[0])
// 	for i := 0; i < len(t); i++ {
// 		fmt.Println(t[i])
// 		c.Set(t[i])
// 	}

// 	return c
// }

// ref https://tutorialedge.net/golang/parsing-json-with-golang/

func (l *Lru_cache) Display() {
	node := l.last
	for node != nil {
		fmt.Printf("%+v ->", node.id)
		node = node.prev
	}
	fmt.Println()
}

// func Display(node *Node) {
// 	for node != nil {
// 		fmt.Printf("%v ->", node.id)
// 		node = node.next
// 	}
// 	fmt.Println()
// }

// func main() {
// 	// db, _ = sql.Open("mysql", "root:62011212@tcp(127.0.0.1:3306)/prodj")

// 	// defer profile.Start(profile.MemProfile).Stop()

// 	c := Cache_cons(10)

// 	temp := Data{Name: "pune", Quantity: 20, Price: 100}
// 	temp2 := Data{Name: "pune2", Quantity: 20, Price: 100}
// 	temp3 := Data{Name: "pune3", Quantity: 20, Price: 100}

// 	c.Set(1, temp)
// 	c.Set(1, temp2)
// 	c.Set(1, temp3)
// 	// c.Set(1, temp3)
// 	c.Display()
// 	fmt.Println(c.GetCache(1))
// 	fmt.Println("last: ", c.last)
// 	fmt.Println("head: ", c.head)

// 	// for i := 0; i < 10; i++ {
// 	// 	for j := 0; j < 2; j++ {
// 	// 		// start := time.Now()
// 	// 		c.Set()
// 	// end := time.Since(start)
// 	// fmt.Printf("%v\n", end)

// 	// t := c.cache(i)
// 	// fmt.Println(t)
// 	// fmt.Printf("%T\n", t)
// 	// 	}
// 	// }
// }
