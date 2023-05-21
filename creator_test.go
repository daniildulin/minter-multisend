package multisend

import (
	"github.com/joho/godotenv"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const filePath = "./tmp/data.csv"

func TestCreateFromFileWithDiffCoins(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	txCreator := NewTxCreator(os.Getenv("MINTER_GRPC_HOST"), 1)

	address, mnemonic, txs, err := txCreator.CreateFromFileWithDiffCoins(filePath)

	assert.NoError(t, err)
	assert.NotEmpty(t, address)
	assert.NotEmpty(t, mnemonic)
	assert.NotEmpty(t, txs)
}

func TestCreateFromFile(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}
	// Создаем экземпляр TxCreator
	txCreator := NewTxCreator(os.Getenv("MINTER_GRPC_HOST"), 1)

	coinID := uint64(0)

	address, mnemonic, txs, err := txCreator.CreateFromFile(filePath, coinID)

	assert.NoError(t, err)
	assert.NotEmpty(t, address)
	assert.NotEmpty(t, mnemonic)
	assert.NotEmpty(t, txs)
}

func TestCreateFromFileWithFixedValue(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	txCreator := NewTxCreator(os.Getenv("MINTER_GRPC_HOST"), 1)

	coinID := uint64(0)
	value := int64(100)

	address, mnemonic, txs, err := txCreator.CreateFromFileWithFixedValue(filePath, coinID, value)

	assert.NoError(t, err)
	assert.NotEmpty(t, address)
	assert.NotEmpty(t, mnemonic)
	assert.NotEmpty(t, txs)
}

func TestCreateTxs(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}
	txCreator := NewTxCreator(os.Getenv("MINTER_GRPC_HOST"), 1)

	recipients := []Recipient{
		{
			Address: "Mx885f69ec994bfbcd05630b9e751b69799ce1f8bd",
			Value:   big.NewInt(100),
			CoinID:  1,
		},
		{
			Address: "Mxaeac266a4533cb0b4255ea2922f997353a18b2e8",
			Value:   big.NewInt(200),
			CoinID:  1,
		},
	}

	txs, err := txCreator.CreateTxs(recipients)

	assert.NoError(t, err)
	assert.NotEmpty(t, txs)
}

func TestGenerateTxFromMap100(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}
	txCreator := NewTxCreator(os.Getenv("MINTER_GRPC_HOST"), 1)

	recipients := []Recipient{
		{
			Address: "Mx885f69ec994bfbcd05630b9e751b69799ce1f8bd",
			Value:   big.NewInt(100),
			CoinID:  1,
		},
		{
			Address: "Mxaeac266a4533cb0b4255ea2922f997353a18b2e8",
			Value:   big.NewInt(200),
			CoinID:  1,
		},
	}

	nonce := uint64(1)

	txString, err := txCreator.generateTxFromMap100(recipients, nonce)

	assert.NoError(t, err)
	assert.NotEmpty(t, txString)
}

func TestReadCsvFile(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	txCreator := NewTxCreator(os.Getenv("MINTER_GRPC_HOST"), 1)

	records, err := txCreator.readCsvFile(filePath)

	assert.NoError(t, err)
	assert.NotEmpty(t, records)
}

func TestNewTxCreator(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}
	txCreator := NewTxCreator(os.Getenv("MINTER_GRPC_HOST"), 1)

	assert.NotNil(t, txCreator)
	assert.NotNil(t, txCreator.mntClient)
	assert.NotNil(t, txCreator.wallet)
}

func TestNewTxCreatorFromMnemonic(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	txCreator := NewTxCreatorFromMnemonic(os.Getenv("MINTER_GRPC_HOST"), os.Getenv("MINTER_MNEMONIC"), 1)

	assert.NotNil(t, txCreator)
	assert.NotNil(t, txCreator.mntClient)
	assert.NotNil(t, txCreator.wallet)
}
