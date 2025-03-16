package myUtil

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	blockchainclient "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainClient"
	"github.com/pborgen/liquidityFinder/myConst"
)

func MyAsync[T any](f func() T) chan T {
	ch := make(chan T)
	go func() {
		ch <- f()
	}()
	return ch
}

func MyAsyncWithError[T any](f func() (T, error)) chan T {
	ch := make(chan T)
	go func() {
		result, err := f()
		if err != nil {
			panic(err)
		}
		ch <- result
	}()
	return ch
}

func Wei2EthAsString(amount *big.Int) string {

	compact_amount := big.NewInt(0)
	reminder := big.NewInt(0)
	divisor := big.NewInt(1e18)
	compact_amount.QuoRem(amount, divisor, reminder)
	return fmt.Sprintf("%v.%018s", compact_amount.String(), reminder.String())
}

func IsTransactionSuccess(receipt *types.Receipt) (bool, error) {

	if receipt.Status == types.ReceiptStatusSuccessful {
		return true, nil
	} else {
		return false, nil
	}
}

func FromAddressFromTransaction(transaction *types.Transaction) (common.Address, error) {

	from, err := types.Sender(types.NewEIP155Signer(transaction.ChainId()), transaction)
	if err != nil {
		return common.Address{}, err
	}

	return from, nil
}

func IsWPLS(address common.Address) bool {
	return address == common.HexToAddress(myConst.WPLS_ADDRESS_STRING)
}

func SortTokens(token0 common.Address, token1 common.Address) (common.Address, common.Address) {
    // Check for identical addresses
    if token0 == token1 {
        return common.Address{}, common.Address{}
    }

	    // Check for zero address
	if token0 == (common.Address{}) {
		return common.Address{}, common.Address{}
	}

	// Convert address to a big.Int
	tokenAAddressBigInt := int32(binary.BigEndian.Uint32(token0.Bytes()))
	tokenBAddressBigInt := int32(binary.BigEndian.Uint32(token1.Bytes()))

	if tokenAAddressBigInt < tokenBAddressBigInt {
		return token0, token1
	} else {
		return token0, token1
	}
    // Sort tokens (tokenA < tokenB)
    // if tokenA.Hex() < tokenB.Hex() {
    //     token0 = tokenA
    //     token1 = tokenB
    // } else {
    //     token0 = tokenB
    //     token1 = tokenA
    // }

}

func HasDuplicates[T comparable](arr []T) bool {
	seen := make(map[T]bool)
	for _, num := range arr {
		if seen[num] {
			return true // Duplicate found
		}
		seen[num] = true
	}
	return false // No duplicates found
}


func FormatBigIntWithDecimals(value *big.Int, decimals int) string {

	// Convert Wei to Ether
	amountEther := new(big.Float).SetInt(value)
	amountEther.Quo(amountEther, big.NewFloat(0).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)))

	// Format the amount
	formattedAmount := fmt.Sprintf("%.4f", amountEther) // Adjust precision as needed

	formattedAmount, err := addCommasToNumberString(formattedAmount)
	if err != nil {
		return "error"
	}

	return formattedAmount
}

func GetRevertReason(from common.Address, tx *types.Transaction) (string, error) {

	client := blockchainclient.GetHttpClient()

	msg := ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}

	res, err := client.CallContract(context.Background(), msg, nil)

	if err != nil {
		if strings.Contains(err.Error(), "reverted") {
			return err.Error(), nil
		}
		return "", err
	}

	return string(res), nil
}

func parseRevertReason(input []byte) (string, error) {
	if len(input) < 4 {
		return "", fmt.Errorf("invalid input")
	}

	// methodID := input[:4]
	inputData := input[4:]

	parsedRevert, err := abi.JSON(strings.NewReader(abiRevert))
	if err != nil {
		return "", err
	}

	var reason string
	err = parsedRevert.UnpackIntoInterface(&reason, "Error", inputData)
	if err != nil {
		return "", err
	}

	return reason, nil
}

const abiRevert = `[{ "name": "Error", "type": "function", "inputs": [ { "type": "string" } ] }]`

func addCommasToNumberString(numberStr string) (string, error) {
	// Split the number into integer and fractional parts
	parts := strings.Split(numberStr, ".")
	integerPart := parts[0]
	var fractionalPart string
	if len(parts) > 1 {
		fractionalPart = parts[1]
	}

	// Format the integer part with commas
	formattedIntegerPart := formatWithCommas(integerPart)

	// Recombine the integer and fractional parts
	if fractionalPart != "" {
		return formattedIntegerPart + "." + fractionalPart, nil
	}
	return formattedIntegerPart, nil
}

func formatWithCommas(integerStr string) string {
	n := len(integerStr)
	if n <= 3 {
		return integerStr
	}

	// Reverse the string to make it easier to insert commas
	reversed := reverseString(integerStr)
	var result strings.Builder

	// Insert commas every three digits
	for i, char := range reversed {
		if i > 0 && i%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(char)
	}

	// Reverse the result to get the final formatted string
	return reverseString(result.String())
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func ReadFileToString(fileName string) (string, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func IsWithinNumeric78Range(value *big.Int) bool {
	if value == nil {
		return false
	}

	// NUMERIC(78,0) max value is 10^78 - 1
	maxValue := new(big.Int)
	maxValue.Exp(big.NewInt(10), big.NewInt(78), nil)
	maxValue.Sub(maxValue, big.NewInt(1))

	// NUMERIC(78,0) min value is -(10^78 - 1)
	minValue := new(big.Int)
	minValue.Neg(maxValue)

	// Check if value is within range
	return value.Cmp(maxValue) <= 0 && value.Cmp(minValue) >= 0
}

// Check if address exists in slice
func ContainsAddress(addresses []common.Address, search common.Address) bool {
	for _, addr := range addresses {
		if addr == search {
			return true
		}
	}
	return false
}
