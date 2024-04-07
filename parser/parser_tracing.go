package parser

import (
	"fmt"
	"strings"
)

var traceLevel int = 0

const traceIdentPlaceholder string = "\t"

func identLevel() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel)
}

func tracePrint(msg string) {
	fmt.Printf("%s%s\n", identLevel(), msg)
}

func incIdent() {
	traceLevel += 1
}

func decIdent() {
	traceLevel -= 1
}

func trace(msg string) string {
	incIdent()
	tracePrint("BEGIN " + msg)
	return msg
}

func untrace(msg string) {
	tracePrint("END " + msg)
	decIdent()
}
