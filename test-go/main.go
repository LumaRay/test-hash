package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const HASH_SIZE = 32

// const MAX_HASHES = 8 * 1024 * 1024
// const ROOT_NODE_ADDRESS_BITS = 3*8 + 2
// const CHILD_NODE_ADDRESS_BITS = 1*8 - 6

const MAX_HASHES = 32 * 1024 * 1024
const ROOT_NODE_ADDRESS_BITS = 3*8 + 4
const CHILD_NODE_ADDRESS_BITS = 1*8 - 4

// const MAX_HASHES = 64 * 1024 * 1024
// const ROOT_NODE_ADDRESS_BITS = 3*8 + 4
// const CHILD_NODE_ADDRESS_BITS = 1*8 - 4

type HashNode[T any] struct {
	mpObject     *T
	mpHash       *[HASH_SIZE]uint8
	mChildNodes  []*HashNode[T]
	mAddressBits int
}

type HashTree[T any] struct {
	mRootNode HashNode[T]
}

var hashes [MAX_HASHES][HASH_SIZE]uint8

func (this *HashNode[T]) Init() {
	this.mChildNodes = nil
	this.mAddressBits = CHILD_NODE_ADDRESS_BITS
}

func HashNodeCreate[T any](pHash *[HASH_SIZE]uint8, pObject *T) *HashNode[T] {
	pNode := new(HashNode[T])
	pNode.Init()
	pNode.mpObject = pObject
	pNode.mpHash = pHash
	return pNode
}

// func testEq(a, b []uint8) bool {
// 	if len(a) != len(b) {
// 		return false
// 	}
// 	for i := range a {
// 		if a[i] != b[i] {
// 			return false
// 		}
// 	}
// 	return true
// }

func (this *HashNode[T]) SetObject(pHash *[HASH_SIZE]uint8, pObject *T, iByteOffset int) {
	if this.mChildNodes == nil || len(this.mChildNodes) == 0 {
		//this->mChildNodes = calloc((uint64_t)1 << this->mAddressBits, sizeof(struct HashNode));
		this.mChildNodes = make([]*HashNode[T], uint64(1)<<this.mAddressBits)
	}
	if this.mpHash != nil && this.mpObject != nil {
		pExistingHash := this.mpHash
		iExistingHashAddress := ^(math.MaxUint32 << this.mAddressBits) & binary.LittleEndian.Uint32(pExistingHash[iByteOffset:])
		this.mChildNodes[iExistingHashAddress] = HashNodeCreate(this.mpHash, this.mpObject)
		this.mpHash = nil
		this.mpObject = nil
	}
	iHashAddress := ^(math.MaxUint32 << this.mAddressBits) & binary.LittleEndian.Uint32(pHash[iByteOffset:])
	ppChildNode := &this.mChildNodes[iHashAddress]
	if *ppChildNode == nil || iByteOffset == HASH_SIZE-1 {
		*ppChildNode = HashNodeCreate(pHash, pObject)
	} else {
		(**ppChildNode).SetObject(pHash, pObject, iByteOffset+int(math.Ceil(float64(this.mAddressBits)/8.0)))
	}
}

func (this *HashNode[T]) GetObject(pHash *[HASH_SIZE]uint8, iByteOffset int) *T {
	if this.mpHash != nil && this.mpObject != nil && bytes.Equal(this.mpHash[iByteOffset:], pHash[iByteOffset:]) {
		return this.mpObject
	}
	iHashAddress := ^(math.MaxUint32 << this.mAddressBits) & binary.LittleEndian.Uint32(pHash[iByteOffset:])
	ppChildNode := &this.mChildNodes[iHashAddress]
	if *ppChildNode == nil || iByteOffset+int(math.Ceil(float64(this.mAddressBits)/8.0)) >= HASH_SIZE-1 {
		fmt.Println("Hash overflow!")
		return nil
	}
	return (**ppChildNode).GetObject(pHash, iByteOffset+int(math.Ceil(float64(this.mAddressBits)/8.0)))
}

func (this *HashTree[T]) Init() {
	this.mRootNode.mAddressBits = ROOT_NODE_ADDRESS_BITS
}

func (this *HashTree[T]) Add(hash *[HASH_SIZE]uint8, object *T) {
	hashCopy := *hash
	objectCopy := *object
	this.mRootNode.SetObject(&hashCopy, &objectCopy, 0)
}

func (this *HashTree[T]) Get(hash *[HASH_SIZE]uint8) *T {
	return this.mRootNode.GetObject(hash, 0)
}

func generateRandomHashes() {
	rand.Seed(time.Now().Unix())
	for idxHash := 0; idxHash < MAX_HASHES; idxHash++ {
		for idxByte := 0; idxByte < HASH_SIZE; idxByte += 8 {
			binary.LittleEndian.PutUint64(hashes[idxHash][idxByte:], rand.Uint64())
		}
	}
}

func MakeHashTree(pHashTree *HashTree[int32]) {
	for idxHash := 0; idxHash < MAX_HASHES; idxHash++ {
		hash := hashes[idxHash]
		in := int32(hash[0])
		pHashTree.Add(&hash, &in)
	}
}

func TestHashTree(pHashTree *HashTree[int32]) {
	for idxHash := 0; idxHash < MAX_HASHES; idxHash++ {
		hash := hashes[idxHash]
		//b := pHashTree.Get(&hash)
		//hashByte := hash[0]
		if *pHashTree.Get(&hash) != int32(hash[0]) {
			fmt.Println("TestHashTree!")
		}
	}
}

func MakeStdUMap(hashMap *map[[HASH_SIZE]uint8]int32) {
	for idxHash := 0; idxHash < MAX_HASHES; idxHash++ {
		hash := hashes[idxHash]
		(*hashMap)[hash] = int32(hash[0])
	}
}

func TestStdUMap(hashMap *map[[HASH_SIZE]uint8]int32) {
	for idxHash := 0; idxHash < MAX_HASHES; idxHash++ {
		hash := hashes[idxHash]
		if (*hashMap)[hash] != int32(hash[0]) {
			fmt.Println("TestStdUMap!")
		}
	}
}

func TestStdUMapWrongGet(hashMap *map[[HASH_SIZE]uint8]int32) {
	var hash [HASH_SIZE]uint8
	rand.Seed(time.Now().Unix())
	for idxHash := 0; idxHash < MAX_HASHES; idxHash++ {
		for idxByte := 0; idxByte < HASH_SIZE; idxByte += 8 {
			binary.LittleEndian.PutUint64(hash[idxByte:], rand.Uint64())
		}
		_, ok := (*hashMap)[hash]
		if ok {
			fmt.Println("TestStdUMapWrongGet!")
		}
	}
}

func main() {
	fmt.Print("Generating hashes...")
	generateRandomHashes()
	fmt.Println("Done")

	// var hashTree HashTree[int32]
	// hashTree.Init()

	// start := time.Now()
	// MakeHashTree(&hashTree)
	// fmt.Println("MakeHashTree", time.Since(start).Milliseconds())

	// start = time.Now()
	// TestHashTree(&hashTree)
	// fmt.Println("TestHashTree", time.Since(start).Milliseconds())

	hashMap := make(map[[HASH_SIZE]uint8]int32)

	start := time.Now()
	MakeStdUMap(&hashMap)
	fmt.Println("MakeStdUMap", time.Since(start).Milliseconds())

	start = time.Now()
	TestStdUMap(&hashMap)
	fmt.Println("TestStdUMap", time.Since(start).Milliseconds())

	TestStdUMapWrongGet(&hashMap)
}
