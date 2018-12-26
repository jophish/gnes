package gnes

var mapperMap = map[uint32]func(*cartInfo) (mapper, error){
	1: newMapper_MMC1,
}

func numberToMapper(mapper uint32, info *cartInfo) (mapper, error) {
	if mapFunc, ok := mapperMap[mapper]; ok {
		newMapper, err := mapFunc(info)
		if err != nil {
			return nil, err
		}
		return newMapper, nil
	} else {
		return nil, &gError1{err_MAPPER_UNSUPPORTED, uint64(mapper)}
	}
}

type mapper interface {
	write(val uint8, addr uint16) error
	read(addr uint16) (uint8, error)
}
