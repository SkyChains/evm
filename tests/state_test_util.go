// (c) 2021-2024, Lux Partners Limited.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package tests

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/skychains/evm/core"
	"github.com/skychains/evm/core/rawdb"
	"github.com/skychains/evm/core/state"
	"github.com/skychains/evm/core/state/snapshot"
	"github.com/skychains/evm/core/types"
	"github.com/skychains/evm/core/vm"
	"github.com/skychains/evm/ethdb"
	"github.com/skychains/evm/params"
	"github.com/skychains/evm/trie"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

// StateTest checks transaction processing without block context.
// See https://github.com/ethereum/EIPs/issues/176 for the test format specification.
type StateTest struct {
	json stJSON
}

// StateSubtest selects a specific configuration of a General State Test.
type StateSubtest struct {
	Fork  string
	Index int
}

func (t *StateTest) UnmarshalJSON(in []byte) error {
	return json.Unmarshal(in, &t.json)
}

type stJSON struct {
	Env  stEnv                    `json:"env"`
	Pre  core.GenesisAlloc        `json:"pre"`
	Tx   stTransaction            `json:"transaction"`
	Out  hexutil.Bytes            `json:"out"`
	Post map[string][]stPostState `json:"post"`
}

type stPostState struct {
	Root            common.UnprefixedHash `json:"hash"`
	Logs            common.UnprefixedHash `json:"logs"`
	TxBytes         hexutil.Bytes         `json:"txbytes"`
	ExpectException string                `json:"expectException"`
	Indexes         struct {
		Data  int `json:"data"`
		Gas   int `json:"gas"`
		Value int `json:"value"`
	}
}

//go:generate go run github.com/fjl/gencodec -type stEnv -field-override stEnvMarshaling -out gen_stenv.go
type stEnv struct {
	Coinbase   common.Address `json:"currentCoinbase"   gencodec:"required"`
	Difficulty *big.Int       `json:"currentDifficulty" gencodec:"required"`
	Random     *big.Int       `json:"currentRandom"     gencodec:"optional"`
	GasLimit   uint64         `json:"currentGasLimit"   gencodec:"required"`
	Number     uint64         `json:"currentNumber"     gencodec:"required"`
	Timestamp  uint64         `json:"currentTimestamp"  gencodec:"required"`
	BaseFee    *big.Int       `json:"currentBaseFee"  gencodec:"optional"`
}

