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
		return nil, err_MAPPER_UNSUPPORTED
	}
}

type mapper interface {
	write(val, addr uint8) error
	read(addr uint8) (uint8, error)
}
