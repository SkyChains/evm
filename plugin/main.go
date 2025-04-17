// (c) 2021-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"fmt"

	"github.com/skychains/chain/version"
	"github.com/skychains/evm/plugin/evm"
	"github.com/skychains/evm/plugin/runner"
)

func main() {

	versionString := fmt.Sprintf("EVM/%s [Luxd=%s, rpcchainvm=%d]", evm.Version, version.Current, version.RPCChainVMProtocol)
	runner.Run(versionString)
}
