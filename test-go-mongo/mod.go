package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

func FillHashesCollection(ctx *context.Context, collection *mongo.Collection) {
	for idxHash := 0; idxHash < MAX_HASHES; idxHash++ {
		hash := hashes[idxHash]
		collection.InsertOne(*ctx, bson.M{"hash": hash, "val": int(hash[0])})
		//(*hashMap)[hash] = int32(hash[0])
	}
}

func TestHashesCollection(ctx *context.Context, collection *mongo.Collection) {
	var result struct {
		//Hash   [HASH_SIZE]uint8 `bson:"hash"`
		Val int `bson:"val"`
	}
	for idxHash := 0; idxHash < MAX_HASHES; idxHash++ {
		hash := hashes[idxHash]
		//var result3 bson.M
		err := collection.FindOne(*ctx, bson.M{"hash": hash}).Decode(&result)
		if err == mongo.ErrNoDocuments {
			// Do something when no record was found
			fmt.Println("record does not exist 2")
		} else if err != nil {
			log.Fatal(err)
		} else {
			//if int(result3["val"]) != int(hash[0]) {
			if result.Val != int(hash[0]) {
				//fmt.Println("TestHashes!", result3["val"], int(hash[0]))
				fmt.Println("TestHashes!", result.Val, int(hash[0]))
			}
		}
		//if (*hashMap)[hash] != int32(hash[0]) {
	}
}

func main() {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb:///127.0.0.1:27017"))

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	fmt.Print("Generating hashes...")
	generateRandomHashes()
	fmt.Println("Done")

	collection := client.Database("testing").Collection("numbers")

	collection.DeleteMany(ctx, bson.M{})
	defer collection.DeleteMany(ctx, bson.M{})

	mod := mongo.IndexModel{
		Keys: bson.M{
			"hash": 1, // index in ascending order
		},
		// create UniqueIndex option
		Options: options.Index().SetUnique(true),
	}

	ind, err1 := collection.Indexes().CreateOne(ctx, mod)

	// Check if the CreateOne() method returned any errors
	if err1 != nil {
		fmt.Println("Indexes().CreateOne() ERROR:", err1)
		return //os.Exit(1) // exit in case of error
	} else {
		// API call returns string of the index name
		fmt.Println("CreateOne() index:", ind)
		//fmt.Println("CreateOne() type:", reflect.TypeOf(ind), "\n")
	}

	start := time.Now()
	FillHashesCollection(&ctx, collection)
	fmt.Println("FillHashesCollection", time.Since(start).Milliseconds())

	start = time.Now()
	TestHashesCollection(&ctx, collection)
	fmt.Println("TestHashesCollection", time.Since(start).Milliseconds())

	/*
		//res, err := collection.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
		_, _ = collection.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
		// id := res.InsertedID

		_, _ = collection.InsertOne(ctx, bson.M{"hash": hashes[0], "val": int(hashes[0][0])})

		var result struct {
			Value float64
		}
		filter := bson.D{{"name", "pi"}}
		err = collection.FindOne(ctx, filter).Decode(&result)
		if err == mongo.ErrNoDocuments {
			// Do something when no record was found
			fmt.Println("record does not exist")
		} else if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("record =", result.Value)
		}

		//var result2 struct {
		//	Val int
		//}
		var result3 bson.M
		filter2 := bson.M{"hash": hashes[0]}
		err = collection.FindOne(ctx, filter2).Decode(&result3)
		if err == mongo.ErrNoDocuments {
			// Do something when no record was found
			fmt.Println("record does not exist 2")
		} else if err != nil {
			log.Fatal(err)
		} else {
			//fmt.Println("record2 =", result2.Value)
			fmt.Println("record2 =", result3["val"])
		}
	*/

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}
