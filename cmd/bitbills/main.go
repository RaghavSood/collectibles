package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/RaghavSood/btcsupply/address"
	"github.com/RaghavSood/collectibles/bitcoinrpc"
	"github.com/RaghavSood/collectibles/bitcoinrpc/types"
	"github.com/RaghavSood/collectibles/electrum"
	ctypes "github.com/RaghavSood/collectibles/types"
	"github.com/rs/zerolog"
	"golang.org/x/exp/maps"
)

const (
	ONE       = 100000000
	FIVE      = 500000000
	TEN       = 1000000000
	TWENTY    = 2000000000
	MAX_DEPTH = 7 // Maximum depth for recursive scanning
)

var addressList = []string{"1FhTe1bMtoKHDbNw13v3BQ9sb2kFTagaRH"}

var addressGraph = make(map[string]map[string]map[string]struct{})
var visitedAddresses = make(map[string]struct{})

var selectedAddresses = make(map[string]bool)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func main() {
	eclient, err := electrum.NewElectrum()
	if err != nil {
		fmt.Println("Error creating Electrum client", err)
		return
	}

	bclient := bitcoinrpc.NewRpcClient(os.Getenv("BITCOIND_HOST"), os.Getenv("BITCOIND_USER"), os.Getenv("BITCOIND_PASS"))

	for _, addr := range addressList {
		script, err := address.AddressToScript(addr)
		if err != nil {
			fmt.Printf("Error converting address %s to script: %s\n", addr, err)
			return
		}

		txids, heights, err := eclient.GetScriptHistory(script)
		if err != nil {
			fmt.Println("Error getting balance for", addr, err)
			return
		}

		// Check all transaction IDs
		for i, txid := range txids {
			// BitBills ceased production on May 15, 2012
			if heights[i] >= 185000 {
				continue
			}

			var tx types.TransactionDetail

			tx, err = bclient.GetTransaction(txid)
			if err != nil {
				fmt.Println("Error getting transaction", txid, err)
				continue // Skip this transaction and continue with the next
			}

			spendingTransaction := checkInputs(tx.Vin, addr)
			if !spendingTransaction {
				continue // Skip this transaction and continue with the next
			}

			outputs := categorizeOutputs(tx.Vout, addr)

			// Print categorized addresses for the root address
			printCategorizedOutputs(addr, txid, outputs)

			// If there are categorized outputs, find the change output
			if len(outputs) > 0 {
				changeAddress := findChangeOutput(tx.Vout, outputs)
				if changeAddress != "" {
					trackAddress(addr, txid, changeAddress)
					fmt.Printf("Change output found: %s -> %s (Transaction ID: %s)\n", addr, changeAddress, txid)
					// Recursively scan the change address
					scanChangeAddress(changeAddress, bclient, eclient, 1) // Start at depth 1
				}
			}
		}
	}

	addressesFound := maps.Keys(selectedAddresses)
	for _, addr := range addressesFound {
		fmt.Println(addr)
	}

	fmt.Println("Address graph:")
	for inputAddr, txs := range addressGraph {
		fmt.Printf("  %s:\n", inputAddr)
		for txid, outputs := range txs {
			fmt.Printf("    %s:\n", txid)
			for outputAddr := range outputs {
				fmt.Printf("      %s\n", outputAddr)
			}
		}
	}
}

func checkInputs(vins []types.Vin, addr string) bool {
	for _, vin := range vins {
		if vin.Prevout.ScriptPubKey.Address == addr {
			return true
		}
	}
	return false
}

func categorizeOutputs(vouts []types.Vout, addr string) map[int64][]string {
	categories := make(map[int64][]string)

	for _, vout := range vouts {
		if vout.ScriptPubKey.Address == addr {
			continue
		}
		value := ctypes.FromBTCString(ctypes.BTCString(vout.Value)).Int64()
		categories[value] = append(categories[value], vout.ScriptPubKey.Address)

		if value == ONE || value == FIVE || value == TEN || value == TWENTY {
			selectedAddresses[vout.ScriptPubKey.Address] = true
		}
	}
	return categories
}

func findChangeOutput(vouts []types.Vout, categorized map[int64][]string) string {
	for _, vout := range vouts {
		value := ctypes.FromBTCString(ctypes.BTCString(vout.Value)).Int64()
		// If the output value is not in the categorized values, it is the change output
		if value != ONE && value != FIVE && value != TEN && value != TWENTY {
			return vout.ScriptPubKey.Address
		}
	}
	return ""
}

func trackAddress(inputAddr, txid, outputAddr string) {
	if addressGraph[inputAddr] == nil {
		addressGraph[inputAddr] = make(map[string]map[string]struct{})
	}
	if addressGraph[inputAddr][txid] == nil {
		addressGraph[inputAddr][txid] = make(map[string]struct{})
	}
	addressGraph[inputAddr][txid][outputAddr] = struct{}{}
}

func scanChangeAddress(addr string, bclient *bitcoinrpc.RpcClient, eclient *electrum.Electrum, depth int) {
	// Check if the address has already been visited
	if _, visited := visitedAddresses[addr]; visited {
		return
	}
	// Mark the address as visited
	visitedAddresses[addr] = struct{}{}

	// Stop if the maximum depth is reached
	if depth > MAX_DEPTH {
		return
	}

	script, err := address.AddressToScript(addr)
	if err != nil {
		fmt.Printf("Error converting change address %s to script: %s\n", addr, err)
		return
	}

	txids, heights, err := eclient.GetScriptHistory(script)
	if err != nil {
		fmt.Println("Error getting balance for change address", addr, err)
		return
	}

	for i, txid := range txids {
		// BitBills ceased production on May 15, 2012
		if heights[i] >= 185000 {
			continue
		}

		var tx types.TransactionDetail

		tx, err = bclient.GetTransaction(txid)
		if err != nil {
			fmt.Println("Error getting transaction", txid, err)
			continue
		}

		spendingTransaction := checkInputs(tx.Vin, addr)
		if !spendingTransaction {
			continue // Skip this transaction and continue with the next
		}

		outputs := categorizeOutputs(tx.Vout, addr)

		// Print categorized addresses for the change address
		printCategorizedOutputs(addr, txid, outputs)

		if len(outputs) > 0 {
			changeAddress := findChangeOutput(tx.Vout, outputs)
			if changeAddress != "" {
				trackAddress(addr, txid, changeAddress)
				fmt.Printf("Change output found: %s -> %s (Transaction ID: %s)\n", addr, changeAddress, txid)
				// Recursively scan the change address, increasing the depth
				scanChangeAddress(changeAddress, bclient, eclient, depth+1)
			}
		}
	}
}

func printCategorizedOutputs(addr string, txid string, categorized map[int64][]string) {
	fmt.Printf("Categorized outputs for address %s (Transaction ID: %s):\n", addr, txid)
	for value, addrs := range categorized {
		fmt.Printf("  Value: %.8f BTC, Addresses: %s\n", float64(value)/1e8, strings.Join(addrs, ", "))
	}
}
