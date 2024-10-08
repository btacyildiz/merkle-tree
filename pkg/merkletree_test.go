package merkletree

import (
	"crypto/sha256"
	"gotest.tools/v3/assert"
	"math/rand"
	"merkletree/pkg/util"
	"strconv"
	"testing"
)

func Test_InitHash(t *testing.T) {
	type testCase struct {
		name          string
		inputHashes   []string
		expectedCount int
	}

	testCases := []testCase{
		{
			name:          "even number of elements",
			inputHashes:   generateRandomHashes(6),
			expectedCount: merkleTreeElementCount(6),
		},
		{
			name:          "odd number of elements",
			inputHashes:   generateRandomHashes(5),
			expectedCount: merkleTreeElementCount(5),
		},
		{
			name:          "odd number more items",
			inputHashes:   generateRandomHashes(345),
			expectedCount: merkleTreeElementCount(345),
		},
		{
			name:          "even number more items",
			inputHashes:   generateRandomHashes(320),
			expectedCount: merkleTreeElementCount(320),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			merkleTree := &Data{}
			err := merkleTree.Init(tt.inputHashes)
			assert.NilError(t, err)
			assert.Equal(t, tt.expectedCount, len(merkleTree.hashArray))
			isValid, err := merkleTree.VerifyTree()
			assert.NilError(t, err)
			assert.Equal(t, true, isValid)
		})
	}
}

func TestGenerateProof(t *testing.T) {
	type testCase struct {
		name          string
		leafHashCount int
		leafIndex     int
		expected      []ProofItem
		expectError   bool
	}
	testCases := []testCase{
		{
			name:          "invalid index",
			leafHashCount: 3,
			leafIndex:     3,
			expected:      nil,
			expectError:   true,
		},
		{
			name:          "empty",
			leafHashCount: 0,
			leafIndex:     3,
			expected:      nil,
			expectError:   true,
		},
		{
			name:          "odd number - mid",
			leafHashCount: 5,
			leafIndex:     3,
			expected:      []ProofItem{{true, 8}, {true, 3}, {false, 2}},
		},
		{
			name:          "odd number - end",
			leafHashCount: 5,
			leafIndex:     4,
			expected:      []ProofItem{{false, 10}, {false, 5}, {true, 1}},
		},
		{
			name:          "odd number - start",
			leafHashCount: 5,
			leafIndex:     0,
			expected:      []ProofItem{{false, 7}, {false, 4}, {false, 2}},
		},
		{
			name:          "even number - mid",
			leafHashCount: 6,
			leafIndex:     4,
			expected:      []ProofItem{{false, 11}, {false, 5}, {true, 1}},
		},
		{
			name:          "even number - end",
			leafHashCount: 6,
			leafIndex:     5,
			expected:      []ProofItem{{true, 10}, {false, 5}, {true, 1}},
		},
		{
			name:          "even number - start",
			leafHashCount: 5,
			leafIndex:     0,
			expected:      []ProofItem{{false, 7}, {false, 4}, {false, 2}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			merkleTree := Data{}
			err := merkleTree.Init(generateRandomHashes(tt.leafHashCount))
			assert.NilError(t, err)
			proof, err := merkleTree.GenerateProof(tt.leafIndex)
			assert.Equal(t, tt.expectError, err != nil)
			assert.DeepEqual(t, tt.expected, proof)
		})

	}
}

func TestData_UpdateLeaf(t *testing.T) {
	tests := []struct {
		name            string
		leafCount       int
		updateLeafIndex int
		updateLeafHash  string
		wantErr         bool
	}{
		{
			name:            "out of bound",
			leafCount:       3,
			updateLeafIndex: 3,
			updateLeafHash:  generateRandomHashes(1)[0],
			wantErr:         true,
		},
		{
			name:            "empty merkle",
			leafCount:       0,
			updateLeafIndex: 3,
			updateLeafHash:  generateRandomHashes(1)[0],
			wantErr:         true,
		},
		{
			name:            "update mid",
			leafCount:       5,
			updateLeafIndex: 3,
			updateLeafHash:  generateRandomHashes(1)[0],
		},
		{
			name:            "update end",
			leafCount:       5,
			updateLeafIndex: 4,
			updateLeafHash:  generateRandomHashes(1)[0],
		},
		{
			name:            "update start",
			leafCount:       5,
			updateLeafIndex: 0,
			updateLeafHash:  generateRandomHashes(1)[0],
		},
		{
			name:            "big tree",
			leafCount:       145,
			updateLeafIndex: 100,
			updateLeafHash:  generateRandomHashes(1)[0],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merkleTree := Data{}
			err := merkleTree.Init(generateRandomHashes(tt.leafCount))
			assert.NilError(t, err)
			err = merkleTree.UpdateLeaf(tt.updateLeafIndex, tt.updateLeafHash)
			assert.Equal(t, tt.wantErr, err != nil)
			isValid, err := merkleTree.VerifyTree()
			assert.NilError(t, err)
			assert.Equal(t, true, isValid)
		})
	}
}

func TestData_VerifyLeaf(t *testing.T) {

	tests := []struct {
		name            string
		leafCount       int
		verifyLeafIndex int
		// verifyLeafHash or verifyLeafHashIndex can be provided
		verifyLeafHash      string
		verifyLeafHashIndex int
		isValid             bool
		wantErr             bool
	}{
		{
			name:            "out of bound",
			leafCount:       3,
			verifyLeafIndex: -1,
			wantErr:         true,
		},
		{
			name:           "empty merkle",
			leafCount:      0,
			verifyLeafHash: "bla",
			wantErr:        true,
		},
		{
			name:                "should verify",
			leafCount:           5,
			verifyLeafIndex:     4,
			verifyLeafHashIndex: 4,
			isValid:             true,
		},
		{
			name:            "should be invalid",
			leafCount:       5,
			verifyLeafIndex: 1,
			verifyLeafHash:  generateRandomHashes(1)[0],
			isValid:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merkleTree := Data{}
			err := merkleTree.Init(generateRandomHashes(tt.leafCount))
			var hashToVerify string
			if tt.verifyLeafHash == "" {
				hashToVerify = merkleTree.hashArray[merkleTree.getHashListIndex(tt.verifyLeafHashIndex)]
			} else {
				hashToVerify = tt.verifyLeafHash
			}
			assert.NilError(t, err)
			isValid, err := merkleTree.VerifyLeaf(tt.verifyLeafIndex, hashToVerify)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func generateRandomHashes(n int) []string {
	var generatedHashes []string
	for i := 0; i < n; i++ {
		bytes := sha256.Sum256([]byte(strconv.FormatInt(int64(rand.Uint64()), 10)))
		generatedHashes = append(generatedHashes, util.ByteToHex(bytes[:]))
	}
	return generatedHashes
}
