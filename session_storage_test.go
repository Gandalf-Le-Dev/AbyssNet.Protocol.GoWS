package gows

import (
	"testing"

	"github.com/Gandalf-Le-Dev/abyssnet.protocol.gows/internal"
	"github.com/dolthub/maphash"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	var as = assert.New(t)
	var m1 = make(map[string]any)
	var m2 = newSmap()
	var count = internal.AlphabetNumeric.Intn(1000)
	for i := 0; i < count; i++ {
		var key = string(internal.AlphabetNumeric.Generate(10))
		var val = internal.AlphabetNumeric.Uint32()
		m1[key] = val
		m2.Store(key, val)
	}

	var keys = make([]string, 0)
	for k := range m1 {
		keys = append(keys, k)
	}
	for i := 0; i < len(keys)/2; i++ {
		delete(m1, keys[i])
		m2.Delete(keys[i])
	}

	for k, v := range m1 {
		v1, ok := m2.Load(k)
		as.Equal(true, ok)
		as.Equal(v, v1)
	}
	as.Equal(len(m1), m2.Len())
}

func TestSliceMap(t *testing.T) {
	var as = assert.New(t)
	var m = newSmap()
	m.Store("hong", 1)
	m.Store("mei", 2)
	m.Store("ming", 3)
	{
		v, _ := m.Load("hong")
		as.Equal(1, v)
	}
	{
		m.Delete("hong")
		v, ok := m.Load("hong")
		as.Equal(false, ok)
		as.Nil(v)

		m.Store("hong", 4)
		v, _ = m.Load("hong")
		as.Equal(4, v)
	}
}

func TestMap_Range(t *testing.T) {
	var as = assert.New(t)
	var m1 = make(map[any]any)
	var m2 = newSmap()
	var count = 1000
	for i := 0; i < count; i++ {
		var key = string(internal.AlphabetNumeric.Generate(10))
		var val = internal.AlphabetNumeric.Uint32()
		m1[key] = val
		m2.Store(key, val)
	}

	{
		var keys []any
		m2.Range(func(key string, value any) bool {
			v, ok := m1[key]
			as.Equal(true, ok)
			as.Equal(v, value)
			keys = append(keys, key)
			return len(keys) < 100
		})
		as.Equal(100, len(keys))
	}

	{
		var keys []any
		m2.Range(func(key string, value any) bool {
			v, ok := m1[key]
			as.Equal(true, ok)
			as.Equal(v, value)
			keys = append(keys, key)
			return true
		})
		as.Equal(1000, len(keys))
	}
}

func TestConcurrentMap(t *testing.T) {
	var as = assert.New(t)
	var m1 = make(map[string]any)
	var m2 = NewConcurrentMap[string, uint32](5)
	var count = internal.AlphabetNumeric.Intn(1000)
	for i := 0; i < count; i++ {
		var key = string(internal.AlphabetNumeric.Generate(10))
		var val = internal.AlphabetNumeric.Uint32()
		m1[key] = val
		m2.Store(key, val)
	}

	var keys = make([]string, 0)
	for k := range m1 {
		keys = append(keys, k)
	}
	for i := 0; i < len(keys)/2; i++ {
		delete(m1, keys[i])
		m2.Delete(keys[i])
	}

	for k, v := range m1 {
		v1, ok := m2.Load(k)
		as.Equal(true, ok)
		as.Equal(v, v1)
	}
	as.Equal(len(m1), m2.Len())
}

func TestConcurrentMap_Range(t *testing.T) {
	var as = assert.New(t)
	var m1 = make(map[any]any)
	var m2 = NewConcurrentMap[string, uint32](13)
	var count = 1000
	for i := 0; i < count; i++ {
		var key = string(internal.AlphabetNumeric.Generate(10))
		var val = internal.AlphabetNumeric.Uint32()
		m1[key] = val
		m2.Store(key, val)
	}

	{
		var keys []any
		m2.Range(func(key string, value uint32) bool {
			v, ok := m1[key]
			as.Equal(true, ok)
			as.Equal(v, value)
			keys = append(keys, key)
			return len(keys) < 100
		})
		as.Equal(100, len(keys))
	}

	{
		var keys []any
		m2.Range(func(key string, value uint32) bool {
			v, ok := m1[key]
			as.Equal(true, ok)
			as.Equal(v, value)
			keys = append(keys, key)
			return true
		})
		as.Equal(1000, len(keys))
	}
}

func TestHash(t *testing.T) {
	var h = maphash.NewHasher[string]()
	for i := 0; i < 1000; i++ {
		var a = string(internal.AlphabetNumeric.Generate(16))
		var b = string(internal.AlphabetNumeric.Generate(16))
		assert.Equal(t, h.Hash(a), h.Hash(a))
		assert.NotEqual(t, h.Hash(a), h.Hash(b))
	}
}
