package nodes

type NodesMap map[string]INodeConn

type INodeConn interface {
	GetAddress() string
	GetAuthKey() []byte
	GetEncrKey() []byte
}
