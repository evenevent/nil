package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/NilFoundation/nil/nil/cmd/nil_block_generator/internal/commands"
	"github.com/NilFoundation/nil/nil/common/check"
	"github.com/NilFoundation/nil/nil/common/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func main() {
	check.PanicIfErr(execute())
}

func execute() error {
	rootCmd := &cobra.Command{
		Use:   os.Args[0],
		Short: "Run Block Generator CLI Tool",
	}

	logging.SetupGlobalLogger("info")
	logger := logging.NewLogger("block_generator_cli")

	initCmd := buildInitCmd(logger)
	rootCmd.AddCommand(initCmd)

	showContractsCmd := buildShowContractsCmd()
	rootCmd.AddCommand(showContractsCmd)

	showCallsCmd := buildShowCallsCmd()
	rootCmd.AddCommand(showCallsCmd)

	addContractCmd, err := buildAddContractCmd(logger)
	if err != nil {
		return err
	}
	callContractCmd, err := buildCallContractCmd()
	if err != nil {
		return err
	}
	getBlockCmd := buildGetBlockCmd(logger)
	rootCmd.AddCommand(addContractCmd, callContractCmd, getBlockCmd)

	return rootCmd.Execute()
}

func buildInitCmd(logger zerolog.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize new wallet",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := commands.CleanFiles()
			if err != nil {
				return err
			}
			err = commands.RunNilNode()
			if err != nil {
				return err
			}

			walletAdr, hexKey, err := commands.CreateNewSmartAccount(logger)
			if err != nil {
				return err
			}
			return commands.InitConfig(walletAdr, hexKey)
		},
	}
	return cmd
}

func buildAddContractCmd(logger zerolog.Logger) (*cobra.Command, error) {
	var contractName string
	var contractPath string
	cmd := &cobra.Command{
		Use:   "add-contract",
		Short: "Deploy new contract",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := commands.RunNilNode()
			if err != nil {
				return err
			}

			cfg, err := commands.ReadConfigFromFile()
			if err != nil {
				return err
			}
			adr, err := commands.DeployContract(contractPath, cfg.PrivateKey, logger)
			if err != nil {
				return err
			}
			return commands.AddContract(contractName, contractPath, adr)
		},
	}

	const contractNameFlag = "contract-name"
	cmd.Flags().StringVar(&contractName, contractNameFlag, contractName, "name of contract")
	if err := cmd.MarkFlagRequired(contractNameFlag); err != nil {
		return nil, err
	}

	const contractPathFlag = "contract-path"
	cmd.Flags().StringVar(&contractPath, contractPathFlag, contractPath, "path to contract code")
	if err := cmd.MarkFlagRequired(contractPathFlag); err != nil {
		return nil, err
	}

	return cmd, nil
}

func buildCallContractCmd() (*cobra.Command, error) {
	var contractName string
	var method string
	var argsCmd string
	var count int
	cmd := &cobra.Command{
		Use:   "call-contract",
		Short: "Call method of deployed contract",
		RunE: func(cmd *cobra.Command, args []string) error {
			callArgs := strings.Fields(argsCmd)
			return commands.AddCall(contractName, method, callArgs, count)
		},
	}

	const contractNameFlag = "contract-name"
	cmd.Flags().StringVar(&contractName, contractNameFlag, contractName, "name of contract")
	if err := cmd.MarkFlagRequired(contractNameFlag); err != nil {
		return nil, err
	}

	const methodFlag = "method"
	cmd.Flags().StringVar(&method, methodFlag, method, "method to be called")
	if err := cmd.MarkFlagRequired(methodFlag); err != nil {
		return nil, err
	}

	const argsFlag = "args"
	cmd.Flags().StringVar(&argsCmd, argsFlag, argsCmd, "method arguments")
	if err := cmd.MarkFlagRequired(argsFlag); err != nil {
		return nil, err
	}

	const countFlag = "count"
	cmd.Flags().IntVar(&count, countFlag, count, "number of calls")
	if err := cmd.MarkFlagRequired(countFlag); err != nil {
		return nil, err
	}

	return cmd, nil
}

func buildGetBlockCmd(logger zerolog.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-block",
		Short: "Generate block from the current config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := commands.ReadConfigFromFile()
			if err != nil {
				return err
			}
			err = commands.RunNilNode()
			if err != nil {
				return err
			}

			blockHash, err := commands.CallContract(cfg.WalletAdr, cfg.PrivateKey, cfg.Calls, logger)
			if err != nil {
				return err
			}
			fmt.Printf("Hash of the block: %s\n", blockHash)
			return nil
		},
	}
	return cmd
}

func buildShowContractsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-contracts",
		Short: "Print deployed contracts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.ShowContracts(os.Stdout)
		},
	}
	return cmd
}

func buildShowCallsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-calls",
		Short: "Print calls of contracts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.ShowCalls(os.Stdout)
		},
	}
	return cmd
}
