package object

import (
	"testing"
)

func TestStringHashKey(t *testing.T) {
    hello1 := &String{Value: "Hello World"}
    hello2 := &String{Value: "Hello World"}
    diff1 := &String{Value: "My name is johnny"}
    diff2 := &String{Value: "My name is johnny"}
    if hello1.HashKey() != hello2.HashKey() {
        t.Errorf("strings with same content have different hash keys")
    }
    if diff1.HashKey() != diff2.HashKey() {
        t.Errorf("strings with same content have different hash keys")
    }
    if hello1.HashKey() == diff1.HashKey() {
        t.Errorf("strings with different content have same hash keys")
    }
}

func TestIntegerHashKey(t *testing.T) {
    one1 := &Integer{Value: 1}
    one2 := &Integer{Value: 1}
    meaningOfLife1 := &Integer{Value: 42}
    meaningOfLife2 := &Integer{Value: 42}
    if one1.HashKey() != one2.HashKey() {
        t.Errorf("Integers with same content have different hash keys")
    }
    if meaningOfLife1.HashKey() != meaningOfLife2.HashKey() {
        t.Errorf("Integers with same content have different hash keys")
    }
    if one1.HashKey() == meaningOfLife1.HashKey() {
        t.Errorf("Integers with different content have same hash keys")
    }
}

func TestBooleanHashKey(t *testing.T) {
	true1 := &Boolean{Value: true}
	true2 := &Boolean{Value: true}
	false1 := &Boolean{Value: false}
	false2 := &Boolean{Value: false}
    if true1.HashKey() != true2.HashKey() {
        t.Errorf("Booleans with same content have different hash keys")
    }
    if false1.HashKey() != false2.HashKey() {
        t.Errorf("Booleans with same content have different hash keys")
    }
    if true1.HashKey() == false1.HashKey() {
        t.Errorf("Booleans with different content have same hash keys")
    }
}
