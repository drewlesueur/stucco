// stucco: Library for a Single Threaded Un-explicit Cooperative Coroutine DSL
// Single Threaded: Like JavaScript Everything within a single isntance of stucco, runs in a single conceptual thread. (Go coroutine)
// Un-explicit:  Like Go, all Coroutines/Concurrency are implied/automatic and synchronous. No meed to "await" like in JavaScript, no "function colors"
// Cooperative: at the library level, pausing only happens at known points.
// Coroutines: like goroutines
// DSL: embeddable in existing Go code with unique programming language-like api.

// different from stucco3, this uses idea of cc library
// usiing locks instead of goroutines.
// allows for idiomatic go functions to be interjected.
// all the locking is abstracted
// the elegance of async code in stucco3 and linescript4 etc was to just return nil, and wait for the next coroutine in the channel (or queue)
// for languages that don't have green threads that's a nice way to implement it.
// but Go has green threads already, so I can get the same essential thing with lock boundaries under the hood.

package stucco

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
	"os"
	"sync"
)

type State struct {
	GlobalNumbers []float64
	Vals          *List
	VarsStack     []*Record
	Vars          *Record
	CodeStack     [][]any
	Code          []any
	IStack        []int
	I             int
	CallingParent *State
	LexicalParent *State
	AsyncParent   *State
	Tag           string
	Mu sync.Mutex
}


type Callback struct {
	State        *State
	ReturnValues []any
	Vars         *Record
}

type ScopedBlock struct {
	Code          []any
	LexicalParent *State
}

var GlobalState *State

func init() {
	GlobalState = New()
}




func New() *State {
	_ = log.Println

	s := &State{
		Vals:        NewList(),
		Vars:        NewRecord(),
		Code:        nil,
		I:           -1,
		Mu: sync.Mutex{},
	}
	
	
	s.Mu.Lock() // only unlocks during async actions
	ret := s
	// standard library
	s.E(`
        (
            newScope
            .b as
            .a as
            b valOnly
            a valOnly
            dropScope
        ) .swap as
        (
            newScope
            .a as
            dropScope
        ) .drop as
        // (
        //     .a as
        // ) enclose .drop as
        (
            newScope
            .a as
            a valOnly
            a valOnly
            dropScope
        ) .dup as
        
        ( swap as ) .var as
        ( swap to ) .let as
        
        
    
        .globalCounter 0 var
        .if (
            swap
            guard
            do
        ) var
    
        .loop (
            // globalCounter 1 + .globalCounter to
            // "i_" globalCounter ++
            newScope
            .block as
            .max as
            0 .i as
            // .name "who??" var
            (
                i 1 + .i to
                // "increasing i to " i ++ log
                // "and max is " max ++ log
                i max <= guard
                // i max guardLte
                i block valOnly
                downScope
                do
                upScope
                // "length of scope is " scopeLen ++ say
                repeat
            ) do
            dropScope
        ) var
        
        
        .and (
            swap
            do
            dup
            guard
            drop
            do
        ) var
        
        // .or (
        //     swap
        //
        // )
        // .enclose (
        //     .theBlock as
        //     {
        //         .Code theBlock valOnly
        //         .LexicalParent __state .LexicalParent at
        //     }
        // ) var
	`)
	return ret
}

func (s *State) BeDone() * State {
	if len(s.CodeStack) > 0 {
	    s.PopCodeStack()
	    return s
	}
	return s.CallingParent
}
func (s *State) PopCodeStack() {
    s.Code = s.CodeStack[len(s.CodeStack)-1]
    s.I = s.IStack[len(s.IStack)-1]
    
    s.CodeStack = s.CodeStack[0:len(s.CodeStack)-1]
    s.IStack = s.IStack[0:len(s.IStack)-1]
}

// func R(code ...any) *State {
// 	return GlobalState.R(code...)
// }
func E(codeStringsAndFuncs ...any) *State {
	return GlobalState.E(codeStringsAndFuncs...)
}

