# Poryscript Playground

See it in action here: http://www.huderlem.com/poryscript-playground/

This is an online playground for [Poryscript](https://github.com/huderlem/poryscript), a higher-level scripting language for Gen 3 pokemon decompilation projects.

# Getting Started with Development

These instructions will get you setup and working with the Poryscript Playground locally.

## Building and Running Locally

First, install [Go](http://golang.org).  Poryscript Playground, uses Go's experiemental WebAssembly support to interface with the JavaScript on the webpage.

Clone and navigate to the Poryscript working directory. To build, you must specify that WebAssembly is the target.
```
GOARCH=wasm GOOS=js go build -o main.wasm main.go
```
If you're on Windows, you might be better off running them in separate commands:
```
set GOARCH=wasm
set GOOS=js
go build -o main.wasm main.go
```

This will create a `main.wasm` file which contains the Poryscript compiler logic.

Next, run a local web server in the current directory. Any web server will work.  I use `goexec`:
```
# install goexec: go get -u github.com/shurcooL/goexec
goexec "http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))"
```

Finally, visit `http://localhost:8080/`, and you should see the page up and running. The `main.wasm` payload is automatically loaded when the page is loaded. Look in the browser's debugging console to see if there are any errors reported.
