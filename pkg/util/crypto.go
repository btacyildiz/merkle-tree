package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func MerkleHash(data string) (string, error) {
	byteArr, err := HexToByte(data)
	if err != nil {
		return "", fmt.Errorf("unaable to convert to byteArray %w", err)
	}
	h := sha256.New()
	n, err := h.Write(byteArr)
	if err != nil {
		return "", fmt.Errorf("unable to write data to hash function %w", err)
	}
	if n != len(byteArr) {
		return "", fmt.Errorf("written %d bytes but expected %d bytes", n, len(data))
	}
	return ByteToHex(h.Sum(nil)), nil
}

func ByteToHex(data []byte) string {
	return hex.EncodeToString(data)
}

func HexToByte(hexStr string) ([]byte, error) {
	return hex.DecodeString(hexStr)
}
