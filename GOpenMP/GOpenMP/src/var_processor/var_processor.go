/*
 =================================================================================
 Name        : for_parallel_processor.go
 Author      : Enrique Madridejos Zamorano
 Version     :
 Copyright   : Apache Licence Version 2.0
 Description : Módulo para tratamiento de declaraciones de variables.
 			   NOTA IMPORTANTE: Sólo admite variables declaradas IMPLÍCITAMENTE. 
 =================================================================================
 */

package var_processor

import (
	//"fmt"
	"go/token"
	. "goprep"

)

// Stack of bools, model with a slice.
type braceStack []bool

func (stack braceStack) Empty() bool {
	return len(stack) == 0
}
func (stack braceStack) Peek() bool {
	return stack[len(stack)-1]
}
func (stack *braceStack) Push(i bool) {
	(*stack) = append((*stack), i)
}
func (stack *braceStack) Pop() bool {
	i := (*stack)[len(*stack)-1]
	(*stack) = (*stack)[:len(*stack)-1]
	return i
}

type Variable struct { // Estructura para variables inicializadas
	Ident string // Identificador
	Type  string // Tipo
	Ini   bool   // ¿Está inicializada?
}

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

// Funcion que trata la declaracion de un interface.
func interface_declare(tok Token, in chan Token, out chan string, sync chan interface{}) string {
	var str string = "interface"
	passToken(tok, out, sync)
	tok = <-in
	str = str + tok.Str
	passToken(tok, out, sync)
	fin := false
	for !fin {
		tok = <-in
		if tok.Token == token.RBRACE {
			str = str + tok.Str
			passToken(tok, out, sync)
			fin = true
		} else {
			str = str + tok.Str
			passToken(tok, out, sync)
		}
	}
	return str
}

// Funcion que trata la declaracion de un slice. 
func slice_declare(tok Token, in chan Token, out chan string, sync chan interface{}) string {
	var str string = "["
	passToken(tok, out, sync)
	fin := false
	for !fin {
		tok = <-in
		if tok.Token == token.IDENT {
			str = str + tok.Str
			passToken(tok, out, sync)
			fin = true
		} else {
			str = str + tok.Str
			passToken(tok, out, sync)
		}
	}
	return str
}

// Funcion que trata la declaracion de un puntero.
func pointer_declare(tok Token, in chan Token, out chan string, sync chan interface{}) string {
	var str string = "*"
	passToken(tok, out, sync)
	fin := false
	for !fin {
		tok = <-in
		if tok.Token == token.IDENT {
			str = str + tok.Str
			passToken(tok, out, sync)
			fin = true
		} else {
			str = str + tok.Str
			passToken(tok, out, sync)
		}
	}
	return str
}

// Funcion que trata la declaracion de un map.
func map_declare(tok Token, in chan Token, out chan string, sync chan interface{}) string {
	var str string = "map"
	passToken(tok, out, sync)
	fin_brack := false
	for !fin_brack {
		tok = <-in
		if tok.Token == token.RBRACK {
			str = str + tok.Str
			passToken(tok, out, sync)
			tok = <-in
			str = str + tok.Str
			passToken(tok, out, sync)
			fin_brack = true
		} else {
			str = str + tok.Str
			passToken(tok, out, sync)
		}
	}
	return str
}

// Funcion que trata la declaracion de un struct.
func struct_declare(tok Token, in chan Token, out chan string, sync chan interface{}) string {
	var str string = "struct"
	passToken(tok, out, sync)
	fin_brace := false
	for !fin_brace {
		tok = <-in
		if tok.Token == token.RBRACE {
			str = str + tok.Str
			passToken(tok, out, sync)
			fin_brace = true
		} else {
			str = str + tok.Str + " "
			passToken(tok, out, sync)
		}
	}
	return str
}

// Funcion que trata la declaracion de una funcion como variable.
func func_declare(tok Token, in chan Token, out chan string, sync chan interface{}) string {
	var str string = "func"
	passToken(tok, out, sync)
	fin_paren := false
	for !fin_paren {
		tok = <-in
		if tok.Token == token.RPAREN {
			str = str + tok.Str
			passToken(tok, out, sync)
			tok = <-in
			str = str + tok.Str
			passToken(tok, out, sync)
			fin_paren = true
		} else {
			str = str + tok.Str + " "
			passToken(tok, out, sync)
		}
	}
	return str
}

// Funcion que trata la declaracion de un canal.
func channel_declare(tok Token, in chan Token, out chan string, sync chan interface{}) string {
	var str string = "chan"
	passToken(tok, out, sync)

	return str
}

// Funcion para asignar los campos de las variables de la lista de variables.
func assign(identList []string, tipe string, ini bool) []Variable {
	var varAux Variable
	var res []Variable
	for i := range identList {
		varAux.Ident = identList[i]
		varAux.Type = tipe
		varAux.Ini = ini
		res = append(res, varAux)
	}
	return res
}

