package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
)

const HASH_SIZE = 32

//const MAX_HASHES = 32 * 1024 * 1024
const MAX_HASHES = 32 * 1024

var hashes [MAX_HASHES][HASH_SIZE]uint8

func generateRandomHashes() {
	rand.Seed(time.Now().Unix())
	for idxHash := 0; idxHash < MAX_HASHES; idxHash++ {
		for idxByte := 0; idxByte < HASH_SIZE; idxByte += 8 {
			binary.LittleEndian.PutUint64(hashes[idxHash][idxByte:], rand.Uint64())
		}
	}
}

func SetHashes(rdb *redis.Client) {
	for idxHash := 0; idxHash < MAX_HASHES; idxHash++ {
		hash := hashes[idxHash]
		err := rdb.Set(ctx, string(hash[:]), int(hash[0]), 0).Err()
		if err != nil {
			panic(err)
		}
		//(*hashMap)[hash] = int32(hash[0])
	}
}

func TestHashes(rdb *redis.Client) {
	for idxHash := 0; idxHash < MAX_HASHES; idxHash++ {
		hash := hashes[idxHash]
		val, err := rdb.Get(ctx, string(hash[:])).Int()
		if err != nil {
			panic(err)
		}
		//if (*hashMap)[hash] != int32(hash[0]) {
		if val != int(hash[0]) {
			fmt.Println("TestHashes!")
		}
	}
}

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	/*var bs [HASH_SIZE]uint8;// = make([]uint8, HASH_SIZE)
	  //err := rdb.Set(ctx, "key", "value", 0).Err()
	  err := rdb.Set(ctx, string(bs[:]), "value", 0).Err()
	  if err != nil {
	      panic(err)
	  }

	  val, err := rdb.Get(ctx, "key").Result()
	  if err != nil {
	      panic(err)
	  }
	  fmt.Println("key", val)

	  val2, err := rdb.Get(ctx, "key2").Result()
	  if err == redis.Nil {
	      fmt.Println("key2 does not exist")
	  } else if err != nil {
	      panic(err)
	  } else {
	      fmt.Println("key2", val2)
	  }
	  // Output: key value
	  // key2 does not exist*/

	fmt.Print("Generating hashes...")
	generateRandomHashes()
	fmt.Println("Done")

	start := time.Now()
	SetHashes(rdb)
	fmt.Println("SetHashes", time.Since(start).Milliseconds())

	start = time.Now()
	TestHashes(rdb)
	fmt.Println("TestHashes", time.Since(start).Milliseconds())
}