func (s *State) R(globalNumberNameMap map[string]int, code ...any) *State {
	freshState := &State{
		I:             -1,
		Code:          code,
		Vals:          s.Vals,
		Vars:          s.Vars, // since it's global, we reuse global vars
		LexicalParent: s,
		CallingParent: nil,
		Mu: s.Mu,
		GlobalNumbers: make([]float64, len(globalNumberNameMap)),
	}
	s = freshState
	for {
		if s == nil {
			break
		}
		s.I++
		if s.I >= len(s.Code) {
			s = s.BeDone()
			continue
		}
		v := s.Code[s.I]
		switch v := v.(type) {
		case func(s *State) *State:
			s = v(s)
		case func():
			v()
		case func() any:
			s.Push(v())
		case func(any) any:
			s.Push(v(s.Pop()))
		case func(any):
			v(s.Pop())
		case func(string) string:
			s.Push(v(toStringInternal(s.Pop())))
		case func(string, string) string:
			b := toStringInternal(s.Pop())
			a := toStringInternal(s.Pop())
			s.Push(v(a, b))
		default:
			s.Push(v)
		}
	}
	return freshState
}
// TODO: consolidate this with R
func ExecBlock(s *State, block []any) *State {
	origState := s
	for i := 0; i < len(block); i++ {
		v := block[i]
		switch v := v.(type) {
		case func(s *State) *State:
			s = v(s)
		case func():
			v()
		case func() any:
			s.Push(v())
		case func(any) any:
			s.Push(v(s.Pop()))
		case func(any):
			v(s.Pop())
		case func(string) string:
			s.Push(v(toStringInternal(s.Pop())))
		case func(string, string) string:
			b := toStringInternal(s.Pop())
			a := toStringInternal(s.Pop())
			s.Push(v(a, b))
		default:
			s.Push(v)
		}
		if s != origState {
		    return s
		}
	}
	return s
}

func AddBuiltin(name string, f any) {
	var normalized func(s *State) *State

	switch v := f.(type) {
	case func(s *State) *State:
        normalized = v
	case func():
	    normalized = func(s *State) *State {
			v()
	        return s
	    }
	case func() any:
	    normalized = func(s *State) *State {
			s.Push(v())
	        return s
	    }
	case func(any) any:
	    normalized = func(s *State) *State {
			s.Push(v(s.Pop()))
	        return s
	    }
	case func(any):
	    normalized = func(s *State) *State {
			v(s.Pop())
	        return s
	    }
	case func(string) string:
	    normalized = func(s *State) *State {
			s.Push(v(toStringInternal(s.Pop())))
	        return s
	    }
	case func(string, string) string:
	    normalized = func(s *State) *State {
			b := toStringInternal(s.Pop())
			a := toStringInternal(s.Pop())
			s.Push(v(a, b))
	        return s
	    }
	default:
		panic("unexpected type")
	}

	Builtins[name] = normalized
}

func (s *State) E(codeStringsAndFuncs ...any) *State {
	globalNumberNameMap := map[string]int{}
	
	tokens := []any{}
	
	for _, cf := range codeStringsAndFuncs {
		switch cf := cf.(type) {
		case string:
			// code = append(code, Parse(cf, globalNumberNameMap)...)
			tokens = append(tokens, Tokenize(cf)...)
		default:
			// code = append(code, cf)
			tokens = append(tokens, cf)
		}
	}
	code := Parse(tokens, globalNumberNameMap)
	return s.R(globalNumberNameMap, code...)
	
}

var internedStrings = map[string]string{}

func interned(v string) string {
	if s, ok := internedStrings[v]; ok {
		return s
	}
	internedStrings[v] = v
	return v
}

