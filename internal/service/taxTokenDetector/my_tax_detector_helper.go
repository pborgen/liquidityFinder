package taxTokenDetector

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/pborgen/liquidityFinder/internal/myConfig"
)

type FunctionFromJson struct {
	Selector        string `json:"selector"`
	BytecodeOffset  int    `json:"bytecode_offset"`
	Arguments       string `json:"arguments"`
	StateMutability string `json:"state_mutability"`
}

type FunctionSignatureFromJson struct {
	Selector string `json:"selector"`
	Name     string `json:"name"`
}

type StorageRecordFromJson struct {
	Slot        string `json:"slot"`
	Offset      int    `json:"offset"`
	Type        string `json:"type"`
	Reads       []string `json:"reads"`
	Writes      []string `json:"writes"`
}


type ByteCodeDataFromJson struct {
	Functions []FunctionFromJson `json:"functions"`
	StorageRecord   []StorageRecordFromJson `json:"storage"`
}

type MyTaxDetectorHelper struct {
	FunctionSignaturesMap map[string]string
}

const (
	myTransferFromFunctionSelector = "23b872dd"
	myTransferFunctionSelector    = "a9059cbb"
)

var functionSignaturesList = []FunctionSignatureFromJson{
	{"07ab1f8a", "ReflectionFee()"},
	{"098729db", "ReflectionsFromContractBalance(uint256)"},
	{"0a3cd7d2", "ReflectionReward(uint256)"},
	{"0eb8f9de", "ReflectionDOGE(uint256)"},
	{"0ed95d6d", "ReflectionSFee()"},
	{"10cd3ff2", "ReflectionRewardsBake(uint256)"},
	{"164f3986", "ReflectionRATE()"},
	{"17918c7f", "ReflectionADA(uint256)"},
	{"18eab04f", "ReflectionsAvax()"},
	{"1948363b", "Reflections()"},
	{"1d872333", "ReflectionSellFee()"},
	{"20b997af", "Reflection(address,uint256)"},
	{"22e74946", "ReflectionsBurn(address[])"},
	{"245506c1", "ReflectionFeeSettings()"},
	{"2b2ba854", "Reflection()"},
	{"2ee63882", "ReflectionRewardsCake(uint256)"},
	{"33f10eaa", "Reflection(address[])"},
	{"3406a0ca", "ReflectionTokens(uint256)"},
	{"40b25e5b", "Reflections(address,uint256)"},
	{"4685e527", "ReflectionFeeFromToken(uint256,bool)"},
	{"4a4ca7b7", "Reflection_TOKEN()"},
	{"4c77a166", "ReflectionRedistribution(address,uint256)"},
	{"4cc7b2c2", "ReflectionRate()"},
	{"535d5efa", "Reflections(address[])"},
	{"57e54478", "ReflectionsTx(address,uint256)"},
	{"5dcf60d9", "ReflectionsClaim(address,uint256)"},
	{"63decc93", "Reflection_Fee()"},
	{"650ec0bf", "ReflectionsShare(address,uint256)"},
	{"68d5d84a", "ReflectionsGPT4(address,uint256)"},
	{"6c469882", "ReflectionsVaults(uint256)"},
	{"71d3951f", "ReflectionWallet()"},
	{"73fdc769", "ReflectionsFeeFromToken(uint256,bool)"},
	{"7b1d6071", "ReflectionsMarketing(address,uint256)"},
	{"7d0cc617", "ReflectionContract()"},
	{"87b0c5d1", "Reflections(address[],bool)"},
	{"8edbfd9a", "ReflectionFeeSell()"},
	{"8f329884", "ReflectionsAllowed(address[])"},
	{"93c93792", "Reflections(uint256)"},
	{"96d39e69", "ReflectionAddress()"},
	{"9838d194", "ReflectionRewardsFee()"},
	{"9e01e1b1", "Reflections(address,bool)"},
	{"a1c6f281", "Reflections(address)"},
	{"a8c07e10", "ReflectionRewardsAda(uint256)"},
	{"b4df808a", "ReflectionTaxFee()"},
	{"badc96f0", "Reflection(uint256)"},
	{"bdad86c9", "ReflectionC(uint256)"},
	{"bdb6ebb8", "ReflectionRateFromWallet()"},
	{"bff48c5d", "Reflectionnumbers(bool)"},
	{"c4fe93c7", "ReflectionsToken(address[])"},
	{"c531096c", "ReflectionFeeonBuy()"},
	{"cc490302", "ReflectionToken()"},
	{"e25aab9b", "ReflectionFeeRewards(uint256)"},
	{"edb46998", "ReflectionFeeonSell()"},
	{"f3362572", "ReflectionPaused()"},
	{"f8725843", "ReflectionToken(uint256)"},
	{"f8d26425", "ReflectionFeeBuy()"},
	{"fb35c8b4", "ReflectionBuyFee()"},
	{"fc669093", "ReflectionPercentage()"},
	{"00281dc1", "reflectionFeeOnSelling()"},
	{"02551b28", "reflectionFrompanthertoken(uint256,bool)"},
	{"0296b2ad", "reflection(address,address,uint256)"},
	{"0479dc66", "reflectionFromToken(uint256,uint8)"},
	{"07f04f27", "reflectionFromTokenTransfer(uint256,bool)"},
	{"0826bf0e", "reflectionFromPOPO(address,address,uint256,uint256)"},
	{"0a0dc697", "reflectionFromMIDIT(uint256,bool)"},
	{"0db246ff", "reflectionontimer()"},
	{"0e222dc8", "reflectionOf(address)"},
	{"11329b71", "reflections2()"},
	{"128dc0a0", "reflectionWithdrawn(uint256)"},
	{"1597dd5e", "reflectionRoyalty()"},
	{"18411e9a", "reflectionFromTokenSell(uint256,bool)"},
	{"1a842872", "reflectionsClaimed(address)"},
	{"1c36cccd", "reflectionTokensAddresses(uint256)"},
	{"1d1fb0c9", "reflectionPercentage()"},
	{"20b7be1e", "reflectionTaxOf(address)"},
	{"216f3a0b", "reflectionFromToken(address,uint256)"},
	{"26d4585f", "reflectionEnable()"},
	{"27e1519a", "reflectionMintBalance()"},
	{"290772b6", "reflectionDivisor()"},
	{"2e0d9936", "reflectionDistributor()"},
	{"2f239880", "reflectionsfee()"},
	{"30c84cd6", "reflectionFrompolymaxclubtoken(uint256,bool)"},
	{"314a5184", "reflectionTokenBalanceOf(address)"},
	{"316dd9ee", "reflectionFromROSE(uint256,bool)"},
	{"33ae05f3", "reflectiontoken(address,uint256)"},
	{"348a5b73", "reflectionCheck()"},
	{"34a4a860", "reflectionFromTokenBuy(uint256,bool)"},
	{"376de5b7", "reflectionFromToken(address,uint256,bool)"},
	{"38f32480", "reflectionMultiplier()"},
	{"3d064e1f", "reflectionFromMastersCoin(address,address,uint256,uint256)"},
	{"409d0566", "reflectionFeeBuy()"},
	{"40d604a2", "reflectionsFee()"},
	{"42b2bf9e", "reflectionsVaultInLocker(address)"},
	{"47883741", "reflectionFromPolyField(uint256,bool)"},
	{"4e3a4ba5", "reflections4()"},
	{"4edc3ad6", "reflectionBuyFee()"},
	{"4f01e881", "reflectionsPaidToHolderAccount()"},
	{"4f61ded4", "reflectionTransfer(address,uint256)"},
	{"5025de6b", "reflectionFromMadridOX(address,address,uint256,uint256)"},
	{"534207d5", "reflectionFromAgoraSpace(uint256,bool)"},
	{"5a630ca0", "reflection(uint256,address)"},
	{"5c24f9c0", "reflectionPercentCalculate(address)"},
	{"5cad33b9", "reflectionFeesdisabled()"},
	{"6173e1f6", "reflections(address,address[])"},
	{"63b24be3", "reflectionsPerShareAccuracyFactor()"},
	{"6415bac3", "reflectionEnabled()"},
	{"65b23e07", "reflectionFeedisabled()"},
	{"69d2567b", "reflectionFromPolyOctian(uint256,bool)"},
	{"6ef83671", "reflectiontoken(uint256)"},
	{"6f2240af", "reflectionsOwned(address)"},
	{"6f8df198", "reflectionFromtokenswap(uint256,bool)"},
	{"71fc79c1", "reflectionBurn(uint256)"},
	{"72414ddc", "reflectionFrompolythoreum(uint256,bool)"},
	{"73dc6da5", "reflectionFromBabyRydell(address,address,uint256,uint256)"},
	{"794ae2ef", "reflectionsBalance()"},
	{"7f6a5177", "reflectionReleased()"},
	{"88ccb44b", "reflectionFromMaterial(address,address,uint256,uint256)"},
	{"9077ac59", "reflectionsPerShare()"},
	{"94316c7b", "reflectionPenaltyFee_sec()"},
	{"9d8c36e7", "reflectionTenthPercent()"},
	{"9f113d1b", "reflections(address)"},
	{"a06108d7", "reflectionsEarnInLocker(address)"},
	{"a0e085be", "reflectionFromPassWay(address,address,uint256,uint256)"},
	{"a3b86b64", "reflectionFrompikachutoken(uint256,bool)"},
	{"a67513bd", "reflectionStep()"},
	{"a9223fe1", "reflectionFromTokenWithDeductTransferFee(uint256,uint256)"},
	{"ac19ad9b", "reflectionBase()"},
	{"ac45ed1c", "reflectionFeeSell()"},
	{"b3c5fb73", "reflectionReleaseTimestamp()"},
	{"bb054705", "reflectionPercentage(uint256)"},
	{"bb3e3872", "reflectionUpdate(address)"},
	{"bc6a9e01", "reflectionToken()"},
	{"bf266469", "reflectionFromTokenInTiers(uint256,uint256,bool)"},
	{"bfc49e4b", "reflectionFromPolyElephant(uint256,bool)"},
	{"c07ea818", "reflectionsToPay(address)"},
	{"c3c5bd63", "reflectionFromspidertoken(uint256,bool)"},
	{"c9140c04", "reflectionRewardBTCB()"},
	{"cd8307c9", "reflections_exclude_Wallet(address)"},
	{"cdd597a2", "reflectionSleepage()"},
	{"cff83b67", "reflectionTriggerTime()"},
	{"d23e25e8", "reflectionFeeReceiver()"},
	{"d61b5b3a", "reflectionpoolEligible(address)"},
	{"d63ace43", "reflectionBasis()"},
	{"dbf08001", "reflectionRewardsFee()"},
	{"dd073829", "reflectionFeeOnBuying()"},
	{"dde508d2", "reflectionsB(address)"},
	{"e2708c01", "reflectionFromDuaLipaCoin(address,address,uint256,uint256)"},
	{"e60a54b3", "reflectionsWithdrawn(address)"},
	{"ef0b7329", "reflections3()"},
	{"ef656ce3", "reflectionFromStudent(address,address,uint256,uint256)"},
	{"f039526f", "reflectionFromElonMuskCoin(address,address,uint256,uint256)"},
	{"f03d7939", "reflectionOwner()"},
	{"0xf404bce2", "reflectionFromPolyRock(uint256,bool)"},
	{"f594aaf9", "reflectionTriggerAmount()"},
	{"f88a2fad", "reflections_include_Wallet(address)"},
	{"faaf6287", "reflections1()"},
	{"fcef8867", "reflectionFee(uint256)"},
}

