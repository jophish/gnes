package gnes

type mapper_MMC1 struct {
}

func (*mapper_MMC1) write(val, addr uint8) error {
	return nil
}

func (*mapper_MMC1) read(addr uint8) (uint8, error) {
	return 0, nil
}

func newMapper_MMC1() mapper {
	return &mapper_MMC1{}
}
