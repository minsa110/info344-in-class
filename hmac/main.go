package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
)

const usage = `
usage:
	hmac sign|verify <key> <value>
`

func main() {
	if len(os.Args) < 4 ||
		(os.Args[1] != "sign" && os.Args[1] != "verify") {
		fmt.Println(usage)
		os.Exit(1)
	}

	cmd := os.Args[1]
	key := os.Args[2]
	value := os.Args[3]

	switch cmd {
	case "sign":
		v := []byte(value)
		h := hmac.New(sha256.New, []byte(key))
		h.Write(v)
		sig := h.Sum(nil)

		buf := make([]byte, len(v)+len(sig))
		copy(buf, v) // target <-- source
		copy(buf[len(v):], sig)
		fmt.Println(base64.URLEncoding.EncodeToString(buf)) // base64 encoding

	case "verify":
		buf, err := base64.URLEncoding.DecodeString(value)
		if err != nil {
			fmt.Printf("Error decoding: %v\n", err)
			os.Exit(1) // DO NOT USE THIS OR LOG.FATAL ON WEB SERVER ***
		}

		v := buf[:len(buf)-sha256.Size]   // first part of buffer
		sig := buf[len(buf)-sha256.Size:] // second part of buffer

		// how do we know if the signature is valid?
		h := hmac.New(sha256.New, []byte(key)) // has to be the same as "sign"
		h.Write(v)
		sig2 := h.Sum(nil)
		if hmac.Equal(sig, sig2) {
			fmt.Println("Signature is valid!")
		} else {
			fmt.Println("DANGER DANGER INVALID SIGNATURE!!")
		}
	}

}
