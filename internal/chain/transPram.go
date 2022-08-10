package chain

import "github.com/centrifuge/go-substrate-rpc-client/v4/types"

var (
	//trade
	BuySpaceTransactionName   = "FileBank.buy_space"
	UploadFileTransactionName = "FileBank.upload"
	DeleteFileTransactionName = "FileBank.delete_file"
	UploadDeclaration         = "FileBank.upload_declaration"

	//find
	PurchasedSpaceChainModule = "FileBank"
	PurchasedPackage          = "PurchasedPackage"

	FindPriceChainModule  = "FileBank"
	FindPriceModuleMethod = "UnitPrice"

	FindFileChainModule  = "FileBank"
	FindFileModuleMethod = []string{"File", "UserHoldFileList"}

	FindSchedulerInfoModule = "FileMap"
	FindSchedulerInfoMethod = "SchedulerMap"
)

type CessInfo struct {
	RpcAddr               string
	IdentifyAccountPhrase string
	PublicKeyOfIdentify   []byte
	TransactionName       string
	ChainModule           string
	ChainModuleMethod     string
}

type SpacePackage struct {
	Space           types.U128
	Used_space      types.U128
	Remaining_space types.U128
	Tenancy         types.U32
	Package_type    types.U8
	Start           types.U32
	Deadline        types.U32
	State           types.Bytes
}

//---FileMetaInfo
type FileMetaInfo struct {
	FileSize  types.U64
	Index     types.U32
	FileState types.Bytes
	Users     []types.AccountID
	Names     []types.Bytes
	ChunkInfo []ChunkInfo
}

type ChunkInfo struct {
	MinerId   types.U64
	ChunkSize types.U64
	BlockNum  types.U32
	ChunkId   types.Bytes
	MinerIp   types.Bytes
	MinerAcc  types.AccountID
}

type FileList struct {
	Fileid types.Bytes8 `json:"fileid"`
}
type SchedulerInfo struct {
	Ip             types.Bytes     `json:"ip"`
	Owner          types.AccountID `json:"stash_user"`
	ControllerUser types.AccountID `json:"controller_user"`
}

type UserFileList struct {
	File_hash types.Bytes
	File_size types.U64
}
