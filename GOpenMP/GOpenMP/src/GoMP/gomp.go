package main

import (
	"fmt"
	"go/token"
	. "goprep"
	"os"
	. "pragma_processor"
	"strings"
	. "var_processor"
	//. "gomp_lib"
	"runtime"
	//"strconv"
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

// Eliminate black spaces in a given string.
func noSpaces(str string) string {
	return strings.Replace(str, " ", "", -1)
}

// Check is a token is a "pragma gomp"
func isPragma(token Token) bool {
	res := strings.HasPrefix(noSpaces(token.Str), "//pragmagomp")
	return res
}

// Numero de declaraciones de variables. Solo para testeo.
var num_dec int = 0

// Lista de variables declaradas.
var varList []Variable

// Numero de pragmas. Solo para testeo.
var num_prag int = 0

// Información del pragma.
var pragma Pragma

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

// Funcion barrier. Añade una barrera cuando sea necesario.
func barrier(numParallel int) (string, int) {
	if numParallel == 0 {
		numParallel++
		return "var _barrier = make(chan bool)\n", numParallel
	} else {
		return "", numParallel
	}
}

// Funcion routineNum. Trata el token Gomp_get_routine_num()
func subs_Gomp_get_routine_num(in chan Token, out chan string, sync chan interface{}) {
	out <- "_routine_num"
	sync <- nil
	tok := <-in
	// Parentesis. Esto quizá pueda eliminarse y dejar que el compiladro se encargue de indicar el error.
	if tok.Token != token.LPAREN {
		panic("Error Parentesis Izquierdo")
	}
	eliminateToken(out, sync)
	tok = <-in
	if tok.Token != token.RPAREN {
		panic("Error Parentesis Derecho")
	}
	eliminateToken(out, sync)
}

// Funcion que ignora Gomp_set_num_routine() donde corresponda.
func ign_Gomp_set_num_routine(in chan Token, out chan string, sync chan interface{}) {
	eliminateToken(out, sync)
	fin := false
	for !fin {
		tok := <-in
		if tok.Token == token.RPAREN {
			eliminateToken(out, sync)
			fin = true
		} else {
			eliminateToken(out, sync)
		}
	}
}

// Funcion para tratar la clausula "num_threads".
func set_num_threads(pragma Pragma) string {
	if pragma.Num_threads == "" {
		return "Gomp_get_num_routines()"
	} else {
		return pragma.Num_threads
	}
}

func notFound(a string, varList []Variable) (bool, string) {
	var res bool = true
	var v string = ""
	if len(varList) == 0 {
		res = false
		v = a
	} else {
		for i := range varList {
			if a == varList[i].Ident {
				res = false
				break
			}
			if i == len(varList)-1 {
				v = a
			}
		}
	}
	return res, v
}

// Funcion que comprueba las variables de una clausula default(none) con las variables declaradas.
func var_not_prev_declare(pragma Pragma, varList []Variable) (bool, string) {
	var res bool
	var v string
	if len(pragma.Variable_List) == 0 {
		res = true
		v = ""
	} else {
		for i := range pragma.Variable_List {
			res, v = notFound(pragma.Variable_List[i], varList)
			if res {
				break
			}
		}
	}
	return res, v
}

func declare(ident string, varList []Variable) string { // CORREGIR!!!!
	var res string
	for i := range varList {
		if ident == varList[i].Ident {
			if varList[i].Type == "no_type" {
				res = varList[i].Ident + "= reflect.Zero(reflect.TypeOf(" + varList[i].Ident + ")).Interface()"
			} else {
				res = varList[i].Ident + " " + varList[i].Type
			}
			break
		}
	}
	return res
}

// Funcion que crea el string para re-inicializacion de variables private. WARNING: Variables declaradas de forma implicita!!!
func declareList(pragma Pragma, varList []Variable) string {
	var res string
	if len(pragma.Private_List) == 0 {
		res = ""
	} else {
		res = declare(pragma.Private_List[0], varList)
		for i := 1; i < len(pragma.Private_List); i++ {
			res = res + ";\n" + declare(pragma.Private_List[i], varList)
		}
	}
	return res
}

// Programa principal.
func main() {
	// Establecemos GOMAXPROCS
	_numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(_numCPUs)
	p := PipeInit(os.Stdin)
	//Lines(p) // Muestra las lineas
	Link(func(in chan Token,
		tOut chan Token,
		out chan string,
		sync chan interface{}) {
		var numParallel int = 0 // Inicializa el numero de regiones paralelas
		//var listAux []Variable
		//var tipe, bar string = "", ""
		//var ini bool = false
		var bar string = ""
		for tok := range in {

			switch { // Tratamiento de Tokens

			case tok.Str == "var": // Tratamiento para declaración de variables.
				num_dec++ // Numero de declaraciones de variables (para testeo).
				passToken(tok, out, sync)
				tok = <-in
				if tok.Token == token.LPAREN {
					// Declaracion simple
					varList = Var_concat(varList, Var_simple_processor(tok, in, out, sync))
					continue
				} else {
					// Declaracion multiple
					varList = Var_concat(varList, Var_multi_processor(tok, in, out, sync))
					continue
				}

			case isPragma(tok): // Reconocedor de "pragma gomp"

				num_prag++
				fmt.Println("Numero de pragmas actual: ", num_prag, "\n")
				fmt.Println("Pragma: ", tok.Str, "\n") // Recordar retirar los fmt.PrintLn
				pragma = ProcessPragma(tok.Str)
				fmt.Println("Información del pragma: ", pragma, "\n")

				switch pragma.Type { // Tratamiento de pragmas por tipo
				case 0: // PRAGMA PARALLEL
					//s, b := (*BoolStack)(nil), false
					var b bool
					var s braceStack
					endParallel := false

					// Comprobar clausula default
					if pragma.Default == NONE {
						def_cond, def_var := var_not_prev_declare(pragma, varList)
						if def_cond {
							panic("Error: variable " + def_var + " no declarada previamente")
						}
					}

					bar, numParallel = barrier(numParallel)
					out <- bar + "for i := 0; i < " + set_num_threads(pragma) + "; i++{\n" + "go func(_routine_num int)"
					sync <- nil

					tok = <-in
					// init LBRACE
					if tok.Token != token.LBRACE {
						panic("Error: Falta la llave de inicio del pragma")
					}
					//s = Push(s, true)
					s.Push(true) // Llave de apertura de bloque Parallel

					//VARIABLES PRIVATE
					privateList := declareList(pragma, varList)
					fmt.Println("Variables privadas:\n", privateList, "\n")
					out <- " {" + "var (" + privateList + ") \n"
					sync <- nil

					// Tratamiento del contenido del Parallel
					for !endParallel {
						tok = <-in
						switch {
						case tok.Token == token.LBRACE:
							// An lbrace not associated with parallel
							//s = Push(s, false)
							s.Push(false)
							passToken(tok, out, sync)
						case tok.Token == token.RBRACE:
							//s, b = Pop(s)
							b = s.Pop()
							if b {
								// End the parallel
								out <- " _barrier <- true\n" + "}(i)\n" + "}\n" + "for i := 0; i < " + set_num_threads(pragma) + "; i++{\n" + "<-_barrier\n" + "}\n"
								sync <- nil
								endParallel = true
							} else {
								// An rbrace not associated with parallel
								passToken(tok, out, sync)
							}
						case tok.Str == "Gomp_get_routine_num":
							subs_Gomp_get_routine_num(in, out, sync)
						case tok.Str == "Gomp_set_num_routines":
							ign_Gomp_set_num_routine(in, out, sync)
						default:
							// Ignore
							passToken(tok, out, sync)
						}
					}
					continue
				case 1: // PRAGMA PRIVATE_FOR
					panic("Error: Pragma Parallel_For en proceso...")
				case 3: // PRAGMA THREADPRIVATE
					passToken(tok, out, sync)
					continue
					// TO DO: Resto de tratamiento de pragmas
				}
			default:
				// Ignore
				passToken(tok, out, sync)
				continue
			}
		}
		close(tOut)
	})(p)

	PipeEnd(p, os.Stdout)
	fmt.Println("Numero de declaraciones en el código: ", num_dec, "\n")
	fmt.Println("Lista de variables declaradas en el código: ", varList, "\n")
}
