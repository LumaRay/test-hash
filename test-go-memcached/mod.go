package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

const HASH_SIZE = 32

//const MAX_HASHES = 32 * 1024 * 1024
const MAX_HASHES = 32 * 1024

var hashes [MAX_HASHES][HASH_SIZE]uint8

func generateRandomHashes() {
	rand.Seed(time.Now().Unix())
	for idxHash := 0; idxHash < MAX_HASHES; idxHash++ {
		for idxByte := 0; idxByte < HASH_SIZE; idxByte += 8 {
			binary.LittleEndian.PutUint64(hashes[idxHash][idxByte:], (rand.Uint64()|0x3030303030303030)&0x7E7E7E7E7E7E7E7E)
		}
	}
}

func SetHashes(mc *memcache.Client) {
	for idxHash := 0; idxHash < MAX_HASHES; idxHash++ {
		hash := hashes[idxHash]
		mc.Set(&memcache.Item{Key: string(hash[:]), Value: []byte{hash[0]}})
		//(*hashMap)[hash] = int32(hash[0])
	}
}

func TestHashes(mc *memcache.Client) {
	for idxHash := 0; idxHash < MAX_HASHES; idxHash++ {
		hash := hashes[idxHash]
		it, _ := mc.Get(string(hash[:]))
		//if (*hashMap)[hash] != int32(hash[0]) {
		if it.Value[0] != hash[0] {
			fmt.Println("TestHashes!")
		}
	}
}

func main() {
	//mc := memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212")
	mc := memcache.New("127.0.0.1:11211")
	mc.DeleteAll()
	mc.FlushAll()
	//mc.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})
	//it, err := mc.Get("foo")
	//fmt.Println(it, err)

	fmt.Print("Generating hashes...")
	generateRandomHashes()
	fmt.Println("Done")

	//mc.Set(&memcache.Item{Key: string(hashes[0][:]), Value: []byte {hashes[0][0]}})
	//it, err := mc.Get(string(hashes[0][:]))
	//fmt.Println(it, err)

	start := time.Now()
	SetHashes(mc)
	fmt.Println("SetHashes", time.Since(start).Milliseconds())

	start = time.Now()
	TestHashes(mc)
	fmt.Println("TestHashes", time.Since(start).Milliseconds())
}