func Tokenize(codeString string) []any {
	tokens := []any{}
	lines := strings.Split(codeString, "\n")
	for lineI := 0; lineI < len(lines); lineI++ {
		line := lines[lineI]
		words := strings.Split(line, " ")
		parseState := "normal"
		stringStart := 0
	wordLoop:
		for i := 0; i < len(words); i++ {
			w := words[i]
			if parseState == "normal" {
				if strings.HasPrefix(w, "//") {
					break
				}
				if strings.HasPrefix(w, "#") {
					break
				}
				if w == "string:" {
					tokens = append(tokens, "."+strings.Join(words[i+1:], " "))
					break
				}
				if w == "beginString" {
					indent := strings.Join(words[0:i], " ") + " "
					startLineI := lineI
					search := indent + "endString"
					for lineI = lineI; lineI < len(lines); lineI++ {
						line := lines[lineI]
						index := strings.Index(line, search)
						if index != -1 {
							words = strings.Split(line[index+len(search):], " ")
							i = -1
							contents := lines[startLineI+1 : lineI]
							trimIndent := indent + "    "
							for contentsI, l := range contents {
								contents[contentsI] = strings.TrimPrefix(l, trimIndent)
							}
							tokens = append(tokens, "."+strings.Join(contents, "\n"))

							continue wordLoop
						}
					}
					panic("no corresponding endString found")
				}
				if strings.HasPrefix(w, `"`) {
					parseState = "string"
					stringStart = i
					if len(w) != 1 {
						i--
					}
					continue
				}
				tokens = append(tokens, w)
			} else if parseState == "string" {
				// if strings.HasSuffix(w, `"`)  && len(w) != 1 {
				if strings.HasSuffix(w, `"`) {
					parseState = "normal"
					// log.Println(ToJsonF(words[stringStart:i+1]))
					// log.Println("#skyblue", stringStart, i+1)
					theString := strings.Join(words[stringStart:i+1], " ")
					// log.Println("#skyblue", theString)
					tokens = append(tokens, "."+theString[1:len(theString)-1])
					continue
				}
			}
		}
		tokens = append(tokens, "\n")
	}
	return tokens
}

func Parse(tokens []any, globalNumberNameMap map[string]int) ([]any) {

	codeStack := [][]any{}
	code := []any{}
    
	for _, t := range tokens {
		switch t := t.(type) {
		case string:
			switch t {
			case "(":
				codeStack = append(codeStack, code)
				code = []any{}
			// case ")":
			// 	parentCode := codeStack[len(codeStack)-1]
			// 	// you could also put a type there and to the addingnon switch
			// 	parentCode = append(parentCode, B(code...))
			// 	// parentCode = append(parentCode, code, MakeBlock)
			// 	code = parentCode
			// 	codeStack = codeStack[0 : len(codeStack)-1]
			case ")":
				parentCode := codeStack[len(codeStack)-1]
				parentCode = append(parentCode, code)
				code = parentCode
				codeStack = codeStack[0 : len(codeStack)-1]
			default:
				if strings.TrimSpace(t) == "" {
					continue
				}
				if strings.HasPrefix(t, ".") {
					code = append(code, interned(t[1:]))
				} else if strings.HasPrefix(t, "$") {
					if v, ok := globalNumberNameMap[t]; ok {
						code = append(code, v)
					} else {
					    v := len(globalNumberNameMap)
					    globalNumberNameMap[t] = v
						code = append(code, v)
					}
				} else if f, ok := Builtins[t]; ok {
					code = append(code, f)
				} else if t == "true" {
					code = append(code, true)
				} else if t == "false" {
					code = append(code, false)
				} else if t == "nil" {
					code = append(code, nil)
					// } else if i, err := strconv.Atoi(t); err == nil {
					// 	code = append(code, i)
				} else if t == "valOnly" {
					code[len(code)-1] = Get // instead of A (access)
				} else if t == "varName" {
					code = code[0:len(code)-1]
				} else if f, err := strconv.ParseFloat(t, 64); err == nil {
					code = append(code, f)
				} else {
				    if v, ok := globalNumberNameMap["$" + t]; ok {
						code = append(code, v, GetGlobalNumber)
				    } else {
						code = append(code, interned(t), A)
				    }

					// This is actually slightly slower?
					// code = append(code, func(s *State) *State {
					// 	v := s.Get(interned(t))
					// 	s.Push(v)
					// 	if _, ok := v.(*Block); ok {
					// 		return Call(s)
					// 	}
					// 	return s
					// })
				}
			}
		// case func(s *State) *State:
		default:
			code = append(code, t)
		}
	}
	return code
}

