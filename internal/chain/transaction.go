package chain

import (
	"fmt"
	"math/big"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

//BuySpaceOnChain means initiating a transaction to purchase data on the chain
func (ci *CessInfo) BuySpaceOnChain(Quantity, Duration, Expected int) error {
	var (
		err         error
		accountInfo types.AccountInfo
	)
	api.getSubstrateApiSafe()
	defer func() {
		api.releaseSubstrateApi()
		err := recover()
		if err != nil {
			fmt.Printf("[Error]Recover UserHoldSpaceDetails panic :%s\n", err)
		}
	}()
	keyring, err := signature.KeyringPairFromSecret(ci.IdentifyAccountPhrase, 0)
	if err != nil {
		return errors.Wrap(err, "KeyringPairFromSecret err")
	}

	meta, err := api.r.RPC.State.GetMetadataLatest()
	if err != nil {
		return errors.Wrap(err, "GetMetadataLatest err")
	}

	c, err := types.NewCall(meta, ci.TransactionName,
		types.NewU128(*big.NewInt(int64(Quantity))),
		types.NewU128(*big.NewInt(int64(Duration))),
		types.NewU128(*big.NewInt(int64(Expected))))
	if err != nil {
		return errors.Wrap(err, "NewCall err")
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		return errors.Wrap(err, "NewExtrinsic err")
	}

	genesisHash, err := api.r.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return errors.Wrap(err, "GetBlockHash err")
	}

	rv, err := api.r.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return errors.Wrap(err, "GetRuntimeVersionLatest err")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", keyring.PublicKey)
	if err != nil {
		return errors.Wrap(err, "CreateStorageKey err")
	}

	keye, err := types.CreateStorageKey(meta, "System", "Events", nil)
	if err != nil {
		return errors.Wrap(err, "CreateStorageKey Events err")
	}

	ok, err := api.r.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return errors.Wrap(err, "GetStorageLatest err")
	}
	if !ok {
		return errors.New("GetStorageLatest return value is empty")
	}

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(keyring, o)
	if err != nil {
		return errors.Wrap(err, "Sign err")
	}

	// Do the transfer and track the actual status
	sub, err := api.r.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return errors.Wrap(err, "SubmitAndWatchExtrinsic err")
	}
	defer sub.Unsubscribe()

	timeout := time.After(10 * time.Second)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := MyEventRecords{}
				h, err := api.r.RPC.State.GetStorageRaw(keye, status.AsInBlock)
				if err != nil {
					return err
				}
				err = types.EventRecordsRaw(*h).DecodeEventRecords(meta, &events)
				if err != nil {
					return err
				}
				if events.FileBank_BuySpace != nil {
					return nil
				} else {
					return errors.New("Buy space on chain fail!")
				}
			}
		case <-timeout:
			return errors.Errorf("[%v] tx timeout", ci.TransactionName)
		default:
			time.Sleep(time.Second)
		}
	}
}

