package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"
)

func generateUniqueID() (string, error) {
	/*
		Generate a somewhat unique string.
	*/
	c := 10
	b := make([]byte, c)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	bs := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(bs, time.Now().UnixNano())

	sha := sha256.New()
	sha.Write(b)
	sha.Write(bs)

	return fmt.Sprintf("%x", sha.Sum([]byte{})), nil
}