func NewMyTaxDetectorHelper() (*MyTaxDetectorHelper, error) {
	
	functionSignaturesMap := make(map[string]string)
	for _, signature := range functionSignaturesList {
		functionSignaturesMap[signature.Selector] = signature.Name
	}

    return &MyTaxDetectorHelper{
		FunctionSignaturesMap: functionSignaturesMap,
	}, nil
}

func (h *MyTaxDetectorHelper) IsTaxToken(bytecode string) (bool, error) {
	byteCodeData, err := h.getByteCodeData(bytecode)
	if err != nil {
		panic(err)
	}

	readsCounter := 0
	writesCounter := 0

	// Check how many storage slots are used by the transfer functions
	for _, storage := range byteCodeData.StorageRecord {
		reads := storage.Reads
		
		for _, read := range reads {

			_, exists := h.FunctionSignaturesMap[read]
			if exists {
				return true, nil
			}

			if read == myTransferFromFunctionSelector || read == myTransferFunctionSelector {
				readsCounter++
			}
		}

		writes := storage.Writes
		for _, write := range writes {

			_, exists := h.FunctionSignaturesMap[write]
			if exists {
				return true, nil
			}

			if write == myTransferFromFunctionSelector || write == myTransferFunctionSelector {
				writesCounter++
			}
		}
	}

	total := readsCounter + writesCounter

	if total > 10 {
		return true, nil
	} else {
		return false, nil
	}

}

