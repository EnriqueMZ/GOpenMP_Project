/*
 ============================================================================
 Name        : gomp.go
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
               
 Description : Main GOpenMP pre-processor module
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
	"gomp_lib"
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

// Check is a token is a "pragma gomp".
func isPragma(token Token) bool {
	res := strings.HasPrefix(noSpaces(token.Str), "//pragmagomp")
	return res
}

// Number of variable declarations. Test only.
var num_dec int = 0

// Variable declarations list.
var varGlobalList []Variable

// Number of pragmas. Test only.
var num_prag int = 0

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

// Barrier function. Add a barrier when necessary.
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
		panic("Error: Invalid operator in redution clause.")
	}
	return str
}

// Function that constructs a barrier for a reduction variable.
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

// Function to retrive the type of a varible marked as "reduction" from the Variable list.
// Launch panic if the the variable has not been previously initialized.
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
		panic("Variable \"" + id + "\" in reduction clause not previously initialized.")
	}
	return typ
}

// Function that constructs barriers for all variables in a reduction clause.
func barrier_single_reduction(numBarrier int, clause Reduction_Type, varGlobalList []Variable, varLocalList []Variable) (string, string, string, string, int) {
	var dcls, sends, var_dcls, rcvs string
	var numB int = numBarrier
	opr := stringOperator(clause.Operator)
	for i := range clause.Variables {
		typ := search_typ(clause.Variables[i], varGlobalList, varLocalList)
		dcl, send, var_dcl, rcv := barrier_variable(numB, clause.Variables[i], typ, opr)
		dcls = dcls + dcl
		sends = sends + send
		var_dcls = var_dcls + var_dcl
		rcvs = rcvs + rcv
		numB++
	}
	return dcls, sends, var_dcls, rcvs, numB
}

// Function that constructs barriers for all variables in a reduction clauses list.
func barrier_list_reduction(numBarrier int, reductionList []Reduction_Type, varGlobalList []Variable, varLocalList []Variable) (string, string, string, string, int) {
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
			dcls_aux, sends_aux, var_dcls_aux, rcvs_aux, numB_aux := barrier_single_reduction(numBarrier, reductionList[i], varGlobalList, varLocalList)
			dcls = dcls + dcls_aux
			sends = sends + sends_aux
			var_dcls = var_dcls + var_dcls_aux
			rcvs = rcvs + rcvs_aux
			numBarrier = numB_aux
		}
	}
	return dcls, sends, var_dcls, rcvs, numBarrier
}

// Funtion routineNum. It process token Gomp_get_routine_num()
func subs_Gomp_get_routine_num(in chan Token, out chan string, sync chan interface{}) {
	out <- "_routine_num"
	sync <- nil
	tok := <-in
	// Brackets. This can be eliminate, and leave the error treating to the compiler.
	if tok.Token != token.LPAREN {
		panic("Error: Left bracket lost in Gomp_get_routine_num.")
	}
	eliminateToken(out, sync)
	tok = <-in
	if tok.Token != token.RPAREN {
		panic("Error: Right bracket lost in Gomp_get_routine_num.")
	}
	eliminateToken(out, sync)
}

// Function that ingnore Gomp_set_num_routine() where applicable.
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

// Function that compares variables from a "default(none)" clause with declared variables.
func var_not_prev_declare(pragma Pragma, varGlobalList []Variable, varLocalList []Variable) (bool, string) {
	var res bool
	var v string
	if len(pragma.Variable_List) == 0 {
		res = true
		v = ""
	} else {
		for i := range pragma.Variable_List {
			res, v = notFound(pragma.Variable_List[i], varGlobalList)
			if res {
				break
			}
		}
		for i := range pragma.Variable_List {
			res, v = notFound(pragma.Variable_List[i], varLocalList)
			if res {
				break
			}
		}
	}
	return res, v
}

func declare(ident string, varGlobalList []Variable, varLocalList []Variable) string { // WARNING: May need correction.
	var res string
	var enc bool = false
	for i := range varGlobalList {
		if ident == varGlobalList[i].Ident {
			enc = true
			if varGlobalList[i].Type == "no_type" {
				res = varGlobalList[i].Ident + "= reflect.Zero(reflect.TypeOf(" + varGlobalList[i].Ident + ")).Interface()"
			} else {
				res = varGlobalList[i].Ident + " " + varGlobalList[i].Type
			}
			break
		}
	}
	for i := range varLocalList {
		if ident == varLocalList[i].Ident {
			enc = true
			if varLocalList[i].Type == "no_type" {
				res = varLocalList[i].Ident + "= reflect.Zero(reflect.TypeOf(" + varLocalList[i].Ident + ")).Interface()"
			} else {
				res = varLocalList[i].Ident + " " + varLocalList[i].Type
			}
			break
		}
	}
	if !enc {
		panic("Variable \"" + ident + "\" in private clause not previously initialized.")
	}
	return res
}


// Funcion that writes the string code for private variables re-initialization.
// WARNING: Implicitly declared variables.
func declareList(pragma Pragma, varGlobalList []Variable, varLocalList []Variable) string {
	var res string
	if len(pragma.Private_List) == 0 {
		res = ""
	} else {
		res = declare(pragma.Private_List[0], varGlobalList, varLocalList)
		for i := 1; i < len(pragma.Private_List); i++ {
			res = res + ";\n" + declare(pragma.Private_List[i], varGlobalList, varLocalList)
		}
	}
	return res
}

// Function that rewrites the code if it is a pragma.
func pragma_rewrite(tok Token, in chan Token, out chan string, sync chan interface{}, num_prag int, in_parallel bool, routine_num string, for_threads string, numBarriers int, varLocalList []Variable) (int, int) {
	num_prag++
	fmt.Print("\n")
	fmt.Println("  Number of current pragmas: ", num_prag, "\n")
	fmt.Println("  Pragma: ", tok.Str, "\n")
	pragma := ProcessPragma(tok.Str)
	fmt.Println("  Pragma information: ", pragma, "\n")

	switch pragma.Type { // Pragma type treatment
	case 0: // PARALLEL PRAGMA
		var b bool
		var s braceStack
		in_parallel = true
		endParallel := false
		routine_num = "_routine_num"     // String with goroutine number ID variable.
		for_threads = pragma.Num_threads // String with threads number in Parallel block.
		fmt.Println("  Pragma type: PARALLEL \n")
		
		// Check Default clause
		if pragma.Default == NONE {
			def_cond, def_var := var_not_prev_declare(pragma, varGlobalList, varLocalList)
			if def_cond {
				panic("Error: Variable \"" + def_var + "\" not previously initialized.")
			}
		}

		// REDUCTION VARIABLES
		dcls, sends, var_dcls, rcvs, numB := barrier_list_reduction(numBarriers, pragma.Reduction_List, varGlobalList, varLocalList)
		numBarriers = numB
		out <- dcls + "for _i := 0; _i < " + pragma.Num_threads + "; _i++{\n" + "go func(_routine_num int)"
		sync <- nil

		tok = <-in
		// init LBRACE
		if tok.Token != token.LBRACE {
			panic("Error: Missing init brace in pragma.")
		}
		s.Push(true) // Init brace in Parallel block.

		// PRIVATE VARIABLES
		privateList := declareList(pragma, varGlobalList, varLocalList)
		fmt.Println("  Private variables in pragma Parallel:", privateList)

		// Private and reduction variables redeclarations.
		out <- " {" + "var (" + privateList + ") \n" + var_dcls
		sync <- nil

		// Parallel content treatment
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
				num_prag, numBarriers = pragma_rewrite(tok, in, out, sync, num_prag, in_parallel, routine_num, for_threads, numBarriers, varLocalList)
			default:
				// Ignore
				passToken(tok, out, sync)
			}
		}
		in_parallel = false
	case 1: // PARALLEL_FOR PRAGMA
		var b bool
		var s braceStack
		var iterations string = "0" // Parallelized loop iterations. Test only.
		var ini, var_index, assign string
		fmt.Println("  Pragma type: PARALLEL_FOR \n")

		// Check Default clause
		if pragma.Default == NONE {
			def_cond, def_var := var_not_prev_declare(pragma, varGlobalList, varLocalList)
			if def_cond {
				panic("Error: Variable \"" + def_var + "\" not previously initialized.")
			}
		}
		// REDUCTION VARIABLES
		dcls, sends, var_dcls, rcvs, numB := barrier_list_reduction(numBarriers, pragma.Reduction_List, varGlobalList, varLocalList)
		numBarriers = numB
		out <- dcls // Changes pragma by channels declaration.
		sync <- nil

		tok = <-in // "for" token
		fmt.Println("  Variables declared before Parallel For block:", varGlobalList, varLocalList)
		iterations, ini, var_index, assign, tok = For_parallel_declare(tok, in, out, sync, varGlobalList, varLocalList)
		fmt.Println("  Parallelized loop iterations:", iterations)

		// PRIVATE VARIABLES
		privateList := declareList(pragma, varGlobalList, varLocalList)
		fmt.Println("  Private variables in Parallel For pragma:", privateList)

		// Goroutines start. Varibles re-declatations.
		out <- tok.Str + "\n" + "go func(_routine_num int) {\n" + "var (" + privateList + ") \n" + var_dcls + "for " + var_index + " " + assign + " _routine_num + " + ini + "; " + var_index + " <" + iterations + "; " + var_index + " += _numCPUs {\n"
		sync <- nil

		// init LBRACE
		s.Push(true) // Init brace in Parallel For block.
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
			case tok.Str == "Gomp_get_routine_num":
				subs_Gomp_get_routine_num(in, out, sync)
			case tok.Str == "Gomp_set_num_routines":
				ign_Gomp_set_num_routine(in, out, sync)
			default:
				// Ignore
				passToken(tok, out, sync)
			}
		}
	case 2: // FOR PRAGMA
		var b bool
		var s braceStack
		var iterations string = "0" // Parallelized loop iterations. Test only.
		fmt.Println("  Pragma type: FOR \n")

		// Check Default clause
		if pragma.Default == NONE {
			def_cond, def_var := var_not_prev_declare(pragma, varGlobalList, varLocalList)
			if def_cond {
				panic("Error: Variable \"" + def_var + "\" not previously initialized.")
			}
		}

		/*
			// REDUCTION VARIABLES
			dcls, sends, var_dcls, rcvs, numB := barrier_list_reduction(numBarriers, pragma.Reduction_List, varList)
			numBarriers = numB

			out <- dcls
			sync <- nil
		*/

		eliminateToken(out, sync) // Remove pragma from code

		tok = <-in // "for" token
		fmt.Println("  Variables declared before For block:", varGlobalList, varLocalList)
		tok = For_declare(tok, in, out, sync, varGlobalList, varLocalList, routine_num, for_threads)
		fmt.Println("Parallelized loop iterations:", iterations)

		// PRIVATE VARIABLES
		privateList := declareList(pragma, varGlobalList, varLocalList)
		fmt.Println("  Private variables in For pragma:", privateList)

		// Goroutines start. Varibles re-declatations.
		//out <- tok.Str + "\n" + "var (" + privateList + ") \n" + var_dcls + "for _i := _routine_num; _i <" + iteraciones + "; _i += _numCPUs {\n"
		out <- tok.Str + "\n" + "var (" + privateList + ") \n"
		sync <- nil

		// init LBRACE
		s.Push(true) // Init brace in For block.
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
			case tok.Str == "Gomp_get_routine_num":
				subs_Gomp_get_routine_num(in, out, sync)
			case tok.Str == "Gomp_set_num_routines":
				ign_Gomp_set_num_routine(in, out, sync)
			default:
				// Ignore
				passToken(tok, out, sync)
			}
		}
	case 3: // THREADPRIVATE PRAGMA
		fmt.Println("  Pragma type: THREADPRIVATE \n")
		eliminateToken(out, sync)
		// TO DO: Rest of pragma treatment
	}
	return num_prag, numBarriers
}

