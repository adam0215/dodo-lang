package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is Dodo"}
	diff2 := &String{Value: "My name is Dodo"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with the same content produce different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with the same content produce different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with the different content produce same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	true1 := &Boolean{Value: true}
	true2 := &Boolean{Value: true}
	diff1 := &Boolean{Value: false}
	diff2 := &Boolean{Value: false}

	if true1.HashKey() != true2.HashKey() {
		t.Errorf("strings with the same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with the same content have different hash keys")
	}

	if true1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with the different content have same hash keys")
	}
}

func TestIntegerHashKey(t *testing.T) {
	fifteen1 := &Integer{Value: 15}
	fifteen2 := &Integer{Value: 15}
	diff1 := &Integer{Value: 420}
	diff2 := &Integer{Value: 420}

	if fifteen1.HashKey() != fifteen2.HashKey() {
		t.Errorf("strings with the same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with the same content have different hash keys")
	}

	if fifteen1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with the different content have same hash keys")
	}
}
