package client

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestNewClient1(t *testing.T) {
	c, err := New("localhost:5432")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	fmt.Println("SET => foo")
	if err := c.Set(context.Background(), "foo", "bar"); err != nil {
		log.Fatal(err)
	}

	val, err := c.Get(context.Background(), "foo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET =>", val)
}


func TestNewClient(t *testing.T) {
	c, err := New("localhost:5432")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	for i := 0; i < 1; i++ {
		fmt.Println("SET => foo")
		if err := c.Set(context.Background(), "foo", "bar"); err != nil {
			log.Fatal(err)
		}

		fmt.Println("SET => psk")
		if err := c.Set(context.Background(), "psk", "ssk"); err != nil {
			log.Fatal(err)
		}

		val, err := c.Get(context.Background(), "foo")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("GET =>", val)
	}
}