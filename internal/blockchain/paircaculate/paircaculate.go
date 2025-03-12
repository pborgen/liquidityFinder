package paircaculate

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// More complete version with configurable init code hash:
type PairCalculator struct {
    Factory      common.Address
    InitCodeHash []byte
}

func NewPairCalculator(factory common.Address, initCodeHash string) *PairCalculator {
    return &PairCalculator{
        Factory:      factory,
        InitCodeHash: common.FromHex(initCodeHash),
    }
}

func (pc *PairCalculator) PairFor(tokenA common.Address, tokenB common.Address) common.Address {
    // Sort tokens
    var token0, token1 common.Address
    if tokenA.Hex() < tokenB.Hex() {
        token0 = tokenA
        token1 = tokenB
    } else {
        token0 = tokenB
        token1 = tokenA
    }

    // Inner hash
    innerPacked := append(token0.Bytes(), token1.Bytes()...)
    innerHash := crypto.Keccak256(innerPacked)

    // Final packed data
    finalPacked := []byte{255} // hex'ff'
    finalPacked = append(finalPacked, pc.Factory.Bytes()...)
    finalPacked = append(finalPacked, innerHash...)
    finalPacked = append(finalPacked, pc.InitCodeHash...)

    // Final hash to address
    finalHash := crypto.Keccak256(finalPacked)
    return common.BytesToAddress(finalHash[12:])
}
