// object/object.go

package object

import (
	"bytes"
	"container/list"
	"fmt"
	"hash/fnv"
	"math"
	"renelle/ast"
	"strconv"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	FLOAT_OBJ        = "FLOAT"
	STRING_OBJ       = "STRING"
	BOOLEAN_OBJ      = "BOOLEAN"
	ATOM_OBJ         = "ATOM"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	TUPLE_OBJ        = "TUPLE"
	MAP_OBJ          = "MAP"
	SLICE_OBJ        = "SLICE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}

type Pair struct {
	Key   Object
	Value Object
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashTable struct {
	Buckets []*list.List
	Size    int
	Length  int
}

func NewHashTable(size int) *HashTable {
	return &HashTable{
		Buckets: make([]*list.List, size),
		Size:    size,
		Length:  0,
	}
}

func (h *HashKey) Equals(other *HashKey) bool {
	return h.Type == other.Type && h.Value == other.Value
}

func (h *HashTable) Put(pair Pair) {
	hashKey := pair.Key.(Hashable).HashKey()
	index := int(hashKey.Value % uint64(h.Size))
	if h.Buckets[index] == nil {
		h.Buckets[index] = list.New()
	} else {
		// If the key already exists, update its value
		for e := h.Buckets[index].Front(); e != nil; e = e.Next() {
			if Equals(pair.Key, e.Value.(Pair).Key) {
				e.Value = pair
				return
			}
		}
	}
	h.Length++
	// If the key does not exist, add a new entry
	h.Buckets[index].PushBack(pair)
}

func (h *HashTable) Get(key Object) (Object, bool) {
	hashKey := key.(Hashable).HashKey()
	index := int(hashKey.Value % uint64(h.Size))
	if h.Buckets[index] == nil {
		return nil, false
	}
	for e := h.Buckets[index].Front(); e != nil; e = e.Next() {
		if Equals(key, e.Value.(Pair).Key) {
			return e.Value.(Pair).Value, true
		}
	}
	return nil, false
}

func (h *HashTable) Keys() []Object {
	var keys []Object
	for _, bucket := range h.Buckets {
		if bucket != nil {
			for e := bucket.Front(); e != nil; e = e.Next() {
				keys = append(keys, e.Value.(Pair).Key)
			}
		}
	}
	return keys
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Float struct {
	Value float64
}

func (f *Float) Inspect() string  { return fmt.Sprintf("%g", f.Value) }
func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) HashKey() HashKey {
	return HashKey{Type: f.Type(), Value: math.Float64bits(f.Value)}
}

type String struct {
	Value string
}

func (s *String) Inspect() string  { return "\"" + s.Value + "\"" }
func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) HashKey() HashKey {
	hasher := fnv.New64a()
	hasher.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: hasher.Sum64()}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	}
	return HashKey{Type: b.Type(), Value: value}
}

type Atom struct {
	Value string
}

func (a *Atom) Inspect() string  { return ":" + a.Value }
func (a *Atom) Type() ObjectType { return ATOM_OBJ }
func (a *Atom) HashKey() HashKey {
	hasher := fnv.New64a()
	hasher.Write([]byte(a.Value))
	return HashKey{Type: a.Type(), Value: hasher.Sum64()}
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

type Error struct {
	Message string
	Line    int
	Column  int
}

func (e *Error) Inspect() string {
	return fmt.Sprintf("Line: %d, Column %d: ERROR: %s", e.Line, e.Column, e.Message)
}
func (e *Error) Type() ObjectType { return ERROR_OBJ }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n\t")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()

}
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) HashKey() HashKey {
	hasher := fnv.New64a()
	hasher.Write([]byte(fmt.Sprintf("%p", f)))
	return HashKey{Type: f.Type(), Value: hasher.Sum64()}
}

type BuiltinFunction func(ctx *EvalContext, args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) HashKey() HashKey {
	hasher := fnv.New64a()
	hasher.Write([]byte(fmt.Sprintf("%p", b)))
	return HashKey{Type: b.Type(), Value: hasher.Sum64()}
}

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	const maxElements = 3
	var out bytes.Buffer

	elements := []string{}
	for i, el := range ao.Elements {
		if i > maxElements {
			elements = append(elements, "...")
			break
		}
		elements = append(elements, el.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, " "))
	out.WriteString("]")

	return out.String()
}

func (ao *Array) HashKey() HashKey {
	hasher := fnv.New64a()
	for _, item := range ao.Elements {
		itemHashKey := item.(Hashable).HashKey()
		hasher.Write([]byte(strconv.FormatUint(itemHashKey.Value, 10)))
	}
	return HashKey{Type: ao.Type(), Value: hasher.Sum64()}
}

type Tuple struct {
	Elements []Object
}

func (to *Tuple) Type() ObjectType { return TUPLE_OBJ }
func (to *Tuple) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range to.Elements {
		elements = append(elements, el.Inspect())
	}

	out.WriteString("(")
	out.WriteString(strings.Join(elements, " "))
	out.WriteString(")")

	return out.String()
}

func (to *Tuple) HashKey() HashKey {
	hasher := fnv.New64a()
	for _, item := range to.Elements {
		itemHashKey := item.(Hashable).HashKey()
		hasher.Write([]byte(strconv.FormatUint(itemHashKey.Value, 10)))
	}
	return HashKey{Type: to.Type(), Value: hasher.Sum64()}
}

