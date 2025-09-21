package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

/*
	Pre-requisites for using net/rpc library:
	- all the functions needed to be called remotely need to be a method.
	- all the methods need to be exported (upper case names)>
	- all methods should have two arguments.
	- all method's second argument should be a pointer to a data type that shall be
	  populated as the expected response type.
	- all methods should have error as their return type.
*/

type Item struct {
	Title string
	Body  string
}

type Api struct{}

var database []Item

func (api *Api) GetDB(title string, reply *[]Item) error {
	*reply = database
	return nil
}

func (api *Api) GetItemByName(title string, reply *Item) error {
	for _, item := range database {
		if item.Title == title {
			*reply = item
		}
	}
	return nil
}

func (api *Api) CreateItem(item Item, reply *Item) error {
	database = append(database, item)

	for _, dbItem := range database {
		if dbItem.Title == item.Title {
			*reply = dbItem
		}
	}
	return nil
}

func (api *Api) EditItem(item Item, reply *Item) error {
	for i, dbItem := range database {
		if dbItem.Title == item.Title {
			database[i] = item
			*reply = database[i]
		}
	}
	return nil
}

func (api *Api) DeleteItem(item Item, reply *Item) error {
	for idx, dbItem := range database {
		if dbItem.Title == item.Title && dbItem.Body == item.Body {
			database = append(database[:idx], database[idx+1:]...)
			*reply = dbItem
		}
	}
	return nil
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	api := new(Api)
	err := rpc.RegisterName("Remote", api)
	if err != nil {
		log.Println(fmt.Errorf("failed to register API. err=%s", err))
	}
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Println(fmt.Errorf("listener error. err=%s", err))
	}
	log.Printf("starting a server on port=%s", listener.Addr())
	log.Fatal(http.Serve(listener, nil))
}
