package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type Item struct {
	Title string
	Body  string
}

func main() {
	var reply Item
	var db []Item

	client, err := rpc.DialHTTP("tcp", ":4000")
	if err != nil {
		log.Fatalf("client error. err=%s", err)
	}

	client.Call("Remote.GetDB", "", &db)
	fmt.Printf("initial db: %+v\n", db)

	a := Item{"first", "item one"}
	b := Item{"second", "item two"}
	c := Item{"three", "item three"}

	client.Call("Remote.CreateItem", a, &reply)
	client.Call("Remote.CreateItem", b, &reply)
	client.Call("Remote.CreateItem", c, &reply)

	client.Call("Remote.GetDB", "", &db)
	fmt.Printf("db after creating items: %+v\n", db)

	// delete
	client.Call("Remote.DeleteItem", b, &reply)

	client.Call("Remote.GetDB", "", &db)
	fmt.Printf("db after deleting b item: %+v\n", db)

	// Edit an item
	var newItemReply Item
	newItem := Item{
		Title: "first",
		Body:  "this is a new one",
	}
	client.Call("Remote.EditItem", newItem, &newItemReply)
	client.Call("Remote.GetDB", "", &db)
	fmt.Printf("db after editing c item: %+v\n", db)

	// Get items by name
	var x, y, z Item
	client.Call("Remote.GetItemByName", "first", &x)
	client.Call("Remote.GetItemByName", "three", &y)
	client.Call("Remote.GetItemByName", "four", &z)

	fmt.Println("Retrieved items:", x, y, z)
	fmt.Printf("db after getting items: %+v\n", db)
}

// func fodder() {
// 	fmt.Printf("initial db: %+v\n", database)
//
// 	a := Item{"first", "item one"}
// 	b := Item{"second", "item two"}
// 	c := Item{"three", "item three"}
//
// 	api.CreateItem(a)
// 	api.CreateItem(b)
// 	api.CreateItem(c)
// 	fmt.Printf("db after creating items: %+v\n", database)
//
// 	DeleteItem(b)
// 	fmt.Printf("db after deleting b item: %+v\n", database)
//
// 	EditItem("three", Item{
// 		Title: "four",
// 		Body:  "item three replaced",
// 	})
// 	fmt.Printf("db after editing c item: %+v\n", database)
//
// 	x := GetItemByName("first")
// 	y := GetItemByName("three")
// 	z := GetItemByName("four")
//
// 	fmt.Println(x, y, z)
// }
