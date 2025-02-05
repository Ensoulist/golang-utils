package rand

import (
	"fmt"
	"math/rand/v2"
)

func NewRand(typ RandType, param any) (*rand.Rand, error) {
	facotry, ok := registerdSources[typ]
	if !ok {
		return nil, fmt.Errorf("no rand type: %v", typ)
	}
	source, err := facotry(param)
	if err != nil {
		return nil, err
	}
	return rand.New(source), nil
}

type FactoryFunc func(param any) (rand.Source, error)

var registerdSources = map[RandType]FactoryFunc{}

func RegisterSource(typ RandType, factoryFunc FactoryFunc) {
	registerdSources[typ] = factoryFunc
}
