package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/api/gmail/v1"
)

// CreateHash produces 32 byte hash using simple MD5 hash
func CreateHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// EncryptFile encrypts an email message using given hash key. The file name is the message ID
func EncryptFile(msg *gmail.Message, hashStr string) {
	// f, _ := os.Create(filename)
	f, err := os.OpenFile(msg.Id, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)		
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to save encrypted item (Email ID - %s): %v", msg.Id, err)
	}

	data, err := json.Marshal(msg)
    if err != nil {
        fmt.Println(err)
        return
    }
	
	// json.NewEncoder(f).Encode(msg)
	f.Write(encrypt(data, hashStr))
}

// DecryptFile decrypts the contents of an encrypted file and stores into 
// a new file (of same name suffixed with "_d")
func DecryptFile(filename string, hashStr string) {
	data, _ := ioutil.ReadFile(filename)
	decryptedFileName := filename + "_d"
	f, err := os.OpenFile(decryptedFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)		
	defer f.Close()

	if err != nil {
		log.Fatalf("Unable to save decrypted item (Email ID - %s): %v", filename, err)
	}

	f.Write(decrypt(data, hashStr))
}

// encrypt a byte array with a given hashkey
func encrypt(data []byte, hashStr string) []byte {
	block, _ := aes.NewCipher([]byte(hashStr))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

// decrypt a byte array with a given hashkey
func decrypt(data []byte, hashStr string) []byte {
	key := []byte(hashStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return plaintext
}
