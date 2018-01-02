package ethclient

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

type KeyGenerator struct {
	keyFolder string
}

func NewKeyGenerator(folder string) *KeyGenerator {
	return &KeyGenerator{folder}
}

func (kg *KeyGenerator) GenerateKey(password string) (common.Address, error) {
	ks := keystore.NewKeyStore(kg.keyFolder, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		return common.Address{}, err
	} else {
		address := account.Address
		if err := os.Rename(account.URL.Path, kg.keyFilePath(address.Hex())); err != nil {
			panic("rename key file failed:" + err.Error())
		}
		return address, nil
	}
}

func (kg *KeyGenerator) keyFilePath(address string) string {
	keyFileName := strings.Join([]string{address, "key"}, ".")
	return filepath.Join(kg.keyFolder, keyFileName)
}

func (kg *KeyGenerator) GetAccount(address common.Address, password string) (*Account, error) {
	return NewAccount(kg.keyFilePath(address.Hex()), password)
}

func (kg *KeyGenerator) ListAddress() []common.Address {
	keyFiles, err := ioutil.ReadDir(kg.keyFolder)
	if err != nil {
		return nil
	}

	var addresses []common.Address
	for _, f := range keyFiles {
		if f.IsDir() == false {
			fname := f.Name()
			if strings.HasSuffix(fname, ".key") {
				addresses = append(addresses, common.HexToAddress(strings.TrimSuffix(fname, ".key")))
			}
		}
	}
	return addresses
}
