/*
 ============================================================================================
 Name        : pragma_processor.go
 Author      : Enrique Madridejos Zamorano
 Version     :
 Copyright   : Apache Licence Version 2.0
 Description : Module that handles "pragma gomp" expresions from the original code.
 ============================================================================================
 */

package pragma_processor

import (
	"runtime"
	"strings"
	"unicode/utf8"
)
// ==========================================================================================
// Private functions
// ==========================================================================================

/*	
	Function explode splits s into an array of UTF-8 sequences,
	one per Unicode character (still strings) up to a maximum of n (n < 0 means no limit).
	Invalid UTF-8 sequences become correct encodings of U+FFF8.
	
	Import from golang source (all rights reserved)
*/
func explode(s string, n int) []string {
	if n == 0 {
		return nil
	}
	l := utf8.RuneCountInString(s)
	if n <= 0 || n > l {
		n = l
	}
	a := make([]string, n)
	var size int
	var ch rune
	i, cur := 0, 0
	for ; i+1 < n; i++ {
		ch, size = utf8.DecodeRuneInString(s[cur:])
		if ch == utf8.RuneError {
			a[i] = string(utf8.RuneError)
		} else {
			a[i] = s[cur : cur+size]
		}
		cur += size
	}
	// add the rest, if there is any
	if cur < len(s) {
		a[i] = s[cur:]
	}
	return a
}
/*
	Generic split: splits after each instance of sep,
	including sepSave bytes of sep in the subarrays.
	
	Import from golang source (all rights reserved)
*/
func genSplit(s, sep string, sepSave, n int) []string {
	if n == 0 {
		return nil
	}
	if sep == "" {
		return explode(s, n)
	}
	if n < 0 {
		n = strings.Count(s, sep) + 1
	}
	c := sep[0]
	start := 0
	a := make([]string, n)
	na := 0
	for i := 0; i+len(sep) <= len(s) && na+1 < n; i++ {
		if s[i] == c && (len(sep) == 1 || s[i:i+len(sep)] == sep) {
			a[na] = s[start : i+sepSave]
			na++
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	a[na] = s[start:]
	return a[0 : na+1]
}
/*
	Function that splits a string into different parts, separated by a given character.
	Additionally, it separates the first element of the resulting slice.
	Return the first element and the slice.
*/
func splitFirstBySep(s, sep string) (string, []string) {
	return genSplit(s, sep, 0, 2)[0], genSplit(s, sep, 0, 2)[1:]
}

// Function that removes whitespaces within a string. 
func noSpaces(str string) string {
	return strings.Replace(str, " ", "", -1)
}

// Function that removes commas within a string
func noCommas(str string) string {
	return strings.Replace(str, ",", "", -1)
}

// Function that searches repeated elements within a slice of strings.
// Return the first elemento found repeted, if exists.
func repeatIn(a []string) (bool, string) {
	var res bool = false
	var elem = ""
	for i := 0; i < len(a)-1; i++ {
		for j := i + 1; j < len(a); j++ {
			if a[i] == a[j] {
				res = true
				elem = a[i]
			}
		}
	}
	return res, elem
}

// Function that searches repeated elements between two slice of strings.
// Return the first elemento found repeted, if exists.

func repeat(a, b []string) (bool, string) {
	var res bool = false
	var rep string
	for i := range b {
		for j := range a {
			if a[j] == b[i] {
				res = true
				rep = a[j]
			}
		}
	}
	return res, rep
}

// Function that concatenates two slices of strings.
func concat(a, b []string) []string {
	for i := range b {
		a = append(a, b[i])
	}
	return a
}

// ==========================================================================================

// Pragma types

type Pragma_Type int

const (
	PARALLEL Pragma_Type = iota
	PARALLEL_FOR
	FOR
	THREADPRIVATE
	// Add as needed
)

type Clause_Type int

const (
	DEFAULT Clause_Type = iota
	IF
	NUM_THREADS
	SHARED
	PRIVATE
	FIRSTPRIVATE
	COPYIN
	REDUCTION
	LASTPRIVATE
	SCHEDULE
	COLLAPSE
	ORDERED
	NOWAIT
	// Add as needed
)

type Spec interface {
	// Spec can be *ListSpec, *DefaultSpec, *EscalarSpec, *ReductionSpec
}

type ListSpec struct {
	Variables []string // Variables list
}

type Default_Type int

const (
	SHA Default_Type = iota
	NONE
)

type DefaultSpec struct {
	Type Default_Type
}

type EscalarSpec struct {
	Value int
}

type Red_Operator int

const (
	SUMA Red_Operator = iota
	PRODUCTO
	RESTA
	ANDB
	ORB
	XORB
	AND
	OR
)

type Reduction_Type struct {
	Operator  Red_Operator // Operator
	Variables []string     // Variables which affect the operator
}

type Clause struct {
	Type  Clause_Type
	Specs []Spec
}

type Pragma struct {
	Type            Pragma_Type 	// Pragma type
	Default         Default_Type
	Def_Num_threads int 			// Number of threads by default (goroutines)
	Num_threads     string
	If_content      string
	Variable_List   []string 		// Complete list of variables within the pragma (for searching)
	Shared_List     []string
	Private_List    []string
	First_List      []string
	Last_List       []string
	Copyin_List     []string
	Reduction_List  []Reduction_Type
	// Add as needed
}

// Group clauses valid for each type of pragma
var (
	parallelClauses    = []string{"if", "num_threads", "default", "shared", "private", "firstprivate", "copyin", "reduction"}
	parallForClauses   = []string{"if", "num_threads", "default", "shared", "private", "firstprivate", "copyin", "reduction", "lastprivate", "schedule", "collapse", "ordered"}
	forClauses         = []string{"private", "firstprivate", "lastprivate", "reduction", "schedule", "collapse", "ordered", "nowait"}
	thPrivateClauses   = []string{""}
	reductionOperators = []string{"+", "*", "-", "&", "|", "^", "&&", "||"}
	// Add as needed
)

// Function that identifies a pragma gomp.
// Additionaly, it returns pragma content.
func isPragmaGomp(pragma string) (bool, []string) {
	var res bool = true
	pragmaAux := strings.TrimSpace(pragma) // Removes whitespaces at the beginning and the end of the pragma.
	prag, auxList := splitFirstBySep(pragmaAux, " ")
	aux := strings.TrimSpace(auxList[0])
	gomp, list := splitFirstBySep(aux, " ")
	if !(noSpaces(prag) == "//pragma" && noSpaces(gomp) == "gomp") {
		res, list = false, nil
	}
	return res, list
}

/*
	Function that identifies and return the type of pragma from pragma's body.
	Additionally, it returns a slices containing clauses from the pragma. 
	Launch "panic" if the type is incorrect.
*/
func pragmaType(body string) (Pragma_Type, []string) {
	var typID Pragma_Type = -1
	var res []string
	bodyAux := strings.TrimSpace(body)
	typ, auxList := splitFirstBySep(bodyAux, " ")
	res = auxList
	switch typ {
	case "parallel":
		if len(auxList) == 0 {
			typID = 0
			break
		}
		subTyp, list := splitFirstBySep(auxList[0], " ")
		if subTyp == "for" {
			typID = 1
			res = list
		} else {
			typID = 0
		}
	case "for":
		typID = 2
	case "threadprivate":
		typID = 3
	// Add cases as needed
	default:
		panic("Error: Pragma_Gomp type incorrect")
	}
	return typID, res
}

// Function that checks whether a clause belongs to a given list (it is valid for that list).
// Ignore whitespaces within the clause.
func validClause(clause string, group []string) bool {
	var res bool = false
	for i := range group {
		if noSpaces(clause) == group[i] {
			res = true
			break
		}
	}
	return res
}

// Sorting function clauses, according to its string.
func clauseType(clause string) Clause_Type { // Perhaps redundant?
	var typ Clause_Type
	switch clause {
	case "default":
		typ = 0
	case "if":
		typ = 1
	case "num_threads":
		typ = 2
	case "shared":
		typ = 3
	case "private":
		typ = 4
	case "firstprivate":
		typ = 5
	case "copyin":
		typ = 6
	case "reduction":
		typ = 7
	case "lastprivate":
		typ = 8
	case "schedule":
		typ = 9
	case "collapse":
		typ = 10
	case "ordered":
		typ = 11
	case "nowait":
		typ = 12
	default:
		typ = 13 // List clause
		// Add as needed
	}
	return typ
}

func splitClauses(prgTyp Pragma_Type, clauses string) (Clause_Type, string, []string) {
	var group []string
	switch prgTyp { // Determines the corresponding clauses group.
	case 0:
		group = parallelClauses
	case 1:
		group = parallForClauses
	case 2:
		group = forClauses
	case 3:
		group = thPrivateClauses
	}
	clauseAll, res := splitFirstBySep(clauses, "(")
	clauseAux := noSpaces(clauseAll)
	clause := noCommas(clauseAux)
	if !validClause(clause, group) {
		panic("Error: Invalidad clause \"" + clause + "\" inside this pragma")
	}
	return clauseType(clause), clause, res
}

// Functions for clauses content.

func clauseContentList(body string) ([]string, []string) {
	contAll, res := splitFirstBySep(body, ")")
	contS := noSpaces(contAll)
	cont := strings.Split(contS, ",")
	if cont[0] == "" {
		panic("No specified variables inside the clause")
	}
	rep, elem := repeatIn(cont)
	if rep {
		panic("Variable \"" + elem + "\" repeated in clause content")
	}
	return cont, res
}

func clauseContentDef(body string) (Default_Type, []string) {
	var typ Default_Type
	contAll, res := splitFirstBySep(body, ")")
	cont := noSpaces(contAll)
	switch cont {
	case "shared":
		typ = 0
	case "none":
		typ = 1
	default:
		panic("Error: Invalid argument \"" + cont + "\" in clause Default")
	}
	return typ, res
}

func clauseContentStr(body string) (string, []string) {
	contAll, res := splitFirstBySep(body, ")")
	return contAll, res
}

func validOperator(operator string) bool {
	var res bool = false
	for i := range reductionOperators {
		if noSpaces(operator) == reductionOperators[i] {
			res = true
			break
		}
	}
	return res
}

func operatorType(op string) Red_Operator {
	var typ Red_Operator
	switch op {
	case "+":
		typ = 0
	case "*":
		typ = 1
	case "-":
		typ = 2
	case "&":
		typ = 3
	case "|":
		typ = 4
	case "^":
		typ = 5
	case "&&":
		typ = 6
	case "||":
		typ = 7
	default:
		panic("Error: Invalid operator \"" + op + "\" in clause Reduction")
	}
	return typ
}

func clauseContentRed(body string) (Red_Operator, []string, []string) {
	var typOp Red_Operator
	contAll, res := splitFirstBySep(body, ")")
	cont := noSpaces(contAll)
	op, listAux := splitFirstBySep(cont, ":")
	if !validOperator(op) {
		panic("Error: Invalid operator \"" + op + "\" in clause Reduction")
	}
	typOp = operatorType(op)
	list := strings.Split(listAux[0], ",")
	return typOp, list, res
}

func ProcessPragma(pragma string) Pragma {
	var (
		PragmaProcesed                                                            Pragma
		oneDefault, oneIf, oneNThreads, oneThPrivate                              bool = false, false, false, false
		prgTyp                                                                    Pragma_Type
		defTyp                                                                    Default_Type
		clauseType                                                                Clause_Type
		ifCont, nThreadsCont                                                      string
		clauseList, contentList, variableList, privateList, sharedList, firstList []string
		reductionList                                                             []Reduction_Type
		operator                                                                  Red_Operator
	)
	cond, body := isPragmaGomp(pragma)
	if !cond {
		panic("Error: Not a pragma gomp\n")
	}
	prgTyp, clauseList = pragmaType(body[0])
	PragmaProcesed.Type = prgTyp
	switch prgTyp { // Pragma initializer
	case 0:
		PragmaProcesed.Default = SHA
		PragmaProcesed.Def_Num_threads = runtime.NumCPU()
		PragmaProcesed.Num_threads = "_numCPUs"
		PragmaProcesed.Variable_List = variableList
		PragmaProcesed.Shared_List = sharedList
		PragmaProcesed.Private_List = privateList
		PragmaProcesed.First_List = firstList
		PragmaProcesed.Reduction_List = reductionList
	case 1:
		PragmaProcesed.Default = SHA
		PragmaProcesed.Def_Num_threads = runtime.NumCPU()
		PragmaProcesed.Num_threads = "_numCPUs"
		PragmaProcesed.Variable_List = variableList
		PragmaProcesed.Shared_List = sharedList
		PragmaProcesed.Private_List = privateList
		PragmaProcesed.First_List = firstList
		PragmaProcesed.Reduction_List = reductionList
	case 2:
		PragmaProcesed.Default = SHA
		PragmaProcesed.Def_Num_threads = runtime.NumCPU()
		PragmaProcesed.Num_threads = "_numCPUs"
		PragmaProcesed.Variable_List = variableList
		PragmaProcesed.Shared_List = sharedList
		PragmaProcesed.Private_List = privateList
		PragmaProcesed.First_List = firstList
		PragmaProcesed.Reduction_List = reductionList
	case 3:
		PragmaProcesed.Variable_List = variableList
	default:
		panic("Pragma type still not recognized by the program. In process...")
		//TO DO: // Add remaining pragmas
	}
	for {
		if len(clauseList) == 0 || clauseList[0] == "" { // No clauses, or already finished processing the string of clauses.
			break
		}
		clauseType, _, clauseList = splitClauses(prgTyp, clauseList[0])
		switch clauseType {
		case 0: // DEFAULT
			if oneDefault == true {
				panic("Error: Can not declare several additional Default clauses inside a pragma")
			}
			oneDefault = true
			defTyp, clauseList = clauseContentDef(clauseList[0])
			PragmaProcesed.Default = defTyp
		case 1: // IF
			if oneIf == true {
				panic("Error: Can not declare several additional If clauses inside a pragma")
			}
			oneIf = true
			ifCont, clauseList = clauseContentStr(clauseList[0])
			PragmaProcesed.If_content = ifCont
		case 2: // NUM_THREADS
			if oneNThreads == true {
				panic("Error: Can not declare several additional Num_threads clauses inside a pragma")
			}
			oneNThreads = true
			nThreadsCont, clauseList = clauseContentStr(clauseList[0])
			PragmaProcesed.Num_threads = nThreadsCont
		case 3: // SHARED
			contentList, clauseList = clauseContentList(clauseList[0])
			res, rep := repeat(contentList, PragmaProcesed.Variable_List)
			if res {
				panic("Error: Variable \"" + rep + "\" is repeted in several clauses data")
			}
			PragmaProcesed.Variable_List = concat(contentList, PragmaProcesed.Variable_List)
			PragmaProcesed.Shared_List = concat(contentList, PragmaProcesed.Shared_List)
		case 4: // PRIVATE
			contentList, clauseList = clauseContentList(clauseList[0])
			res, rep := repeat(contentList, PragmaProcesed.Variable_List)
			if res {
				panic("Error: Variable \"" + rep + "\" is repeted in several clauses data")
			}
			PragmaProcesed.Variable_List = concat(contentList, PragmaProcesed.Variable_List)
			PragmaProcesed.Private_List = concat(contentList, PragmaProcesed.Private_List)
		case 5: // FIRSTPRIVATE
			contentList, clauseList = clauseContentList(clauseList[0])
			res, rep := repeat(contentList, PragmaProcesed.Variable_List)
			if res {
				panic("Error: Variable \"" + rep + "\" is repeted in several clauses data")
			}
			PragmaProcesed.Variable_List = concat(contentList, PragmaProcesed.Variable_List)
			PragmaProcesed.First_List = concat(contentList, PragmaProcesed.First_List)
		case 6: // COPYING
			panic("Copyin clause not yet implemented") // To do along with Threadprivate pragma
		case 7: // REDUCTION
			operator, contentList, clauseList = clauseContentRed(clauseList[0])
			res, rep := repeat(contentList, PragmaProcesed.Variable_List)
			if res {
				panic("Error: Variable \"" + rep + "\" is repeted in several clauses data")
			}
			PragmaProcesed.Variable_List = concat(contentList, PragmaProcesed.Variable_List)
			var reductionClause Reduction_Type
			reductionClause.Operator = operator
			reductionClause.Variables = contentList
			PragmaProcesed.Reduction_List = append(PragmaProcesed.Reduction_List, reductionClause)
			
		//TO DO: Add remaining cases
		
		case 13:
			if prgTyp == 3 && oneThPrivate == true {
				panic("Error: Can not declare more than one list in pragma Threadprivate")
			}
			oneThPrivate = true
			contentList, clauseList = clauseContentList(clauseList[0])
			PragmaProcesed.Variable_List = concat(contentList, PragmaProcesed.Variable_List)
		default:
			panic("In process...")
		}
	}
	//TO DO: Add remaining cases
	
	return PragmaProcesed
}
