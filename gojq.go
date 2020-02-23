package gojq

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// JNode is the node for tree representation of underlying JSON.
type JNode struct {
	value   interface{}
	key     string
	parent  *JNode
	parray  bool
	isArray bool
	Err     error
}

// FromBytes constructs value of JNode from a JSON bytes string.
func FromBytes(jsonStr []byte) *JNode {
	jn := &JNode{}
	err := json.Unmarshal(jsonStr, &jn.value)
	if err != nil {
		jn.Err = err
	}
	return jn
}

// FromInterface creates jq object around the given value.
func FromInterface(val interface{}) *JNode {
	return &JNode{value: val}
}

// String returns string representation of JNode's contents.
func (jn *JNode) String() string {
	bjson, err := json.Marshal(jn.value)
	if err != nil {
		jn.Err = err
		return fmt.Sprintf("[JNode Error: %v]", err)
	}
	return fmt.Sprintf("%s", bjson)
}

// As unmarshals value JSON into the given variable.
func (jn *JNode) As(v interface{}) error {
	if jn.Err != nil {
		return jn.Err
	}
	bstr, err := json.Marshal(jn.value)
	if err != nil {
		return err
	}
	return json.Unmarshal(bstr, v)
}

// Select walks down the JSON tree to the node at specified path.
func (jn *JNode) Select(path ...string) *JNode {
	ret, err := recursiveWalk(jn, path)
	if err != nil {
		jn.Err = err
		return jn
	}
	return ret
}

func recursiveWalk(jn *JNode, path []string) (*JNode, error) {
	if len(path) == 0 {
		return jn, nil
	}
	key := path[0]
	if len(path) == 1 {
		var val interface{}
		var arr bool
		switch jn.value.(type) {
		case string, float64, bool, nil: // simple value
			return nil, fmt.Errorf("Can't index into primitive value [%s]", key)
		case []interface{}: // array indexing
			arr = true
			idx, err := strconv.Atoi(key)
			if err != nil {
				return nil, err
			}
			val = jn.value.([]interface{})[idx]
		case map[string]interface{}: // map walk
			val = jn.value.(map[string]interface{})[key]
			arr = false
		}

		ret := &JNode{
			value:   val,
			key:     key,
			parent:  jn,
			parray:  jn.isArray,
			isArray: arr,
		}
		return ret, nil
	}
	rest := path[1:]
	switch jn.value.(type) {
	case []interface{}:
		idx, err := strconv.Atoi(key)
		if err != nil {
			return nil, err
		}
		node := JNode{
			value:  jn.value.([]interface{})[idx],
			key:    key,
			parent: jn,
			parray: jn.isArray,
		}
		return recursiveWalk(&node, rest)
	case map[string]interface{}:
		node := JNode{
			value:  jn.value.(map[string]interface{})[key],
			key:    key,
			parent: jn,
			parray: false,
		}
		return recursiveWalk(&node, rest)
	}
	// if the value is a primitive
	if len(rest) > 0 {
		return jn, fmt.Errorf("Can not index into primitive value [%s]", rest[0])
	}
	return jn, nil
}

// Set sets the value of the current JNode in the json tree
func (jn *JNode) Set(val interface{}) *JNode {
	if jn.parent != nil {
		// update parent
		switch jn.parent.value.(type) {
		case []interface{}:
			idx, err := strconv.Atoi(jn.key)
			if err != nil {
				jn.Err = err
				return jn
			}
			jn.parent.value.([]interface{})[idx] = val
		case map[string]interface{}:
			jn.parent.value.(map[string]interface{})[jn.key] = val
		}
	} else {
		// root, update just its value
		jn.value = val
	}
	return jn
}

// Exists checks to see if a property is present
func (jn *JNode) Exists(path ...string) bool {
	q, err := recursiveWalk(jn, path)
	if err != nil {
		jn.Err = err
		return false
	}
	if q.value != nil {
		return true
	}
	return false
}

// Iterator returns an iterator object for an array property
func (jn *JNode) Iterator() *JNIterator {
	// TODO make sure it is an array
	return &JNIterator{i: -1, data: jn.value.([]interface{})}
}

//JNIterator is iterator for array fields in a JNode
type JNIterator struct {
	i      int
	data   []interface{}
	parent *JNode
}

// Next advances the iterator and returns if there is a next element
func (it *JNIterator) Next() bool {
	it.i++
	if it.i < len(it.data) {
		return true
	}
	return false
}

// Value returns value of the current element
func (it *JNIterator) Value() interface{} {
	return it.data[it.i]
}

// Node returns the current element wrapped in a JNode
func (it *JNIterator) Node() *JNode {
	return &JNode{
		value:  it.data[it.i],
		key:    strconv.Itoa(it.i),
		parent: it.parent,
		parray: true,
	}
}

// Count returns the number of elements in the iterator
func (it *JNIterator) Count() int {
	return len(it.data)
}
