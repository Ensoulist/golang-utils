package rand

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func testRand(all, max int, r *rand.Rand) {
	m := make(map[int64]int)
	for i := 0; i < all; i++ {
		n := r.Uint64N(uint64(max))
		m[int64(n)]++
	}
	fmt.Println(m)
}

func testDistribution(all int, r *rand.Rand) {
	m := make(map[int64]int)
	for i := 0; i < all; i++ {
		n := r.Uint64()
		m[int64(n)]++
	}
	fmt.Println(m)
}

func TestChaCha8(t *testing.T) {
	r, err := NewRand(ChaCha8, nil)
	if err != nil {
		t.Fatal(err)
	}
	testRand(100, 10, r)
	testRand(10000, 10, r)
	testRand(100000, 10, r)
}

func TestPCG(t *testing.T) {
	r, err := NewRand(PCG, nil)
	if err != nil {
		t.Fatal(err)
	}
	testRand(100, 10, r)
	testRand(10000, 10, r)
	testRand(100000, 10, r)
}

func TestLCG(t *testing.T) {
	r, err := NewRand(LCG, nil)
	if err != nil {
		t.Fatal(err)
	}
	testRand(100, 10, r)
	testRand(10000, 10, r)
	testRand(100000, 10, r)
}

func TestLCG2(t *testing.T) {
	r, err := NewRand(LCG2, nil)
	if err != nil {
		t.Fatal(err)
	}
	testRand(100, 10, r)
	testRand(10000, 10, r)
	testRand(100000, 10, r)
}

func TestZipf(t *testing.T) {
	r, err := NewRand(Zipf, &ZipfParam{
		S:    1.5,
		V:    1.0,
		IMax: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	testDistribution(100, r)
	testDistribution(10000, r)
	testDistribution(100000, r)
}

func TestExponential(t *testing.T) {
	r, err := NewRand(Exponential, &ExponentialParam{
		Lambda: 1.0,
	})
	if err != nil {
		t.Fatal(err)
	}
	testDistribution(100, r)
	testDistribution(10000, r)
	testDistribution(100000, r)
}

func TestCotoutine(t *testing.T) {
	coFunc := func(i int, yield func(int) int) int {
		fmt.Println(i)
		j := yield(2)
		fmt.Println(j)
		return 100
	}
	co := New(coFunc)
	fmt.Println(co(1))
	fmt.Println(co(3))
}

func New[In, Out any](f func(in In, yield func(Out) In) Out) (resume func(In) Out) {
	cin := make(chan In)
	cout := make(chan Out)
	resume = func(in In) Out {
		cin <- in
		return <-cout
	}
	yield := func(out Out) In {
		cout <- out
		return <-cin
	}
	go func() { cout <- f(<-cin, yield) }()
	return resume
}
