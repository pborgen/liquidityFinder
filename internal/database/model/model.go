package model

import (

	"github.com/pborgen/liquidityFinder/internal/database/model/dex"
	"github.com/pborgen/liquidityFinder/internal/database/model/erc20"
	"github.com/pborgen/liquidityFinder/internal/database/model/network"
	"github.com/pborgen/liquidityFinder/internal/types"
)

type MyModel interface {
	dex.ModelDex | erc20.ModelERC20 | network.NetworkModel | types.ModelPair 
}
