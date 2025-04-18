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

package filters

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	"github.com/skychains/evm/consensus/dummy"
	"github.com/skychains/evm/core"
	"github.com/skychains/evm/core/rawdb"
	"github.com/skychains/evm/core/types"
	"github.com/skychains/evm/params"
	"github.com/skychains/evm/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func makeReceipt(addr common.Address) *types.Receipt {
	receipt := types.NewReceipt(nil, false, 0)
	receipt.Logs = []*types.Log{
		{Address: addr},
	}
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	return receipt
}

func BenchmarkFilters(b *testing.B) {
	var (
		db, _   = rawdb.NewLevelDBDatabase(b.TempDir(), 0, 0, "", false)
		_, sys  = newTestFilterSystem(b, db, Config{})
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		addr2   = common.BytesToAddress([]byte("jeff"))
		addr3   = common.BytesToAddress([]byte("ethereum"))
		addr4   = common.BytesToAddress([]byte("random addresses please"))

		gspec = &core.Genesis{
			Config:  params.TestChainConfig,
			Alloc:   core.GenesisAlloc{addr1: {Balance: big.NewInt(1000000)}},
			BaseFee: big.NewInt(1),
		}
	)
	defer db.Close()
	_, chain, receipts, err := core.GenerateChainWithGenesis(gspec, dummy.NewFaker(), 100010, 10, func(i int, gen *core.BlockGen) {
		switch i {
		case 2403:
			receipt := makeReceipt(addr1)
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction(999, common.HexToAddress("0x999"), big.NewInt(999), 999, gen.BaseFee(), nil))
		case 1034:
			receipt := makeReceipt(addr2)
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction(999, common.HexToAddress("0x999"), big.NewInt(999), 999, gen.BaseFee(), nil))
		case 34:
			receipt := makeReceipt(addr3)
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction(999, common.HexToAddress("0x999"), big.NewInt(999), 999, gen.BaseFee(), nil))
		case 99999:
			receipt := makeReceipt(addr4)
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction(999, common.HexToAddress("0x999"), big.NewInt(999), 999, gen.BaseFee(), nil))
		}
	})
	require.NoError(b, err)
	// The test txs are not properly signed, can't simply create a chain
	// and then import blocks. TODO(rjl493456442) try to get rid of the
	// manual database writes.
	gspec.MustCommit(db)
	for i, block := range chain {
		rawdb.WriteBlock(db, block)
		rawdb.WriteCanonicalHash(db, block.Hash(), block.NumberU64())
		rawdb.WriteHeadBlockHash(db, block.Hash())
		rawdb.WriteReceipts(db, block.Hash(), block.NumberU64(), receipts[i])
	}
	b.ResetTimer()

	filter, err := sys.NewRangeFilter(0, int64(rpc.LatestBlockNumber), []common.Address{addr1, addr2, addr3, addr4}, nil)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		logs, _ := filter.Logs(context.Background())
		if len(logs) != 4 {
			b.Fatal("expected 4 logs, got", len(logs))
		}
	}
}

