package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/zsrv/supermicro-product-key/pkg/nonjson"
	"os"
	"text/tabwriter"
	"time"
)

var swSKU string
var swDisplayName string
var swID []byte

var softwareVersion string
var invoiceNumber string
var creationDate string
var expirationDate string
var property []byte

func init() {
	rootCmd.AddCommand(nonjsonCmd)
	nonjsonCmd.AddCommand(nonjsonDecodeCmd)
	nonjsonCmd.AddCommand(nonjsonEncodeCmd)
	nonjsonCmd.AddCommand(nonjsonBruteForceCmd)
	nonjsonCmd.AddCommand(nonjsonListSKUCmd)

	nonjsonEncodeCmd.Flags().StringVar(&swSKU, "sku", "", "license SKU (one of sku, display-name, id are required)")
	nonjsonEncodeCmd.Flags().StringVar(&swDisplayName, "display-name", "", "software display name (one of sku, display-name, id are required)")
	nonjsonEncodeCmd.Flags().BytesHexVar(&swID, "id", []byte{}, "software ID in hex, e.g. '02' (one of sku, display-name, id are required)")
	nonjsonEncodeCmd.MarkFlagsMutuallyExclusive("sku", "display-name", "id")

	nonjsonEncodeCmd.Flags().StringVar(&softwareVersion, "software-version", "", "software version (default 'none')")
	nonjsonEncodeCmd.Flags().StringVar(&invoiceNumber, "invoice-number", "", "invoice number (default 'none')")
	nonjsonEncodeCmd.Flags().StringVar(&creationDate, "creation-date", "", "creation date in RFC3339 format, e.g. '1970-01-01T00:00:00Z' (default is the current date and time)")
	nonjsonEncodeCmd.Flags().StringVar(&expirationDate, "expiration-date", "", "expiration date in RFC3339 format (default '1970-01-01T00:00:00Z' - no expiration)")
	nonjsonEncodeCmd.Flags().BytesHexVar(&property, "property", []byte{}, "property in hex, e.g. '01EE02FF' (default null)")
}

var nonjsonCmd = &cobra.Command{
	Use:   "nonjson",
	Short: "Non-JSON product key operations",
}

var nonjsonDecodeCmd = &cobra.Command{
	Use:   "decode MAC_ADDRESS PRODUCT_KEY",
	Short: "Decode a non-JSON product key",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		macAddress := args[0]
		encodedProductKey := args[1]

		productKey, err := nonjson.DecodeProductKey(encodedProductKey, macAddress)
		if err != nil {
			return errors.WithMessage(err, "failed to decode product key")
		}

		s, _ := json.MarshalIndent(productKey, "", "\t")
		fmt.Println(string(s))
		return nil
	},
}

var nonjsonEncodeCmd = &cobra.Command{
	Use:   "encode {--sku sku | --display-name name | --id id} [--software-version version] [--invoice-number invoice] [--creation-date date] [--expiration-date date] [--property property] MAC_ADDRESS",
	Short: "Encode a non-JSON product key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		macAddress := args[0]

		pk := nonjson.NewProductKey()

		var swidSet bool
		switch {
		case swSKU != "":
			swid, err := nonjson.SoftwareIdentifiers.BySKU(swSKU)
			if err != nil {
				return err
			}
			pk.SoftwareIdentifier = *swid
			swidSet = true
		case swDisplayName != "":
			swid, err := nonjson.SoftwareIdentifiers.ByDisplayName(swDisplayName)
			if err != nil {
				return err
			}
			pk.SoftwareIdentifier = *swid
			swidSet = true
		case len(swID) > 0:
			if len(swID) != 1 {
				return errors.New("software id must have a length of 1 byte")
			}

			swid, err := nonjson.SoftwareIdentifiers.ByID(swID[0])
			if err != nil {
				return err
			}
			pk.SoftwareIdentifier = *swid
			swidSet = true
		}
		if swidSet != true {
			return errors.New("one of sku, display-name, id are required")
		}

		if softwareVersion != "" {
			pk.SoftwareVersion = softwareVersion
		}

		if invoiceNumber != "" {
			pk.InvoiceNumber = invoiceNumber
		}

		if creationDate != "" {
			dt, err := time.Parse(time.RFC3339, creationDate)
			if err != nil {
				return errors.WithMessage(err, "failed to parse creation date")
			}
			pk.CreationDate = dt
		}

		if expirationDate != "" {
			dt, err := time.Parse(time.RFC3339, expirationDate)
			if err != nil {
				return errors.WithMessage(err, "failed to parse expiration date")
			}
			pk.ExpirationDate = dt
		}

		if len(property) > 0 {
			pk.Property = property
		}

		encoded, err := pk.Encode(macAddress)
		if err != nil {
			return errors.WithMessage(err, "failed to encode product key")
		}

		fmt.Println(encoded)
		return nil
	},
}

var nonjsonBruteForceCmd = &cobra.Command{
	Use:   "bruteforce PRODUCT_KEY",
	Short: "Find the MAC address associated with a non-JSON product key by brute force",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		productKey := args[0]

		fmt.Println("searching for mac address ...")

		mac, err := nonjson.BruteForceProductKeyMACAddress(productKey)
		if err != nil {
			return err
		}

		fmt.Printf("found match! mac = '%s'\n", mac)
		return nil
	},
}

var nonjsonListSKUCmd = &cobra.Command{
	Use:   "listswid",
	Short: "Get a list of software identifiers that can be used in product keys",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 3, 1, 2, ' ', 0)
		fmt.Fprintf(w, "License SKU\tDisplay Name\tID\n")
		fmt.Fprintf(w, "-----------\t------------\t--\n")
		for _, swid := range nonjson.SoftwareIdentifiers.List() {
			fmt.Fprintf(w, "%v\t%v\t%v\n", swid.SKU, swid.DisplayName, swid.ID)
		}
		w.Flush()
	},
}
