// (c) 2019-2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package message

import (
	"github.com/skychains/chain/codec"
	"github.com/skychains/chain/codec/linearcodec"
	"github.com/skychains/chain/utils/units"
	"github.com/skychains/chain/utils/wrappers"
)

const (
	Version        = uint16(0)
	maxMessageSize = 1 * units.MiB
)

var (
	Codec           codec.Manager
	CrossChainCodec codec.Manager
)

func init() {
	Codec = codec.NewManager(maxMessageSize)
	c := linearcodec.NewDefault()

	errs := wrappers.Errs{}
	errs.Add(
		// Gossip types
		c.RegisterType(TxsGossip{}),

		// Types for state sync frontier consensus
		c.RegisterType(SyncSummary{}),

		// state sync types
		c.RegisterType(BlockRequest{}),
		c.RegisterType(BlockResponse{}),
		c.RegisterType(LeafsRequest{}),
		c.RegisterType(LeafsResponse{}),
		c.RegisterType(CodeRequest{}),
		c.RegisterType(CodeResponse{}),

		// Warp request types
		c.RegisterType(MessageSignatureRequest{}),
		c.RegisterType(BlockSignatureRequest{}),
		c.RegisterType(SignatureResponse{}),

		Codec.RegisterCodec(Version, c),
	)

	if errs.Errored() {
		panic(errs.Err)
	}

	CrossChainCodec = codec.NewManager(maxMessageSize)
	ccc := linearcodec.NewDefault()

	errs = wrappers.Errs{}
	errs.Add(
		// CrossChainRequest Types
		ccc.RegisterType(EthCallRequest{}),
		ccc.RegisterType(EthCallResponse{}),

		CrossChainCodec.RegisterCodec(Version, ccc),
	)

	if errs.Errored() {
		panic(errs.Err)
	}
}
