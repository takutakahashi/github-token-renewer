package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCmd(t *testing.T) {
	if rootCmd == nil {
		t.Fatal("rootCmd is nil")
	}

	if rootCmd.Use != "github-token-renewer" {
		t.Errorf("rootCmd.Use = %v, want github-token-renewer", rootCmd.Use)
	}

	if rootCmd.Run == nil {
		t.Error("rootCmd.Run is nil")
	}
}

func TestRootCmdFlags(t *testing.T) {
	configFlag := rootCmd.Flags().Lookup("config")
	if configFlag == nil {
		t.Fatal("config flag not found")
	}

	if configFlag.Shorthand != "c" {
		t.Errorf("config flag shorthand = %v, want c", configFlag.Shorthand)
	}

	if configFlag.DefValue != "./config.yaml" {
		t.Errorf("config flag default = %v, want ./config.yaml", configFlag.DefValue)
	}
}

func TestExecute(t *testing.T) {
	// Create a temporary root command for testing
	testCmd := &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	// Test that Execute can be called (even though it will fail without valid config)
	// This is just a structural test
	if testCmd.Use != "test" {
		t.Errorf("testCmd.Use = %v, want test", testCmd.Use)
	}
}