//UploadFileMetaInformation means upload file metadata to the chain
func (ci *CessInfo) UploadFileMetaInformation(fileid, filename, filehash string, ispublic bool, backups uint8, filesize uint64, downloadfee *big.Int, WalletAddress string) (string, error) {
	var (
		err         error
		accountInfo types.AccountInfo
	)
	api.getSubstrateApiSafe()
	defer func() {
		api.releaseSubstrateApi()
		err := recover()
		if err != nil {
			fmt.Printf("[Error]Recover UploadFileMetaInformation panic :%s\n", err)
		}
	}()
	keyring, err := signature.KeyringPairFromSecret(ci.IdentifyAccountPhrase, 0)
	if err != nil {
		return "", errors.Wrap(err, "KeyringPairFromSecret err")
	}

	meta, err := api.r.RPC.State.GetMetadataLatest()
	if err != nil {
		return "", errors.Wrap(err, "GetMetadataLatest err")
	}

	c, err := types.NewCall(
		meta,
		ci.TransactionName,
		types.NewBytes([]byte(WalletAddress)),
		types.NewBytes([]byte(filename)),
		types.NewBytes([]byte(fileid)),
		types.NewBytes([]byte(filehash)),
		types.NewBool(ispublic),
		types.NewU8(backups),
		types.NewU64(filesize),
		types.NewU128(*downloadfee),
	)
	if err != nil {
		return "", errors.Wrap(err, "NewCall err")
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		return "", errors.Wrap(err, "NewExtrinsic err")
	}

	genesisHash, err := api.r.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return "", errors.Wrap(err, "GetBlockHash err")
	}

	rv, err := api.r.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return "", errors.Wrap(err, "GetRuntimeVersionLatest err")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", keyring.PublicKey)
	if err != nil {
		return "", errors.Wrap(err, "CreateStorageKey err")
	}

	keye, err := types.CreateStorageKey(meta, "System", "Events", nil)
	if err != nil {
		return "", errors.Wrap(err, "CreateStorageKey Events err")
	}

	ok, err := api.r.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return "", errors.Wrap(err, "GetStorageLatest err")
	}
	if !ok {
		return "", errors.New("GetStorageLatest return value is empty")
	}

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(keyring, o)
	if err != nil {
		return "", errors.Wrap(err, "Sign err")
	}

	// Do the transfer and track the actual status
	sub, err := api.r.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return "", errors.Wrap(err, "SubmitAndWatchExtrinsic err")
	}
	defer sub.Unsubscribe()

	timeout := time.After(10 * time.Second)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := MyEventRecords{}
				h, err := api.r.RPC.State.GetStorageRaw(keye, status.AsInBlock)
				if err != nil {
					return "", err
				}
				err = types.EventRecordsRaw(*h).DecodeEventRecords(meta, &events)
				if err != nil {
					return "", err
				}
				if events.FileBank_FileUpload != nil {
					return "success!", nil
				} else {
					return "fail", errors.New("upload file fail")
				}
			}
		case <-timeout:
			return "", errors.New("upload file meta info timeout,please check your Internet!")
		default:
			time.Sleep(time.Second)
		}
	}
}

func (ci *CessInfo) DeleteFileOnChain(fileid string) error {
	var (
		err         error
		accountInfo types.AccountInfo
	)
	api.getSubstrateApiSafe()
	defer func() {
		api.releaseSubstrateApi()
		err := recover()
		if err != nil {
			fmt.Printf("[Error]Recover UserHoldSpaceDetails panic :%s\n", err)
		}
	}()
	keyring, err := signature.KeyringPairFromSecret(ci.IdentifyAccountPhrase, 0)
	if err != nil {
		return errors.Wrap(err, "KeyringPairFromSecret err")
	}

	meta, err := api.r.RPC.State.GetMetadataLatest()
	if err != nil {
		return errors.Wrap(err, "GetMetadataLatest err")
	}

	c, err := types.NewCall(meta, ci.TransactionName, types.NewBytes([]byte(fileid)))
	if err != nil {
		return errors.Wrap(err, "NewCall err")
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		return errors.Wrap(err, "NewExtrinsic err")
	}

	genesisHash, err := api.r.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return errors.Wrap(err, "GetBlockHash err")
	}

	rv, err := api.r.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return errors.Wrap(err, "GetRuntimeVersionLatest err")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", keyring.PublicKey)
	if err != nil {
		return errors.Wrap(err, "CreateStorageKey Account err")
	}

	keye, err := types.CreateStorageKey(meta, "System", "Events", nil)
	if err != nil {
		return errors.Wrap(err, "CreateStorageKey Events err")
	}

	ok, err := api.r.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return errors.Wrap(err, "GetStorageLatest err")
	}
	if !ok {
		return errors.New("GetStorageLatest return value is empty")
	}

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(keyring, o)
	if err != nil {
		return errors.Wrap(err, "Sign err")
	}

	// Do the transfer and track the actual status
	sub, err := api.r.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return errors.Wrap(err, "SubmitAndWatchExtrinsic err")
	}
	defer sub.Unsubscribe()

	timeout := time.After(10 * time.Second)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := MyEventRecords{}
				h, err := api.r.RPC.State.GetStorageRaw(keye, status.AsInBlock)
				if err != nil {
					return err
				}
				err = types.EventRecordsRaw(*h).DecodeEventRecords(meta, &events)
				if err != nil {
					errors.Wrap(err, "DecodeEvent err")
					return err
				}
				if events.FileBank_DeleteFile != nil {
					return nil
				} else {
					return errors.Errorf("Delete file info on chain fail!")
				}
			}
		case <-timeout:
			return errors.Errorf("[%v] tx timeout", ci.TransactionName)
		default:
			time.Sleep(time.Second)
		}
	}
}

