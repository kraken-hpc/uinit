package uinit

import "fmt"

var _ KeyValue = (*SimpleKV)(nil)

// NewSimpleKV creates a new, initialized SimpleKV
func NewSimpleKV() (kv *SimpleKV) {
	kv = &SimpleKV{}
	kv.Init()
	return
}

// SimpleKV just provides an interface to a map
type SimpleKV struct {
	store map[string]string
}

// Init performs initializes internal data for the store
func (kv *SimpleKV) Init() {
	kv.store = make(map[string]string)
}

// Set sets a named value in a map (will override)
func (kv *SimpleKV) Set(name, val string) (err error) {
	kv.store[name] = val
	return
}

// Get tries to get a named value, returns an error if it doesn't exist
func (kv *SimpleKV) Get(name string) (val string, err error) {
	var ok bool
	if val, ok = kv.store[name]; !ok {
		return "", fmt.Errorf("no such key: %s", name)
	}
	return
}

// GetDefault tries to get a named value, sets to dval if it doesn't exist
func (kv *SimpleKV) GetDefault(name, dval string) string {
	if val, err := kv.Get(name); err != nil {
		return val
	}
	return dval
}

// Clear clears the kv store
func (kv *SimpleKV) Clear() (err error) {
	kv.Init()
	return
}
