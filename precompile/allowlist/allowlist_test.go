// (c) 2021-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import (
	"testing"

	"github.com/skychains/evm/core/state"
	"github.com/skychains/evm/precompile/contract"
	"github.com/skychains/evm/precompile/modules"
	"github.com/skychains/evm/precompile/precompileconfig"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ precompileconfig.Config = &dummyConfig{}
	_ contract.Configurator   = &dummyConfigurator{}

	dummyAddr = common.Address{1}
)

type dummyConfig struct {
	precompileconfig.Upgrade
	AllowListConfig
}

func (d *dummyConfig) Key() string      { return "dummy" }
func (d *dummyConfig) IsDisabled() bool { return false }
func (d *dummyConfig) Verify(chainConfig precompileconfig.ChainConfig) error {
	return d.AllowListConfig.Verify(chainConfig, d.Upgrade)
}

func (d *dummyConfig) Equal(cfg precompileconfig.Config) bool {
	other, ok := (cfg).(*dummyConfig)
	if !ok {
		return false
	}
	return d.AllowListConfig.Equal(&other.AllowListConfig)
}

type dummyConfigurator struct{}

func (d *dummyConfigurator) MakeConfig() precompileconfig.Config {
	return &dummyConfig{}
}

func (d *dummyConfigurator) Configure(
	chainConfig precompileconfig.ChainConfig,
	precompileConfig precompileconfig.Config,
	state contract.StateDB,
	blockContext contract.ConfigurationBlockContext,
) error {
	cfg := precompileConfig.(*dummyConfig)
	return cfg.AllowListConfig.Configure(chainConfig, dummyAddr, state, blockContext)
}

func TestAllowListRun(t *testing.T) {
	dummyModule := modules.Module{
		Address:      dummyAddr,
		Contract:     CreateAllowListPrecompile(dummyAddr),
		Configurator: &dummyConfigurator{},
		ConfigKey:    "dummy",
	}
	RunPrecompileWithAllowListTests(t, dummyModule, state.NewTestStateDB, nil)
}

func BenchmarkAllowList(b *testing.B) {
	dummyModule := modules.Module{
		Address:      dummyAddr,
		Contract:     CreateAllowListPrecompile(dummyAddr),
		Configurator: &dummyConfigurator{},
		ConfigKey:    "dummy",
	}
	BenchPrecompileWithAllowList(b, dummyModule, state.NewTestStateDB, nil)
}