func (h *MyTaxDetectorHelper) getByteCodeData(bytecode string) (ByteCodeDataFromJson, error) {

	baseDir := myConfig.GetInstance().GetBaseDir()
	baseDir = baseDir + "/apps/python/parsebytecode" // Change this to your actual directory

	err := os.Chdir(baseDir)
	if err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
		return ByteCodeDataFromJson{}, err
	}

    // Define the Python script and arguments
	cmd := exec.Command("bash", "-c", "source venv/bin/activate && python3 parsebytecode.py " + bytecode)
 
    // Run the command and capture the output
    output, err := cmd.Output()
    if err != nil {
        fmt.Println("Error:", err)
        return ByteCodeDataFromJson{}, err
    }

	var data ByteCodeDataFromJson

	err = json.Unmarshal(output, &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return ByteCodeDataFromJson{}, err
	}

    // Print the output
    return data, nil
}

func (h *MyTaxDetectorHelper) countStorageSlotsForSignature(signature string, StorageRecord []StorageRecordFromJson) int {
	countReads := 0
	countWrites := 0
	for _, record := range StorageRecord {
		for _, read := range record.Reads {
			if read == signature {
				countReads++
			}
		}
		for _, write := range record.Writes {
			if write == signature {
				countWrites++
			}
		}
	}
	return countReads + countWrites
}
