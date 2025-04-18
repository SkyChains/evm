// Copyright (C) 2023-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"os"
	"testing"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	// Import the solidity package, so that ginkgo maps out the tests declared within the package
	"github.com/skychains/evm/tests/precompile/solidity"
)

func TestE2E(t *testing.T) {
	if basePath := os.Getenv("TEST_SOURCE_ROOT"); basePath != "" {
		os.Chdir(basePath)
	}
	gomega.RegisterFailHandler(ginkgo.Fail)
	solidity.RegisterAsyncTests()
	ginkgo.RunSpecs(t, "evm precompile ginkgo test suite")
}
