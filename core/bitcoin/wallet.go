package bitcoin

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"

	"github.com/indexone/niter/core/crypto"
	"github.com/indexone/niter/core/utils"
)

// The address struct should represent all kinds of addresses in Bitcoin.
// For now, it only represents the P2WPKH and P2WSH addresses.
type Address struct {
	commitment []byte
	params     chaincfg.Params
}

func (a *Address) Serialize() (string, error) {
	data := utils.ToIntArray(a.commitment)
	addr, err := SegwitAddrEncode(a.params.Bech32HRPSegwit, 0, data)
	if err != nil {
		return "", err
	}

	return addr, nil
}

func LoadAddress(addr string) (*Address, error) {
	_, data, err := SegwitAddrDecode("bc", addr)
	if err != nil {
		return nil, err
	}
	commitment := utils.ToByteArray(data)
	return &Address{commitment: commitment, params: chaincfg.MainNetParams}, nil
}

func (a *Address) IsWitnessV0KeyHash() bool {
	return len(a.commitment) == 20
}

func (a *Address) IsWitnessV0Script() bool {
	return len(a.commitment) == 32
}

// Wallet represents a Bitcoin wallet.
//
// NOTE: This is the simple version of a wallet, it only contains a private key.
// For a more secure wallet, the HDWallet should be used.
type Wallet struct {
	privateKey *secp256k1.PrivateKey
	params     *chaincfg.Params
}

func LoadWallet(wif string, chainParams *chaincfg.Params) (*Wallet, error) {
	wifKey, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		privateKey: wifKey.PrivKey,
		params:     chainParams,
	}, nil
}

// Returns the Wallet Import Format (WIF) string for the wallet's private key.
func (w *Wallet) WIF() (string, error) {
	wif, err := btcutil.NewWIF(w.privateKey, w.params, true)
	if err != nil {
		return "", err
	}

	return wif.String(), nil
}

func (w *Wallet) getPublicKey() *secp256k1.PublicKey {
	return w.privateKey.PubKey()
}

func (w *Wallet) getPublicKeyCommitment() []byte {
	return crypto.Hash160(w.getPublicKey().SerializeCompressed())
}

func (w *Wallet) Address() *Address {
	publicKeyHash := w.getPublicKeyCommitment()
	return &Address{commitment: publicKeyHash, params: *w.params}
}

// GenerateWallet generates a new wallet with a random private key.
// TODO: Map network to chain config better
func GenerateWallet(chainParams *chaincfg.Params) (*Wallet, error) {
	privateKey, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	return &Wallet{
		privateKey: privateKey,
		params:     chainParams,
	}, nil
}

func GetLockingScriptAddress(commitmentHash []byte, chainParams *chaincfg.Params) (*Address, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_0)
	builder.AddData(commitmentHash)
	builder.AddOp(txscript.OP_EQUAL)

	script, _ := builder.Script()

	scriptCommitmentStr := utils.Hash(script)
	scriptCommitment, _ := hex.DecodeString(scriptCommitmentStr)

	scriptPubKey, err := txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(scriptCommitment).Script()
	if err != nil {
		return nil, err
	}

	return &Address{commitment: scriptPubKey, params: *chainParams}, nil
}