// Funcion publica que concatena dos slices de Variable.
func Var_concat(a, b []Variable) []Variable {
	for i := range b {
		a = append(a, b[i])
	}
	return a
}

// Funcion publica que almacena las variables de declaraciones simples.
func Var_simple_processor(tok Token, in chan Token, out chan string, sync chan interface{}) []Variable {
	var listAux []Variable
	var tipe string = ""
	var ini bool = false
	var identList []string
	var fin_paren bool = false

	passToken(tok, out, sync)
	for !fin_paren {
		tok = <-in
		if tok.Token == token.RPAREN {
			fin_paren = true
		} else {
			identList = append(identList, tok.Str)
			passToken(tok, out, sync)
			finIn := false
			for !finIn {
				tok = <-in
				switch tok.Token {
				case token.COMMA: // SEGUIR A PARTIR DE AQUI
					passToken(tok, out, sync)
					tok = <-in
					identList = append(identList, tok.Str)
					passToken(tok, out, sync)
				case token.IDENT:
					tipe = tipe + tok.Str
					passToken(tok, out, sync)
				case token.MUL:
					// Tratamiento de punteros
					tipe = pointer_declare(tok, in, out, sync)
				case token.INTERFACE:
					// Tratamiento declaracion interface
					tipe = interface_declare(tok, in, out, sync)
				case token.LBRACK:
					// Tratamiento declaracion slice
					tipe = slice_declare(tok, in, out, sync)
				case token.MAP:
					// Tratamiento declaracion map
					tipe = map_declare(tok, in, out, sync)
				case token.STRUCT:
					//tratamiento declaracion struct
					tipe = struct_declare(tok, in, out, sync)
				case token.FUNC:
					// Tratamiento declaracion funcion
					tipe = func_declare(tok, in, out, sync)
				case token.PERIOD:
					// Tipos compuestos
					tipe = tipe + tok.Str
					passToken(tok, out, sync)
				case token.CHAN:
					// Tratamiento declaracion chan
					tipe = tipe + tok.Str + " "
					passToken(tok, out, sync)
				case token.ASSIGN:
					ini = true
					passToken(tok, out, sync)
					for {
						tok = <-in
						if tok.Token == token.SEMICOLON { // REVISAR AQUI
							passToken(tok, out, sync)
							break
						} else {
							passToken(tok, out, sync)
						}
					}
					finIn = true
				case token.SEMICOLON:
					passToken(tok, out, sync)
					finIn = true
				}
			}
		}
		if fin_paren {
			passToken(tok, out, sync)
			break
		} else {
			listAux = assign(identList, tipe, ini)
			tipe = ""
			ini = false
			identList = nil
		}
	}
	return listAux
}

// Funcion publica que almacena las variables de declaraciones multiples.
func Var_multi_processor(tok Token, in chan Token, out chan string, sync chan interface{}) []Variable {
	var listAux []Variable
	var tipe string = ""
	var ini bool = false
	var identList []string
	identList = append(identList, tok.Str)
	passToken(tok, out, sync)
	fin := false
	for !fin {
		tok = <-in
		switch tok.Token {
		case token.COMMA:
			passToken(tok, out, sync)
			tok = <-in
			identList = append(identList, tok.Str)
			passToken(tok, out, sync)
		case token.IDENT:
			tipe = tipe + tok.Str
			passToken(tok, out, sync)
		case token.MUL:
			// Tratamiento de punteros
			tipe = pointer_declare(tok, in, out, sync)
		case token.INTERFACE:
			// Tratamiento declaracion interface
			tipe = interface_declare(tok, in, out, sync)
		case token.LBRACK:
			// Tratamiento declaracion slice
			tipe = slice_declare(tok, in, out, sync)
		case token.MAP:
			// Tratamiento declaracion map
			tipe = map_declare(tok, in, out, sync)
		case token.STRUCT:
			// Tratamiento declaracion struct
			tipe = struct_declare(tok, in, out, sync)
		case token.FUNC:
			// Tratamiento declaracion funcion
			tipe = func_declare(tok, in, out, sync)
		case token.PERIOD:
			// Tipos compuestos
			tipe = tipe + tok.Str
			passToken(tok, out, sync)
		case token.CHAN:
			// Tratamiento declaracion chan
			tipe = tipe + tok.Str + " "
			passToken(tok, out, sync)
		case token.ASSIGN:
			ini = true
			passToken(tok, out, sync)
			fin = true
		case token.SEMICOLON:
			passToken(tok, out, sync)
			fin = true
		}
	}
	listAux = assign(identList, tipe, ini)
	tipe = ""
	ini = false
	identList = nil
	return listAux
}
