package main

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123"
	dbname   = "test"
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

func FillHashTable(db *sql.DB) {
	for idxHash := 0; idxHash < MAX_HASHES; idxHash++ {
		hash := hashes[idxHash]
		bitString := hashToBits(&hash)
		sqlStatement := `
		INSERT INTO test (hash, val)
		VALUES ($1, $2)`
		_, err := db.Exec(sqlStatement, bitString, int32(hash[0]))
		if err != nil {
			panic(err)
		}
		//(*hashMap)[hash] = int32(hash[0])
	}
}

func TestHashTable(db *sql.DB) {
	for idxHash := 0; idxHash < MAX_HASHES; idxHash++ {
		hash := hashes[idxHash]
		bitString := hashToBits(&hash)
		sqlStatement := `SELECT val FROM test WHERE hash=$1`
		var resValue int32 = 0
		//_, err := db.Exec(sqlStatement, bitString).Scan(&resValue)
		err := db.QueryRow(sqlStatement, bitString).Scan(&resValue)
		if err != nil {
			panic(err)
		}
		//if (*hashMap)[hash] != int32(hash[0]) {
		if resValue != int32(hash[0]) {
			fmt.Println("TestHashTable!")
		}
	}
}

func hashToBits(hash *[HASH_SIZE]uint8) string {
	//bitString := "B'"
	bitString := "B"
	for _, n := range hash {
		bitString += fmt.Sprintf("%08b", n)
	}
	//bitString += "'::bit(256)"
	return bitString
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	defer db.Exec("TRUNCATE TABLE test")
	_, err = db.Exec("TRUNCATE TABLE test")
	if err != nil {
		panic(err)
	}

	fmt.Print("Generating hashes...")
	generateRandomHashes()
	fmt.Println("Done")

	/*sqlStatement := `
	INSERT INTO test (hash, val)
	VALUES ($1, $2)`
	err = db.QueryRow(sqlStatement, "124::bit(256)", 123)
	if err != nil {
		panic(err)
	}*/

	/*sqlStatement := `
	INSERT INTO test (hash, val)
	VALUES (125::bit(256), 123)`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}*/

	/*var bs [HASH_SIZE]uint8;// = make([]uint8, HASH_SIZE)
	bitString := hashToBits(&bs)
	sqlStatement := `
	INSERT INTO test (hash, val)
	VALUES ($1, $2)`
	_, err = db.Exec(sqlStatement, bitString, "123")
	if err != nil {
		panic(err)
	}

	sqlStatement = `
	SELECT * FROM test WHERE hash=(124::bit(256))`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}*/

	start := time.Now()
	FillHashTable(db)
	fmt.Println("FillHashTable", time.Since(start).Milliseconds())

	start = time.Now()
	TestHashTable(db)
	fmt.Println("TestHashTable", time.Since(start).Milliseconds())
}
