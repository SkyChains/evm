// Code generated
// This file is a generated precompile contract config with stubbed abstract functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package nativeminter

import (
	"math/big"

	"github.com/skychains/evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// NativeCoinMintedEventGasCost is the gas cost of the NativeCoinMinted event.
	// It is the base gas cost + the gas cost of the topics (signature, sender, recipient)
	// and the gas cost of the non-indexed data (32 bytes for amount).
	NativeCoinMintedEventGasCost = contract.LogGas + contract.LogTopicGas*3 + contract.LogDataGas*common.HashLength
)

// PackNativeCoinMintedEvent packs the event into the appropriate arguments for NativeCoinMinted.
// It returns topic hashes and the encoded non-indexed data.
func PackNativeCoinMintedEvent(sender common.Address, recipient common.Address, amount *big.Int) ([]common.Hash, []byte, error) {
	return NativeMinterABI.PackEvent("NativeCoinMinted", sender, recipient, amount)
}

// UnpackNativeCoinMintedEventData attempts to unpack non-indexed [dataBytes].
func UnpackNativeCoinMintedEventData(dataBytes []byte) (*big.Int, error) {
	var eventData = struct {
		Amount *big.Int
	}{}
	err := NativeMinterABI.UnpackIntoInterface(&eventData, "NativeCoinMinted", dataBytes)
	return eventData.Amount, err
}
