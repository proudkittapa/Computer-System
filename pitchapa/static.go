// package main

// import (
// 	"net/http"
// 	"os"
// )

// func main(){
// 	dir, _:= os.Getwd()
// 	http.ListenAndServe(";3000", http.FileServer(http.Dir(dir)))
// }

/*
	This is a LRU cache
*/

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

const SIZE = 5 // size of cache

type Node struct {
	Val   [16]byte
	Left  *Node
	Right *Node
}

// double linked list
type Queue struct {
	Head   *Node
	Tail   *Node
	Length int
}

// maps string to node in Queue
type Hash (map[[16]byte]*Node)

// type hash map[int]byte

type Cache struct {
	Queue Queue
	Hash  Hash
}

func NewCache() Cache {
	return Cache{Queue: NewQueue(), Hash: Hash{}}
}

func NewQueue() Queue {
	head := &Node{}
	tail := &Node{}
	head.Right = tail
	tail.Left = head

	return Queue{Head: head, Tail: tail}
}

func (c *Cache) Check(str [16]byte) {
	node := &Node{}
	if val, ok := c.Hash[str]; ok {
		node = c.Remove(val)
	} else {
		node = &Node{Val: str}
	}

	c.Add(node)
	c.Hash[str] = node
}

func (c *Cache) Remove(n *Node) *Node {
	fmt.Printf("remove: %s\n", n.Val)
	left := n.Left
	right := n.Right
	left.Right = right
	right.Left = left
	c.Queue.Length -= 1

	delete(c.Hash, n.Val)
	return n
}

func (c *Cache) Add(n *Node) {
	fmt.Printf("add: %s\n", n.Val)
	tmp := c.Queue.Head.Right
	c.Queue.Head.Right = n
	n.Left = c.Queue.Head
	n.Right = tmp
	tmp.Left = n

	c.Queue.Length++
	if c.Queue.Length > SIZE {
		c.Remove(c.Queue.Tail.Left)
	}
	// if
	// 	print("This value in this cache already!")

}

func (c *Cache) Display() {
	c.Queue.Display()
}

func (q *Queue) Display() {
	node := q.Head.Right
	fmt.Printf("%d - [", q.Length)
	for i := 0; i < q.Length; i++ {
		fmt.Printf("{%s}", node.Val)
		if i < q.Length-1 {
			fmt.Printf(" <--> ")
		}
		node = node.Right
	}
	fmt.Println("]")
}

// func main() {
// 	fmt.Println("TEST CACHE")
// 	cache := NewCache()
// 	for _, word := range [16]byte("shirt") {
// 		cache.Check(word)
// 		cache.Display()q
// 	}

// }

func main() {
	f, err := os.Open("index.html")

	if err != nil {
		fmt.Println("File reading error", err)

	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	chunksize := 512
	reader := bufio.NewReader(f)
	part := make([]byte, chunksize)
	buffer := bytes.NewBuffer(make([]byte, 0))
	var bufferLen int
	for {
		count, err := reader.Read(part)
		if err != nil {
			break
		}
		bufferLen += count
		buffer.Write(part[:count])
	}
	fmt.Println("home")
	return buffer.String()
}
