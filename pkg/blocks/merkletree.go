package blocks

import (
	"crypto/sha256"

	"github.com/GhostNet-Dev/GhostNet-Core/libs/container"
	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
)

func CreateMerkleRoot(txHashList [][]byte) []byte {
	depth := GetDepth(len(txHashList))
	if depth == 0 {
		return make([]byte, gbytes.HashSize)
	}
	hashList := container.NewQueue()
	for _, hash := range txHashList {
		hashList.Push(hash)
	}

	hash := sha256.New()

	for hashList.Count > 1 {
		size := hashList.Count
		for i := uint32(0); i < size; i += 2 {
			var left []byte
			if size-(i+1) > 0 {
				left = hashList.Pop().([]byte)
				right := hashList.Pop().([]byte)
				left = append(left, right...)
			} else {
				left = hashList.Pop().([]byte)
				left = append(left, left...)
			}
			hash.Write(left)
			hashList.Push(hash.Sum(nil))
			hash.Reset()
		}
	}

	return hashList.Pop().([]byte)
}

func GetDepth(count int) uint32 {
	return uint32(count)
}
