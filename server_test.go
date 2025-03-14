package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func waitForServer(addr string, timeout time.Duration) error {
	start := time.Now()
	for {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			conn.Close()
			return nil
		}
		if time.Since(start) > timeout {
			return fmt.Errorf("server did not start within %s", timeout)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func TestOfficialRedisClient(t *testing.T) {
	listenAddr := ":5001"
	server := NewServer(Config{
		ListenAddr: listenAddr,
	})
	go func() {
		log.Fatal(server.Start())
	}()
	if err := waitForServer(listenAddr, 5*time.Second); err != nil {
		t.Fatal(err)
	}

	// Redis client creation.
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost%s", ":5001"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test key-value operations.
	testCases := map[string]string{
		"foo":  "bar",
		"a":    "gg",
		"your": "mom",
		"step": "dad",
	}
	for key, val := range testCases {
		if err := rdb.Set(context.Background(), key, val, 0).Err(); err != nil {
			t.Fatal(err)
		}
		newVal, err := rdb.Get(context.Background(), key).Result()
		if err != nil {
			t.Fatal(err)
		}
		if newVal != val {
			t.Fatalf("expected %s but got %s", val, newVal)
		}
	}
}

func TestTTL(t *testing.T) {
	listenAddr := ":5001"
	server := NewServer(Config{
		ListenAddr: listenAddr,
	})
	go func() {
		log.Fatal(server.Start())
	}()
	if err := waitForServer(listenAddr, 5*time.Second); err != nil {
		t.Fatal(err)
	}

	// Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost%s", ":5001"),
		Password: "",
		DB:       0,
	})

	// Set a key with TTL of 1 second
	key := "temp"
	val := "expires_soon"
	ttl := 10 * time.Second

	if err := rdb.Set(context.Background(), key, val, ttl).Err(); err != nil {
		t.Fatal(err)
	}

	// Immediately check if key exists
	newVal, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		t.Fatal(err)
	}
	if newVal != val {
		t.Fatalf("expected %s but got %s", val, newVal)
	}

	// Wait for TTL to expire
	for range ttl + 1 {
		_, err := rdb.Get(context.Background(), key).Result()
		if err != nil { // Key expired
			break
		}
		time.Sleep(500 * time.Millisecond) // Check every 500ms
	}

	// Verify key is expired
	_, err = rdb.Get(context.Background(), key).Result()
	if err == nil {
		t.Fatalf("expected key %s to be expired, but it still exists", key)
	}
}
