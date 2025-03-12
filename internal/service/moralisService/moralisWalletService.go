package moralisService

import (
	"encoding/json"
	"io"

	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pborgen/liquidityFinder/internal/myConfig"
)

var apiKey = ""
var baseUrl = ""
var cfg = myConfig.GetInstance()

func init() {
	apiKey = cfg.GetMoralisApiKey()
	baseUrl = cfg.GetMoralisBaseUrl()
}

type MoralisWallet struct {
	Page        int `json:"page"`
	Page_size   int `json:"page_size"`
	MoralisWalletResult   []MoralisWalletResult  `json:"result"`
}

type MoralisWalletResult struct {
	IsNativeToken 		bool `json:"native_token"`
	Token_address 		string `json:"token_address"`
	Token_name    		string `json:"name"`
	Token_symbol  		string `json:"symbol"`
	Token_decimals 		int    `json:"decimals"`
	Total_Usd_value     float64 `json:"usd_value"`
	Token_usd_price     float64 `json:"usd_price"`
	Token_balance     	string `json:"balance"`
}

func GetWalletBalances(walletAddress common.Address) (*MoralisWallet, error) {

  	url := baseUrl + "/wallets/" + walletAddress.String() + "/tokens?chain=pulse"

	payload := strings.NewReader("")

	req, _ := http.NewRequest("GET", url, payload)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var jsonData MoralisWallet
	if err := json.Unmarshal(body, &jsonData); err != nil {
		return nil, err
	}

	return &jsonData, nil

}