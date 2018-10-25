package main

import (
	"fmt"
	"reflect"
)

type X struct{}
type Y struct{}

type str string

func main() {
	foo := "foo"
	explain(foo, "string")
	explain(&foo, "pointer to string")

	explain(X{}, "X struct")
	explain(&X{}, "pointer to X struct")

	var bar str = "bar"
	explain(bar, "str type")
	explain(&bar, "pointer to str type")

	explain([]string{"foo","bar"}, "slice of strings")
	explain([]str{"foo","bar"}, "slice of str")

	explain(42.2, "42.2")
}

func explain(val interface{}, desc string){
	v := reflect.ValueOf(val)
	fmt.Printf("%s\n", desc)
	fmt.Printf("  Kind = %s\n", v.Type().Kind().String())
	fmt.Printf("  Name = %s\n", v.Type().Name())
}
