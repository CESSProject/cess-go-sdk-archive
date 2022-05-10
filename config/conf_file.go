package config

type CessConf struct {
	ChainData ChainData `yaml:"chainData"`
}

type ChainData struct {
	CessRpcAddr           string `yaml:"cessRpcAddr"`
	FaucetAddress         string `yaml:"faucetAddress"`
	IdAccountPhraseOrSeed string `yaml:"idAccountPhraseOrSeed"`
	AccountPublicKey      string `yaml:"accountPublicKey"`
	WalletAddress         string `yaml:"walletAddress"`
}