var Builtins = map[string]func(*State) *State{
	"newScope": NewScope,
	"downScope": DownScope,
	"upScope": UpScope,
	"dropScope": DropScope,
	"scopeLen": func(s *State) *State {
	    s.Push(len(s.VarsStack))
	    return s
	},
	"as":       As,
	"to":       To,
	"at":       At,
	"setField": SetAt,
	"setIndex": SetAt,
	"slice":    Slice,
	"[":        BeginList,
	"]":        EndList,
	"{":        BeginRecord,
	"}":        EndRecord,
	"sleepMs":  SleepMS,
	"get":      Get,
	"get$":      GetGlobalNumber,
	"as$":      func(s *State) *State {
	    s.GlobalNumbers[s.Pop().(int)] = s.Pop().(float64)
	    return s
	},
	"enclose":  Enclose,
	"call":     Call,
	"do":       Do,
	"drop":     Drop,
	// "swap":     Swap,
	// "dup":     Dup,
	"is":     Is,
	"loopGo": func(s *State) *State {
	    origState := s
	    block := s.Pop().([]any)
	    times := int(s.Pop().(float64))
	    for i := 0; i < times; i++ {
	        // _ = block
	        // _ = origState
	        s.Push(float64(i))
	        s = ExecBlock(s, block)
	        if s != origState {
	            return s
	        }
	    }
	    return s
	},
	"+":            Plus,
	"-":            Minus,
	"*":            Times,
	"/":            Divide,
	">":            GreaterThan,
	"<":            LessThan,
	">=":            GreaterThanOrEqualTo,
	"<=":            LessThanOrEqualTo,
	"==":            EqualTo,
	"say":          Say,
	"log":          Log,
	"addrOfString": AddrOfString,
	"repeat":       Repeat,
	"exit":       Exit,
	"tag":        Tag,
	"breakTag":     BreakTag,
	"break":        Break,
	"guard":        Guard,
	"guardLte":     GuardLTE,
	"return":        Return,
	// "delay":        Delay,
	"nowMs":        NowMS,
	"now":          NowMS,
	"++":           CC,
	"push":         Push,
	"pushTo":       PushTo,
	// "if":       If,
	// "ifelse":       IfElse,
	"len":       Len,
	"__state":       func(s *State) *State {
	    s.Push(s)
	    return s
	},
	"split": func(s *State) *State {
	    v := s.Pop().(string)
	    str := s.Pop().(string)
	    theList := NewListFromStringSlice(strings.Split(str, v))
	    s.Push(theList)
	    return s
	},
}

func GetGlobalNumber(s *State) *State {
    s.Push(s.GlobalNumbers[s.Pop().(int)])
    return s
}

func Len(s *State) *State {
	list := s.Pop().(*List)
	s.Push(float64(list.Length()))
	return s
}
func NowMS(s *State) *State {
	s.Push(float64(time.Now().UnixMilli()))
	return s
}
func CC(s *State) *State {
	b := toStringInternal(s.Pop())
	a := toStringInternal(s.Pop())
	s.Push(a + b)
	return s
}
func Swap(s *State) *State {
	b := s.Pop()
	a := s.Pop()
	s.Push(b)
	s.Push(a)
	return s
}
func Dup(s *State) *State {
	a := s.Pop()
	s.Push(a)
	s.Push(a)
	return s
}
func Is(s *State) *State {
	b := s.Pop()
	a := s.Pop()
	s.Push(a == b)
	return s
}
func Tag(s *State) *State {
	s.Tag = s.Pop().(string)
	return s
}

func (s *State) ChildState(block any) *State {
    switch block := block.(type) {
    case *ScopedBlock:
        return &State{
		    Vals:          s.Vals,
		    Vars:          nil,
		    Code:          block.Code,
		    CallingParent: s,
		    LexicalParent: block.LexicalParent,
		    I:             -1,
            Mu: s.Mu,
            GlobalNumbers: s.GlobalNumbers,
        }
    case []any:
        panic("should not get here 1")
	    s.CodeStack = append(s.CodeStack, s.Code)
	    s.IStack = append(s.IStack, s.I)
	    s.Code = block
	    s.I = -1
	    return s
    }
    return s
}
func (s *State) AsyncChildState(block any) *State {
    switch block := block.(type) {
    case *ScopedBlock:
        return &State{
	        Vals:          NewList(),
	        Vars:          nil,
	        Code:          block.Code,
	        CallingParent: nil,
	        LexicalParent: block.LexicalParent,
	        AsyncParent:   s,
	        I:             -1,
            Mu: s.Mu,
            GlobalNumbers: s.GlobalNumbers,
        }
    case []any:
        panic("should not get here 2")
	    s.CodeStack = append(s.CodeStack, s.Code)
	    s.IStack = append(s.IStack, s.I)
	    s.Code = block
	    return s
    }
    return s
}

