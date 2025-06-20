package dc

type dcState struct {
	iteration  uint64
	generators []IGenerator
}

func NewDCState(i uint64, g ...IGenerator) IDCState {
	return &dcState{
		iteration:  i,
		generators: g,
	}
}

func (p *dcState) Iteration() uint64 {
	return p.iteration
}

func (p *dcState) Generate() byte {
	sum := byte(0)
	for _, g := range p.generators {
		sum ^= g.Generate(p.iteration)
	}
	p.iteration++
	return sum
}
