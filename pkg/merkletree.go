package merkletree

import (
	"fmt"
	"merkletree/pkg/util"
)

type Data struct {
	// Array represents merkle tree as a slice in the format of binary tree.
	// First element always represents root hash value.
	// Consecutive values are constructed binary tree format which leads to
	// individual hashes of individual hash of the elements.
	// e.g. hash(abcd), hash(ab), hash(cd), hash(a), hash(b), hash(c) hash(d)
	hashArray []string
	itemCount int
}

// ProofItem is used to construct merkle proof
// during verification indexed hash value will be merged with current value
// according to the IsLeft value;
// left -> hash(left + current) | right -> hash(current + right)
type ProofItem struct {
	IsLeft bool
	Index  int
}

type merklePath struct {
	ProofItem
	ParentIndex int
}

// Init initialises merkle tree, with leaf hashes
// leaf hash values should be in hex string format
func (d *Data) Init(hashes []string) error {
	d.itemCount = len(hashes)
	d.hashArray = []string{}
	d.hashArray = append(d.hashArray, hashes...)
	var currentArr = hashes
	for len(currentArr) > 1 {
		var localArr []string
		var curIndex = 0
		for curIndex < len(currentArr) {
			var hash string
			var err error
			if curIndex+1 < len(currentArr) {
				hash, err = util.MerkleHash(currentArr[curIndex] + currentArr[curIndex+1])
			} else {
				// if there is odd number of elements we concat and hash with self
				hash, err = util.MerkleHash(currentArr[curIndex] + currentArr[curIndex])
			}
			if err != nil {
				return fmt.Errorf("error while init - calculating merkle hash %w", err)
			}
			localArr = append(localArr, hash)
			curIndex += 2
		}
		d.hashArray = append(localArr, d.hashArray...)
		currentArr = localArr
	}
	return nil
}

// GenerateProof iterates over path to the root, and generates proof items
func (d *Data) GenerateProof(leafIndex int) ([]ProofItem, error) {
	var proofIndexes []ProofItem
	err := d.generatePath(leafIndex, func(merklePath merklePath) error {
		proofIndexes = append(proofIndexes, ProofItem{
			IsLeft: merklePath.IsLeft,
			Index:  merklePath.Index,
		})
		return nil
	})
	return proofIndexes, err
}

