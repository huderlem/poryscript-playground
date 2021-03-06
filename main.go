package main

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/huderlem/poryscript/types"

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

	genElement := document.Call("getElementById", "gen-picker")
	if genElement == js.Undefined() {
		fmt.Println("ERROR: couldn't get gen picker input element")
		return nil
	}
	gen := types.GEN3
	if genElement.Get("value").String() == "2" {
		gen = types.GEN2
	}

	compileTimeSwitchesElement := document.Call("getElementById", "switches-text")
	if compileTimeSwitchesElement == js.Undefined() {
		fmt.Println("ERROR: couldn't get compile-time switches input element")
		return nil
	}
	compileSwitches := make(map[string]string)
	for _, option := range strings.Split(compileTimeSwitchesElement.Get("value").String(), " ") {
		parts := strings.SplitN(option, "=", 2)
		if len(parts) != 2 {
			setErrorText(fmt.Sprintf("Error: invalid compile-time switch %s\n", option))
			return nil
		}
		compileSwitches[parts[0]] = parts[1]
	}

	parser := parser.New(lexer.New(inputScript, gen), gen, "", compileSwitches)
	program, err := parser.ParseProgram()
	if err != nil {
		setErrorText(err.Error())
		return nil
	}

	optimize := optimizeElement.Get("checked").Bool()
	emitter := emitter.New(program, gen, optimize)
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