// func If(s *State) *State {
// 	block := s.Pop()
// 	cond := s.Pop().(bool)
// 	if cond {
// 	    return s.ChildState(block)
// 	}
// 	return s
// }

// func IfElse(s *State) *State {
// 	elseBlock := s.Pop()
// 	ifBlock := s.Pop()
// 	cond := s.Pop().(bool)
// 	if cond {
// 	    return s.ChildState(ifBlock)
// 	} else {
// 	    return s.ChildState(elseBlock)
// 	}
// 	return s
// }

func BreakTag(s *State) *State {
	tag := s.Pop().(string)
	for s != nil {
		if s.Tag == tag {
			break
		}
		s = s.LexicalParent
	}
	return s
}

func Break(s *State) *State {
	s.PopCodeStack()
	return s
}
func Guard(s *State) *State {
	cond := s.Pop().(bool)
	if !cond {
	    // return s.BeDone()
	    s.I = len(s.Code)
	}
	return s
}
func GuardLTE(s *State) *State {
	b := s.Pop().(float64)
	a := s.Pop().(float64)
	if a > b {
	    s.I = len(s.Code)
	}
	return s
}
func Return(s *State) *State {
    s.I = len(s.Code)-1
    return s
    // return s.CallingParent
}
func Repeat(s *State) *State {
	s.I = -1
	return s
}
func Exit(s *State) *State {
    os.Exit(0)
    return s
}


type Waiter struct {
    State *State
}
// func Delay(s *State) *State {
//     block := s.Pop()
//     ms := int(s.Pop().(float64))
//     w := &Waiter{
//         State: s,
//     }
//     s.Push(w)
//     go func () {
//         fmt.Println("sleeping", ms)
//         time.Sleep(time.Duration(ms) * time.Millisecond)
// 		s.AddCallback(&Callback{
// 			State: s.AsyncChildState(block),
// 		})
//     }()
//     return s
// }


func AddrOfString(s *State) *State {
	v := s.Pop().(string)
	dataPtr := (*reflect.StringHeader)(unsafe.Pointer(&v)).Data
	s.Push(uintptr(dataPtr))
	return s
}

// #### **Option B: []byte conversion (careful with empty strings!)**
// func AddrOfString(s *State) *State {
//     v := s.Pop().(string)
//     if len(v) == 0 {
//         s.Push(uintptr(0))
//     } else {
//         s.Push(uintptr(unsafe.Pointer(&([]byte(v))[0])))
//     }
//     return s
// }

func B(code ...any) func(state *State) *State {
	return func(state *State) *State {
		state.Push(&ScopedBlock{
			Code:          code,
			LexicalParent: state,
		})
		return state
	}
}
func Enclose(state *State) *State {
	state.Push(&ScopedBlock{
		Code:          state.Pop().([]any),
		LexicalParent: state,
	})
	return state
}

func C(code ...any) []any {
    return code
}

func Call(s *State) *State {
	b := s.Pop()
	return s.ChildState(b)
}
func Do(s *State) *State {
	switch code := s.Pop().(type) {
	case []any:
		s.CodeStack = append(s.CodeStack, s.Code)
		s.IStack = append(s.IStack, s.I)
		s.Code = code
		s.I = -1
	default:
        s.Push(code)
    }
	return s
}

// todo: you could always have the vars in the stack? might simplify
func NewScope(s *State) *State {
	s.VarsStack = append(s.VarsStack, s.Vars)
	s.Vars = NewRecord()
	s.VarsStack = append(s.VarsStack, s.Vars)
	return s
}
func DownScope(s *State) *State {
	s.Vars = s.VarsStack[len(s.VarsStack)-2]
	return s
}
func UpScope(s *State) *State {
	s.Vars = s.VarsStack[len(s.VarsStack)-1]
	return s
}
func DropScope(s *State) *State {
	s.Vars = s.VarsStack[len(s.VarsStack)-2]
	s.VarsStack = s.VarsStack[0:len(s.VarsStack)-2]
	return s
}

