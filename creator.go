package multisend

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/MinterTeam/minter-go-sdk/v2/api/grpc_client"
	"github.com/MinterTeam/minter-go-sdk/v2/transaction"
	"github.com/MinterTeam/minter-go-sdk/v2/wallet"
	"math/big"
	"os"
	"regexp"
)

type TxCreator struct {
	wallet    *wallet.Wallet
	mntClient *grpc_client.Client
	chainID   transaction.ChainID
}

// CreateFromFileWithDiffCoins creates transactions from csv file with different coins
// Return address, mnemonic, list of transactions and error
func (svc *TxCreator) CreateFromFileWithDiffCoins(filePath string) (address, mnemonic string, txs []string, err error) {
	data, err := svc.readCsvFile(filePath)
	if err != nil {
		return "", "", nil, err
	}

	floatRegex := regexp.MustCompile(`^[-+]?[0-9]*\.[0-9]+([eE][-+]?[0-9]+)?$`)

	recipients := make([]Recipient, 0)
	for i := range data {
		if floatRegex.MatchString(data[i][1]) {
			bigVal, ok := big.NewFloat(0).SetString(data[i][1])
			if !ok {
				return "", "", nil, fmt.Errorf("failed to parse value %s", data[i][1])
			}
			bigVal = bigVal.Mul(bigVal, big.NewFloat(1e18))
			v, _ := bigVal.Int(nil)

			coinID, ok := big.NewInt(0).SetString(data[i][2], 10)
			if !ok {
				return "", "", nil, fmt.Errorf("failed to parse coin id %s", data[i][2])
			}

			recipients = append(recipients, Recipient{
				Address: data[i][0],
				Value:   v,
				CoinID:  coinID.Uint64(),
			})
		} else {
			v, ok := big.NewInt(0).SetString(data[i][1], 10)
			if !ok {
				return "", "", nil, fmt.Errorf("failed to parse value %s", data[i][1])
			}
			v = v.Mul(v, big.NewInt(1e18))

			coinID, ok := big.NewInt(0).SetString(data[i][2], 10)
			if !ok {
				return "", "", nil, fmt.Errorf("failed to parse coin id %s", data[i][2])
			}

			recipients = append(recipients, Recipient{
				Address: data[i][0],
				Value:   v,
				CoinID:  coinID.Uint64(),
			})
		}
	}

	txs, err = svc.CreateTxs(recipients)
	if err != nil {
		return "", "", nil, err
	}
	address = svc.wallet.Address
	mnemonic = svc.wallet.Mnemonic

	return
}

// CreateFromFile creates transactions from csv file
// Return address, mnemonic, list of transactions and error
func (svc *TxCreator) CreateFromFile(filePath string, coinID uint64) (address, mnemonic string, txs []string, err error) {
	data, err := svc.readCsvFile(filePath)
	if err != nil {
		return "", "", nil, err
	}

	floatRegex := regexp.MustCompile(`^[-+]?[0-9]*\.[0-9]+([eE][-+]?[0-9]+)?$`)

	recipients := make([]Recipient, 0)
	for i := range data {
		if floatRegex.MatchString(data[i][1]) {
			bigVal, ok := big.NewFloat(0).SetString(data[i][1])
			if !ok {
				return "", "", nil, fmt.Errorf("failed to parse value %s", data[i][1])
			}
			bigVal = bigVal.Mul(bigVal, big.NewFloat(1e18))
			v, _ := bigVal.Int(nil)

			recipients = append(recipients, Recipient{
				Address: data[i][0],
				Value:   v,
				CoinID:  coinID,
			})
		} else {
			v, ok := big.NewInt(0).SetString(data[i][1], 10)
			if !ok {
				return "", "", nil, fmt.Errorf("failed to parse value %s", data[i][1])
			}
			v = v.Mul(v, big.NewInt(1e18))
			recipients = append(recipients, Recipient{
				Address: data[i][0],
				Value:   v,
				CoinID:  coinID,
			})
		}
	}

	txs, err = svc.CreateTxs(recipients)
	if err != nil {
		return "", "", nil, err
	}
	address = svc.wallet.Address
	mnemonic = svc.wallet.Mnemonic

	return
}

