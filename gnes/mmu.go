package gnes

type mmu struct {
	mapper mapper
}

func newMmu(mapperNum uint32) (*mmu, error) {
	mmu := &mmu{}
	mapper, err := numberToMapper(mapperNum)
	if err != nil {
		return nil, err
	}
	mmu.mapper = mapper
	return mmu, nil
}

func (*mmu) read(addr uint8) (uint8, error) {
	return 0, nil
}

func (*mmu) write(val, addr uint8) error {
	return nil
}