func Drop(s *State) *State {
	s.Pop()
	return s
}

func Push(s *State) *State {
    v := s.Pop()
    l := s.Pop().(*List)
    l.Push(v)
    return s
}
func PushTo(s *State) *State {
    l := s.Pop().(*List)
    v := s.Pop()
    l.Push(v)
    return s
}

func (s *State) Push(v any) {
	s.Vals.Push(v)
}
func (s *State) Pop() any {
	return s.Vals.Pop()
}

func Get(s *State) *State {
	s.Push(s.Get(s.Pop().(string)))
	return s
}

// A for Access
func A(s *State) *State {
	v := s.Get(s.Pop().(string))
	s.Push(v)
	switch v.(type) {
	case *ScopedBlock:
		return Call(s)
	case []any:
		return Do(s)
	}
	return s
}

func (state *State) Get(varName string) any {
	parent, v := state.findParentAndValue(varName)
	if parent == nil {
		panic(fmt.Sprintf("var not found: %q", varName))
	}
	return v
}

func (state *State) findParentAndValue(varName string) (*State, any) {
	scopesUp := 0
	for state != nil {
		v, ok := state.Vars.GetHas(varName)
		if ok {
			return state, v
		}
		state = state.LexicalParent
		scopesUp++
	}
	return nil, nil
}


func To(s *State) *State {
	varName := s.Pop().(string)
	v := s.Pop()
	s.Let(varName, v)
	return s
}

func (state *State) Let(varName string, v any) {
	parent, _ := state.findParentAndValue(varName)
	if parent == nil {
		panic("var not found " + varName)
	}
	parent.Vars.Set(varName, v)
}

func As(s *State) *State {
	varName := s.Pop().(string)
	v := s.Pop()
	s.Var(varName, v)
	return s
}

func (state *State) Var(varNameAny any, v any) {
	varName := varNameAny.(string)
	if state.Vars == nil {
	    state.Vars = NewRecord()
	}
	state.Vars.Set(varName, v)
}

// Slice is for 1 indexed slice
func Slice(state *State) *State {
	endInt := int(state.Pop().(float64))
	startInt := int(state.Pop().(float64))
	s := state.Pop()
	switch s := s.(type) {
	case *List:
		state.Push(s.Slice(startInt, endInt))
		return state
	case string:
		if len(s) == 0 {
			state.Push("")
			return state
		}
		if startInt < 0 {
			startInt = len(s) + startInt + 1
		}
		if startInt <= 0 {
			startInt = 1
		}
		if startInt > len(s) {
			state.Push("")
			return state
		}
		if endInt < 0 {
			endInt = len(s) + endInt + 1
		}
		if endInt <= 0 {
			state.Push("")
			return state
		}
		if endInt > len(s) {
			endInt = len(s)
		}
		if startInt > endInt {
			state.Push("")
			return state
		}
		state.Push(s[startInt-1 : endInt])
		return state
	}
	state.Push(nil)
	return state
}

func At(s *State) *State {
	key := s.Pop()
	v := s.Pop()
	switch v := v.(type) {
	case *List:
		s.Push(v.Get(int(key.(float64))))
		return s
	case map[string]any:
		s.Push(v[key.(string)])
		return s
	case string:
		indexInt := int(key.(float64))
		if indexInt == 0 {
			s.Push("")
			return s
		}
		if indexInt < 1 {
			s.Push(string(v[len(v)+indexInt]))
			return s
		}
		if indexInt-1 >= len(v) {
			s.Push("")
			return s
		}
		s.Push(string(v[indexInt-1]))
		return s
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Pointer {
			rv = rv.Elem()
		}
		if rv.Kind() == reflect.Struct {
			fieldName := key.(string)
			fv := rv.FieldByName(fieldName)
			if fv.IsValid() {
				s.Push(fv.Interface())
			}
		}
		return s
	}
}

