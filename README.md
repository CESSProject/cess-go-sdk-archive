# CESS-GO-SDK

## Overview

CESS-GO-SDK (hereinafter referred to as sdk) is suitable for operating systems of go version go1.17.2 and above. It is an API abstracted from the CESS client project cess-portal for external calls. Developers can use sdk to quickly access to the CESS storage system, the functions it provides include: file upload, file download, file deletion, file decryption, file encryption, search for specified file information, search for user file lists, search for real-time storage prices in the CESS system, and search for user-purchased space , CESS system storage space to purchase and obtain CESS transaction coins. The implementation of the sdk does not require any token. Anyone accessing the sdk only needs to pass in the parameters required by the API to obtain the response of the CESS system. Please keep the network unobstructed when using the sdk.



## Project structure

```shell
cess-go-sdk
├─config
├─internal
│  ├─chain    ##Code files for interacting with the CESS chain
│  └─rpc	  ##Code files that interact with scheduler
├─module
│  └─result  ##sdk return data structure
├─sdk		##sdk code file
├─test		##sdk code test file
└─tools		##tool
```



## Function introduction

Method receiver:

| method receiver |                        Include method                        |                             Link                             |
| :-------------: | :----------------------------------------------------------: | :----------------------------------------------------------: |
|     FileSDK     | FileUpload、FileDownload、FileDelete、FileDncrypt、FileEncrypt | https://github.com/CESSProject/cess-go-sdk/blob/5f5b41204890cbb0da2134c421ace8d546d535c1/sdk/file_sdk.go#L23 |
|    QuerySDK     |  QueryPurchasedSpace、QueryPrice、QueryFile、QueryFileList   | https://github.com/CESSProject/cess-go-sdk/blob/5f5b41204890cbb0da2134c421ace8d546d535c1/sdk/query_sdk.go#L13 |
|   PurchaseSDK   |                 ObtainFromFaucet、Expansion                  | https://github.com/CESSProject/cess-go-sdk/blob/5f5b41204890cbb0da2134c421ace8d546d535c1/sdk/purchase_sdk.go#L12 |



### File upload

* Example link:https://github.com/CESSProject/cess-go-sdk/blob/d5d1c99f54f09166a1c5bf6c6ebab1724747556c/test/file_upload_test.go#L8

* Method explanation: Uploading files to the CESS system needs to consume the purchased storage space

* Parameter explanation:

  Value method receiver:

  |    parameter name     | Parameter explanation |  type  |
  | :-------------------: | :-------------------: | :----: |
  |      CessRpcAddr      |  CESS chain address   | string |
  | IdAccountPhraseOrSeed |    wallet mnemonic    | string |
  |     WalletAddress     |    wallet address     | string |

  Method parameters:

  | parameter name |                    Parameter explanation                     |     type      |
  | :------------: | :----------------------------------------------------------: | :-----------: |
  |   blocksize    |                      CESS chain address                      | sdk.BlockSize |
  |      path      |                       wallet mnemonic                        |    string     |
  |    backups     |                        wallet address                        |    string     |
  |   privatekey   | File encryption password, if the length is zero, it is a public file |    string     |

* Return parameter:

  | parameter name | Parameter explanation |  type  |
  | :------------: | :-------------------: | :----: |
  |     fileid     |    file unique id     | string |
  |      err       |     return error      | error  |

### File download

* Example link:https://github.com/CESSProject/cess-go-sdk/blob/d5d1c99f54f09166a1c5bf6c6ebab1724747556c/test/file_download_test.go#L8

* Method Explanation:Restore and download uploaded files from CESS system to local

* Parameter explanation:

  Value method receiver:

  | parameter name | Parameter explanation |  type  |
  | :------------: | :-------------------: | :----: |
  |  CessRpcAddr   |  CESS chain address   | string |
  | WalletAddress  |    wallet address     | string |

  Method parameters:

  | parameter name | Parameter explanation  |  type  |
  | :------------: | :--------------------: | :----: |
  |     fileid     |        file id         | string |
  |  installpath   | Download to local path | string |
  
