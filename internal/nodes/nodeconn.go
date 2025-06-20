package nodes

var (
	_ INodeConn = &nodeConn{}
)

type nodeConn struct {
	Addr string `json:"addr"`
	Pasw string `json:"pasw"`

	authKey []byte
	encrKey []byte
}

func (p *nodeConn) GetAddress() string {
	return p.Addr
}

func (p *nodeConn) GetAuthKey() []byte {
	return p.authKey
}

func (p *nodeConn) GetEncrKey() []byte {
	return p.encrKey
}