// MAIN PROGRAM.
func main() {
	_fIn,_ := os.Open(os.Args[1]) 
	p := PipeInit(_fIn)
	//Lines(p) // Show pre-processor lines
	Link(func(in chan Token,
		tOut chan Token,
		out chan string,
		sync chan interface{}) {
		var numFunc int = 0          // Initializes function counter of the original code.
		var numBarriers int = 0      // Inicializes barrier counter.
		var in_parallel bool = false // Inside parallel block?
		var routine_num string = "0" // Goroutine ID string
		var for_threads = "1" 		 // String with for-loop thread number.
		
		fmt.Print("\n")
		fmt.Println("  ===============================================")
		fmt.Print("\n")
		fmt.Println("  GOPENMP_PREPROCESSOR")
		fmt.Println("  Go/OpenMP version")
		fmt.Print("\n")
		fmt.Println("  Number of processors available = ", gomp_lib.Gomp_get_num_procs())
		fmt.Println("  Number of threads =              ", gomp_lib.Gomp_get_num_routines())
		fmt.Print("\n")
		fmt.Println("  ===============================================")
		fmt.Print("\n")
		fmt.Println("  Start preprocessign...")
		fmt.Print("\n")
		
		for tok := range in {
			
			switch { // Token treatment

			case tok.Token == token.IMPORT: // Import treatment
				Imports_declare(tok, in, out, sync)
				continue

			case tok.Token == token.FUNC:   // Variable _numCPUs treatment
				var b bool
				var s braceStack
				var varLocalList []Variable // Local variable list of a function.
				if numFunc == 0 {           // First function declare in the original code.
					numFunc++
					out <- "var _numCPUs = runtime.NumCPU()\n" + "func _init_numCPUs(){\n" + "runtime.GOMAXPROCS(_numCPUs)\n" + "}\n" + tok.Str
					sync <- nil
					tok = <-in
					fmt.Println("  Now entering the first function:", tok.Str)
					if tok.Str == "main" { // First function is the "main" funtion.
						for tok.Token != token.LPAREN {
							passToken(tok, out, sync)
							tok = <-in
						}
						varLocalList = Var_argument_processor(tok, in, out, sync, varLocalList)
						fmt.Println("  Argument list of the function:", varLocalList)
						tok = <-in
						for tok.Token != token.LBRACE {
							passToken(tok, out, sync)
							tok = <-in
						}
						// Initializes number of CPUs
						out <- tok.Str + "\n" + "_init_numCPUs()\n"
						sync <- nil
						// init LBRACE
						s.Push(true) // Init brace in function.
						endFunc := false
						for !endFunc {
							tok = <-in
							switch {
							case isPragma(tok): // "pragma gomp" recognizer
								num_prag, numBarriers = pragma_rewrite(tok, in, out, sync, num_prag, in_parallel, routine_num, for_threads, numBarriers, varLocalList)
							case tok.Str == "var": // Variable declaration treatment
								num_dec++ 		   // Variable declaration counter (test only)
								passToken(tok, out, sync)
								tok = <-in
								fmt.Println("  Local variable:", tok.Str)
								if tok.Token == token.LPAREN {
									// Simple declaration
									varLocalList = Var_concat(varLocalList, Var_simple_processor(tok, in, out, sync))
								} else {
									// Multiple declaration
									varLocalList = Var_concat(varLocalList, Var_multi_processor(tok, in, out, sync))
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
									fmt.Println("  Variable list declared in this function: ", varLocalList, "\n")
								} else {
									// An rbrace not associated with parallel
									passToken(tok, out, sync)
								}
							default: // Ignore
								passToken(tok, out, sync)
							}
						}
						continue
					} else { // First function is not the "main" funtion.
						for tok.Token != token.LPAREN {
							passToken(tok, out, sync)
							tok = <-in
						}
						varLocalList = Var_argument_processor(tok, in, out, sync, varLocalList)
						fmt.Println("  Argument list of the function:", varLocalList)
						tok = <-in
						for tok.Token != token.LBRACE {
							passToken(tok, out, sync)
							tok = <-in
						}
						passToken(tok, out, sync)
						// init LBRACE
						s.Push(true) // Init brace in function.
						endFunc := false
						for !endFunc {
							tok = <-in
							switch {
							case isPragma(tok): // "pragma gomp" recognizer
								num_prag, numBarriers = pragma_rewrite(tok, in, out, sync, num_prag, in_parallel, routine_num, for_threads, numBarriers, varLocalList)
							case tok.Str == "var": // Variable declaration treatment
								num_dec++ 		   // Variable declaration counter (test only)
								passToken(tok, out, sync)
								tok = <-in
								fmt.Println("  Local variable:", tok.Str)
								if tok.Token == token.LPAREN {
									// Simple declaration
									varLocalList = Var_concat(varLocalList, Var_simple_processor(tok, in, out, sync))
								} else {
									// Multiple declaration
									varLocalList = Var_concat(varLocalList, Var_multi_processor(tok, in, out, sync))
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
									fmt.Println("  Variable list declared in this function: ", varLocalList, "\n")
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
				} else { // Not first function declare in the original code.
					passToken(tok, out, sync)
					tok = <-in
					fmt.Println("  Now entering the function:", tok.Str)
					if tok.Str == "main" { // "main" function.
						for tok.Token != token.LPAREN {
							passToken(tok, out, sync)
							tok = <-in
						}
						varLocalList = Var_argument_processor(tok, in, out, sync, varLocalList)
						fmt.Println("  Argument list of the function:", varLocalList)
						tok = <-in
						for tok.Token != token.LBRACE {
							passToken(tok, out, sync)
							tok = <-in
						}
						// Initializes number of CPUs
						out <- tok.Str + "\n" + "_init_numCPUs()\n"
						sync <- nil
						// init LBRACE
						s.Push(true) // Init brace in function.
						endFunc := false
						for !endFunc {
							tok = <-in
							switch {
							case isPragma(tok): // "pragma gomp" recognizer
								num_prag, numBarriers = pragma_rewrite(tok, in, out, sync, num_prag, in_parallel, routine_num, for_threads, numBarriers, varLocalList)
							case tok.Str == "var": // Variable declaration treatment
								num_dec++ 		   // Variable declaration counter (test only)
								passToken(tok, out, sync)
								tok = <-in
								fmt.Println("  Local variable:", tok.Str)
								if tok.Token == token.LPAREN {
									// Simple declaration
									varLocalList = Var_concat(varLocalList, Var_simple_processor(tok, in, out, sync))
								} else {
									// Multiple declaration
									varLocalList = Var_concat(varLocalList, Var_multi_processor(tok, in, out, sync))
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
									fmt.Println("  Variable list declared in this function: ", varLocalList, "\n")
								} else {
									// An rbrace not associated with parallel
									passToken(tok, out, sync)
								}
							default: // Ignore
								passToken(tok, out, sync)
							}
						}
						continue
					} else { // Another function
						for tok.Token != token.LPAREN {
							passToken(tok, out, sync)
							tok = <-in
						}
						varLocalList = Var_argument_processor(tok, in, out, sync, varLocalList)
						fmt.Println("  Argument list of the function:", varLocalList)
						tok = <-in
						for tok.Token != token.LBRACE {
							passToken(tok, out, sync)
							tok = <-in
						}
						passToken(tok, out, sync)
						// init LBRACE
						s.Push(true) // Init brace in function
						endFunc := false
						for !endFunc {
							tok = <-in
							switch {
							case isPragma(tok): // "pragma gomp" recognizer
								num_prag, numBarriers = pragma_rewrite(tok, in, out, sync, num_prag, in_parallel, routine_num, for_threads, numBarriers, varLocalList)
							case tok.Str == "var": // Variable declaration treatment
								num_dec++ 		   // Variable declaration counter (test only)
								passToken(tok, out, sync)
								tok = <-in
								fmt.Println("  Local variable:", tok.Str)
								if tok.Token == token.LPAREN {
									// Simple declaration
									varLocalList = Var_concat(varLocalList, Var_simple_processor(tok, in, out, sync))
								} else {
									// Multiple declaration
									varLocalList = Var_concat(varLocalList, Var_multi_processor(tok, in, out, sync))
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
									fmt.Println("  Variable list declared in this function: ", varLocalList, "\n")
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
			case tok.Str == "var": // Variable declaration treatment
				num_dec++ 		   // Variable declaration counter (test only)
				passToken(tok, out, sync)
				tok = <-in
				fmt.Println("  Global variable:", tok.Str)
				if tok.Token == token.LPAREN {
					// Simple declaration
					varGlobalList = Var_concat(varGlobalList, Var_simple_processor(tok, in, out, sync))
					continue
				} else {
					// Multiple declaration
					varGlobalList = Var_concat(varGlobalList, Var_multi_processor(tok, in, out, sync))
					continue
				}
			default: // Ignore
				passToken(tok, out, sync)
				continue
			}
		}
		close(tOut)
	})(p)
	
	_fOut,_ := os.Create(os.Args[2])
	PipeEnd(p, _fOut)
	
	fmt.Println("  Number of variable declarations in original code: ", num_dec, "\n")
	fmt.Println("  Variable declared list in original code: ", varGlobalList, "\n")
	fmt.Println("  Preprocessing finished. \n")
}
