package ConfigLayouts

//Server file layout
type Server struct {
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