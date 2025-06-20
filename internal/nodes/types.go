package nodes

type Nodes map[string]*NodeConn

type NodeConn struct {
	Addr string `json:"addr"`
	Pasw string `json:"pasw"`
}