type Slice struct {
	Start Object
	End   Object
}

func (s *Slice) Type() ObjectType { return "SLICE" }
func (s *Slice) Inspect() string {
	return fmt.Sprintf("%s::%s", s.Start.Inspect(), s.End.Inspect())
}

type Map struct {
	Store *HashTable
}

func (m *Map) Type() ObjectType { return MAP_OBJ }
func (m *Map) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, bucket := range m.Store.Buckets {
		if bucket == nil {
			continue
		}
		for e := bucket.Front(); e != nil; e = e.Next() {
			pair := e.Value.(Pair)
			elements = append(elements, pair.Key.Inspect()+" = "+pair.Value.Inspect())
		}
	}

	out.WriteString("{")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")

	return out.String()
}

func (m *Map) HashKey() HashKey {
	hasher := fnv.New64a()
	for _, bucket := range m.Store.Buckets {
		if bucket != nil {

			for e := bucket.Front(); e != nil; e = e.Next() {
				pair := e.Value.(Pair)
				key, keyOk := pair.Key.(Hashable)
				value, valueOk := pair.Value.(Hashable)
				if keyOk && valueOk {
					keyHash := key.HashKey().Value
					valueHash := value.HashKey().Value
					hasher.Write([]byte(fmt.Sprintf("%s%d%s%d", key.HashKey().Type, keyHash, value.HashKey().Type, valueHash)))
				}
			}
		}
	}
	return HashKey{Type: m.Type(), Value: hasher.Sum64()}
}

func (m *Map) Get(key Object) (Object, bool) {
	hashKey := key.(Hashable).HashKey()
	index := int(hashKey.Value % uint64(len(m.Store.Buckets)))
	if m.Store.Buckets[index] == nil {
		return nil, false
	}
	for e := m.Store.Buckets[index].Front(); e != nil; e = e.Next() {
		if Equals(key, e.Value.(Pair).Key) {
			return e.Value.(Pair).Value, true
		}
	}
	return nil, false
}
func (m *Map) Put(key, value Object) {
	m.Store.Put(Pair{Key: key, Value: value})
}

func (m *Map) Copy(newItems int) *Map {
	hashTableSize := m.Store.Size + newItems
	hashTable := NewHashTable(hashTableSize)
	for _, key := range m.Store.Keys() {
		value, _ := m.Store.Get(key)
		hashTable.Put(Pair{Key: key, Value: value})
	}
	return &Map{Store: hashTable}
}

func (m *Map) Keys() []Object {
	return m.Store.Keys()
}

type Env interface {
	Get(name string) (Object, bool)
	Set(name string, val Object) Object
}

type Module struct {
	Name        string
	Environment Env
}

func (m *Module) Type() ObjectType { return "MODULE" }
func (m *Module) Inspect() string {
	return fmt.Sprintf("module %s", m.Name)
}

func Equals(a, b Object) bool {
	switch a := a.(type) {
	case *Integer:
		b, ok := b.(*Integer)
		return ok && a.Value == b.Value
	case *Float:
		b, ok := b.(*Float)
		return ok && a.Value == b.Value
	case *String:
		b, ok := b.(*String)
		return ok && a.Value == b.Value
	case *Boolean:
		b, ok := b.(*Boolean)
		return ok && a.Value == b.Value
	case *Atom:
		b, ok := b.(*Atom)
		return ok && a.Value == b.Value
	case *ReturnValue:
		b, ok := b.(*ReturnValue)
		return ok && Equals(a.Value, b.Value)
	case *Array:
		b, ok := b.(*Array)
		if !ok || len(a.Elements) != len(b.Elements) {
			return false
		}
		for i, el := range a.Elements {
			if !Equals(el, b.Elements[i]) {
				return false
			}
		}
		return true
	case *Function:
		b, ok := b.(*Function)
		return ok && a == b
	case *Builtin:
		b, ok := b.(*Builtin)
		return ok && a == b
	case *Tuple:
		b, ok := b.(*Tuple)
		if !ok || len(a.Elements) != len(b.Elements) {
			return false
		}
		for i, el := range a.Elements {
			if !Equals(el, b.Elements[i]) {
				return false
			}
		}
		return true

	case *Map:
		b, ok := b.(*Map)
		if !ok || len(a.Store.Buckets) != len(b.Store.Buckets) {
			return false
		}
		for i, bucketA := range a.Store.Buckets {
			bucketB := b.Store.Buckets[i]
			if (bucketA == nil && bucketB != nil) || (bucketA != nil && bucketB == nil) {
				return false
			}
			if bucketA != nil && bucketB != nil && bucketA.Len() != bucketB.Len() {
				return false
			}
			if bucketA != nil && bucketB != nil {
				for eA, eB := bucketA.Front(), bucketB.Front(); eA != nil && eB != nil; eA, eB = eA.Next(), eB.Next() {
					pairA := eA.Value.(Pair)
					pairB := eB.Value.(Pair)
					if !Equals(pairA.Key, pairB.Key) || !Equals(pairA.Value, pairB.Value) {
						return false
					}
				}
			}
		}
		return true
	default:
		return false
	}
}