//go:generate go run github.com/fjl/gencodec -type stTransaction -field-override stTransactionMarshaling -out gen_sttransaction.go
type stTransaction struct {
	GasPrice             *big.Int            `json:"gasPrice"`
	MaxFeePerGas         *big.Int            `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *big.Int            `json:"maxPriorityFeePerGas"`
	Nonce                uint64              `json:"nonce"`
	To                   string              `json:"to"`
	Data                 []string            `json:"data"`
	AccessLists          []*types.AccessList `json:"accessLists,omitempty"`
	GasLimit             []uint64            `json:"gasLimit"`
	Value                []string            `json:"value"`
	PrivateKey           []byte              `json:"secretKey"`
}

// GetChainConfig takes a fork definition and returns a chain config.
// The fork definition can be
// - a plain forkname, e.g. `Byzantium`,
// - a fork basename, and a list of EIPs to enable; e.g. `Byzantium+1884+1283`.
func GetChainConfig(forkString string) (baseConfig *params.ChainConfig, eips []int, err error) {
	var (
		splitForks            = strings.Split(forkString, "+")
		ok                    bool
		baseName, eipsStrings = splitForks[0], splitForks[1:]
	)

	// NOTE: this is added to support mapping geth fork names to
	// evm fork names.
	forkAliases := map[string]string{
		"Berlin":       "Pre-EVM",
		"London":       "EVM",
		"ArrowGlacier": "EVM",
		"GrayGlacier":  "EVM",
		"Merge":        "EVM",
	}
	if alias, ok := forkAliases[baseName]; ok {
		baseName = alias
	}

	if baseConfig, ok = Forks[baseName]; !ok {
		return nil, nil, UnsupportedForkError{baseName}
	}
	for _, eip := range eipsStrings {
		if eipNum, err := strconv.Atoi(eip); err != nil {
			return nil, nil, fmt.Errorf("syntax error, invalid eip number %v", eipNum)
		} else {
			if !vm.ValidEip(eipNum) {
				return nil, nil, fmt.Errorf("syntax error, invalid eip number %v", eipNum)
			}
			eips = append(eips, eipNum)
		}
	}
	return baseConfig, eips, nil
}

// Subtests returns all valid subtests of the test.
func (t *StateTest) Subtests() []StateSubtest {
	var sub []StateSubtest
	for fork, pss := range t.json.Post {
		for i := range pss {
			sub = append(sub, StateSubtest{fork, i})
		}
	}
	return sub
}

// checkError checks if the error returned by the state transition matches any expected error.
// A failing expectation returns a wrapped version of the original error, if any,
// or a new error detailing the failing expectation.
// This function does not return or modify the original error, it only evaluates and returns expectations for the error.
func (t *StateTest) checkError(subtest StateSubtest, err error) error {
	expectedError := t.json.Post[subtest.Fork][subtest.Index].ExpectException
	if err == nil && expectedError == "" {
		return nil
	}
	if err == nil && expectedError != "" {
		return fmt.Errorf("expected error %q, got no error", expectedError)
	}
	if err != nil && expectedError == "" {
		return fmt.Errorf("unexpected error: %w", err)
	}
	if err != nil && expectedError != "" {
		// Ignore expected errors (TODO MariusVanDerWijden check error string)
		return nil
	}
	return nil
}

// Run executes a specific subtest and verifies the post-state and logs
func (t *StateTest) Run(subtest StateSubtest, vmconfig vm.Config, snapshotter bool) (*snapshot.Tree, *state.StateDB, error) {
	snaps, statedb, root, err := t.RunNoVerify(subtest, vmconfig, snapshotter)
	if checkedErr := t.checkError(subtest, err); checkedErr != nil {
		return snaps, statedb, checkedErr
	}
	// The error has been checked; if it was unexpected, it's already returned.
	if err != nil {
		// Here, an error exists but it was expected.
		// We do not check the post state or logs.
		return snaps, statedb, nil
	}
	post := t.json.Post[subtest.Fork][subtest.Index]
	// N.B: We need to do this in a two-step process, because the first Commit takes care
	// of suicides, and we need to touch the coinbase _after_ it has potentially suicided.
	if root != common.Hash(post.Root) {
		return snaps, statedb, fmt.Errorf("post state root mismatch: got %x, want %x", root, post.Root)
	}
	if logs := rlpHash(statedb.Logs()); logs != common.Hash(post.Logs) {
		return snaps, statedb, fmt.Errorf("post state logs hash mismatch: got %x, want %x", logs, post.Logs)
	}
	return snaps, statedb, nil
}

// RunNoVerify runs a specific subtest and returns the statedb and post-state root
func (t *StateTest) RunNoVerify(subtest StateSubtest, vmconfig vm.Config, snapshotter bool) (*snapshot.Tree, *state.StateDB, common.Hash, error) {
	config, eips, err := GetChainConfig(subtest.Fork)
	if err != nil {
		return nil, nil, common.Hash{}, UnsupportedForkError{subtest.Fork}
	}
	vmconfig.ExtraEips = eips
	block := t.genesis(config).ToBlock()
	snaps, statedb := MakePreState(rawdb.NewMemoryDatabase(), t.json.Pre, snapshotter)

	var baseFee *big.Int
	if config.IsEVM(0) {
		baseFee = t.json.Env.BaseFee
		if baseFee == nil {
			// Retesteth uses `0x10` for genesis baseFee. Therefore, it defaults to
			// parent - 2 : 0xa as the basefee for 'this' context.
			baseFee = big.NewInt(0x0a)
		}
	}
	post := t.json.Post[subtest.Fork][subtest.Index]
	msg, err := t.json.Tx.toMessage(post, baseFee)
	if err != nil {
		return nil, nil, common.Hash{}, err
	}

	// Try to recover tx with current signer
	if len(post.TxBytes) != 0 {
		var ttx types.Transaction
		err := ttx.UnmarshalBinary(post.TxBytes)
		if err != nil {
			return nil, nil, common.Hash{}, err
		}

		if _, err := types.Sender(types.LatestSigner(config), &ttx); err != nil {
			return nil, nil, common.Hash{}, err
		}
	}

	// Prepare the EVM.
	txContext := core.NewEVMTxContext(msg)
	context := core.NewEVMBlockContext(block.Header(), nil, &t.json.Env.Coinbase)
	context.GetHash = vmTestBlockHash
	context.BaseFee = baseFee
	if config.IsEVM(0) && t.json.Env.Random != nil {
		context.Difficulty = big.NewInt(0)
	}
	evm := vm.NewEVM(context, txContext, statedb, config, vmconfig)
	// Execute the message.
	snapshot := statedb.Snapshot()
	gaspool := new(core.GasPool)
	gaspool.AddGas(block.GasLimit())
	_, err = core.ApplyMessage(evm, msg, gaspool)
	if err != nil {
		statedb.RevertToSnapshot(snapshot)
	}
	// Add 0-value mining reward. This only makes a difference in the cases
	// where
	// - the coinbase suicided, or
	// - there are only 'bad' transactions, which aren't executed. In those cases,
	//   the coinbase gets no txfee, so isn't created, and thus needs to be touched
	statedb.AddBalance(block.Coinbase(), new(big.Int))
	// Commit block
	statedb.Commit(config.IsEIP158(block.Number()), false)
	// And _now_ get the state root
	root := statedb.IntermediateRoot(config.IsEIP158(block.Number()))
	return snaps, statedb, root, err
}

func MakePreState(db ethdb.Database, accounts core.GenesisAlloc, snapshotter bool) (*snapshot.Tree, *state.StateDB) {
	sdb := state.NewDatabaseWithConfig(db, &trie.Config{Preimages: true})
	statedb, _ := state.New(types.EmptyRootHash, sdb, nil)
	for addr, a := range accounts {
		statedb.SetCode(addr, a.Code)
		statedb.SetNonce(addr, a.Nonce)
		statedb.SetBalance(addr, a.Balance)
		for k, v := range a.Storage {
			statedb.SetState(addr, k, v)
		}
	}
	// Commit and re-open to start with a clean state.
	root, _ := statedb.Commit(false, false)

	var snaps *snapshot.Tree
	if snapshotter {
		snapconfig := snapshot.Config{
			CacheSize:  1,
			NoBuild:    false,
			AsyncBuild: false,
			SkipVerify: true,
		}
		snaps, _ = snapshot.New(snapconfig, db, sdb.TrieDB(), common.Hash{}, root)
	}
	statedb, _ = state.New(root, sdb, snaps)
	return snaps, statedb
}

func (t *StateTest) genesis(config *params.ChainConfig) *core.Genesis {
	genesis := &core.Genesis{
		Config:     config,
		Coinbase:   t.json.Env.Coinbase,
		Difficulty: t.json.Env.Difficulty,
		GasLimit:   t.json.Env.GasLimit,
		Number:     t.json.Env.Number,
		Timestamp:  t.json.Env.Timestamp,
		Alloc:      t.json.Pre,
	}
	if t.json.Env.Random != nil {
		// Post-Merge
		genesis.Mixhash = common.BigToHash(t.json.Env.Random)
		genesis.Difficulty = big.NewInt(0)
	}
	return genesis
}

func (tx *stTransaction) toMessage(ps stPostState, baseFee *big.Int) (*core.Message, error) {
	// Derive sender from private key if present.
	var from common.Address
	if len(tx.PrivateKey) > 0 {
		key, err := crypto.ToECDSA(tx.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %v", err)
		}
		from = crypto.PubkeyToAddress(key.PublicKey)
	}
	// Parse recipient if present.
	var to *common.Address
	if tx.To != "" {
		to = new(common.Address)
		if err := to.UnmarshalText([]byte(tx.To)); err != nil {
			return nil, fmt.Errorf("invalid to address: %v", err)
		}
	}

	// Get values specific to this post state.
	if ps.Indexes.Data > len(tx.Data) {
		return nil, fmt.Errorf("tx data index %d out of bounds", ps.Indexes.Data)
	}
	if ps.Indexes.Value > len(tx.Value) {
		return nil, fmt.Errorf("tx value index %d out of bounds", ps.Indexes.Value)
	}
	if ps.Indexes.Gas > len(tx.GasLimit) {
		return nil, fmt.Errorf("tx gas limit index %d out of bounds", ps.Indexes.Gas)
	}
	dataHex := tx.Data[ps.Indexes.Data]
	valueHex := tx.Value[ps.Indexes.Value]
	gasLimit := tx.GasLimit[ps.Indexes.Gas]
	// Value, Data hex encoding is messy: https://github.com/ethereum/tests/issues/203
	value := new(big.Int)
	if valueHex != "0x" {
		v, ok := math.ParseBig256(valueHex)
		if !ok {
			return nil, fmt.Errorf("invalid tx value %q", valueHex)
		}
		value = v
	}
	data, err := hex.DecodeString(strings.TrimPrefix(dataHex, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid tx data %q", dataHex)
	}
	var accessList types.AccessList
	if tx.AccessLists != nil && tx.AccessLists[ps.Indexes.Data] != nil {
		accessList = *tx.AccessLists[ps.Indexes.Data]
	}
	// If baseFee provided, set gasPrice to effectiveGasPrice.
	gasPrice := tx.GasPrice
	if baseFee != nil {
		if tx.MaxFeePerGas == nil {
			tx.MaxFeePerGas = gasPrice
		}
		if tx.MaxFeePerGas == nil {
			tx.MaxFeePerGas = new(big.Int)
		}
		if tx.MaxPriorityFeePerGas == nil {
			tx.MaxPriorityFeePerGas = tx.MaxFeePerGas
		}
		gasPrice = math.BigMin(new(big.Int).Add(tx.MaxPriorityFeePerGas, baseFee),
			tx.MaxFeePerGas)
	}
	if gasPrice == nil {
		return nil, errors.New("no gas price provided")
	}

	msg := &core.Message{
		From:       from,
		To:         to,
		Nonce:      tx.Nonce,
		Value:      value,
		GasLimit:   gasLimit,
		GasPrice:   gasPrice,
		GasFeeCap:  tx.MaxFeePerGas,
		GasTipCap:  tx.MaxPriorityFeePerGas,
		Data:       data,
		AccessList: accessList,
	}
	return msg, nil
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func vmTestBlockHash(n uint64) common.Hash {
	return common.BytesToHash(crypto.Keccak256([]byte(big.NewInt(int64(n)).String())))
}
