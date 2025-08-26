#!/usr/local/bin/linescript4
"closeloop"
string
    SEE FILE BLOCK: stucco.go @@ func Parse(codeString string)
    
    update the default case so it handles numbers
    numbers with no decimals are ints
    numbers with decimals are float64s
    handle "true" and "false" and "nil"
    
end
drop string
end
execBashStdinStream
say