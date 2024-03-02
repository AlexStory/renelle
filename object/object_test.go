// object/object_test.go

package object

import (
	"testing"
)

func TestMap(t *testing.T) {
	m := &Map{Store: NewHashTable(10)}

	key1 := &String{Value: "key1"}
	value1 := &Integer{Value: 1}
	m.Put(key1, value1)

	key2 := &String{Value: "key2"}
	value2 := &Integer{Value: 2}
	m.Put(key2, value2)

	v, ok := m.Get(key1)
	if !ok || v.Inspect() != value1.Inspect() {
		t.Errorf("Expected %s, got %s", value1.Inspect(), v.Inspect())
	}

	v, ok = m.Get(key2)
	if !ok || v.Inspect() != value2.Inspect() {
		t.Errorf("Expected %s, got %s", value2.Inspect(), v.Inspect())
	}
}

func TestHashIndependence(t *testing.T) {
	m := &Map{
		Store: NewHashTable(1000),
	}

	// Integer and String with similar value
	integer123 := &String{Value: "Integer 123"}
	m.Put(&Integer{Value: 123}, integer123)

	string123 := &String{Value: "String 123"}
	m.Put(&String{Value: "123"}, string123)

	val, ok := m.Get(&Integer{Value: 123})
	if !ok || (ok && val != integer123) {
		t.Errorf("Expected 'Integer 123', got '%v'", val)
	}

	val, ok = m.Get(&String{Value: "123"})
	if !ok || (ok && val != string123) {
		t.Errorf("Expected 'String 123', got '%v'", val)
	}

	// Float and Integer with similar value
	float456 := &String{Value: "Float 456.0"}
	m.Put(&Float{Value: 456.0}, float456)

	integer456 := &String{Value: "Integer 456"}
	m.Put(&Integer{Value: 456}, integer456)

	val, ok = m.Get(&Integer{Value: 456})
	if !ok || (ok && val != integer456) {
		t.Errorf("Expected 'Integer 456', got '%v'", val)
	}

	val, ok = m.Get(&Float{Value: 456.0})
	if !ok || (ok && val != float456) {
		t.Errorf("Expected 'Float 456.0', got '%v'", val)
	}

	// Boolean and String similar value
	booleanTrue := &String{Value: "Boolean true"}
	m.Put(&Boolean{Value: true}, booleanTrue)

	stringTrue := &String{Value: "String true"}
	m.Put(&String{Value: "true"}, stringTrue)

	val, ok = m.Get(&Boolean{Value: true})
	if !ok || (ok && val != booleanTrue) {
		t.Errorf("Expected 'Boolean true', got '%v'", val)
	}

	val, ok = m.Get(&String{Value: "true"})
	if !ok || (ok && val != stringTrue) {
		t.Errorf("Expected 'String true', got '%v'", val)
	}
}
