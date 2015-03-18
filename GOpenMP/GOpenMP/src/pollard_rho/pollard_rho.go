package pollard_rho 

import (
	"math/big"
)

func factorial(n *big.Int) (result *big.Int) {
	result = new(big.Int)

	switch n.Cmp(&big.Int{}) {
	case -1, 0:
		result.SetInt64(1)
	default:
		result.Set(n)
		var one big.Int
		one.SetInt64(1)
		result.Mul(result, factorial(n.Sub(n, &one)))
	}
	return
}

func f (x *big.Int, k *big.Int, a *big.Int, n *big.Int) *big.Int {
	// F(x) = x^2k + a mod n
	var res big.Int
	var pow big.Int
	var pow_x big.Int
	var add big.Int
	
	b := big.NewInt(0)
	c := big.NewInt(2)
	
	pow.Mul(c, k)
	//fmt.Println("2*k=", pow)
	pow_x.Exp(x, &pow, b)
	//fmt.Println("x^2k=", pow_x)
	add.Add(&pow_x, a)
	//fmt.Println("x^2k + a=", add)
	res.Mod(&add, n)
	//fmt.Println("x^2k + a mod=", res)
	
	return &res
	}

func Pollard_rho(n *big.Int, B *big.Int, a *big.Int, s *big.Int) (*big.Int, int64) {
	var count int64 = 0
	var U = s
	var V = s
	var g = big.NewInt(1)
	var one = big.NewInt(1)
	
	k := factorial(B)
	
	for g.Cmp(one) == 0 {
		
		U = f(U, k, a, n)
		V = f(V, k, a, n)
		V = f(V, k, a, n)
		
		g.Abs(g.Sub(U, V))
		g.GCD(nil, nil, g, n)
		
		if g.Cmp(n) == 0 {
			U.Add(U, one)
			V.Add(V, one)
			g = big.NewInt(1)
			}
		count++
		}
	return g, count
	}