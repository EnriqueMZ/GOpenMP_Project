/*
 =========================================================================================================
 Name        : for_processor.go
 Author      : Enrique Madridejos Zamorano
 Version     :
 Copyright   : Apache Licence Version 2.0
 Description : Module that handles loops inside a pragma for.
               Loop valid structure:
               
               for  init-expr , var relop b , incr-expr
               
               Where,
               
               - “init-expr”: initialization of “var” variable (loop variable), by an integer expression.
               - “relop”: valid operators <, <=, >, >=.
               - “b”: integer expression.
               - “incr-expr”: Increase or decrease of “var”, in a integer number,
                  using a standard operator (++, --,  +=, -=), or by the form “var = var + incr”. 
 ==========================================================================================================
*/

package for_processor

import (
	"go/token"
	. "goprep"
	. "var_processor"
)

// Private token work functions.

// Funtion that let a token pass.
func passToken(tok Token, out chan string, sync chan interface{}) {
	out <- tok.Str
	sync <- nil
}

// Funtion that eliminate a token.
func eliminateToken(out chan string, sync chan interface{}) {
	out <- ""
	sync <- nil
}

// Function that determines valid logical operators.
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

// Function to get the type of a variable declared as reduction.
// Launch panic if variable not previously declared.
func search_typ(id string, varGlobalList []Variable, varLocalList []Variable) string {
	var typ string = "error"
	for i := range varGlobalList {
		if id == varGlobalList[i].Ident {
			typ = varGlobalList[i].Type
			break
		}
	}
	for i := range varLocalList {
		if id == varLocalList[i].Ident {
			typ = varLocalList[i].Type
			break
		}
	}
	if typ == "error" {
		panic("Variable \"" + id + "\" in clause not previously declared.")
	}
	return typ
}

// Function that removes a variable from a private varible list.
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

// Function that adds a variable to a private variable list.
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

// Function that process a loop declaration include in a pragma for.
func For_declare(tok Token, in chan Token, out chan string, sync chan interface{}, varGlobalList []Variable, varLocalList []Variable, routine_num string, for_threads string) Token {
	var num_iter, ini, fin, inc, steps, var_indice, aux, assign string
	var err bool
	if tok.Token != token.FOR {
		panic("Error: It must start with keyword \"for\".")
	}
	passToken(tok, out, sync)
	tok = <-in
	var_indice = tok.Str
	// Rewrite the loop
	eliminateToken(out, sync)
	tok = <-in
	if tok.Token != token.DEFINE && tok.Token != token.ASSIGN {
		panic("Error: Loop variable must be defined implicitly.")
	}else{
		assign = tok.Str
		}
	eliminateToken(out, sync)
	tok = <-in
	if tok.Token != token.INT {
		panic("Error: Variable \"" + tok.Str + "\" must be defined as an integer.")
	}
	ini = tok.Str
	// Rewrite the loop
	eliminateToken(out, sync)
	tok = <-in
	if tok.Token != token.SEMICOLON {
		panic("Error: Wait a semicolon.")
	}
	eliminateToken(out, sync)
	tok = <-in
	aux = tok.Str
	if aux != var_indice {
		panic("Error: It must use the same variable in the for declarion.")
	}
	// Rewrite the loop
	eliminateToken(out, sync)
	tok = <-in
	err, inc = logic_operator(tok)
	if err {
		panic("Invalid logical operator.")
	}
	// Rewrite the loop
	eliminateToken(out, sync)
	tok = <-in
	if tok.Token == token.INT {
		fin = tok.Str
	} else {
		if tok.Token == token.IDENT {
			typ := search_typ(tok.Str, varGlobalList, varLocalList)
			if typ != "int" {
				panic("Error: Variable \"" + tok.Str + "\" must be defined as an integer.")
			} else {
				fin = tok.Str
			}
		} else {
			panic("Error: Variable \"" + tok.Str + "\" must be defined as an integer.")
		}
	}
	// Rewrite the loop
	eliminateToken(out, sync)
	tok = <-in
	if tok.Token != token.SEMICOLON {
		panic("Error: Wait a semicolon.")
	}
	eliminateToken(out, sync)
	tok = <-in
	aux = tok.Str
	if aux != var_indice {
		panic("Error: It must use the same variable in the for declarion.")
	}
	eliminateToken(out, sync)
	tok = <-in
	switch tok.Token {
	case token.INC, token.DEC:
		steps = "1"
	case token.ADD_ASSIGN, token.SUB_ASSIGN:
		if tok.Token != token.INT {
			panic("Error: Variable \"" + tok.Str + "\" must be an integer.")
		}
		steps = tok.Str
	}
	num_iter = "(" + fin + " + " + inc + ") / " + steps // String "(fin + inc) / steps"
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
