// This is little program to help me understand what sort of type conversions
// are allowed. I wrote this to help me write the ValueIs() implementation.
package main

import (
	"log"
	"reflect"
	"unsafe"
)

type s struct {
	value string
}

func main() {
	val := "foo"
	vals := []interface{}{
		true,
		int(0),
		int8(0),
		int16(0),
		int32(0),
		int64(0),
		uint(0),
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),
		float32(0),
		float64(0),
		complex(float32(1), float32(1)),
		complex(float64(1), float64(1)),
		[1]string{"str"},
		make(chan int),
		map[string]string{"foo": "bar"},
		&val,
		[]string{"str"},
		"string",
		s{"value"},
		unsafe.Pointer(&val),
	}

	for _, v1 := range vals {
		for _, v2 := range vals {
			t1 := reflect.TypeOf(v1)
			t2 := reflect.TypeOf(v2)

			log.Printf("%s  -> %s ? %v", t1.Kind(), t2.Kind(), t1.ConvertibleTo(t2))
		}
	}
}
