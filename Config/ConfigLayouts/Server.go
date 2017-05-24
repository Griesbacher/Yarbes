package ConfigLayouts

//Server file layout
type Server struct {
	TLS struct {
		Enable    bool
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
		Enabled        bool
		RPCInterface   string
		HTTPInterface  string
		HtpasswdPath   string
		InfluxAddress  string
		InfluxDatabase string
		InfluxUsername string
		InfluxPassword string
		RuleSystem     string
	}
	Proxy struct {
		Enabled      bool
		RPCInterface string
		ModuleFolder string
	}
}
