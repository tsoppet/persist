package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type obj struct {
	Name string
	Number int
	When time.Time
}

// Passphrase for encrypting/decrypting the data, this should be stored as a secret
var passphrase = "000102030405060708090A0B0C0D0E0FF0E0D0C0B0A090807060504030201000"

// Main converts saves encrypts and writes an object to file.tmp then reads the encrypted data from
// file.tmp, decrypts it and marshals it into json format and prints it out to confirm it is the same
// data that was originally written.
func main() {
	o := &obj{
		Name:   "Mat",
		Number: 47,
		When:   time.Now(),
	}

	salt, _ := RandomSalt(32)		// should be a large random number stored in a secret

	// save object `o` to file.tmp
	if err := Save("./file.tmp", o, passphrase, salt); err != nil {
		log.Fatalln(err)
	}
	// load it back
	var o2 obj
	if err := Load("./file.tmp", &o2, passphrase, salt); err != nil {
		log.Fatalln(err)
	}

	data, _ := json.Marshal(o2)
	fmt.Printf("%s\n", data)
	// o and o2 are now the same
	// check out file.tmp - you'll see the encrypted data
}