func (ci *CessInfo) UploadDeclaration(filehash, filename string) (string, error) {
	api.getSubstrateApiSafe()
	defer func() {
		api.releaseSubstrateApi()
		if err := recover(); err != nil {
			fmt.Printf("[Error] Upload Declaration panic :%s\n", err)
		}
	}()

	var txhash string
	var accountInfo types.AccountInfo

	meta, err := api.r.RPC.State.GetMetadataLatest()
	if err != nil {
		return txhash, errors.Wrap(err, "GetMetadataLatest err")
	}

	c, err := types.NewCall(meta, ci.TransactionName, types.NewBytes([]byte(filehash)), types.NewBytes([]byte(filename)))
	if err != nil {
		return txhash, errors.Wrap(err, "[NewCall]")
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		return txhash, errors.Wrap(err, "[NewExtrinsic]")
	}

	genesisHash, err := api.r.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return txhash, errors.Wrap(err, "GetBlockHash")
	}

	rv, err := api.r.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return txhash, errors.Wrap(err, "GetRuntimeVersionLatest err")
	}

	keyring, err := signature.KeyringPairFromSecret(ci.IdentifyAccountPhrase, 0)
	if err != nil {
		return txhash, errors.Wrap(err, "KeyringPairFromSecret err")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", keyring.PublicKey, nil)
	if err != nil {
		return txhash, errors.Wrap(err, "[CreateStorageKey]")
	}

	ok, err := api.r.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return txhash, errors.Wrap(err, "[GetStorageLatest]")
	}

	if !ok {
		return txhash, errors.New("GetStorageLatest return value is empty")
	}

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(keyring, o)
	if err != nil {
		return txhash, errors.Wrap(err, "[Sign]")
	}

	// Do the transfer and track the actual status
	sub, err := api.r.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return txhash, errors.Wrap(err, "[SubmitAndWatchExtrinsic]")
	}
	defer sub.Unsubscribe()
	timeout := time.After(10 * time.Second)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := MyEventRecords{}
				txhash, _ = types.EncodeToHexString(status.AsInBlock)
				keye, err := types.CreateStorageKey(meta, "System", "Events", nil)
				if err != nil {
					return txhash, errors.Wrap(err, "CreateStorageKey Events err")
				}
				h, err := api.r.RPC.State.GetStorageRaw(keye, status.AsInBlock)
				if err != nil {
					return txhash, errors.Wrap(err, "GetStorageRaw")
				}

				err = types.EventRecordsRaw(*h).DecodeEventRecords(meta, &events)
				if err != nil {
					errors.Wrap(err, "DecodeEvent err")
					return txhash, err
				}

				if events.FileBank_UploadDeclaration != nil {
					return txhash, nil
				} else {
					return txhash, errors.Errorf("Delete file info on chain fail!")
				}
			}
		case err = <-sub.Err():
			return txhash, errors.Wrap(err, "<-sub")
		case <-timeout:
			return txhash, errors.Errorf("[%v] tx timeout", ci.TransactionName)
		}
	}
}
