package stucco

import (
	"fmt"
	// "testing"
	// "time"
)

func ExampleE_say() {
	E(`
	    "hello" say
	`)
	// Output:
	// hello
}

func ExampleE_block() {
	E(`
	    ( "hello world 1" say ) do
	`)
	// Output:
	// hello world 1
}

func ExampleE_if() {
	E(`
	    true ( "hello world 2" say ) if
	    false ( "goodbye world" say ) if
	`)
	// Output:
	// hello world 2
}

func ExampleE_varName() {
	E(`
	    YoDude varName say
	`)
	// Output:
	// YoDude
}
func ExampleE_valOnly() {
	E(`
	    .stuff ( "hello world 3" say ) var
	    stuff
	    stuff valOnly drop
	    stuff

	`)
	// Output:
	// hello world 3
	// hello world 3
}
func ExampleE_state() {
	E(`
	    __state .I at say
	`)
	// Output:
	// 2
}

func ExampleE_loop() {
	E(`
	    10 ( say ) loop
	`)
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
	// 10
}

func ExampleE_nested_loop() {
	E(`
	    (
            .name "Drew" var
            4 (
                name swap ++ say
                3 (
                    "..." swap ++ name ++ say
                ) loop
            ) loop
        ) enclose call
	`)
	// Output:
	// Drew1
	// ...1Drew
	// ...2Drew
	// ...3Drew
	// Drew2
	// ...1Drew
	// ...2Drew
	// ...3Drew
	// Drew3
	// ...1Drew
	// ...2Drew
	// ...3Drew
	// Drew4
	// ...1Drew
	// ...2Drew
	// ...3Drew
}

func ExampleE_sleep() {
	E(`
	    nowMs
	    150 sleepMs
	    "slept" say
	    nowMs swap - .diff as
        diff 5 / round 5 * say
        ` /*, func(s *State) *State {s.Push(s.Vals); return s},` say*/ + `
	`)
	// Output:
	// slept
	// 150
}

// TODO: nested Go funcs
func ExampleE_and() {
	// <-E(`
	E(`
	    true 2 and say
	    true ( 3 ) and say
	    false 100 and say
	    false ( 300 ) and say
	`)
	// Output:
	// 2
	// 3
	// false
	// false
}
func ExampleE_string() {
	E(`
	    " " "x" ++ say
	`)
	// Output:
	//  x
}
func ExampleE_loopWithGoFunc() {
	E(`
	    10 ( `, func(v any) {fmt.Println(v)} ,` ) loop
	`)
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
	// 10
}
func ExampleE_split() {
	E(`
	    "a b c d" " " split say
	`)
	// Output:
    // [
    //     "a",
    //     "b",
    //     "c",
    //     "d"
    // ]
}

func ExampleE_newline() {
	E(`
	    "a" newline ++ "b" ++
	    say
	`)
	// Output:
	// a
	// b
}
func ExampleE_each() {
	E(`
	    [ .a .b .c ] (
	        ": " swap ++ ++
	        say
	    ) each
	`)
	// Output:
    // 1: a
    // 2: b
    // 3: c
}

/*
    (
        (
           300 return
        ) do
        "never get here" say
    ) call

// early return goes to function block
// hmm dynamic function, runs in same scope?
// no separate return flow like that?

.if (
    swap
    guard
    do
) dynamicDef

loop var (
    block var as
    count var as
    i var 1 let
    (
        i count lte guard
        i block do
        repeat
    ) do
) dynamicDef

loop var (
    block var as
    count var as
    i var 1 let
    (
        i count lte guard
        i block do
        repeat
    ) do
) dynamicDef


# no work with early return
loop var (
    // count
    // block
    // i
    0
    (
        1 +
        1 pick 3 pick lte guard
        1 pick 2 pick do
        repeat
    ) do
) dynamicDef



true (
    "yay true" say
) if
"done" say

*/

