package gnes

var mapperMap = map[uint32]func() mapper{
	1: newMapper_MMC1,
}

func numberToMapper(mapper uint32) (mapper, error) {
	if mapFunc, ok := mapperMap[mapper]; ok {
		return mapFunc(), nil
	} else {
		return nil, err_MAPPER_UNSUPPORTED
	}
}

type mapper interface {
	write(val, addr uint8) error
	read(addr uint8) (uint8, error)
}
