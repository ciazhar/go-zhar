package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	// Initialize the Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Adjust if your Redis server is on a different host/port
	})

	// Add members to sets
	rdb.SAdd(ctx, "set1", "a", "b", "c")
	rdb.SAdd(ctx, "set2", "b", "c", "d", "e")

	// Get all members of a set
	members, err := rdb.SMembers(ctx, "set1").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Members of set1: %v\n", members)

	// Perform union operation
	union, err := rdb.SUnion(ctx, "set1", "set2").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Union of set1 and set2: %v\n", union)

	// Perform intersection operation
	intersection, err := rdb.SInter(ctx, "set1", "set2").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Intersection of set1 and set2: %v\n", intersection)

	// Perform difference operation (elements in set1 but not in set2)
	difference, err := rdb.SDiff(ctx, "set1", "set2").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Difference of set1 and set2: %v\n", difference)
}
