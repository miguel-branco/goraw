package main

import (
	"fmt"
    "os"
    . "interpreter"
	. "interpreter/expressions"
    . "interpreter/operators"
    . "interpreter/plugins"
)

// FIXME: Plugin interface shouldn't require a top-level collection value but simply a value.
//        If it happens to be a collection, ok, otherwise, can be a e.g. record value as is the case with JSON.

func main() {
    fmt.Println("goraw interpreter - version 0.0")

    //
/*
    rc := NewRecordConstruction([]AttributeConstruction{AttributeConstruction{"a", NewBoolConst(true)}, AttributeConstruction{"b", NewIntConst(2)}, AttributeConstruction{"c", NewIntConst(5)}})
    rp1 := NewRecordProjection(rc, "b")
    rp2 := NewRecordProjection(rc, "c")

    bop := NewBinaryOperation(Mult, rp1, rp2)
    fmt.Println(bop.Execute(nil))

    b := NewBoolConst(true)
    n := NewNot(b)
    fmt.Println(n.Execute(nil))
*/
    // CSV Scan + Select

    fmt.Println("CSV")

    csv := NewCsvPlugin(os.Args[1])
    csv.Init()

    csvScan := make(chan []Value)
    go Scan(csv, csvScan)

    csvSelect := make(chan []Value)
    go Select(csvScan, NewBinaryOperation(Bool, Eq, NewRecordProjection(String, NewArgument(Record, 0), "_2"), NewStringConst("33")), csvSelect)

    Print(csvSelect)

    csv.Fini()

    // JSON Scan + Unnest

    fmt.Println("JSON")

    json := NewJsonPlugin(os.Args[2])
    json.Init()

    jsonScan := make(chan []Value)
    go Scan(json, jsonScan)

    jsonUnnest := make(chan []Value)
    go Unnest(jsonScan, NewPath(0, []string{"phoneNumbers"}), NewBinaryOperation(Bool, Eq, NewRecordProjection(String, NewArgument(Record, 1), "type"), NewStringConst("fax")), jsonUnnest)

    Print(jsonUnnest)

    json.Fini()

    // CSV Scan + Join

    fmt.Println("CSV")

    csv1 := NewCsvPlugin(os.Args[1])
    csv2 := NewCsvPlugin(os.Args[1])
    
    csvScan1 := make(chan []Value)
    go Scan(csv1, csvScan1)

    csvScan2 := make(chan []Value)
    go Scan(csv2, csvScan2)

    csvJoin := make(chan []Value)
    go Join(csvScan1, csvScan2, NewBinaryOperation(Bool, Eq, NewRecordProjection(String, NewArgument(Record, 0), "_2"),  NewRecordProjection(String, NewArgument(Record, 1), "_2")), csvJoin)

    Print(csvJoin)

    csv1.Fini()
    csv2.Fini()

    // FIXME: Should path use NewArgument and NewRecordProjection that get executed? That drops the need for BindPath

}
