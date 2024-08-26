## merkle-tree 

This package provides following merkle tree operations via storing merkle tree data in an array setting;
Simply initialise with following way and run the operations; 


### Tests 

```bash

go test -cover -v  ./...

```

### Init 

Init function accepts list of leaf hashes in hex string format. 

```go 
    mData := &merkletree.Data{}
	err := mData.Init([]string{
	    "91aaa84a4eff79bee548dc37e77535ea5ef0021ec894ea0af6f656cad6c42a2e",
	    "fd1255bffb6dd11fd9d99c476f6c3c82a17c3ce59365eea67fd26bc7ea521b67",
	    "4e578867a315e03da16380602aa34e028a3e6bdca9471ef3e9d0326639adf68d",
	    "51c1aa3bfaf2b382cf0142f8d52640bf3dbfe501b42a535392c1c1403c9d769d",
	    "9a286af55b24fbc423e8fea1bf778905b079d0cefa0010b87840a2263700d2b6",
    })
```

### Generate Proof 

For given `leafIndex` generates proof for it until to root hash value.
Proofs are returned as an array of `ProofItem` format.

### Verify Leaf 

Given leaf hash and index, verify its integrity via computing corresponding root hash and 
comparing with the real root item. 

### Update Leaf

Updates leaf hash and corresponding path to the merkle root.

### Verify Tree

Verifies merkle tree from leafs to the root.