// CreateFromFileWithFixedValue creates transactions from csv file with fixed pay value
// Return list of transactions and error
func (svc *TxCreator) CreateFromFileWithFixedValue(filePath string, coinID uint64, value int64) (address, mnemonic string, txs []string, err error) {
	data, err := svc.readCsvFile(filePath)
	if err != nil {
		return "", "", nil, err
	}

	recipients := make([]Recipient, 0)
	for i := range data {
		recipients = append(recipients, Recipient{
			Address: data[i][0],
			Value:   transaction.BipToPip(big.NewInt(value)),
			CoinID:  coinID,
		})
	}

	txs, err = svc.CreateTxs(recipients)
	if err != nil {
		return "", "", nil, err
	}
	address = svc.wallet.Address
	mnemonic = svc.wallet.Mnemonic

	return
}

// CreateTxs creates transactions from a list of recipients
func (svc *TxCreator) CreateTxs(recipients []Recipient) ([]string, error) {
	i := 0
	count := 0
	recipientsChunks := make(map[int][]Recipient)
	for j := range recipients {
		if recipientsChunks[i] == nil {
			recipientsChunks[i] = make([]Recipient, 0)
		}
		recipientsChunks[i] = append(recipientsChunks[i], recipients[j])
		count++
		if count > 99 {
			count = 0
			i++
		}
	}

	nonce, err := svc.mntClient.Nonce(svc.wallet.Address)
	if err != nil {
		return nil, err
	}

	txs := make([]string, 0)
	for _, chunk := range recipientsChunks {
		txString, err := svc.generateTxFromMap100(chunk, nonce)
		if err != nil {
			return nil, err
		}
		txs = append(txs, txString)
		nonce++
	}

	return txs, nil
}

// CreateTxs creates transaction from a list of recipients less or equal 100
func (svc *TxCreator) generateTxFromMap100(recipients []Recipient, nonce uint64) (string, error) {
	if len(recipients) > 100 {
		return "", errors.New("list is greater than 100")
	}

	sendData := transaction.NewMultisendData()
	for i := range recipients {
		d, err := transaction.NewSendData().
			SetCoin(recipients[i].CoinID).
			SetValue(recipients[i].Value).
			SetTo(recipients[i].Address)
		if err != nil {
			return "", err
		}
		sendData.AddItem(d)
	}

	tx, err := transaction.NewBuilder(svc.chainID).NewTransaction(sendData)
	if err != nil {
		panic(err)
	}
	gp, err := svc.mntClient.MinGasPrice()
	if err != nil {
		return "", err
	}

	tx.SetNonce(nonce).SetGasPrice(uint8(gp.MinGasPrice)).SetGasCoin(0)
	//tx.SetPayload([]byte(""))

	signedTx, err := tx.Sign(svc.wallet.PrivateKey)
	if err != nil {
		return "", err
	}

	txString, err := signedTx.Encode()
	if err != nil {
		return "", err
	}

	return txString, err
}

func (svc *TxCreator) readCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func NewTxCreator(minterGRPCHost string, chainID uint8) *TxCreator {
	w, err := wallet.New()
	if err != nil {
		panic(err)
	}

	nodeAPI, err := grpc_client.New(minterGRPCHost)
	if err != nil {
		panic(err)
	}

	return &TxCreator{
		mntClient: nodeAPI,
		chainID:   transaction.ChainID(chainID),
		wallet:    w,
	}
}

func NewTxCreatorFromMnemonic(minterGRPCHost, mnemonic string, chainID uint8) *TxCreator {
	nodeAPI, err := grpc_client.New(minterGRPCHost)
	if err != nil {
		panic(err)
	}

	seed, err := wallet.Seed(mnemonic)
	if err != nil {
		panic(err)
	}

	w, err := wallet.Create(mnemonic, seed)
	if err != nil {
		panic(err)
	}

	return &TxCreator{
		mntClient: nodeAPI,
		chainID:   transaction.ChainID(chainID),
		wallet:    w,
	}
}
