/*
 ==========================================================================
 Name        : for_processor.go
 Author      : Enrique Madridejos Zamorano
 Version     :
 Copyright   : Apache Licence Version 2.0
 Description : Módulo para tratamiento de bucles dentro de un pragma for
 ==========================================================================
 */

package for_processor

import (
	"go/token"
	. "goprep"
	. "var_processor"
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
func logic_operator(tok Token) (bool, string) {
	var err bool
	var inc string
	switch tok.Token {
	case token.LSS, token.GTR:
		err = false
		inc = "0"
	case token.LEQ, token.GEQ:
		err = false
		inc = "1"
	default:
		err = true
	}
	return err, inc
}

// Función para obtener el tipo de una variable marcada como reduction. Error si no se ha inicializado previamente.
func search_typ(id string, varList []Variable) string {
	var typ string = "error"
	for i := range varList {
		if id == varList[i].Ident {
			typ = varList[i].Type
		}
	}
	if typ == "error" {
		panic("Variable " + id + " en clausula reduction no declarada previamente")
	}
	return typ
}

// Función que elimina una variable de una lista de variables.
func eliminate_element(id string, privateList []string) []string {
	for i := range privateList {
		if id == privateList[i] {
			privateList[i] = privateList[len(privateList)-1]
			privateList = privateList[:len(privateList)-1]
			break
		}
	}
	return privateList
}

// Función que añade una variable a una lista de variables.
func add_element(id string, privateList []string) []string {
	for i := range privateList {
		if id == privateList[i] {
			break
		}else{
			element := id + " int"
			privateList = append(privateList, element)
			}
	}
	return privateList
}

// Funcion que trata la declaracion de un bucle for paralelizado.
func For_declare(tok Token, in chan Token, out chan string, sync chan interface{}, varList []Variable, routine_num string, for_threads string) Token {
	var num_iter, ini, fin, inc, steps, var_indice, aux, assign string
	var err bool
	if tok.Token != token.FOR {
		panic("Error: Debe comenzar con un for")
	}
	passToken(tok, out, sync)
	tok = <-in
	var_indice = tok.Str
	// Reescribe el bucle
	eliminateToken(out, sync)
	tok = <-in
	if tok.Token != token.DEFINE && tok.Token != token.ASSIGN {
		panic("Error: La variable indice debe definirse implicitamente")
	}else{
		assign = tok.Str
		}
	eliminateToken(out, sync)
	tok = <-in
	if tok.Token != token.INT {
		panic("Error: la variable " + tok.Str + " debe definirse como un entero")
	}
	ini = tok.Str
	// Reescribe el bucle
	eliminateToken(out, sync)
	tok = <-in
	if tok.Token != token.SEMICOLON {
		panic("Error: Espera un semicolon")
	}
	eliminateToken(out, sync)
	tok = <-in
	aux = tok.Str
	if aux != var_indice {
		panic("Error: Debe emplear la misma variable en la declaracion del for")
	}
	// Reescribe el bucle
	eliminateToken(out, sync)
	tok = <-in
	err, inc = logic_operator(tok)
	if err {
		panic("Operador lógico no válido")
	}
	// Reescribe el bucle
	eliminateToken(out, sync)
	tok = <-in
	if tok.Token == token.INT {
		fin = tok.Str
	} else {
		if tok.Token == token.IDENT {
			typ := search_typ(tok.Str, varList)
			if typ != "int" {
				panic("Error: la variable " + tok.Str + " debe definirse como un entero")
			} else {
				fin = tok.Str
			}
		} else {
			panic("Error: la variable " + tok.Str + " debe definirse como un entero")
		}
	}
	// Reescribe el bucle
	eliminateToken(out, sync)
	tok = <-in
	if tok.Token != token.SEMICOLON {
		panic("Error: Espera un semicolon")
	}
	eliminateToken(out, sync)
	tok = <-in
	aux = tok.Str
	if aux != var_indice {
		panic("Error: Debe emplear la misma variable en la declaracion del for")
	}
	eliminateToken(out, sync)
	tok = <-in
	switch tok.Token {
	case token.INC, token.DEC:
		steps = "1"
	case token.ADD_ASSIGN, token.SUB_ASSIGN:
		if tok.Token != token.INT {
			panic("Error: Debe definirse como un entero")
		}
		steps = tok.Str
	}
	num_iter = "(" + fin + " + " + inc + ") / " + steps // Cadena "(fin + inc) / steps"
	if for_threads == "1" {
		out <- var_indice + assign + " " + routine_num + "; "+ var_indice + " < " + num_iter + "; " + var_indice + "++"
		sync <- nil
		tok = <-in
	} else {
		out <- var_indice + assign + " " + routine_num + " + " + ini + "; " + var_indice + " < " + num_iter + "; " + var_indice + " += " + for_threads // _i := _routine_num + 0; _i < (n+0)/1; _i += _numCPUs
		sync <- nil
		tok = <-in
	}
	return tok
}
