package Config

type ServerConfig struct {
	RuleSystem struct {
		Enabled      bool
		ModulePath   string
		Rulefile     string
		Worker       int
		TLSEnable    bool
		TLSCert      string
		TLSKey       string
		TLSCaCert    string
		RpcInterface string
	}
}

type ClientConfig struct {
	Client struct {
		RpcInterface string
		TLSEnable    bool
		TLSCert      string
		TLSKey       string
		TLSCaCert    string
	}
}
