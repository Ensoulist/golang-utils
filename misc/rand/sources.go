package rand

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

func init() {
	RegisterSource(ChaCha8, createChaCha8)
	RegisterSource(PCG, createPCG)
	RegisterSource(LCG, createLCG)
	RegisterSource(LCG2, createLCG2)
	RegisterSource(Zipf, createZipf)
	RegisterSource(Exponential, createExponential)
}

// /////// ChaCha8 ////////
type ChaCha8Param struct {
	Seed [32]byte
}

func createChaCha8(param any) (rand.Source, error) {
	var ptrSeed *[32]byte
	if param == nil {
		ptrSeed = &[32]byte{}
		s := (*ptrSeed)[0:32]
		_, err := crand.Read(s)
		if err != nil {
			return nil, err
		}
	} else {
		p, ok := param.(*ChaCha8Param)
		if !ok {
			return nil, fmt.Errorf("create ChaCha8, invalid param")
		}
		ptrSeed = &p.Seed
	}
	return rand.NewChaCha8(*ptrSeed), nil
}

// /////// ChaCha8 ////////
type PCGParam struct {
	Seed1, Seed2 uint64
}

func createPCG(param any) (rand.Source, error) {
	var seed1, seed2 uint64
	if param == nil {
		seed1 = uint64(time.Now().UnixNano())
		if err := binary.Read(crand.Reader, binary.LittleEndian, &seed2); err != nil {
			return nil, fmt.Errorf("failed to generate random seed: %v", err)
		}
	} else {
		p, ok := param.(*PCGParam)
		if !ok {
			return nil, fmt.Errorf("create PCG, invalid param")
		}
		seed1, seed2 = p.Seed1, p.Seed2
	}
	return rand.NewPCG(seed1, seed2), nil
}

// /////// LCG ////////
type LCGParam struct {
	Seed uint64
}

func (s *LCGParam) Uint64() uint64 {
	// Parameters for LCG (common values)
	const a = 6364136223846793005
	const c = 1
	s.Seed = (a*s.Seed + c)
	return uint64(s.Seed)
}

func createLCG(param any) (rand.Source, error) {
	if param == nil {
		return &LCGParam{Seed: uint64(time.Now().UnixNano())}, nil
	} else {
		p, ok := param.(*LCGParam)
		if !ok {
			return nil, fmt.Errorf("create LCG, invalid param")
		}
		return p, nil
	}
}

// /////// LCG2 ////////
type LCG2Param struct {
	Seed uint64
}

func (s *LCG2Param) Uint64() uint64 {
	// Parameters for LCG (common values)
	s.Seed = (s.Seed*9031 + 49297) % 233280
	return uint64((float64(s.Seed) / 233280) * math.MaxUint64)
}

func createLCG2(param any) (rand.Source, error) {
	if param == nil {
		return &LCG2Param{Seed: uint64(time.Now().UnixNano())}, nil
	} else {
		p, ok := param.(*LCG2Param)
		if !ok {
			return nil, fmt.Errorf("create LCG2, invalid param")
		}
		return p, nil
	}
}

// /////// Zipf ////////
type ZipfParam struct {
	Rand *rand.Rand
	S, V float64
	IMax uint64
}

func createZipf(param any) (rand.Source, error) {
	var p *ZipfParam
	if param != nil {
		var ok bool
		p, ok = param.(*ZipfParam)
		if !ok {
			p = nil
		}
	}

	var r *rand.Rand
	if p != nil && p.Rand != nil {
		r = p.Rand
	} else {
		r, _ = NewRand(Default, nil)
	}

	var s float64
	if p != nil && p.S > 1 {
		s = p.S
	} else {
		s = 1.5
	}

	var v float64
	if p != nil && p.V >= 1 {
		v = p.V
	} else {
		v = 1
	}

	var imax uint64
	if p != nil && p.IMax > 0 {
		imax = p.IMax
	} else {
		imax = 10000
	}
	return rand.NewZipf(r, s, v, imax), nil
}

// /////// Exponential ////////
type ExponentialParam struct {
	Rand   *rand.Rand
	Lambda float64
}

func (s *ExponentialParam) Uint64() uint64 {
	r := s.Rand.ExpFloat64() / s.Lambda
	return uint64(r)
}

func createExponential(param any) (rand.Source, error) {
	p, ok := param.(*ExponentialParam)
	if !ok || p == nil {
		p = &ExponentialParam{}
	}
	if p.Rand == nil {
		var err error
		p.Rand, err = NewRand(Default, nil)
		if err != nil {
			return nil, err
		}
		if p.Lambda <= 0 {
			p.Lambda = 1
		}
	}
	return p, nil
}
