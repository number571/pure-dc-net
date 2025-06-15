package dc

type dcNet struct {
	iteration  uint64
	generators []IGenerator
}

func NewDCNet(i uint64, g ...IGenerator) IDCNet {
	return &dcNet{
		iteration:  i,
		generators: g,
	}
}

func (p *dcNet) Iteration() uint64 {
	return p.iteration
}

func (p *dcNet) Generate() byte {
	sum := byte(0)
	for _, g := range p.generators {
		sum ^= g.Generate(p.iteration)
	}
	p.iteration++
	return sum
}
