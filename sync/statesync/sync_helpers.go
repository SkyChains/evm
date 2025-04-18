// (c) 2021-2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package statesync

import (
	"github.com/skychains/evm/core/rawdb"
	"github.com/skychains/evm/core/state/snapshot"
	"github.com/skychains/evm/core/types"
	"github.com/skychains/evm/ethdb"
	"github.com/skychains/evm/trie"
	"github.com/ethereum/go-ethereum/common"
)

// writeAccountSnapshot stores the account represented by [acc] to the snapshot at [accHash], using
// SlimAccountRLP format (omitting empty code/storage).
func writeAccountSnapshot(db ethdb.KeyValueWriter, accHash common.Hash, acc types.StateAccount) {
	slimAccount := snapshot.SlimAccountRLP(acc.Nonce, acc.Balance, acc.Root, acc.CodeHash)
	rawdb.WriteAccountSnapshot(db, accHash, slimAccount)
}

// writeAccountStorageSnapshotFromTrie iterates the trie at [storageTrie] and copies all entries
// to the storage snapshot for [accountHash].
func writeAccountStorageSnapshotFromTrie(batch ethdb.Batch, batchSize int, accountHash common.Hash, storageTrie *trie.Trie) error {
	it := trie.NewIterator(storageTrie.NodeIterator(nil))
	for it.Next() {
		rawdb.WriteStorageSnapshot(batch, accountHash, common.BytesToHash(it.Key), it.Value)
		if batch.ValueSize() > batchSize {
			if err := batch.Write(); err != nil {
				return err
			}
			batch.Reset()
		}
	}
	if it.Err != nil {
		return it.Err
	}
	return batch.Write()
}
