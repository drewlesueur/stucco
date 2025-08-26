package main

import (
	"fmt"
	. "github.com/drewlesueur/stucco"
	"time"
)
func main() {
    
    var totF = 0.0
    var idxF = 0.0
    var start time.Time
    
    
	fmt.Println("\n=====fast Go loop float")
	totF = 0
	start = time.Now()
	for idxF = 1; idxF <= 1000000; idxF++ {
		totF += idxF
	}
	fmt.Println(totF)
	fmt.Println(time.Since(start))
	
	E(`
        "" say
        "======= raw loop" say
        .a 0 var
        .start nowMs var
        .i 0 var
        (
            i 1 + .i as
            i 1000000 <= guard
            i a + .a to
            repeat
        ) do
        "a is " a ++ say
        "it took " nowMs start - ++ log
	    
	`)
	
	E(`
        "" say
        "======= raw loop global" say
        .a 0 var
        .start nowMs var
        0 $i as$
        (
            i 1 + $i as$
            i 1000000 <= guard
            i a + .a to
            repeat
        ) do
        "a is " a ++ say
        "it took " nowMs start - ++ log
	    
	`)
	E(`
        "" say
        "======= raw loop global2" say
        0 $a as$
        .start nowMs var
        0 $i as$
        (
            i 1 + $i as$
            i 1000000 <= guard
            i a + $a as$
            repeat
        ) do
        "a is " a ++ say
        "it took " nowMs start - ++ log
	`)

	E(`
        "" say
        "======= drop" say
        .a 0 var
        .start nowMs var
        1000000 (
            drop
        ) loop
        "a is " a ++ say
        "it took " nowMs start - ++ log
	`)
	E(`
        "" say
        "======= suming loop" say
        .a 0 var
        .start nowMs var
        1000000 (
            a + .a to
        ) loop
        "a is " a ++ say
        "it took " nowMs start - ++ log
	`)
	E(`
        "" say
        "======= suming loop Go" say
        .a 0 var
        .start nowMs var
        1000000 (
            a + .a to
        ) loopGo
        "a is " a ++ say
        "it took " nowMs start - ++ log
	`)
	E(`
        "" say
        "======= suming loop Go faster" say
        0 $a as$
        .start nowMs var
        1000000 (
            a + $a as$
        ) loopGo
        "a is " a ++ say
        "it took " nowMs start - ++ log
	    
	`)
	E(`
        "" say
        "======= summing loop actual go" say
        .a 0 var
        .start nowMs var
        `, func(s *State) *State {
            v := 0
            for i := 1; i <= 1000000; i++ {
                v += i
            }
            s.Vars.Set("a", float64(v))
            return s
        }, `
        "a is " a ++ say
        "it took " nowMs start - ++ log
	`)
}

