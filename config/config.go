package config

// Cfg holds the configuration settings for the application.
var Cfg Config

// Config holds the configuration settings for the application.
type Config struct {
	Mode        string
	Listen      string
	DatabaseURL string
	SecretKey   string

	KubesphereAPI KubesphereAPIConfig
	ChainAPI      ChainAPIConfig
}

// KubesphereAPIConfig holds the configuration for the KubeSphere API.
type KubesphereAPIConfig struct {
	URL      string
	UserName string
	Password string
	Cluster  string
}

// ChainAPIConfig holds the configuration for the chain API.
type ChainAPIConfig struct {
	AddressPrefix        string
	RPC                  string
	TokenContractAddress string
	ServiceName          string
	KeyringDir           string
	FaucetGas            string
	FaucetToken          string
	OrderContractAddress string
}
