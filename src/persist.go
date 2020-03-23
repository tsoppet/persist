package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"sync"
)
var lock sync.Mutex
// Marshal is a function that marshals the object into an io.Reader.
// By default, it uses the JSON marshaller then encrypts the data.
var Marshal = func(v interface{}, passphrase string, salt []byte) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	e, err := Encrypt(b, passphrase, salt)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(e), nil
}

// Unmarshal is a function that un-marshals the data from the
// reader into the specified value.
// By default, it decrypts the data then uses the JSON unmarshaller.
var Unmarshal = func(r io.Reader, v interface{}, passphrase string, salt []byte) error {
	b := make([]byte, 256)
	n, err := r.Read(b); if err != nil { return err }
	d, err := Decrypt(b[:n], passphrase, salt); if err != nil { return err }
	return json.Unmarshal(d, v)
}

// Save saves a representation of v to filename
func Save(filename string, v interface{}, passphrase string, salt []byte) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := Marshal(v, passphrase, salt)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, r)
	return err
}

// Load loads the file at path into v.
// Use os.IsNotExist() to see if the returned error is due
// to the file being missing.
func Load(path string, v interface{}, passphrase string, salt []byte) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return Unmarshal(f, v, passphrase, salt)
}