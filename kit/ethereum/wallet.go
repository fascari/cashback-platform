package ethereum

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	bip32 "github.com/tyler-smith/go-bip32"
	bip39 "github.com/tyler-smith/go-bip39"
)

// Wallet holds an Ethereum private key and its derived address.
type Wallet struct {
	privateKey *ecdsa.PrivateKey
	address    common.Address
}

// NewFromMnemonic derives an Ethereum wallet from a BIP-39 mnemonic using the given
// derivation path (e.g. "m/44'/60'/0'/0/0").
func NewFromMnemonic(mnemonic, derivationPath string) (*Wallet, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("invalid mnemonic")
	}

	seed := bip39.NewSeed(mnemonic, "")

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("derive master key: %w", err)
	}

	key, err := deriveKey(masterKey, derivationPath)
	if err != nil {
		return nil, fmt.Errorf("derive key at path %s: %w", derivationPath, err)
	}

	privateKey, err := crypto.ToECDSA(key.Key)
	if err != nil {
		return nil, fmt.Errorf("convert to ecdsa: %w", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return &Wallet{privateKey: privateKey, address: address}, nil
}

// Address returns the Ethereum address derived from the private key.
func (w *Wallet) Address() common.Address {
	return w.address
}

// SignTransaction signs a transaction with the wallet's private key.
func (w *Wallet) SignTransaction(tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	signer := types.NewLondonSigner(chainID)
	signed, err := types.SignTx(tx, signer, w.privateKey)
	if err != nil {
		return nil, fmt.Errorf("sign transaction: %w", err)
	}
	return signed, nil
}

// PrivateKey returns the raw ECDSA private key (handle with care).
func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

// deriveKey traverses the BIP-32 path from the master key.
func deriveKey(master *bip32.Key, path string) (*bip32.Key, error) {
	// path format: m/44'/60'/0'/0/0
	key := master
	for seg := range strings.SplitSeq(strings.TrimPrefix(path, "m/"), "/") {
		hardened := strings.HasSuffix(seg, "'")
		seg = strings.TrimSuffix(seg, "'")

		index := new(uint32)
		if _, err := fmt.Sscan(seg, index); err != nil {
			return nil, fmt.Errorf("parse path segment %q: %w", seg, err)
		}
		if hardened {
			*index += bip32.FirstHardenedChild
		}

		child, err := key.NewChildKey(*index)
		if err != nil {
			return nil, fmt.Errorf("derive child at index %d: %w", *index, err)
		}
		key = child
	}
	return key, nil
}
