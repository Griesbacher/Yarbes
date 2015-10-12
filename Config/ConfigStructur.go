package Config

//Server config file layout
type ServerConfig struct {
	TLS struct {
		Cert      string
		Key       string
		CaCert    string
		BlackList []string
	}
	RuleSystem struct {
		Enabled      bool
		ModulePath   string
		Rulefile     string
		Worker       int
		RPCInterface string
	}
	LogServer struct {
		Enabled      bool
		RPCInterface string
	}
}

//Client config file layout
type ClientConfig struct {
	TLS struct {
		Cert   string
		Key    string
		CaCert string
	}
	Backend struct {
		RPCInterface string
	}
	LogServer struct {
		RPCInterface string
	}
}
