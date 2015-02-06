package pragma_processor

import (
	"runtime"
	"strings"
	"unicode/utf8"
)

// Funciones privadas

// function explode splits s into an array of UTF-8 sequences,
// one per Unicode character (still strings) up to a maximum of n (n < 0 means no limit).
// Invalid UTF-8 sequences become correct encodings of U+FFF8.
// Import from golang source (all rights reserved)
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

// Generic split: splits after each instance of sep,
// including sepSave bytes of sep in the subarrays.
// Import from golang source (all rights reserved)
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

// Funcion que separa el primer elemento del string antes de un separador dado
// Devuelve el elemento y la lista sin el.
func splitFirstBySep(s, sep string) (string, []string) {
	return genSplit(s, sep, 0, 2)[0], genSplit(s, sep, 0, 2)[1:]
}

// Funcion que elimina los espacio en blanco dentro de un string dado.
func noSpaces(str string) string {
	return strings.Replace(str, " ", "", -1)
}

// Funcion que elimina las comas dentro de un string dado.
func noCommas(str string) string {
	return strings.Replace(str, ",", "", -1)
}

// Funcion que busca elementos repetidos dentro de un slice de strings.
// Devuelve el elmento repetido, si existe.

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

// Funcion que busca elementos repetidos entre dos slices de strings.
// Devuelve el elemento repetido, si existe.
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

// Funcion que concatena dos slices de strings.
func concat(a, b []string) []string {
	for i := range b {
		a = append(a, b[i])
	}
	return a
}

// Tipos para pragmas

type Pragma_Type int

const (
	PARALLEL Pragma_Type = iota
	PARALLEL_FOR
	FOR
	THREADPRIVATE
	// Añadir conforme sea necesario
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
	// Añadir conforme sea necesario
)

type Spec interface {
	// Spec puede ser *ListSpec, *DefaultSpec, *EscalarSpec, *ReductionSpec
}

type ListSpec struct {
	Variables []string // Listas de variables
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
	Operator  Red_Operator //Operador
	Variables []string     //Variables afectadas por el operador
}

type Clause struct {
	Type  Clause_Type
	Specs []Spec
}

type Pragma struct {
	Type            Pragma_Type // Tipo del pragma
	Default         Default_Type
	Def_Num_threads int // Numero de hilos por defecto.
	Num_threads     string
	If_content      string
	Variable_List   []string // Lista completa de variables (para búsquedas)
	Shared_List     []string
	Private_List    []string
	First_List      []string
	Last_List       []string
	Copyin_List     []string
	Reduction_List  []Reduction_Type
	// Añadir campos conforme sea necesario
}

// Grupos de clausulas validas para cada tipo de Pragma
var (
	parallelClauses    = []string{"if", "num_threads", "default", "shared", "private", "firstprivate", "copyin", "reduction"}
	parallForClauses   = []string{"if", "num_threads", "default", "shared", "private", "firstprivate", "copyin", "reduction", "lastprivate", "schedule", "collapse", "ordered"}
	forClauses         = []string{"private", "firstprivate", "lastprivate", "reduction", "schedule", "collapse", "ordered", "nowait"}
	thPrivateClauses   = []string{""}
	reductionOperators = []string{"+", "*", "-", "&", "|", "^", "&&", "||"}
	//  Añadir cuando sea necesario
)

// Función que identifica un pragma gomp
// Además, devuelve el contenido del pragma
func isPragmaGomp(pragma string) (bool, []string) {
	var res bool = true
	pragmaAux := strings.TrimSpace(pragma) // Elimina los espacios en blanco al principio y al final del pragma.
	prag, auxList := splitFirstBySep(pragmaAux, " ")
	aux := strings.TrimSpace(auxList[0])
	gomp, list := splitFirstBySep(aux, " ")
	if !(noSpaces(prag) == "//pragma" && noSpaces(gomp) == "gomp") {
		res, list = false, nil
	}
	return res, list
}

// Funcion que identifica el tipo de pragma en el cuerpo de un pragma
// Además, devuelve el string con las clausulas. Lanza "panic" si el tipo no es correcto.
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
	// Añadir casos cuando vaya siendo necesario
	default:
		panic("Error: Tipo de Pragma Gomp incorrecto")
	}
	return typID, res
}

// Funcion que comprueba si una clausula pertenece a una lista dada (es valida para esa lista)
// Obvia los espacios en blanco que pueda tener la clausula
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

// Funcion que clasifica las clausulas en funcion de su string.
func clauseType(clause string) Clause_Type { // ¿REDUNDANTE?
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
		typ = 13 // Clausula lista
		// Añadir cuando sea necesario
	}
	return typ
}

func splitClauses(prgTyp Pragma_Type, clauses string) (Clause_Type, string, []string) {
	var group []string
	switch prgTyp { // Determina el grupo de clausulas correspondiente al Pragma_Type
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
		panic("Error: Clausula no valida para este pragma")
	}
	return clauseType(clause), clause, res
}

// Funciones para contenido de clausulas

func clauseContentList(body string) ([]string, []string) {
	contAll, res := splitFirstBySep(body, ")")
	contS := noSpaces(contAll)
	cont := strings.Split(contS, ",")
	if cont[0] == "" {
		panic("No se han especificado variables en la clausula")
	}
	rep, elem := repeatIn(cont)
	if rep {
		panic("Variable \"" + elem + "\" repetida en lista de datos")
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
		panic("Error: argumento no valido en clausula default")
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
		panic("Error: operador no valido en clausula reduction")
	}
	return typ
}

