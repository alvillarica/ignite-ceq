package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/chaincmd"
	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/services/chain"
)

var moniker string
var seedphrase string
var runbuild bool

func NewChainInit() *cobra.Command {
	c := &cobra.Command{
		Use:   "init",
		Short: "Initialize your chain",
		Args:  cobra.NoArgs,
		RunE:  chainInitHandler,
	}

	flagSetPath(c)
	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().StringVar(&moniker, "moniker", "mynode", "Name of node")
	c.Flags().StringVar(&seedphrase, "seedphrase", "", "Seed phrase for the pioneer")
	c.Flags().BoolVar(&runbuild, "build-chain", true, "Build the chain?")

	return c
}

func chainInitHandler(cmd *cobra.Command, _ []string) error {
	chainOption := []chain.Option{
		chain.LogLevel(logLevel(cmd)),
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
	}

	c, err := newChainWithHomeFlags(cmd, chainOption...)
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

  if runbuild {
    if _, err := c.Build(cmd.Context(), cacheStorage, ""); err != nil {
      return err
    }
  }

	if err := c.Init(cmd.Context(), true, moniker, seedphrase); err != nil {
		return err
	}

	home, err := c.Home()
	if err != nil {
		return err
	}

	fmt.Printf("🗃  Initialized. Checkout your chain's home (data) directory: %s\n", colors.Info(home))

	return nil
}