func TestFilters(t *testing.T) {
	var (
		db, _   = rawdb.NewLevelDBDatabase(t.TempDir(), 0, 0, "", false)
		_, sys  = newTestFilterSystem(t, db, Config{})
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr    = crypto.PubkeyToAddress(key1.PublicKey)

		hash1 = common.BytesToHash([]byte("topic1"))
		hash2 = common.BytesToHash([]byte("topic2"))
		hash3 = common.BytesToHash([]byte("topic3"))
		hash4 = common.BytesToHash([]byte("topic4"))

		gspec = &core.Genesis{
			Config:  params.TestChainConfig,
			Alloc:   core.GenesisAlloc{addr: {Balance: big.NewInt(1000000)}},
			BaseFee: big.NewInt(1),
		}
	)
	defer db.Close()

	_, chain, receipts, err := core.GenerateChainWithGenesis(gspec, dummy.NewFaker(), 1000, 10, func(i int, gen *core.BlockGen) {
		switch i {
		case 1:
			receipt := types.NewReceipt(nil, false, 0)
			receipt.Logs = []*types.Log{
				{
					Address: addr,
					Topics:  []common.Hash{hash1},
				},
			}
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction(1, common.HexToAddress("0x1"), big.NewInt(1), 1, gen.BaseFee(), nil))
		case 2:
			receipt := types.NewReceipt(nil, false, 0)
			receipt.Logs = []*types.Log{
				{
					Address: addr,
					Topics:  []common.Hash{hash2},
				},
			}
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction(2, common.HexToAddress("0x2"), big.NewInt(2), 2, gen.BaseFee(), nil))

		case 998:
			receipt := types.NewReceipt(nil, false, 0)
			receipt.Logs = []*types.Log{
				{
					Address: addr,
					Topics:  []common.Hash{hash3},
				},
			}
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction(998, common.HexToAddress("0x998"), big.NewInt(998), 998, gen.BaseFee(), nil))
		case 999:
			receipt := types.NewReceipt(nil, false, 0)
			receipt.Logs = []*types.Log{
				{
					Address: addr,
					Topics:  []common.Hash{hash4},
				},
			}
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction(999, common.HexToAddress("0x999"), big.NewInt(999), 999, gen.BaseFee(), nil))
		}
	})
	require.NoError(t, err)
	// The test txs are not properly signed, can't simply create a chain
	// and then import blocks. TODO(rjl493456442) try to get rid of the
	// manual database writes.
	gspec.MustCommit(db)
	for i, block := range chain {
		rawdb.WriteBlock(db, block)
		rawdb.WriteCanonicalHash(db, block.Hash(), block.NumberU64())
		rawdb.WriteHeadBlockHash(db, block.Hash())
		rawdb.WriteReceipts(db, block.Hash(), block.NumberU64(), receipts[i])
	}

	filter, err := sys.NewRangeFilter(0, int64(rpc.LatestBlockNumber), []common.Address{addr}, [][]common.Hash{{hash1, hash2, hash3, hash4}})
	require.NoError(t, err)
	logs, _ := filter.Logs(context.Background())
	if len(logs) != 4 {
		t.Error("expected 4 log, got", len(logs))
	}

	for i, tc := range []struct {
		f          *Filter
		wantHashes []common.Hash
	}{
		{
			mustNewRangeFilter(t, sys, 900, 999, []common.Address{addr}, [][]common.Hash{{hash3}}),
			[]common.Hash{hash3},
		}, {
			mustNewRangeFilter(t, sys, 990, int64(rpc.LatestBlockNumber), []common.Address{addr}, [][]common.Hash{{hash3}}),
			[]common.Hash{hash3},
		}, {
			mustNewRangeFilter(t, sys, 1, 10, nil, [][]common.Hash{{hash1, hash2}}),
			[]common.Hash{hash1, hash2},
		}, {
			mustNewRangeFilter(t, sys, 0, int64(rpc.LatestBlockNumber), nil, [][]common.Hash{{common.BytesToHash([]byte("fail"))}}),
			nil,
		}, {
			mustNewRangeFilter(t, sys, 0, int64(rpc.LatestBlockNumber), []common.Address{common.BytesToAddress([]byte("failmenow"))}, nil),
			nil,
		}, {
			mustNewRangeFilter(t, sys, 0, int64(rpc.LatestBlockNumber), nil, [][]common.Hash{{common.BytesToHash([]byte("fail"))}, {hash1}}),
			nil,
		}, {
			mustNewRangeFilter(t, sys, int64(rpc.LatestBlockNumber), int64(rpc.LatestBlockNumber), nil, nil), []common.Hash{hash4},
		}, {
			// Note: modified from go-ethereum since we don't have FinalizedBlock
			mustNewRangeFilter(t, sys, int64(rpc.AcceptedBlockNumber), int64(rpc.LatestBlockNumber), nil, nil), []common.Hash{hash4},
		}, {
			// Note: modified from go-ethereum since we don't have FinalizedBlock
			mustNewRangeFilter(t, sys, int64(rpc.AcceptedBlockNumber), int64(rpc.AcceptedBlockNumber), nil, nil), []common.Hash{hash4},
		}, {
			// Note: modified from go-ethereum since we don't have FinalizedBlock
			mustNewRangeFilter(t, sys, int64(rpc.LatestBlockNumber), -3, nil, nil), []common.Hash{hash4},
		}, {
			// Note: modified from go-ethereum since we don't have SafeBlock
			mustNewRangeFilter(t, sys, int64(rpc.AcceptedBlockNumber), int64(rpc.LatestBlockNumber), nil, nil), []common.Hash{hash4},
		}, {
			// Note: modified from go-ethereum since we don't have SafeBlock
			mustNewRangeFilter(t, sys, int64(rpc.AcceptedBlockNumber), int64(rpc.AcceptedBlockNumber), nil, nil), []common.Hash{hash4},
		}, {
			// Note: modified from go-ethereum since we don't have SafeBlock
			mustNewRangeFilter(t, sys, int64(rpc.LatestBlockNumber), int64(rpc.AcceptedBlockNumber), nil, nil), []common.Hash{hash4},
		},
		{
			mustNewRangeFilter(t, sys, int64(rpc.PendingBlockNumber), int64(rpc.PendingBlockNumber), nil, nil), nil,
		},
	} {
		logs, _ := tc.f.Logs(context.Background())
		var haveHashes []common.Hash
		for _, l := range logs {
			haveHashes = append(haveHashes, l.Topics[0])
		}
		if have, want := len(haveHashes), len(tc.wantHashes); have != want {
			t.Fatalf("test %d, have %d logs, want %d", i, have, want)
		}
		if len(haveHashes) == 0 {
			continue
		}
		if !reflect.DeepEqual(tc.wantHashes, haveHashes) {
			t.Fatalf("test %d, have %v want %v", i, haveHashes, tc.wantHashes)
		}
	}
}

func mustNewRangeFilter(t *testing.T, sys *FilterSystem, begin, end int64, addresses []common.Address, topics [][]common.Hash) *Filter {
	t.Helper()
	f, err := sys.NewRangeFilter(begin, end, addresses, topics)
	require.NoError(t, err)
	return f
}