func (d *Data) VerifyLeaf(leafIndex int, leafHash string) (bool, error) {
	if err := d.checkIndex(leafIndex); err != nil {
		return false, fmt.Errorf("error verifying leaf - checking index - err: %w", err)
	}

	//leafIndexInMerkle := d.getHashListIndex(leafIndex)
	var currentHash = leafHash
	var err error
	var calculatedRootHash string

	err = d.generatePath(leafIndex, func(merklePath merklePath) error {
		if merklePath.IsLeft {
			currentHash = d.hashArray[merklePath.Index] + currentHash
		} else {
			currentHash = currentHash + d.hashArray[merklePath.Index]
		}
		currentHash, err = util.MerkleHash(currentHash)
		if err != nil {
			return fmt.Errorf("verify leaf operation failed for current id %d  - error generating merkle hash %w",
				merklePath.ParentIndex, err)
		}
		if merklePath.ParentIndex == 0 {
			calculatedRootHash = currentHash
		}
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("verify leaf operation is failed err: %w", err)
	}

	if calculatedRootHash != d.hashArray[0] {
		fmt.Println("reCalculated root hash does not match")
		return false, nil
	}
	return true, nil
}

// UpdateLeaf iterates over path to the root, updates the path
func (d *Data) UpdateLeaf(leafIndex int, newHash string) error {

	if err := d.checkIndex(leafIndex); err != nil {
		return fmt.Errorf("error updating leaf - checking index - err: %w", err)
	}

	leafIndexInMerkle := d.getHashListIndex(leafIndex)
	d.hashArray[leafIndexInMerkle] = newHash
	var err error
	var mergedHashes = newHash
	return d.generatePath(leafIndex, func(merklePath merklePath) error {
		if merklePath.IsLeft {
			mergedHashes = d.hashArray[merklePath.Index] + mergedHashes
		} else {
			mergedHashes = mergedHashes + d.hashArray[merklePath.Index]
		}
		mergedHashes, err = util.MerkleHash(mergedHashes)
		if err != nil {
			return fmt.Errorf("update leaf operation for parent index %d  - error generating merkle hash %w",
				merklePath.ParentIndex, err)
		}
		d.hashArray[merklePath.ParentIndex] = mergedHashes
		return nil
	})
}

// VerifyTree verifies whole structure
func (d *Data) VerifyTree() (bool, error) {
	var index = len(d.hashArray) - 1
	var levelLength = d.itemCount
	var start = len(d.hashArray) - d.itemCount
	var end = len(d.hashArray) - 1
	var parentIndex = start - 1
	for levelLength > 0 && index-1 >= 0 {
		var concat string
		if levelLength%2 == 0 {
			concat = d.hashArray[index-1] + d.hashArray[index]
			index -= 2
		} else {
			if index == end {
				concat = d.hashArray[index] + d.hashArray[index]
				index--
			} else {
				concat = d.hashArray[index-1] + d.hashArray[index]
				index -= 2
			}
		}
		calcHash, err := util.MerkleHash(concat)
		if err != nil {
			return false, fmt.Errorf("unable to create hash to verify %w", err)
		}
		if d.hashArray[parentIndex] != calcHash {
			fmt.Printf("comparing hash at index %d is failed\nexpected hash: %s\nfounded hash: %s\n", parentIndex, calcHash, d.hashArray[parentIndex])
			return false, nil
		}
		parentIndex--
		if index < start {
			if levelLength%2 == 0 {
				levelLength = levelLength / 2
			} else {
				levelLength = levelLength/2 + 1
			}
			start, end = start-levelLength, start-1
		}
	}
	return true, nil
}

// generatePath iterate over merkle tree until root is founded
func (d *Data) generatePath(leafIndex int, handler func(merklePath merklePath) error) error {

	if err := d.checkIndex(leafIndex); err != nil {
		return fmt.Errorf("error generating path - checking index - err: %w", err)
	}

	leafIndexInMerkle := d.getHashListIndex(leafIndex)

	if d.itemCount == 1 {
		// for single item, it will be hashed with self to verify
		return handler(merklePath{ProofItem{
			IsLeft: false,
			Index:  1,
		}, 0})
	}
	var curIndex = leafIndexInMerkle
	var start = len(d.hashArray) - d.itemCount
	var end = len(d.hashArray) - 1
	var layerLength = d.itemCount
	for curIndex > 0 {
		isLeft, siblingIndex := getSibling(curIndex, start, end, layerLength)

		if layerLength%2 == 0 {
			layerLength = layerLength / 2
		} else {
			layerLength = layerLength/2 + 1
		}
		tempStart := start
		start, end = start-layerLength, start-1
		curIndex = start + (curIndex-tempStart)/2
		err := handler(merklePath{
			ProofItem:   ProofItem{IsLeft: isLeft, Index: siblingIndex},
			ParentIndex: curIndex,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Data) checkIndex(index int) error {
	if d.itemCount == 0 {
		return fmt.Errorf("merkle tree is empty")
	}
	if index < 0 || index > d.itemCount-1 {
		return fmt.Errorf("given Index: %d should be within 0-%d range", index, d.itemCount-1)
	}
	return nil
}

func (d *Data) getHashListIndex(leafIndex int) int {
	return len(d.hashArray) - d.itemCount + leafIndex
}

func merkleTreeElementCount(n int) int {
	var total = n
	var current = n
	for current != 1 {
		if current%2 == 0 {
			current = current / 2
		} else {
			current = current/2 + 1
		}
		total += current
	}
	return total
}

func getSibling(curIndex, start, end, layerLength int) (bool, int) {
	if layerLength%2 == 0 {
		if (curIndex-start)%2 == 0 {
			return false, curIndex + 1
		}
		return true, curIndex - 1
	}
	if curIndex == end {
		// hash with self
		return false, curIndex
	}
	if (curIndex-start)%2 == 0 {
		return false, curIndex + 1
	}
	return true, curIndex - 1
}
