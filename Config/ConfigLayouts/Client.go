package ConfigLayouts

//Client file layout
type Client struct {
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
	Livestatus struct {
		Enable  bool
		Type    string
		Address string
	}
}