* Return parameter:

  | parameter name | Parameter explanation | type  |
  | :------------: | :-------------------: | :---: |
  |      err       |     return error      | error |

### FileDelete

* Example link:https://github.com/CESSProject/cess-go-sdk/blob/d5d1c99f54f09166a1c5bf6c6ebab1724747556c/test/file_delete_test.go#L8

* Method Explanation:Delete the uploaded files on the CESS system

* Parameter explanation:

  Value method receiver:

  |    parameter name     | Parameter explanation |  type  |
  | :-------------------: | :-------------------: | :----: |
  |      CessRpcAddr      |  CESS chain address   | string |
  | IdAccountPhraseOrSeed |    wallet mnemonic    | string |

  Method parameters:

  | parameter name | Parameter explanation |  type  |
  | :------------: | :-------------------: | :----: |
  |     fileid     |        file id        | string |
  
* Return parameter:

  | parameter name | Parameter explanation | type  |
  | :------------: | :-------------------: | :---: |
  |      err       |     return error      | error |

### File decrypt

* Example link:https://github.com/CESSProject/cess-go-sdk/blob/d5d1c99f54f09166a1c5bf6c6ebab1724747556c/test/file_decrypt_test.go#L8

* Method Explanation: Decrypt the file, the decryption method is AES

* Parameter explanation:

  Method parameters:

  | parameter name |       Parameter explanation       |  type  |
  | :------------: | :-------------------------------: | :----: |
  |  decryptpath   |       file path to decrypt        | string |
  |    savepath    | The path to save after decryption | string |
  |    password    | File password, length: 16, 24, 32 | string |
  
* Return parameter:

  | parameter name | Parameter explanation | type  |
  | :------------: | :-------------------: | :---: |
  |      err       |     return error      | error |

### File encrypt 

* Example link:https://github.com/CESSProject/cess-go-sdk/blob/d5d1c99f54f09166a1c5bf6c6ebab1724747556c/test/file_encrypt_test.go#L8

* Method Explanation:Encrypt the file, the encryption method is AES

* Parameter explanation:

  Method parameters:

  | parameter name |       Parameter explanation       |  type  |
  | :------------: | :-------------------------------: | :----: |
  |  encryptpath   |        Encrypted file path        | string |
  |    savepath    | The path to save after decryption | string |
  |    password    | File password, length: 16, 24, 32 | string |
  
* Return parameter:

  | parameter name | Parameter explanation | type  |
  | :------------: | :-------------------: | :---: |
  |      err       |     return error      | error |

### Query file

* Example link:https://github.com/CESSProject/cess-go-sdk/blob/d5d1c99f54f09166a1c5bf6c6ebab1724747556c/test/query_file_test.go#L8

* Method Explanation:Query the information of the uploaded file based on the unique id of the file

* Parameter explanation:

  Value method receiver:

  | parameter name | Parameter explanation |  type  |
  | :------------: | :-------------------: | :----: |
  |  CessRpcAddr   |  CESS chain address   | string |

  Method parameters:

  | parameter name | Parameter explanation |  type  |
  | :------------: | :-------------------: | :----: |
  |     fileid     |        file id        | string |
  
* return value:

  | parameter name | Parameter explanation  |  type  |
  | :------------: | :--------------------: | :----: |
  |    FileName    |       file name        | string |
  |    FileSize    |   file size, unit b    | int64  |
  |    FileHash    |       file hash        | string |
  |     Public     |  if a public document  |  bool  |
  |    Backups     | number of file backups |  int8  |
  |  Downloadfee   |     download costs     | int64  |
  |      err       |      return error      | error  |

### Query file list

* Example link:https://github.com/CESSProject/cess-go-sdk/blob/d5d1c99f54f09166a1c5bf6c6ebab1724747556c/test/query_file_test.go#L20

* Method Explanation:Query all files uploaded by the user based on user information

* Parameter explanation:

  Value method receiver:

  | parameter name | Parameter explanation |  type  |
  | :------------: | :-------------------: | :----: |
  |  CessRpcAddr   |  CESS chain address   | string |
  
