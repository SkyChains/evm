// (c) 2023-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package predicate

import (
	"bytes"
	"testing"

	"github.com/SkyChains/chain/utils"
	"github.com/SkyChains/evm/params"
	"github.com/stretchr/testify/require"
)

func testPackPredicate(t testing.TB, b []byte) {
	packedPredicate := PackPredicate(b)
	unpackedPredicated, err := UnpackPredicate(packedPredicate)
	require.NoError(t, err)
	require.Equal(t, b, unpackedPredicated)
}

func FuzzPackPredicate(f *testing.F) {
	for i := 0; i < 100; i++ {
		f.Add(utils.RandomBytes(i))
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		testPackPredicate(t, b)
	})
}

func TestUnpackInvalidPredicate(t *testing.T) {
	require := require.New(t)
	// Predicate encoding requires a 0xff delimiter byte followed by padding of all zeroes, so any other
	// excess padding should invalidate the predicate.
	paddingCases := make([][]byte, 0, 200)
	for i := 1; i < 100; i++ {
		paddingCases = append(paddingCases, bytes.Repeat([]byte{0xee}, i))
		paddingCases = append(paddingCases, make([]byte, i))
	}

	for _, l := range []int{0, 1, 31, 32, 33, 63, 64, 65} {
		validPredicate := PackPredicate(utils.RandomBytes(l))

		for _, padding := range paddingCases {
			invalidPredicate := append(validPredicate, padding...)
			_, err := UnpackPredicate(invalidPredicate)
			require.Error(err, "Predicate length %d, Padding length %d (0x%x)", len(validPredicate), len(padding), invalidPredicate)
		}
	}
}

func TestPredicateResultsBytes(t *testing.T) {
	require := require.New(t)
	dataTooShort := utils.RandomBytes(params.DynamicFeeExtraDataSize - 1)
	_, ok := GetPredicateResultBytes(dataTooShort)
	require.False(ok)

	preDUPgradeData := utils.RandomBytes(params.DynamicFeeExtraDataSize)
	_, ok = GetPredicateResultBytes(preDUPgradeData)
	require.False(ok)
	postDUpgradeData := utils.RandomBytes(params.DynamicFeeExtraDataSize + 2)
	resultBytes, ok := GetPredicateResultBytes(postDUpgradeData)
	require.True(ok)
	require.Equal(resultBytes, postDUpgradeData[params.DynamicFeeExtraDataSize:])
}
