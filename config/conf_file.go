package config

type CessConf struct {
	ChainData ChainData `yaml:"chainData"`
}

type ChainData struct {
	CessRpcAddr           string `yaml:"cessRpcAddr"`
	FaucetAddress         string `yaml:"faucetAddress"`
	IdAccountPhraseOrSeed string `yaml:"idAccountPhraseOrSeed"`
	WalletAddress         string `yaml:"walletAddress"`
}
