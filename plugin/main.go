// (c) 2021-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"fmt"

	"github.com/SkyChains/chain/version"
	"github.com/SkyChains/evm/plugin/evm"
	"github.com/SkyChains/evm/plugin/runner"
)

func main() {
<<<<<<< HEAD
<<<<<<< HEAD
	versionString := fmt.Sprintf("Subnet-EVM/%s [Lux Node=%s, rpcchainvm=%d]", evm.Version, version.Current, version.RPCChainVMProtocol)
=======
	versionString := fmt.Sprintf("Subnet-EVM/%s [Luxd=%s, rpcchainvm=%d]", evm.Version, version.Current, version.RPCChainVMProtocol)
>>>>>>> b36c20f (Update executable to luxd)
=======
	versionString := fmt.Sprintf("EVM/%s [Luxd=%s, rpcchainvm=%d]", evm.Version, version.Current, version.RPCChainVMProtocol)
>>>>>>> fd08c47 (Update import path)
	runner.Run(versionString)
}
