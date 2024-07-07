package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/RaghavSood/btcsupply/address"
	"github.com/RaghavSood/collectibles/bitcoinrpc"
	"github.com/RaghavSood/collectibles/bitcoinrpc/types"
	"github.com/RaghavSood/collectibles/electrum"
	"github.com/rs/zerolog"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func main() {
	txCache := make(map[string]types.TransactionDetail)

	eclient, err := electrum.NewElectrum()
	if err != nil {
		fmt.Println("Error creating Electrum client", err)
		return
	}

	bclient := bitcoinrpc.NewRpcClient(os.Getenv("BITCOIND_HOST"), os.Getenv("BITCOIND_USER"), os.Getenv("BITCOIND_PASS"))

	addressListFile := os.Args[1]
	addressFile, err := os.ReadFile(addressListFile)
	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}

	addressList := strings.Split(string(addressFile), "\n")

	fmt.Printf("Found %d addresses\n", len(addressList))

	for _, addr := range addressList {
		if addr == "" {
			continue
		}

		if addr[len(addr)-1] == '\r' {
			addr = addr[:len(addr)-1]
		}

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

		var firstHeight int64
		var firstTx string
		var firstValue string
		txValueFrequecy := make(map[string]int)

		if len(txids) > 0 {
			firstTx = txids[0]
			firstHeight = heights[0]

			var tx types.TransactionDetail
			if _, ok := txCache[firstTx]; !ok {
				tx, err = bclient.GetTransaction(firstTx)
				if err != nil {
					fmt.Println("Error getting transaction", firstTx, err)
					return
				}
				txCache[firstTx] = tx
			} else {
				tx = txCache[firstTx]
			}

			for _, vout := range tx.Vout {
				if script == vout.ScriptPubKey.Hex {
					firstValue = vout.Value.String()
				}
				txValueFrequecy[vout.Value.String()]++
			}
		}

		fmt.Printf("%s,%s,%d,%s,%s\n", addr, firstTx, firstHeight, firstValue, mapSummary(txValueFrequecy))
	}
}

func mapSummary(m map[string]int) string {
	if len(m) == 0 {
		return ""
	}
	var summary []string
	var total int
	for k, v := range m {
		summary = append(summary, fmt.Sprintf("%dx%s", v, k))
		total += v
	}
	return fmt.Sprintf("%d= %s", total, strings.Join(summary, " and "))
}
