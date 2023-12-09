package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/zsrv/supermicro-product-key/pkg/json"
	netinternal "github.com/zsrv/supermicro-product-key/pkg/net"
)

func init() {
	rootCmd.AddCommand(jsonCmd)
	jsonCmd.AddCommand(jsonVerifyCmd)
	jsonCmd.AddCommand(jsonBruteForceCmd)
	jsonCmd.AddCommand(jsonListSKUCmd)
}

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "JSON product key operations",
}

var jsonVerifyCmd = &cobra.Command{
	Use:   "verify MAC_ADDRESS PRODUCT_KEY",
	Short: "VerifyWithKey the signature of a JSON product key",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		macAddress, err := netinternal.ParseMAC(args[0])
		if err != nil {
			return err
		}
		productKey := args[1]

		license, err := json.ParseLicense(productKey)
		if err != nil {
			return err
		}

		if err := license.Verify(macAddress); err != nil {
			return errors.WithMessage(err, "product key verification failed")
		}

		fmt.Println("signature verified ok")
		return nil
	},
}

var jsonBruteForceCmd = &cobra.Command{
	Use:   "bruteforce PRODUCT_KEY",
	Short: "Find the MAC address associated with a JSON product key by brute force",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		productKey := args[0]

		fmt.Println("searching for mac address ...")

		mac, err := json.BruteForceMACAddressFromString(productKey)
		if err != nil {
			return err
		}

		fmt.Printf("found match! mac = '%s'\n", mac)
		return nil
	},
}

var jsonListSKUCmd = &cobra.Command{
	Use:   "listswid",
	Short: "Get a list of software identifiers that can be found in product keys",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 3, 1, 2, ' ', 0)
		fmt.Fprintf(w, "License SKU\tID\n")
		fmt.Fprintf(w, "-----------\t--\n")
		for _, swid := range json.SoftwareIdentifiers.List() {
			fmt.Fprintf(w, "%v\t%v\n", swid.SKU, swid.ID)
		}
		w.Flush()
	},
}
