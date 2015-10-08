package Config

type ServerConfig struct {
	TLS        struct {
				   Cert string
				   Key string
				   CaCert string
			   }
	RuleSystem struct {
				   Enabled      bool
				   ModulePath   string
				   Rulefile     string
				   Worker       int
				   RpcInterface string
			   }
	LogServer  struct {
				   Enabled      bool
				   RpcInterface string
			   }
}

type ClientConfig struct {
	TLS       struct {
				  Cert string
				  Key string
				  CaCert string
			  }
	Backend   struct {
				  RpcInterface string
			  }
	LogServer struct {
				  RpcInterface string
			  }
}
