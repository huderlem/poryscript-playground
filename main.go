package main

import (
	"fmt"
	"syscall/js"

	"github.com/huderlem/poryscript/emitter"
	"github.com/huderlem/poryscript/lexer"
	"github.com/huderlem/poryscript/parser"
)

func compile(this js.Value, inputs []js.Value) interface{} {
	document := js.Global().Get("document")
	if document == js.Undefined() {
		fmt.Println("ERROR: couldn't get document handle")
		return nil
	}
	inputElement := document.Call("getElementById", "inputtext")
	if inputElement == js.Undefined() {
		fmt.Println("ERROR: couldn't get input element")
		return nil
	}
	outputElement := document.Call("getElementById", "outputtext")
	if outputElement == js.Undefined() {
		fmt.Println("ERROR: couldn't get output element")
		return nil
	}
	optimizeElement := document.Call("getElementById", "optimize-checkbox")
	if optimizeElement == js.Undefined() {
		fmt.Println("ERROR: couldn't get optimize element")
		return nil
	}

	scriptText := inputElement.Get("value").String()
	parser := parser.New(lexer.New(scriptText))
	program, err := parser.ParseProgram()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return nil
	}

	optimize := optimizeElement.Get("checked").Bool()
	emitter := emitter.New(program, optimize)
	result, err := emitter.Emit()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return nil
	}

	outputElement.Set("value", result)
	return nil
}

func registerFunctions() {
	js.Global().Set("compile", js.FuncOf(compile))
}

func main() {
	c := make(chan struct{}, 0)
	registerFunctions()
	<-c
}
