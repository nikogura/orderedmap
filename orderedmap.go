// Copyright 2017 Moritz Fain
// Moritz Fain <moritz@mertinkat.net>
package orderedmap

import (
	"bytes"
	"encoding/json"
	"github.com/ghodss/yaml"
)

type KeyType string
type ValueType interface{}

type KV struct {
	Key   KeyType
	Value ValueType
}

type OrderedMap struct {
	kvList    []*KV
	idxLookup map[KeyType]int
}

// Creates new ordered map
func NewOrderedMap(kvList ...*KV) (om *OrderedMap) {
	om = &OrderedMap{
		idxLookup: make(map[KeyType]int),
	}

	for i := 0; i < len(kvList); i++ {
		om.Set(kvList[i].Key, kvList[i].Value)
	}
	return
}

// Sets value for given key
func (om *OrderedMap) Set(key KeyType, value ValueType) *OrderedMap {
	if idx, ok := om.idxLookup[key]; !ok {
		// insert new key value pair
		om.idxLookup[key] = len(om.kvList)
		om.kvList = append(om.kvList, &KV{key, value})
	} else {
		// update value
		om.kvList[idx].Value = value
	}

	return om
}

// Returns the given key's value or <nil> if key does not exist
func (om *OrderedMap) Get(key KeyType) ValueType {
	if idx, ok := om.idxLookup[key]; ok {
		return om.kvList[idx].Value
	}
	return nil
}

// Checks for existence of a given key
func (om *OrderedMap) Exists(key KeyType) (ok bool) {
	_, ok = om.idxLookup[key]
	return
}

func (obj *OrderedMap) Delete(key KeyType) {
	if idx, ok := obj.idxLookup[key]; ok {
		delete(obj.idxLookup, key)
		obj.kvList[idx] = nil
	}
}

// Returns ordered list of keys
func (om *OrderedMap) GetKeys() (keys []KeyType) {
	for idx := 0; idx < len(om.kvList); idx++ {
		if om.kvList[idx] != nil {
			keys = append(keys, om.kvList[idx].Key)

		}
	}
	return
}

// Returns ordered list of key-value pairs
func (om *OrderedMap) GetList() (kvList []KV) {
	for idx := 0; idx < len(om.kvList); idx++ {
		if om.kvList[idx] != nil {
			kvList = append(kvList, *om.kvList[idx])
		}
	}
	return
}

// Appends the given ordered map to the current one
func (om *OrderedMap) Append(newOm *OrderedMap, overwrite bool) *OrderedMap {
	for _, kv := range newOm.GetList() {
		if !overwrite && om.Exists(kv.Key) {
			continue
		}
		om.Set(kv.Key, kv.Value)
	}
	return om
}

// Returns length
func (om *OrderedMap) Len() int {
	return len(om.idxLookup)
}

// Returns JSON serialized string representation
func (om *OrderedMap) String() string {
	json, _ := om.MarshalJSON()
	return string(json)
}

// Marshal JSON for ordered map
func (om *OrderedMap) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")

	first := true
	for idx := 0; idx < len(om.kvList); idx++ {
		if om.kvList[idx] == nil {
			continue
		}

		if !first {
			buffer.WriteString(",")
		}

		jsonKey, err := json.Marshal(om.kvList[idx].Key)
		if err != nil {
			return nil, err
		}
		jsonValue, err := json.Marshal(om.kvList[idx].Value)
		if err != nil {
			return nil, err
		}
		buffer.Write(jsonKey)
		buffer.WriteByte(58)
		buffer.Write(jsonValue)

		first = false
	}

	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// Marshal JSON for single KV item (convenience only)
func (kv KV) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")

	jsonKey, err := json.Marshal(kv.Key)
	if err != nil {
		return nil, err
	}
	jsonValue, err := json.Marshal(kv.Value)
	if err != nil {
		return nil, err
	}

	buffer.Write(jsonKey)
	buffer.WriteByte(58)
	buffer.Write(jsonValue)

	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// MarshalYAML Marshalls the ordered map to yaml
func (om *OrderedMap) MarshalYAML() ([]byte, error) {
	jsonBytes, err := om.MarshalJSON()
	if err != nil {
		return nil, err
	}

	yamlBytes, err := yaml.JSONToYAML(jsonBytes)
	if err != nil {
		return nil, err
	}

	return yamlBytes, err
}

// MarshalYAML Marshals the ordered map to yaml
func (kv KV) MarshalYAML() ([]byte, error) {
	jsonBytes, err := kv.MarshalJSON()
	if err != nil {
		return nil, err
	}

	yamlBytes, err := yaml.JSONToYAML(jsonBytes)
	if err != nil {
		return nil, err
	}

	return yamlBytes, err

}
