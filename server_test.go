package main

import (
	"context"
	"fmt"
	"log"
	"redis-go/client"
	"sync"
	"testing"
	"time"

)

func TestServerWithMultiClients(t *testing.T) {
	server := NewServer(Config{})
	go func() {
		log.Fatal(server.Start())
	}()
	
	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)
	for i := 0; i < nClients; i++ {
		go func(it int) {
			defer wg.Done()
			c, err := client.New("localhost:5432")
			if err != nil {
				log.Fatal(err)
			}
			defer c.Close()
			key := fmt.Sprintf("client_foo_%d", i)
			value := fmt.Sprintf("client_bar_%d", i)
		
			if err := c.Set(context.Background(), key, value); err != nil {
				log.Fatal(err)
			}
	
			val, err := c.Get(context.Background(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("client %d got this val back => %s\n", i, val)

		}(i)
	}
	wg.Wait()
	time.Sleep(time.Second) // give time to close the connections before terminating

	if len(server.peers) != 0 {
		t.Fatalf("expected 0 peers but got %d", len(server.peers))
	}
}

func TestFooBar(t *testing.T) {
	in := map[string]string {
		"first": "1",
		"second": "2",
	}
	out := writeRespMap(in)
	fmt.Println(out)
}