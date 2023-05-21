## Minter Multisend Transactions Generator

This is a simple utility to generate a multisend transaction for the Minter network from a CSV file.

### CSV file format

The file must contain a list of addresses, the amount of coins to send to each address and coins ID. The file must be in CSV format and have the following structure:

```csv
Mx0000000..., 1.234, 0
Mx0000000..., 234, 0
Mx0000000..., 1.22, 2311
Mx0000000..., 10.2, 2311
```

also you can use a file like this if you want to send the same coin to all addresses

```csv
Mx0000000..., 1.234
Mx0000000..., 234
Mx0000000..., 1.22
Mx0000000..., 10.2
```

or if you want to send the same coin to all addresses with a fixed value a file like this is pretty enough

```csv
Mx0000000...
Mx0000000...
Mx0000000...
Mx0000000...
```

### Usage

1. Add package to your project
```bash
go get github.com/daniildulin/minter-multisend
```
```golang

import (
    "github.com/daniildulin/minter-multisend"
)
```

2. Create a new instance of the generator
```golang
txCreator := multisend.NewTxCreator(minterGRPCAddress)
```
or 

```golang
txCreator := multisend.NewTxCreatorFromMnemonic(minterGRPCAddress, "phrase is here")
```

3. Create a new transactions
```golang
// Getting coins id from the file
// File example: 
// Mx0000000..., 1.234, 0
address, mnemonic, txs, err := txCreator.CreateFromFileWithDiffCoins("path/to/file.csv")
```
or

```golang
// File example: 
// Mx0000000..., 1.234
minterCoinID := uint64(0)
address, mnemonic, txs, err := txCreator.CreateFromFile("path/to/file.csv", minterCoinID)
```

or

```golang   
minterCoinID := uint64(0)
payValue := 10

address, mnemonic, txs, err := txCreator.CreateFromFileWithFixedValue("path/to/file.csv", minterCoinID, payValue)
```