func clauseContentRed(body string) (Red_Operator, []string, []string) {
	var typOp Red_Operator
	contAll, res := splitFirstBySep(body, ")")
	cont := noSpaces(contAll)
	op, listAux := splitFirstBySep(cont, ":")
	if !validOperator(op) {
		panic("Error: operador no valido en clausula reduction")
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
		panic("Error: no es un pragma gomp\n")
	}
	prgTyp, clauseList = pragmaType(body[0])
	PragmaProcesed.Type = prgTyp
	switch prgTyp { // Inicializador de pragmas
	case 0:
		PragmaProcesed.Default = SHA
		PragmaProcesed.Def_Num_threads = runtime.NumCPU()
		PragmaProcesed.Variable_List = variableList
		PragmaProcesed.Shared_List = sharedList
		PragmaProcesed.Private_List = privateList
		PragmaProcesed.First_List = firstList
		PragmaProcesed.Reduction_List = reductionList
	case 1:
		PragmaProcesed.Default = SHA
		PragmaProcesed.Def_Num_threads = runtime.NumCPU()
		PragmaProcesed.Variable_List = variableList
		PragmaProcesed.Shared_List = sharedList
		PragmaProcesed.Private_List = privateList
		PragmaProcesed.First_List = firstList
		PragmaProcesed.Reduction_List = reductionList
	case 3:
		PragmaProcesed.Variable_List = variableList
	default:
		panic("Tipo de pragma aun no reconocido. En proceso.")
		//TO DO: añadir pragmas cuando sea necesario.
	}
	for {
		if len(clauseList) == 0 || clauseList[0] == "" { // Si no hay clausulas, o ya se ha procesado el string de clausulas.
			break
		}
		clauseType, _, clauseList = splitClauses(prgTyp, clauseList[0])
		switch clauseType {
		case 0:
			if oneDefault == true {
				panic("Error: no pueden declararse varias clausulas default en un pragma")
			}
			oneDefault = true
			defTyp, clauseList = clauseContentDef(clauseList[0])
			PragmaProcesed.Default = defTyp
		case 1:
			if oneIf == true {
				panic("Error: no pueden declararse varias clausulas if en un pragma")
			}
			oneIf = true
			ifCont, clauseList = clauseContentStr(clauseList[0])
			PragmaProcesed.If_content = ifCont
		case 2:
			if oneNThreads == true {
				panic("Error: no pueden declararse varias clausulas num_threads en un pragma")
			}
			oneNThreads = true
			nThreadsCont, clauseList = clauseContentStr(clauseList[0])
			PragmaProcesed.Num_threads = nThreadsCont
		case 3:
			contentList, clauseList = clauseContentList(clauseList[0])
			res, rep := repeat(contentList, PragmaProcesed.Variable_List)
			if res {
				panic("Error: la variable " + rep + " se encuentra repetida en varias clausulas de datos")
			}
			PragmaProcesed.Variable_List = concat(contentList, PragmaProcesed.Variable_List)
			PragmaProcesed.Shared_List = concat(contentList, PragmaProcesed.Shared_List)
		case 4:
			contentList, clauseList = clauseContentList(clauseList[0])
			res, rep := repeat(contentList, PragmaProcesed.Variable_List)
			if res {
				panic("Error: la variable \"" + rep + "\" se encuentra repetida en varias clausulas de datos")
			}
			PragmaProcesed.Variable_List = concat(contentList, PragmaProcesed.Variable_List)
			PragmaProcesed.Private_List = concat(contentList, PragmaProcesed.Private_List)
		case 5:
			contentList, clauseList = clauseContentList(clauseList[0])
			res, rep := repeat(contentList, PragmaProcesed.Variable_List)
			if res {
				panic("Error: la variable \"" + rep + "\" se encuentra repetida en varias clausulas de datos")
			}
			PragmaProcesed.Variable_List = concat(contentList, PragmaProcesed.Variable_List)
			PragmaProcesed.First_List = concat(contentList, PragmaProcesed.First_List)
		case 6:
			panic("Clausula copyin aun no implementada") // Hacer cuando enfrentemos directiva threadprivate
		case 7:
			operator, contentList, clauseList = clauseContentRed(clauseList[0])
			res, rep := repeat(contentList, PragmaProcesed.Variable_List)
			if res {
				panic("Error: la variable \"" + rep + "\" se encuentra repetida en varias clausulas de datos")
			}
			PragmaProcesed.Variable_List = concat(contentList, PragmaProcesed.Variable_List)
			var reductionClause Reduction_Type
			reductionClause.Operator = operator
			reductionClause.Variables = contentList
			PragmaProcesed.Reduction_List = append(PragmaProcesed.Reduction_List, reductionClause)

		//TODO Añadir resto de casos

		case 13:
			if prgTyp == 3 && oneThPrivate == true {
				panic("Error: no puede declararse más de una lista en Threadprivate")
			}
			oneThPrivate = true
			contentList, clauseList = clauseContentList(clauseList[0])
			PragmaProcesed.Variable_List = concat(contentList, PragmaProcesed.Variable_List)
		default:
			panic("En proceso...")
		}
	}
	// TODO Añadir resto de casos

	return PragmaProcesed
}
