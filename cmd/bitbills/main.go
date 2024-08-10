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
)

const (
	ONE    = 100000000
	FIVE   = 500000000
	TEN    = 1000000000
	TWENTY = 2000000000
)

var addressList = []string{"1FhTe1bMtoKHDbNw13v3BQ9sb2kFTagaRH"}

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
			height := heights[i]
			var tx types.TransactionDetail

			tx, err = bclient.GetTransaction(txid)
			if err != nil {
				fmt.Println("Error getting transaction", txid, err)
				continue // Skip this transaction and continue with the next
			}

			outputs := categorizeOutputs(tx.Vout, addr)

			fmt.Printf("Address: %s, Transaction ID: %s, Height: %d\n", addr, txid, height)
			for value, addrs := range outputs {
				fmt.Printf("  Value: %.8f BTC, Addresses: %s\n", float64(value)/1e8, strings.Join(addrs, ", "))
			}
		}
	}
}

func categorizeOutputs(vouts []types.Vout, addr string) map[int64][]string {
	categories := make(map[int64][]string)

	for _, vout := range vouts {
		if vout.ScriptPubKey.Address == addr {
			continue
		}
		value := ctypes.FromBTCString(ctypes.BTCString(vout.Value)).Int64()
		categories[value] = append(categories[value], vout.ScriptPubKey.Address)
	}
	return categories
}
