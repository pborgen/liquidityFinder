package myConst

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const (

	WPLS_ADDRESS_STRING = "0xA1077a294dDE1B09bB078844df40758a5D0f9a27"
	WETH_ADDRESS_STRING = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"

	ONE_WPLS_STRING = "1000000000000000000"
	EIGHTEEN_ZEROS_STRING = "000000000000000000"
	NINE_ZEROS_STRING = "000000000"
)

var (
	wpls_address common.Address
	weth_address common.Address
	one_wpls_bigint *big.Int
	one_thousand_wpls_bigint *big.Int
	ten_thousand_wpls_bigint *big.Int
	
	one_hundred_thousand_wpls_bigint *big.Int
	one_million_wpls_bigint *big.Int
	ten_million_wpls_bigint *big.Int
	one_hundred_million_wpls_bigint *big.Int
	one_trillion_wpls_bigint *big.Int

	maxUint256 *big.Int
	ten_percent_of_max_uint256 *big.Int
	ninety_percent_of_max_uint256 *big.Int
)

func init() {
	wpls_address = common.HexToAddress(WPLS_ADDRESS_STRING)
	weth_address = common.HexToAddress(WETH_ADDRESS_STRING)
	one_wpls_bigint, _ = big.NewInt(0).SetString(ONE_WPLS_STRING, 10)
	one_thousand_wpls_bigint = new(big.Int).Mul(one_wpls_bigint, big.NewInt(1000))
	ten_thousand_wpls_bigint = new(big.Int).Mul(one_wpls_bigint, big.NewInt(10000))
	one_hundred_thousand_wpls_bigint = new(big.Int).Mul(one_wpls_bigint, big.NewInt(100000))
	one_million_wpls_bigint = new(big.Int).Mul(one_wpls_bigint, big.NewInt(1000000))
	ten_million_wpls_bigint = new(big.Int).Mul(one_wpls_bigint, big.NewInt(10000000))
	one_hundred_million_wpls_bigint = new(big.Int).Mul(one_wpls_bigint, big.NewInt(100000000))
	one_trillion_wpls_bigint = new(big.Int).Mul(one_wpls_bigint, big.NewInt(1000000000000))

	maxUint256 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

	ten_percent_of_max_uint256 = new(big.Int).Div(maxUint256, big.NewInt(10))
	ninety_percent_of_max_uint256 = new(big.Int).Sub(maxUint256, ten_percent_of_max_uint256)
}

func GetMaxUint256() *big.Int {
	return new(big.Int).Set(maxUint256)
}

func GetTenPercentOfMaxUint256() *big.Int {
	return new(big.Int).Set(ten_percent_of_max_uint256)
}

func GetNinetyPercentOfMaxUint256() *big.Int {
	return new(big.Int).Set(ninety_percent_of_max_uint256)
}

func GetWplsAddress() common.Address {
	return wpls_address
}

func GetOneWplsBigint() *big.Int {
	return new(big.Int).Set(one_wpls_bigint)
}

func GetOneThousandWplsBigint() *big.Int {
	return new(big.Int).Set(one_thousand_wpls_bigint)
}

func GetTenThousandWplsBigint() *big.Int {
	return new(big.Int).Set(ten_thousand_wpls_bigint)
}

func GetOneHundredThousandWplsBigint() *big.Int {
	return new(big.Int).Set(one_hundred_thousand_wpls_bigint)
}

func GetOneMillionWplsBigint() *big.Int {
	return new(big.Int).Set(one_million_wpls_bigint)
}

func GetTenMillionWplsBigint() *big.Int {
	return new(big.Int).Set(ten_million_wpls_bigint)
}

func GetOneHundredMillionWplsBigint() *big.Int {
	return new(big.Int).Set(one_hundred_million_wpls_bigint)
}

func GetOneTrillionWplsBigint() *big.Int {
	return new(big.Int).Set(one_trillion_wpls_bigint)
}

func GetWethAddress() common.Address {
	return weth_address
}