// func ExampleSay() {
// 	<-R(
// 		"hello world", Say,
// 	)
// 	// Output:
// 	// hello world
// }
// 
// func ExamplePlus() {
// 	<-R(
// 		1.0, 2.0, Plus, Say,
// 	)
// 	// Output:
// 	// 3
// }
// 
// func ExamplePlus_newline() {
// 	<-R(
// 		1.0, 2.0, Plus,
// 		Say,
// 	)
// 	// Output:
// 	// 3
// }
// 
// func ExampleVar() {
// 	<-R(
// 		"name", "Drew", Var,
// 		"name", Get, Say,
// 	)
// 	// Output:
// 	// Drew
// }
// func ExampleB() {
// 	<-R(
// 		"sayHi", B("Why Hello", Say), Var,
// 		"sayHi", Get, Call,
// 	)
// 	// Output:
// 	// Why Hello
// }
// 
// func ExampleB_closures() {
// 	<-R(
// 		"increr", B(
// 			"x", 0.0, Var,
// 			B(
// 				"x", Get,
// 				1.0, Plus,
// 				"x", To,
// 				"x", Get,
// 			),
// 		), Var,
// 		"incr", "increr", Get, Call, Var,
// 		"incr", Get, Call, Say,
// 		"incr", Get, Call, Say,
// 		"incr", Get, Call, Say,
// 	)
// 	// Output:
// 	// 1
// 	// 2
// 	// 3
// }
// 
// /*
// 
// def a b c d (
// 
// )
// 
// indent for blocks
// 
// 
// 12 loop i [ say i ]
// 
// 
// 
// x (( 3 + 4 ))
// 
// 
// 
// 
// 3 (
//     "yay" say
// ) loop
// 
// .if (
//     .block as
//     .cond as
//     cond guard
//     block
// ) def
// 
// .if (
//     swap
//     guard
//     call
// ) def
// 
// def if do
//     swap guard call
// end
// 
// 
// def loop do
//     as block
//     as count
//     i = 0
//     do
//         i = i + 1
//         i <= count  guard
//         i block
//         repeat
//     end
// end
// 
// 
// .loop (
//     .block as
//     .count as
//     .i 0 var
//     (
//         i 1 + .i =
//         i count <= guard
//         i block
//         repeat
//     ) run
// ) def
// 
// 
// .loop do
//     .block as
//     .count as
//     .i 0 var
//     do
//         i 1 + .i =
//         i count <= guard
//         i block
//         repeat
//     end run
// end def
// 
// loop var do
//     block var as
//     count var as
//     i var 0 var
//     do
//         i 1 + i var =
//         i count <= guard
//         i block
//         repeat
//     end run
// end def
// 
// 
// .loop
//     .block as
//     .count as
//     .i 0 var
//         i 1 + .i =
//         i count <= guard
//         i block
//         repeat
//     run
// def
// 
// increr var do
//     x var 0 def
//     do
//         x 1 + x var =
//         x
//     end
// end def
// 
// increr var do
//     x var 0 def
//     do
//         x 1 + x var as
//         x
//     end
// end def
// 
// 
// increr do
//     x 0 def
//     do
//         x val 1 + x as
//         x val
//     end
// end
// 
// x 3 is do
//     "x is 3" say
// end if
// 
// .incerer
//     x 0 =
//         x 1 + .x =
//         x
// def
// .incr incerer =
// incr say
// incr say
// 
// 
// .incerer
//     x 0 =
//         x 1 + .x =
//         x
// def
// 
// .incerer ( .x 0 = ( x 1 + .x = x ) ) def
// 
// incerer var ( x var 0 = ( x 1 + x var = x ) ) def
// 
// def increr
//     var x 0
//     func
//         let x x + 1
//     end
// end
// 
// 
// switchNone
//     case x is 3
//     case 
// end
// 
// 
// 
// 
// 
// incerer ( x 0 = ( x$ 1 + x = x$ ) ) def$
// 
// 
// digits = theLines len  toString  len
// 
// .loop do def
// 
// end
// 
// 
// def llmCall provider model prompt
//     switch provider
//     case "ollama"
//         ollamaCall model prompt
//     case "chatgpt"
//         chatGptCall model prompt
//     case "anthropic"
//         anthropicCall model prompt
//     default
//         say "invalid provider" provider
//         exit
//     end
// end
// 
// 
// def llmCall
// 
// 
// end
// 
// 
// llmCall var do
//     provider do
//         "ollama" do
// 
//         case "chatgpt" do
// 
//         end case
// 
//     end switch
// end def
// 
// 
// 
// llmCall var do def
//     provider do switch
//         "ollama" do case
// 
//         "chatgpt" do case
// 
//         "anthropic" do case
//         end
//     end
// end
// 
// 
// 
// switch x
//     case 100
//     
//     case 200
//     
//     case 300
//     
// end
// 
// def case val block
//     var callingParent getCallingParent
//     var cf callingParent lookup caseFunc
//     val cf
//     if
//         block
//         return
//     end
// end
// 
// 
// 
// 
// 3 x is do if
// end
// 
// 
// if x is 3
// 
// else if x is 4
// 
// end
// 
// 
// 
// switch var do
//     block var as
//     val var as
//     caseFunc var ( val is ) def
//     block funcVal scope var at
//     scope caseFunc var caseFunc setAt
//     block
// end def
// 
// case var do
//     block var as
//     caseFunc do
//         block
//         block getFuncVal getLexParent getLexParent setState
//     end if
// end def
// 
// l
// lll
// 
// 
// 
// def if cond block
//     cond and
// end
// 
// 
// 
// 
// */
// 
// 
// func ExampleA_closures2() {
// 	<-R(
// 		"increr", B(
// 			"x", 0.0, Var,
// 			B(
// 				"x", A,
// 				1.0, Plus,
// 				"x", To,
// 				"x", A,
// 			),
// 		), Var,
// 		"incr", "increr", A, Var,
// 		"incr", A, Say,
// 		"incr", A, Say,
// 		"incr", A, Say,
// 	)
// 	// Output:
// 	// 1
// 	// 2
// 	// 3
// }
// 
// func Increr(s *State) *State {
// 	return B(
// 		"x", 0.0, Var,
// 		B(
// 			"x", A,
// 			1.0, Plus,
// 			"x", To,
// 			"x", A,
// 		),
// 	)(s)
// }
// 
// func ExampleA_closures3() {
// 	<-R(
// 		"incr", Increr, Call, Var,
// 		"incr", A, Say,
// 		"incr", A, Say,
// 		"incr", A, Say,
// 	)
// 	// Output:
// 	// 1
// 	// 2
// 	// 3
// }
// 
// func ExampleA_closures4() {
// 	<-E(`
// 	    .increr (
// 	        .x 0 var
// 	        (
// 	            x 1 +
// 	            .x to
// 	            x
// 	        )
// 	    ) var
// 	    .incr increr var
// 	    incr say
// 	    incr say
// 	    incr say
// 	`)
// 	// Output:
// 	// 1
// 	// 2
// 	// 3
// }
// 
// func ExampleB_access() {
// 	<-R(
// 		"x", 300, Var,
// 		"x", A, Say,
// 		100, "x", As,
// 		"x", A, Say,
// 	)
// 	// Output:
// 	// 300
// 	// 100
// }
// func ExampleB_closure() {
// 	<-R(
// 		"x", 300, Var,
// 		B(
// 			"x", A, Say,
// 			"x", 150, Let,
// 			"x", A, Say,
// 		),
// 		Call,
// 		"x", A, Say,
// 	)
// 	// Output:
// 	// 300
// 	// 150
// 	// 150
// }
// 
// func ExampleA_closures5() {
// 	<-E(`
// 	    .increr (
// 	        .x 0 var
// 	        (
// 	            x 1 +
// 	            .x to
// 	            x
// 	        )
// 	    ) var
// 	    .incr increr var
// 	    incr say
// 	    incr say
// 	    incr say
// 	`)
// 	// Output:
// 	// 1
// 	// 2
// 	// 3
// }
// func ExampleE_slice() {
// 	<-E(`
// 	    [ .apple .banana .orange ]
// 	    say
// 	`)
// 	// Output:
// 	// [
// 	//     "apple",
// 	//     "banana",
// 	//     "orange"
// 	// ]
// }
// 
// func ExampleE_map() {
// 	<-E(`
// 	    { .fruit .apple .level 7 }
// 	    say
// 	`)
// 	// Output:
// 	// {
// 	//     "fruit": "apple",
// 	//     "level": 7
// 	// }
// }
// 
// func ExampleE_sleepMS() {
// 	<-E(`
// 	    .Howdy say
// 	    10 sleepMs
// 	    .Duty say
// 	`)
// 	// Output:
// 	// Howdy
// 	// Duty
// }
// 
// func ExampleE_data() {
// 	<-E(`
// 	    [ .apple .pear .banana ] 1 at
// 	    say
// 	`)
// 	// Output:
// 	// apple
// }
// func ExampleE_data2() {
// 	<-E(`
// 	    [ .apple .pear .banana .beans ] 1 2 slice
// 	    say
// 	`)
// 	// Output:
// 	// [
// 	//     "apple",
// 	//     "pear"
// 	// ]
// }
// 
// func ExampleE_data4() {
// 	a := 100
// 	<-E(func() {
// 		a = 30
// 	})
// 	fmt.Println(a)
// 
// 	// Output:
// 	// 30
// }
// 
// func ExampleE_data5() {
// 	v := <-E(`
// 	    20
// 	    25
// 	`)
// 	fmt.Println(v)
// 
// 	// Output:
// 	// 25
// }
// func ExampleE_string() {
// 	<-E(`
// 	    "hello world :)"
// 	    say
// 	`)
// 
// 	// Output:
// 	// hello world :)
// }
// 
// func ExampleE_stringLine() {
// 	<-E(`
//         string: all the way to the end :)
//         say
// 	`)
// 
// 	// Output:
// 	// all the way to the end :)
// }
// 
// func ExampleE_stringMultiLine() {
// 	<-E(`
//         beginString
//             This is
//             a multiline
//             string
//         endString
//         say
// 	`)
// 
// 	// Output:
// 	// This is
// 	// a multiline
// 	// string
// }
// func ExampleE_stringMultiLine2() {
// 	<-E(`
//         beginString
//             This is
//             a multiline
//             string
//         endString say
// 	`)
// 
// 	// Output:
// 	// This is
// 	// a multiline
// 	// string
// }
// func ExampleE_stringMultiLine3() {
// 	<-E(`
//         beginString
//             This is
//             a multiline
//             string
//         endString 1 4 slice say
// 	`)
// 
// 	// Output:
// 	// This
// }
// func ExampleE_quote_no_space() {
// 	<-E(`
// 	    "oneWord" say
// 	`)
// 
// 	// Output:
// 	// oneWord
// }
// 
// func ExampleE_loop() {
// 	<-E(`
// 	    5 (
// 	        say
// 	        "hello" say
// 	    ) loop
// 	`)
// 
// 	// Output:
// 	// 1
// 	// hello
// 	// 2
// 	// hello
// 	// 3
// 	// hello
// 	// 4
// 	// hello
// 	// 5
// 	// hello
// }
// 
// // runs faster when ran with
// // go test -run TestLoop
// func TestLoop(t *testing.T) {
// 	_ = time.Now
// 	a := (<-E(`
// 	    .start nowMs var
// 	    .a 0 var
// 	    10 (
// 	        a + .a to
// 	    ) loop
// 	    "Elapsed: " nowMs start - ++ say
// 	    a say
// 	    a
// 	`)).(float64)
// 
// 	if a != 55 {
// 		t.Fatalf("bad loop count: %v", a)
// 	}
// 
// }
// 
// // func ExampleE_addrOf() {
// // 	<-E(`
// // 	    .yo addrOfString say
// // 	    .yo2 addrOfString say
// // 	    .yo addrOfString say
// // 	    .yo2 addrOfString say
// // 	    "what??!" addrOfString say
// // 	    "what??!" addrOfString say
// // 	`)
// //
// // 	// Output:
// // }
// 
// // func ExampleE_data6() {
// //     E(`
// //         .events [ ] var
// //         .waiters makeList var
// //
// //         .publish (
// //             .event as
// //             events event push
// //             .waiter waiters 1 sub var
// //             .waiters waiters 2 -1 subr var
// //             waiter resume
// //         ) var
// //         .consume (
// //
// //         ) var
// //     `)
// //
// // 	// Output:
// // 	// 25
// // }
// 
// // func ExampleE_data5() {
// //     events := []string{}
// //
// //     publish := func(event string) {
// //         events = append(events, event)
// //
// //     }
// //
// //     consume := func() []string {
// //         E(func() {
// //
// //         }
// //
// //         )
// //         if len(events) > 0 {
// //             ret := events
// //             events = []string{}
// //             return ret
// //         }
// //
// //
// //     }
// //
// //     go func(){
// //         E(func(s *State) *State {
// //             time.Sleep(1 * time.Second)
// //             return nil
// //         }, func() { publish("hello")
// //             publish("world")
// //         })
// //     }()
// //
// //     go func(){
// //         E(func() {
// //             consume()
// //         })
// //     }()
// //
// // 	// Output:
// // 	// 30
// // }


