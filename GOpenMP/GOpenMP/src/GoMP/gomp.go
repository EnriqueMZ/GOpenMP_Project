/*
 ============================================================================
 Name        : gomp.go
 Author      : Enrique Madridejos Zamorano
 Version     :
 Copyright   : Apache Licence Version 2.0
 Description : Módulo principal del preprocesador de texto GOpenMP
 ============================================================================
*/

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
	. "for_parallel_processor"
	. "for_processor"
	. "import_processor"
	// "runtime"
	"strconv"
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
// var pragma Pragma

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

func stringOperator(op Red_Operator) string {
	var str string
	switch op {
	case 0:
		str = "+"
	case 1:
		str = "*"
	case 2:
		str = "-"
	case 3:
		str = "&"
	case 4:
		str = "|"
	case 5:
		str = "^"
	case 6:
		str = "&&"
	case 7:
		str = "||"
	default:
		panic("Error: operador no valido en clausula reduction")
	}
	return str
}

// Funcion que construye una barrera para una variable reduction.
func barrier_variable(numBarrier int, variable string, typ string, opr string) (string, string, string, string) {
	var num string = strconv.Itoa(numBarrier)
	var name, dcl, send, var_dcl, rcv string
	if variable == "nil" {
		name = "_barrier_" + num + "_bool"
		dcl = "var " + name + " = make(chan bool)\n"
		send = name + " <- true\n"
		var_dcl = ""
		rcv = "<- " + name + "\n"
	} else {
		name = "_barrier_" + num + "_" + typ
		dcl = "var " + name + " = make(chan " + typ + ")\n"
		send = name + " <- " + variable + "\n"
		var_dcl = "var " + variable + " " + typ + "\n"
		rcv = variable + " " + opr + "= <- " + name + "\n"
	}
	return dcl, send, var_dcl, rcv

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

// Función que construye barreras para todas la variables de una clausula reduction.
func barrier_single_reduction(numBarrier int, clause Reduction_Type, varList []Variable) (string, string, string, string, int) {
	var dcls, sends, var_dcls, rcvs string
	var numB int = numBarrier
	opr := stringOperator(clause.Operator)
	for i := range clause.Variables {
		typ := search_typ(clause.Variables[i], varList)
		dcl, send, var_dcl, rcv := barrier_variable(numB, clause.Variables[i], typ, opr)
		dcls = dcls + dcl
		sends = sends + send
		var_dcls = var_dcls + var_dcl
		rcvs = rcvs + rcv
		numB++
	}
	return dcls, sends, var_dcls, rcvs, numB
}

// Función que construye barreras para todas las variables de una lista de cluasulas reduction.
func barrier_list_reduction(numBarrier int, reductionList []Reduction_Type, varList []Variable) (string, string, string, string, int) {
	var dcls, sends, var_dcls, rcvs string
	if len(reductionList) == 0 {
		dcls_aux, sends_aux, var_dcls_aux, rcvs_aux := barrier_variable(numBarrier, "nil", "nil", "nil")
		dcls = dcls_aux
		sends = sends_aux
		var_dcls = var_dcls_aux
		rcvs = rcvs_aux
		numBarrier++
	} else {
		for i := range reductionList {
			dcls_aux, sends_aux, var_dcls_aux, rcvs_aux, numB_aux := barrier_single_reduction(numBarrier, reductionList[i], varList)
			dcls = dcls + dcls_aux
			sends = sends + sends_aux
			var_dcls = var_dcls + var_dcls_aux
			rcvs = rcvs + rcvs_aux
			numBarrier = numB_aux
		}
	}
	return dcls, sends, var_dcls, rcvs, numBarrier
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

func declare(ident string, varList []Variable) string { // WARNING: Puede necesitar correción.
	var res string
	var enc bool = false
	for i := range varList {
		if ident == varList[i].Ident {
			enc = true
			if varList[i].Type == "no_type" {
				res = varList[i].Ident + "= reflect.Zero(reflect.TypeOf(" + varList[i].Ident + ")).Interface()"
			} else {
				res = varList[i].Ident + " " + varList[i].Type
			}
			break
		}
	}
	if !enc {
		panic("Variable " + ident + " en clausula private no declarada previamente")
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

// Funcion que reescribe el codigo si se trata de un pragma
func pragma_rewrite(tok Token, in chan Token, out chan string, sync chan interface{}, num_prag int, in_parallel bool, routine_num string, for_threads string, numBarriers int) (int, int) {
	num_prag++
	fmt.Println("Numero de pragmas actual: ", num_prag, "\n")
	fmt.Println("Pragma: ", tok.Str, "\n") // Recordar retirar los fmt.PrintLn
	pragma := ProcessPragma(tok.Str)
	fmt.Println("Información del pragma: ", pragma, "\n")

	switch pragma.Type { // Tratamiento de pragmas por tipo
	case 0: // PRAGMA PARALLEL
		var b bool
		var s braceStack
		in_parallel = true
		endParallel := false
		routine_num = "_routine_num"     // String con el identificador de rutina.
		for_threads = pragma.Num_threads // String con el numero de hilos del Parallel.
		// Comprobar clausula default
		if pragma.Default == NONE {
			def_cond, def_var := var_not_prev_declare(pragma, varList)
			if def_cond {
				panic("Error: variable " + def_var + " no declarada previamente")
			}
		}

		// VARIABLES REDUCTION
		dcls, sends, var_dcls, rcvs, numB := barrier_list_reduction(numBarriers, pragma.Reduction_List, varList)
		numBarriers = numB
		out <- dcls + "for _i := 0; _i < " + pragma.Num_threads + "; _i++{\n" + "go func(_routine_num int)"
		sync <- nil

		tok = <-in
		// init LBRACE
		if tok.Token != token.LBRACE {
			panic("Error: Falta la llave de inicio del pragma")
		}
		s.Push(true) // Llave de apertura de bloque Parallel

		//VARIABLES PRIVATE
		privateList := declareList(pragma, varList)
		fmt.Println("Variables privadas en pragma parallel:", privateList)

		// Redeclaracion de variables private y reduction.
		out <- " {" + "var (" + privateList + ") \n" + var_dcls
		sync <- nil

		// Tratamiento del contenido del Parallel
		for !endParallel {
			tok = <-in
			switch {
			case tok.Token == token.LBRACE:
				// An lbrace not associated with parallel
				s.Push(false)
				passToken(tok, out, sync)
			case tok.Token == token.RBRACE:
				b = s.Pop()
				if b {
					// End the parallel
					out <- sends + "}(_i)\n" + "}\n" + "for _i := 0; _i < " + pragma.Num_threads + "; _i++{\n" + rcvs + "}\n"
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
			case isPragma(tok):
				num_prag, numBarriers = pragma_rewrite(tok, in, out, sync, num_prag, in_parallel, routine_num, for_threads, numBarriers)
			default:
				// Ignore
				passToken(tok, out, sync)
			}
		}
		in_parallel = false
	case 1: // PRAGMA PARALLEL_FOR
		var b bool
		var s braceStack
		var iteraciones string = "0" // Iteraciones del bucle paralelizado. Sólo para testeo.
		var ini, var_indice, assign string

		// Comprobar clausula default
		if pragma.Default == NONE {
			def_cond, def_var := var_not_prev_declare(pragma, varList)
			if def_cond {
				panic("Error: variable " + def_var + " no declarada previamente")
			}
		}
		// VARIABLES REDUCTION
		dcls, sends, var_dcls, rcvs, numB := barrier_list_reduction(numBarriers, pragma.Reduction_List, varList)
		numBarriers = numB
		out <- dcls // Cambia el pragma por la declaracion de canales
		sync <- nil

		tok = <-in // Token "for"
		fmt.Println("Variables declaradas antes del parallel for:", varList)
		iteraciones, ini, var_indice, assign, tok = For_parallel_declare(tok, in, out, sync, varList)
		fmt.Println("Iteraciones del bucle paralelo:", iteraciones)

		// VARIABLES PRIVATE
		privateList := declareList(pragma, varList)
		fmt.Println("Variables privadas en parallel for:", privateList)

		// Lanzamiento de goroutines. Redeclaracion de variables
		out <- tok.Str + "\n" + "go func(_routine_num int) {\n" + "var (" + privateList + ") \n" + var_dcls + "for " + var_indice + " " + assign + " _routine_num + " + ini + "; " + var_indice + " <" + iteraciones + "; "+ var_indice +" += _numCPUs {\n"
		sync <- nil

		// init LBRACE
		s.Push(true) // Llave de apertura de bloque Parallel For
		endParallelFor := false

		for !endParallelFor {
			tok = <-in
			switch {
			case tok.Token == token.LBRACE:
				// An lbrace not associated with parallel
				s.Push(false)
				passToken(tok, out, sync)
			case tok.Token == token.RBRACE:
				b = s.Pop()
				if b {
					// End the parallel for
					out <- tok.Str + "\n" + sends + "}(_i)\n" + "}\n" + "for _i := 0; _i < _numCPUs; _i++{\n" + rcvs + "}\n"
					sync <- nil
					endParallelFor = true
				} else {
					// An rbrace not associated with parallel
					passToken(tok, out, sync)
				}
			//case tok.Str == var_indice: // Variable indice del bucle
			//out <- "_i"
			//sync <- nil
			case tok.Str == "Gomp_get_routine_num":
				subs_Gomp_get_routine_num(in, out, sync)
			case tok.Str == "Gomp_set_num_routines":
				ign_Gomp_set_num_routine(in, out, sync)
			default:
				// Ignore
				passToken(tok, out, sync)
			}
		}
	case 2: // PRAGMA FOR
		var b bool
		var s braceStack
		var iteraciones string = "0" // Iteraciones del bucle paralelizado. Sólo para testeo.

		// Comprobar clausula default
		if pragma.Default == NONE {
			def_cond, def_var := var_not_prev_declare(pragma, varList)
			if def_cond {
				panic("Error: variable " + def_var + " no declarada previamente")
			}
		}

		/*
			// VARIABLES REDUCTION
			dcls, sends, var_dcls, rcvs, numB := barrier_list_reduction(numBarriers, pragma.Reduction_List, varList)
			numBarriers = numB

			out <- dcls // Cambia el pragma por la declaracion de canales
			sync <- nil
		*/

		eliminateToken(out, sync) // Eliminamos el pragma

		tok = <-in // Token "for"
		fmt.Println("Variables declaradas antes del parallel for:", varList)
		tok = For_declare(tok, in, out, sync, varList, routine_num, for_threads)
		fmt.Println("Iteraciones del bucle paralelo:", iteraciones)

		// VARIABLES PRIVATE
		privateList := declareList(pragma, varList)
		fmt.Println("Variables privadas en pragma for:", privateList)

		// Lanzamiento de goroutines. Redeclaracion de variables
		//out <- tok.Str + "\n" + "var (" + privateList + ") \n" + var_dcls + "for _i := _routine_num; _i <" + iteraciones + "; _i += _numCPUs {\n"
		out <- tok.Str + "\n" + "var (" + privateList + ") \n"
		sync <- nil

		// init LBRACE
		s.Push(true) // Llave de apertura de bloque For
		endFor := false

		for !endFor {
			tok = <-in
			switch {
			case tok.Token == token.LBRACE:
				// An lbrace not associated with parallel
				s.Push(false)
				passToken(tok, out, sync)
			case tok.Token == token.RBRACE:
				b = s.Pop()
				if b {
					// End the parallel for
					out <- tok.Str
					sync <- nil
					endFor = true
				} else {
					// An rbrace not associated with parallel
					passToken(tok, out, sync)
				}
			//case tok.Str == var_indice: // Variable indice del bucle
			//out <- "_i"
			//sync <- nil
			case tok.Str == "Gomp_get_routine_num":
				subs_Gomp_get_routine_num(in, out, sync)
			case tok.Str == "Gomp_set_num_routines":
				ign_Gomp_set_num_routine(in, out, sync)
			default:
				// Ignore
				passToken(tok, out, sync)
			}
		}
	case 3: // PRAGMA THREADPRIVATE
		eliminateToken(out, sync)
		// TO DO: Resto de tratamiento de pragmas
	}
	return num_prag, numBarriers
}

// Programa principal.
func main() {
	// Establecemos GOMAXPROCS
	// _numCPUs := runtime.NumCPU()
	// runtime.GOMAXPROCS(_numCPUs)
	p := PipeInit(os.Stdin)
	//Lines(p) // Muestra las lineas
	Link(func(in chan Token,
		tOut chan Token,
		out chan string,
		sync chan interface{}) {
		var numFunc int = 0          // Inicializa el numero de funciones del código original.
		var numBarriers int = 0      // Inicializa el número de barreras
		var in_parallel bool = false // Dentro de una region variable
		var routine_num string = "0" // String con el identificador de rutina
		//var default_threads string = "_numCPUs" // String con el numero de hilos por defecto.
		var for_threads = "1" // String con los hilos de un bucle for
		//var numParallel int = 0 				// Inicializa el numero de regiones paralelas
		//var tipe, bar string = "", ""
		//var ini bool = false
		for tok := range in {

			switch { // Tratamiento de Tokens

			case tok.Token == token.IMPORT: // Tratamiento de import.
				Imports_declare(tok, in, out, sync)
				continue

			case tok.Token == token.FUNC: // Tratamiento de _numCPUs
				var b bool
				var s braceStack
				if numFunc == 0 { // Es la primera funcion del código.
					numFunc++
					out <- "var _numCPUs = runtime.NumCPU()\n" + "func _init_numCPUs(){\n" + "runtime.GOMAXPROCS(_numCPUs)\n" + "}\n" + tok.Str
					sync <- nil
					tok = <-in
					fmt.Println("Entrando en la primera funcion:", tok.Str)
					if tok.Str == "main" { // La primera funcion es "main".
						for tok.Token != token.LBRACE {
							passToken(tok, out, sync)
							tok = <-in
						}
						// Inicializa el numero de CPUs
						out <- tok.Str + "\n" + "_init_numCPUs()\n"
						sync <- nil
						// init LBRACE
						s.Push(true) // Llave de apertura de la funcion
						endFunc := false
						for !endFunc {
							tok = <-in
							switch {
							case isPragma(tok): // Reconocedor de "pragma gomp"
								num_prag, numBarriers = pragma_rewrite(tok, in, out, sync, num_prag, in_parallel, routine_num, for_threads, numBarriers)
							case tok.Str == "var": // Tratamiento para declaración de variables.
								num_dec++ // Numero de declaraciones de variables (para testeo).
								passToken(tok, out, sync)
								tok = <-in
								fmt.Println("Variable local:", tok.Str)
								if tok.Token == token.LPAREN {
									// Declaracion simple
									varList = Var_concat(varList, Var_simple_processor(tok, in, out, sync))
								} else {
									// Declaracion multiple
									varList = Var_concat(varList, Var_multi_processor(tok, in, out, sync))
								}
							case tok.Token == token.LBRACE:
								// An lbrace not associated with parallel
								s.Push(false)
								passToken(tok, out, sync)
							case tok.Token == token.RBRACE:
								b = s.Pop()
								if b {
									// End the parallel for
									out <- tok.Str
									sync <- nil
									endFunc = true
									fmt.Println("Saliendo de la funcion")
								} else {
									// An rbrace not associated with parallel
									passToken(tok, out, sync)
								}
							default: // Ignore
								passToken(tok, out, sync)
							}
						}
						continue
					} else { // La primera funcion no es un "main".						
						for tok.Token != token.LBRACE {
							passToken(tok, out, sync)
							tok = <-in
						}
						passToken(tok, out, sync)
						// init LBRACE
						s.Push(true) // Llave de apertura de la funcion
						endFunc := false
						for !endFunc {
							tok = <-in
							switch {
							case isPragma(tok): // Reconocedor de "pragma gomp"
								num_prag, numBarriers = pragma_rewrite(tok, in, out, sync, num_prag, in_parallel, routine_num, for_threads, numBarriers)
							case tok.Str == "var": // Tratamiento para declaración de variables.
								num_dec++ // Numero de declaraciones de variables (para testeo).
								passToken(tok, out, sync)
								tok = <-in
								fmt.Println("Variable local:", tok.Str)
								if tok.Token == token.LPAREN {
									// Declaracion simple
									varList = Var_concat(varList, Var_simple_processor(tok, in, out, sync))
								} else {
									// Declaracion multiple
									varList = Var_concat(varList, Var_multi_processor(tok, in, out, sync))
								}
							case tok.Token == token.LBRACE:
								// An lbrace not associated with parallel
								s.Push(false)
								passToken(tok, out, sync)
							case tok.Token == token.RBRACE:
								b = s.Pop()
								if b {
									// End the parallel for
									out <- tok.Str
									sync <- nil
									endFunc = true
									fmt.Println("Saliendo de la funcion")
								} else {
									// An rbrace not associated with parallel
									passToken(tok, out, sync)
								}
							default: // Ignore
								passToken(tok, out, sync)
							}
						}
						continue
					}
				} else { // No es la primera funcion del código.
					passToken(tok, out, sync)
					tok = <-in
					fmt.Println("Entrando en la funcion:", tok.Str)
					if tok.Str == "main" { // Funcion "main".
						for tok.Token != token.LBRACE {
							passToken(tok, out, sync)
							tok = <-in
						}
						// Inicializa el numero de CPUs
						out <- tok.Str + "\n" + "_init_numCPUs()\n"
						sync <- nil
						// init LBRACE
						s.Push(true) // Llave de apertura de la funcion
						endFunc := false
						for !endFunc {
							tok = <-in
							switch {
							case isPragma(tok): // Reconocedor de "pragma gomp"
								num_prag, numBarriers = pragma_rewrite(tok, in, out, sync, num_prag, in_parallel, routine_num, for_threads, numBarriers)
							case tok.Str == "var": // Tratamiento para declaración de variables.
								num_dec++ // Numero de declaraciones de variables (para testeo).
								passToken(tok, out, sync)
								tok = <-in
								fmt.Println("Variable local:", tok.Str)
								if tok.Token == token.LPAREN {
									// Declaracion simple
									varList = Var_concat(varList, Var_simple_processor(tok, in, out, sync))
								} else {
									// Declaracion multiple
									varList = Var_concat(varList, Var_multi_processor(tok, in, out, sync))
								}
							case tok.Token == token.LBRACE:
								// An lbrace not associated with parallel
								s.Push(false)
								passToken(tok, out, sync)
							case tok.Token == token.RBRACE:
								b = s.Pop()
								if b {
									// End the parallel for
									out <- tok.Str
									sync <- nil
									endFunc = true
									fmt.Println("Saliendo de la funcion")
								} else {
									// An rbrace not associated with parallel
									passToken(tok, out, sync)
								}
							default: // Ignore
								passToken(tok, out, sync)
							}
						}
						continue
					} else { // Otra funcion.
						for tok.Token != token.LBRACE {
							passToken(tok, out, sync)
							tok = <-in
						}
						passToken(tok, out, sync)
						// init LBRACE
						s.Push(true) // Llave de apertura de la funcion
						endFunc := false
						for !endFunc {
							tok = <-in
							switch {
							case isPragma(tok): // Reconocedor de "pragma gomp"
								num_prag, numBarriers = pragma_rewrite(tok, in, out, sync, num_prag, in_parallel, routine_num, for_threads, numBarriers)
							case tok.Str == "var": // Tratamiento para declaración de variables.
								num_dec++ // Numero de declaraciones de variables (para testeo).
								passToken(tok, out, sync)
								tok = <-in
								fmt.Println("Variable local:", tok.Str)
								if tok.Token == token.LPAREN {
									// Declaracion simple
									varList = Var_concat(varList, Var_simple_processor(tok, in, out, sync))
								} else {
									// Declaracion multiple
									varList = Var_concat(varList, Var_multi_processor(tok, in, out, sync))
								}
							case tok.Token == token.LBRACE:
								// An lbrace not associated with parallel
								s.Push(false)
								passToken(tok, out, sync)
							case tok.Token == token.RBRACE:
								b = s.Pop()
								if b {
									// End the parallel for
									out <- tok.Str
									sync <- nil
									endFunc = true
									fmt.Println("Saliendo de la funcion")
								} else {
									// An rbrace not associated with parallel
									passToken(tok, out, sync)
								}
							default: // Ignore
								passToken(tok, out, sync)
							}
						}
						continue
					}
				}
			case tok.Str == "var": // Tratamiento para declaración de variables.
			num_dec++ // Numero de declaraciones de variables (para testeo).
			passToken(tok, out, sync)
			tok = <-in
			fmt.Println("Variable global:", tok.Str)
			if tok.Token == token.LPAREN {
				// Declaracion simple
				varList = Var_concat(varList, Var_simple_processor(tok, in, out, sync))
				continue
			} else {
				// Declaracion multiple
				varList = Var_concat(varList, Var_multi_processor(tok, in, out, sync))
				continue
			}
			default: // Ignore
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
