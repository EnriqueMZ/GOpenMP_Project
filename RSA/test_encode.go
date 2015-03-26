package main

import (
	"math/big"
	"crypto/rand"
	"crypto/rsa"
	//"io"
	"fmt"
	"os"
	"pollard_rho"
	"pollard_rho_brent"
	//"sieve_github"
)

func check(err error) {
	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}
}

func main() {
	//func GenerateKey(random io.Reader, bits int) (priv *PrivateKey, err error)
	privateKey, err := rsa.GenerateKey(rand.Reader, 64)
	check(err)

	D := privateKey.D //private exponent
	Primes := privateKey.Primes
	PCValues := privateKey.Precomputed

	fmt.Println("Private Key : ", privateKey)
	fmt.Println("Private Exponent : ", D.String())
	fmt.Printf("Primes : %s %s \n", Primes[0].String(), Primes[1].String())
	fmt.Printf("Precomputed Values : Dp[%s] Dq[%s]\n", PCValues.Dp.String(), PCValues.Dq.String())
	fmt.Printf("Precomputed Values : Qinv[%s]", PCValues.Qinv.String())
	fmt.Println()

	var publicKey *rsa.PublicKey
	publicKey = &privateKey.PublicKey
	N := publicKey.N // modulus
	E := publicKey.E // public exponent

	fmt.Println()
	fmt.Println("Public key ", publicKey)
	fmt.Println("Public Exponent : ", E)
	fmt.Println("Modulus : ", N.String())
	fmt.Println()
	
	var B = big.NewInt(1)
	var a = big.NewInt(1)
	var s = big.NewInt(1)
	
	res1, count1 := pollard_rho.Pollard_rho(N, B, a, s)
	fmt.Println("Factor:", res1)
	fmt.Println("Steps:", count1)
	
	res2, count2 := pollard_rho_brent.Pollard_rho_brent(N, B, a, s)
	fmt.Println("Factor:", res2)
	fmt.Println("Steps:", count2)
	
	//res3, res4:= sieve_github.Factorize(N, false)
	//fmt.Println("Factores:", res3, res4)
}