func SetAt(s *State) *State {
	o := s.Pop()
	k := s.Pop()
	v := s.Pop()

	switch o := o.(type) {
	case *List:
		o.Set(int(k.(float64)), v)
	case *Record:
		o.Set(k.(string), v)
	}
	return s
}

func Plus(s *State) *State {
	b := s.Pop().(float64)
	a := s.Pop().(float64)
	s.Push(a + b)
	return s
}

func Minus(s *State) *State {
	b := s.Pop().(float64)
	a := s.Pop().(float64)
	s.Push(a - b)
	return s
}

func Times(s *State) *State {
	b := s.Pop().(float64)
	a := s.Pop().(float64)
	s.Push(a * b)
	return s
}

func Divide(s *State) *State {
	b := s.Pop().(float64)
	a := s.Pop().(float64)
	s.Push(a / b)
	return s
}
func GreaterThan(s *State) *State {
	b := s.Pop().(float64)
	a := s.Pop().(float64)
	s.Push(a > b)
	return s
}

func LessThan(s *State) *State {
	b := s.Pop().(float64)
	a := s.Pop().(float64)
	s.Push(a < b)
	return s
}

func GreaterThanOrEqualTo(s *State) *State {
	b := s.Pop().(float64)
	a := s.Pop().(float64)
	s.Push(a >= b)
	return s
}

func LessThanOrEqualTo(s *State) *State {
	b := s.Pop().(float64)
	a := s.Pop().(float64)
	s.Push(a <= b)
	return s
}

func EqualTo(s *State) *State {
	b := s.Pop()
	a := s.Pop()
	s.Push(a == b)
	return s
}


func Say(s *State) *State {
	v := s.Pop()
	fmt.Printf("%s\n", toStringInternal(v))
	return s
}
func Log(s *State) *State {
	v := s.Pop()
	log.Printf("%s\n", toStringInternal(v))
	return s
}

func BeginList(s *State) *State {
	newState := &State{
		Vals:          NewList(),
		Vars:          nil,
		Code:          s.Code,
		CallingParent: s,
		LexicalParent: s,
		I:             s.I,
 		Mu: s.Mu,
        GlobalNumbers: s.GlobalNumbers,
	}
	return newState
}
func EndList(s *State) *State {
	parentState := s.CallingParent
	parentState.Push(s.Vals)
	parentState.I = s.I
	return parentState
}
func BeginRecord(s *State) *State {
	newState := &State{
		Vals:          NewList(),
		Vars:          NewRecord(),
		Code:          s.Code,
		CallingParent: s,
		LexicalParent: s,
		I:             s.I,
		Mu: s.Mu,
        GlobalNumbers: s.GlobalNumbers,
	}
	return newState
}

func EndRecord(s *State) *State {
	parentState := s.CallingParent
	m := map[string]any{}
	vals := s.Vals
	for i := 0; i < vals.Length(); i += 2 {
		m[vals.Get(i+1).(string)] = vals.Get(i + 2)
	}
	parentState.Push(m)
	parentState.I = s.I
	return parentState
}

func SleepMS(s *State) *State {
	ms := int(s.Pop().(float64))
    s.Mu.Unlock()
	time.Sleep(time.Duration(ms) * time.Millisecond)
    s.Mu.Lock()
	return s
}

// structs (arrays with functions to access)
// parse gets builtins at compile time (good)
// faster maps with string address hashing?
// break level
// Pause
// resume
// go
// wait
// waitloop,
// waitSignal, signal, broadcast
// semaphore, acquire, release
// wait
// waitn
// cancel all


/*
.displayHistoryLite (
    .commandHistory as
    [ ] .ret as
    commandHistory (
        
    ) each
) var

def displayHistoryLite commandHistory
    var ret []
    commandHistory each command
        if command at name, is "llmResponse"
            ret push "llmResponse... " ++ command at content, slice 1 20
            continue
        end
        ret push command at raw
        if command at body, isnt null
            command at body
            pushTo ret
            ret push "END PATCH"
        end
        if command at error
            ret push "    Error: " ++ command at error
        end
        ret push ""
    end
    ret join newline
end

*/