* Return parameter:

  | parameter name | Parameter explanation |  type  |
  | :------------: | :-------------------: | :----: |
  |     FileId     |        file id        | string |
  |      err       |     return error      | error  |

### Query price

* Example link:https://github.com/CESSProject/cess-go-sdk/blob/d5d1c99f54f09166a1c5bf6c6ebab1724747556c/test/query_price_test.go#L8

* Method Explanation:Query the real-time price of the current space of the CESS system, in CESS/G

* Parameter explanation:

  Value method receiver:

  | parameter name | Parameter explanation |  type  |
  | :------------: | :-------------------: | :----: |
  |  CessRpcAddr   |  CESS chain address   | string |
  
* Return parameter:

  | parameter name |    Parameter explanation     |  type   |
  | :------------: | :--------------------------: | :-----: |
  |   spaceprice   | Return price, unit (GB/CESS) | float64 |
  |      err       |         return error         |  error  |

### Query purchased space

* Example link:https://github.com/CESSProject/cess-go-sdk/blob/d5d1c99f54f09166a1c5bf6c6ebab1724747556c/test/query_purchased_space_test.go#L8

* Method Explanation:Query the CESS storage space purchased by the user through user information

* Parameter explanation:

  Value method receiver:

  | parameter name | Parameter explanation |  type  |
  | :------------: | :-------------------: | :----: |
  |  CessRpcAddr   |  CESS chain address   | string |
  | WalletAddress  |    wallet address     | string |
  
* Return parameter:

  | parameter name |     Parameter explanation     |  type  |
  | :------------: | :---------------------------: | :----: |
  | PurchasedSpace |    Purchased space, in MB     | string |
  |   UsedSpace    |       Used space, in MB       | string |
  | RemainingSpace | Remaining unused space, in MB | string |
  |      err       |         return error          | error  |

### Expansion

* Example link:https://github.com/CESSProject/cess-go-sdk/blob/d5d1c99f54f09166a1c5bf6c6ebab1724747556c/test/purchase_expansion_test.go#L8

  * Method Explanation:Purchase storage space for a CESS system

* Parameter explanation:

  Value method receiver:

  |    parameter name     | Parameter explanation |  type  |
  | :-------------------: | :-------------------: | :----: |
  |      CessRpcAddr      |  CESS chain address   | string |
  | IdAccountPhraseOrSeed |    wallet mnemonic    | string |

  Method parameters:

  |       parameter name        |                    Parameter explanation                     | type |
  | :-------------------------: | :----------------------------------------------------------: | :--: |
  | QuantityOfSpaceYouWantToBuy |        The amount of space that needs to be purchased        | int  |
  |       MonthsWantToBuy       |              Need to buy space for a few months              | int  |
  |        ExpectedPrice        | The highest price you can receive when buying space, 0 means accept all prices | int  |

* Return parameter:

  | parameter name | Parameter explanation | type  |
  | :------------: | :-------------------: | :---: |
  |      err       |     return error      | error |

### ObtainFromFaucet

* Example link:https://github.com/CESSProject/cess-go-sdk/blob/d5d1c99f54f09166a1c5bf6c6ebab1724747556c/test/purchase_obtain_from_faucet_test.go#L8

* Method Explanation:Get CESS transaction coins from the faucet

* Parameter explanation:

  Value method receiver:

  | parameter name | Parameter explanation |  type  |
  | :------------: | :-------------------: | :----: |
  | FaucetAddress  |    Faucet address     | string |

  Method parameters:

  | parameter name | Parameter explanation |  type  |
  | :------------: | :-------------------: | :----: |
  | WalletAddress  |    wallet address     | string |
  
* Return parameter:

  | parameter name | Parameter explanation | type  |
  | :------------: | :-------------------: | :---: |
  |      err       |     return error      | error |

## Contribute code

1. Fork
2. Create your feature branch git checkout -b my-new-feature
3. Commit your changes git commit -am 'Added some feature'
4. Commit your changes to the remote git repository git push origin my-new-feature
5. Then go to the my-new-feature branch of the git remote repository on the github website to initiate a Pull Request

If you have any questions, please leave a message in the issue section, thank you.