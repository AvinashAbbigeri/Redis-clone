package client

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNewClient1(t *testing.T) {

	c, err := New("localhost:5001")
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Set(context.TODO(), "foo", "1"); err != nil {
		log.Fatal(err)
	}
	val, err := c.Get(context.TODO(), "foo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(val)
}

func TestNewClient(t *testing.T) {
	c, err := New("localhost:5001")
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second)
	for i := 0; i < 10; i++ {
		fmt.Println("SET =>", fmt.Sprintf("foo_%d", i))
		if err := c.Set(context.TODO(), fmt.Sprintf("foo_%d", i), fmt.Sprintf("bar_%d", i)); err != nil {
			log.Fatal(err)
		}
		val, err := c.Get(context.TODO(), fmt.Sprintf("foo_%d", i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("GET =>", val)
	}
}