func ExampleE_tabs() {
	E(`
    "wow" say
	`)
	// Output:
	// wow
}

func ExampleE_tabs1() {
	E(`
    .wow say
	`)
	// Output:
	// wow
}


func ExampleE_startsWith() {
	E(`
    "abcd" "abc" startsWith say
    "abcd" "abd" startsWith say
	`)
	// Output:
	// true
	// false
}

func ExampleE_incr() {
	E(`
        .increr (
            .x 0 var
            (
                x 1 + .x to
                x
            ) enclose
        ) enclose var

        increr .incr as
        incr say
        incr say
        incr say
        incr say
	`)
	// Output:
	// 1
	// 2
	// 3
	// 4
}

func ExampleE_incr2() {
	E(`
        .increr (
            .x 0 var
            (
                x 1 + .x to
                x
            ) enclose
        ) enclose var

        increr .incrA as
        increr .incrB as
        incrA say
        incrB say
        incrA say
        incrB say
	`)
	// Output:
	// 1
	// 1
	// 2
	// 2
}


func ExampleE_ifElse() {
	E(`

        true ( "it's true" say ) ( "it's false" say ) ifElse
        false ( "it's true" say ) ( "it's false" say ) ifElse
        __vals len say
	`)
	// Output:
	// it's true
	// it's false
    // 0
}


func ExampleE_pick() {
	E(`

        1 2 3
        3 pick
        // 2 3 1
        say say say
	`)
	// Output:
	// 1
	// 3
	// 2
}

