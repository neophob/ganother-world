package main

import "syscall/js"

func addTitle() {
	document := js.Global().Get("document")
	h1 := document.Call("createElement", "h1")
	h1.Set("innerHTML", "Another World")
	document.Get("body").Call("appendChild", h1)
}

func main() {
	println("Hello Another-World")

	addTitle()

	channel := make(chan bool)
	<-channel
}
