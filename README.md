# EVM

[![Build + Test + Release](https://github.com/skychains/evm/actions/workflows/lint-tests-release.yml/badge.svg)](https://github.com/skychains/evm/actions/workflows/lint-tests-release.yml)
[![CodeQL](https://github.com/skychains/evm/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/skychains/evm/actions/workflows/codeql-analysis.yml)

[Lux](https://docs.lux.network/overview/getting-started/lux-platform) is a network composed of multiple blockchains.
Each blockchain is an instance of a Virtual Machine (VM), much like an object in an object-oriented language is an instance of a class.
That is, the VM defines the behavior of the blockchain.

EVM is the [Virtual Machine (VM)](https://docs.lux.network/learn/lux/virtual-machines) that defines the Subnet Contract Chains. EVM is a simplified version of [Coreth VM (C-Chain)](https://github.com/skychains/coreth).

This chain implements the Ethereum Virtual Machine and supports Solidity smart contracts as well as most other Ethereum client functionality.

## Building

The EVM runs in a separate process from the main Luxd process and communicates with it over a local gRPC connection.

### Luxd Compatibility

```text
[v0.1.0] Luxd@v1.7.0-v1.7.4 (Protocol Version: 9)
[v0.1.1-v0.1.2] Luxd@v1.7.5-v1.7.6 (Protocol Version: 10)
[v0.2.0] Luxd@v1.7.7-v1.7.9 (Protocol Version: 11)
[v0.2.1] Luxd@v1.7.10 (Protocol Version: 12)
[v0.2.2] Luxd@v1.7.11-v1.7.12 (Protocol Version: 14)
[v0.2.3] Luxd@v1.7.13-v1.7.16 (Protocol Version: 15)
[v0.2.4] Luxd@v1.7.13-v1.7.16 (Protocol Version: 15)
[v0.2.5] Luxd@v1.7.13-v1.7.16 (Protocol Version: 15)
[v0.2.6] Luxd@v1.7.13-v1.7.16 (Protocol Version: 15)
[v0.2.7] Luxd@v1.7.13-v1.7.16 (Protocol Version: 15)
[v0.2.8] Luxd@v1.7.13-v1.7.18 (Protocol Version: 15)
[v0.2.9] Luxd@v1.7.13-v1.7.18 (Protocol Version: 15)
[v0.3.0] Luxd@v1.8.0-v1.8.6 (Protocol Version: 16)
[v0.4.0] Luxd@v1.9.0 (Protocol Version: 17)
[v0.4.1] Luxd@v1.9.1 (Protocol Version: 18)
[v0.4.2] Luxd@v1.9.1 (Protocol Version: 18)
[v0.4.3] Luxd@v1.9.2-v1.9.3 (Protocol Version: 19)
[v0.4.4] Luxd@v1.9.2-v1.9.3 (Protocol Version: 19)
[v0.4.5] Luxd@v1.9.4 (Protocol Version: 20)
[v0.4.6] Luxd@v1.9.4 (Protocol Version: 20)
[v0.4.7] Luxd@v1.9.5 (Protocol Version: 21)
[v0.4.8] Luxd@v1.9.6-v1.9.8 (Protocol Version: 22)
[v0.4.9] Luxd@v1.9.9 (Protocol Version: 23)
[v0.4.10] Luxd@v1.9.9 (Protocol Version: 23)
[v0.4.11] Luxd@v1.9.10-v1.9.16 (Protocol Version: 24)
[v0.4.12] Luxd@v1.9.10-v1.9.16 (Protocol Version: 24)
[v0.5.0] Luxd@v1.10.0 (Protocol Version: 25)
[v0.5.1] Luxd@v1.10.1-v1.10.4 (Protocol Version: 26)
[v0.5.2] Luxd@v1.10.1-v1.10.4 (Protocol Version: 26)
[v0.5.3] Luxd@v1.10.5-v1.10.8 (Protocol Version: 27)
[v0.5.4] Luxd@v1.10.9-v1.10.12 (Protocol Version: 28)
[v0.5.5] Luxd@v1.10.9-v1.10.12 (Protocol Version: 28)
[v0.5.6] Luxd@v1.10.9-v1.10.12 (Protocol Version: 28)
[v0.5.7] Luxd@v1.10.13-v1.10.14 (Protocol Version: 29)
[v0.5.8] Luxd@v1.10.13-v1.10.14 (Protocol Version: 29)
[v0.5.9] Luxd@v1.10.15-v1.10.17 (Protocol Version: 30)
[v0.5.10] Luxd@v1.10.15-v1.10.17 (Protocol Version: 30)
```

## API

The EVM supports the following API namespaces:

- `eth`
- `personal`
- `txpool`
- `debug`

Only the `eth` namespace is enabled by default.
EVM is a simplified version of [Coreth VM (C-Chain)](https://github.com/skychains/coreth).
Full documentation for the C-Chain's API can be found [here](https://docs.lux.network/apis/node/apis/c-chain).

## Compatibility

The EVM is compatible with almost all Ethereum tooling, including [Remix](https://docs.lux.network/build/dapp/smart-contracts/remix-deploy), [Metamask](https://docs.lux.network/build/dapp/chain-settings), and [Foundry](https://docs.lux.network/build/dapp/smart-contracts/toolchains/foundry).

## Differences Between EVM and Coreth

- Added configurable fees and gas limits in genesis
- Merged Lux hardforks into the single "EVM" hardfork
- Removed Atomic Txs and Shared Memory
- Removed Multicoin Contract and State

## Block Format

To support these changes, there have been a number of changes to the EVM block format compared to what exists on the C-Chain and Ethereum. Here we list the changes to the block format as compared to Ethereum.

### Block Header

- `BaseFee`: Added by EIP-1559 to represent the base fee of the block (present in Ethereum as of EIP-1559)
- `BlockGasCost`: surcharge for producing a block faster than the target rate

## Create an EVM Subnet on a Local Network

### Clone Subnet-evm

First install Go 1.20.10 or later. Follow the instructions [here](https://go.dev/doc/install). You can verify by running `go version`.

Set `$GOPATH` environment variable properly for Go to look for Go Workspaces. Please read [this](https://go.dev/doc/code) for details. You can verify by running `echo $GOPATH`.

As a few software will be installed into `$GOPATH/bin`, please make sure that `$GOPATH/bin` is in your `$PATH`, otherwise, you may get error running the commands below.

Download the `evm` repository into your `$GOPATH`:

```sh
cd $GOPATH
mkdir -p src/github.com/skychains
cd src/github.com/skychains
git clone git@github.com:skychains/evm.git
cd evm
```

This will clone and checkout to `master` branch.

### Run Local Network

To run a local network, it is recommended to use the [cli](https://github.com/skychains/cli#cli) to set up an instance of EVM on an local Lux Network.

There are two options when using the cli:

1. Use an official EVM release: https://docs.lux.network/subnets/build-first-subnet
2. Build and deploy a locally built (and optionally modified) version of EVM: https://docs.lux.network/subnets/create-custom-subnet
