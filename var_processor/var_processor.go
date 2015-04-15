/*
 =======================================================================================================
 Name        : for_parallel_processor.go
 Author      : Enrique Madridejos Zamorano
 Version     :
 Copyright   : Licensed under the Apache License, Version 2.0 (the "License");
   			   you may not use this file except in compliance with the License.
               You may obtain a copy of the License at

               http://www.apache.org/licenses/LICENSE-2.0

               Unless required by applicable law or agreed to in writing, software
               distributed under the License is distributed on an "AS IS" BASIS,
               WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
               See the License for the specific language governing permissions and
               limitations under the License.
               
 Description : Module that handles variable declarations and function arguments from the original code.
 			   IMPORTANT NOTE: Supports only EXPLICITLY declared variables.
 =======================================================================================================
*/

package var_processor

import (
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

type Variable struct { // Initialized variables structure
	Ident string // Identifier
	Type  string // Type
	Ini   bool   // ¿Is the varible inicialized with a value?
}

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

// Function that process an interface declaration.
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

// Function that process a slice declaration.
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

// Function that process a pointer declaration.
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

// Function that process a map declaration.
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

// Function that process a struct declaration.
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

// Function that process a function declaration as a variable.
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

// Function that process a channel declaration.
func channel_declare(tok Token, in chan Token, out chan string, sync chan interface{}) string {
	var str string = "chan"
	passToken(tok, out, sync)
	return str
}

// Function to map the fields of the variables in the variable list.
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

// Function that concatenates two slices of Variable type (public).
func Var_concat(a, b []Variable) []Variable {
	for i := range b {
		a = append(a, b[i])
	}
	return a
}

// Function that stores the variables from simple declarations (public). 
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
				case token.COMMA:
					passToken(tok, out, sync)
					tok = <-in
					identList = append(identList, tok.Str)
					passToken(tok, out, sync)
				case token.IDENT:
					tipe = tipe + tok.Str
					passToken(tok, out, sync)
				case token.MUL:
					// Pointer declaration processing
					tipe = pointer_declare(tok, in, out, sync)
				case token.INTERFACE:
					// Interface declaration processing
					tipe = interface_declare(tok, in, out, sync)
				case token.LBRACK:
					// Slice declaration processing
					tipe = slice_declare(tok, in, out, sync)
				case token.MAP:
					// Map declaration processing
					tipe = map_declare(tok, in, out, sync)
				case token.STRUCT:
					// Struct declaration processing
					tipe = struct_declare(tok, in, out, sync)
				case token.FUNC:
					// Function declaration processing
					tipe = func_declare(tok, in, out, sync)
				case token.PERIOD:
					// Composite types declaration processing
					tipe = tipe + tok.Str
					passToken(tok, out, sync)
				case token.CHAN:
					// Channel declaration processing
					tipe = tipe + tok.Str + " "
					passToken(tok, out, sync)
				case token.ASSIGN:
					ini = true
					passToken(tok, out, sync)
					for {
						tok = <-in
						if tok.Token == token.SEMICOLON { // CHECK THIS POINT
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
// Function that stores the variables from multiple declarations (public). 
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
			// Pointer declaration processing
			tipe = pointer_declare(tok, in, out, sync)
		case token.INTERFACE:
			// Interface declaration processing
			tipe = interface_declare(tok, in, out, sync)
		case token.LBRACK:
			// Slice declaration processing
			tipe = slice_declare(tok, in, out, sync)
		case token.MAP:
			// Map declaration processing
			tipe = map_declare(tok, in, out, sync)
		case token.STRUCT:
			// Struct declaration processing
			tipe = struct_declare(tok, in, out, sync)
		case token.FUNC:
			// Function declaration processing
			tipe = func_declare(tok, in, out, sync)
		case token.PERIOD:
			// Composite types declaration processing
			tipe = tipe + tok.Str
			passToken(tok, out, sync)
		case token.CHAN:
			// Channel declaration processing
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

// Function that stores the arguments from functions declarations (public). 
func Var_argument_processor(tok Token, in chan Token, out chan string, sync chan interface{}, varLocalList []Variable) []Variable {
	var tipe string = ""
	var variable Variable
	var fin_paren bool = false
	passToken(tok, out, sync)
	for !fin_paren {
		tok = <-in
		if tok.Token == token.RPAREN {
			passToken(tok, out, sync)
			break
		} else {
			variable.Ident = tok.Str
			passToken(tok, out, sync)
			finIn := false
			for !finIn {
				tok = <-in
				switch tok.Token {
				case token.IDENT:
					tipe = tipe + tok.Str
					passToken(tok, out, sync)
				case token.MUL:
					// Pointer declaration processing
					tipe = pointer_declare(tok, in, out, sync)
				case token.INTERFACE:
					// Interface declaration processing
					tipe = interface_declare(tok, in, out, sync)
				case token.LBRACK:
					// Slice declaration processing
					tipe = slice_declare(tok, in, out, sync)
				case token.MAP:
					// Map declaration processing
					tipe = map_declare(tok, in, out, sync)
				case token.STRUCT:
					// Struct declaration processing
					tipe = struct_declare(tok, in, out, sync)
				case token.FUNC:
					// Function declaration processing
					tipe = func_declare(tok, in, out, sync)
				case token.PERIOD:
					// Composite types declaration processing
					tipe = tipe + tok.Str
					passToken(tok, out, sync)
				case token.CHAN:
					// Channel declaration processing
					tipe = tipe + tok.Str + " "
					passToken(tok, out, sync)
				case token.COMMA:
					passToken(tok, out, sync)
					finIn = true
				case token.RPAREN:
					finIn = true
					fin_paren = true
				}
			}
		}
		if fin_paren {
			variable.Type = tipe
			variable.Ini = true
			tipe = ""
			varLocalList = append(varLocalList, variable)
			passToken(tok, out, sync)
			break
		} else {
			variable.Type = tipe
			variable.Ini = true
			tipe = ""
			varLocalList = append(varLocalList, variable)
		}
	}
	return varLocalList
}
