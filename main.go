package main

import (
	"fmt"
	"syscall/js"

	"github.com/huderlem/poryscript/emitter"
	"github.com/huderlem/poryscript/lexer"
	"github.com/huderlem/poryscript/parser"
)

func setErrorText(err string) {
	js.Global().Get("document").Call("getElementById", "error-text").Set("innerHTML", err)
}

func compile(this js.Value, inputs []js.Value) interface{} {
	document := js.Global().Get("document")
	window := js.Global().Get("window")
	if window == js.Undefined() {
		fmt.Println("ERROR: couldn't get window handle")
		return nil
	}
	inputScript := window.Get("inputEditor").Call("getValue").String()
	optimizeElement := document.Call("getElementById", "optimize-checkbox")
	if optimizeElement == js.Undefined() {
		fmt.Println("ERROR: couldn't get optimize element")
		return nil
	}

	parser := parser.New(lexer.New(inputScript), "")
	program, err := parser.ParseProgram()
	if err != nil {
		setErrorText(err.Error())
		return nil
	}

	optimize := optimizeElement.Get("checked").Bool()
	emitter := emitter.New(program, optimize)
	resultScript, err := emitter.Emit()
	if err != nil {
		setErrorText(err.Error())
		return nil
	}

	setErrorText("")
	window.Get("outputEditor").Call("setValue", resultScript)
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
