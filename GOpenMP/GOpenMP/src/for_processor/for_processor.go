package for_processor

import (
	"go/token"
	. "goprep"
	"strconv"
)

// Funciones para trabajo con tokens.

// Funcion que deja pasar un token.
func passToken(tok Token, out chan string, sync chan interface{}) {
	out <- tok.Str
	sync <- nil
}

// Funcion que elimina un token.
func eliminateToken(out chan string, sync chan interface{}) {
	out <- ""
	sync <- nil
}

// Funcion que determina operadores lógico válidos
func logic_operator(tok Token) (bool, int) {
	var err bool
	var inc int
	switch tok.Token {
	case token.LSS, token.GTR:
		err = false
		inc = 0
	case token.LEQ, token.GEQ:
		err = false
		inc = 1
	default:
		err = true
	}
	return err, inc
}

// Funcion que trata la declaracion de un interface.
func For_declare(tok Token, in chan Token, out chan string, sync chan interface{}) (string, Token) {
	var num_iter string
	var ini, fin, inc, steps int
	var variable, aux string
	var err bool
	if tok.Token != token.FOR {
		panic("Error: Debe comenzar con un for")
	}
	passToken(tok, out, sync)
	tok = <-in
	variable = tok.Str
	// Reescribe el bucle
	out <- "_i"
	sync <- nil
	tok = <-in
	if tok.Token != token.DEFINE && tok.Token != token.ASSIGN {
		panic("Error: La variable indice debe definirse implicitamente")
	}
	passToken(tok, out, sync)
	tok = <-in
	if tok.Token != token.INT {
		panic("Error: Debe definirse como un entero")
	}
	ini, _ = strconv.Atoi(tok.Str)
	// Reescribe el bucle
	out <- "0"
	sync <- nil
	tok = <-in
	if tok.Token != token.SEMICOLON {
		panic("Error: Espera un semicolon")
	}
	passToken(tok, out, sync)
	tok = <-in
	aux = tok.Str
	if aux != variable {
		panic("Error: Debe emplear la misma variable en la declaracion del for")
	}
	// Reescribe el bucle
	out <- "_i"
	sync <- nil
	tok = <-in
	err, inc = logic_operator(tok)
	if err {
		panic("Operador lógico no válido")
	}
	// Reescribe el bucle
	out <- "<"
	sync <- nil
	tok = <-in
	if tok.Token != token.INT {
		panic("Error: Debe definirse como un entero")
	}
	fin, _ = strconv.Atoi(tok.Str)
	// Reescribe el bucle
	out <- "_numCPUs"
	sync <- nil
	tok = <-in
	if tok.Token != token.SEMICOLON {
		panic("Error: Espera un semicolon")
	}
	passToken(tok, out, sync)
	tok = <-in
	aux = tok.Str
	if aux != variable {
		panic("Error: Debe emplear la misma variable en la declaracion del for")
	}
	out <- "_i"
	sync <- nil
	tok = <-in
	switch tok.Token {
	case token.INC, token.DEC:
		steps = 1
		// Reescribe el bucle
		out <- "++"
		sync <- nil
		tok = <-in
	case token.ADD_ASSIGN, token.SUB_ASSIGN:
		// Reescribe el bucle	
		out <- "++"
		sync <- nil
		tok = <-in
		if tok.Token != token.INT {
			panic("Error: Debe definirse como un entero")
		}
		steps, _ = strconv.Atoi(tok.Str)
		// Reescribe el bucle
		eliminateToken(out, sync)
		tok = <-in
	}
	num_iter = strconv.Itoa(((fin + inc) - ini) / steps)
	return num_iter, tok
	// 	WARNING!!! PARADO HASTA TERMINAR EL PROCESADOR DE IMPORTS
}
