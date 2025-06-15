package syncmap

type ISyncMap interface {
	Store(string, byte)
	Size() int
	Sum() byte
	Clear()
}
