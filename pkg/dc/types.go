package dc

type IDCNet interface {
	Iteration() uint64
	Generate() byte
}

type IGenerator interface {
	Generate(uint64) byte
}

type ITotalizer interface {
	Store(...byte)
	Size() int
	Sum() byte
}
