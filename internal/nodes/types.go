package nodes

type Nodes map[string]*NodeConn

type NodeConn struct {
	Addr string
	Key  []byte
}
