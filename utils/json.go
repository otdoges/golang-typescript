package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"typescript-golang/types"
)

// JSON provides TypeScript-like JSON functionality
type JSONUtils struct{}

// Global JSON utilities instance
var JSON = JSONUtils{}

// Stringify converts object to JSON string (like JSON.stringify())
func (JSONUtils) Stringify(value interface{}, replacer ...func(string, interface{}) interface{}) (string, error) {
	var processedValue interface{}
	
	if len(replacer) > 0 && replacer[0] != nil {
		processedValue = processWithReplacer(value, replacer[0], "")
	} else {
		processedValue = value
	}
	
	bytes, err := json.MarshalIndent(processedValue, "", "  ")
	if err != nil {
		return "", err
	}
	
	return string(bytes), nil
}

// StringifyCompact converts object to compact JSON string
func (JSONUtils) StringifyCompact(value interface{}) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Parse parses JSON string to interface{} (like JSON.parse())
func (JSONUtils) Parse(jsonString string, reviver ...func(string, interface{}) interface{}) (interface{}, error) {
	var result interface{}
	
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, err
	}
	
	if len(reviver) > 0 && reviver[0] != nil {
		result = processWithReviver(result, reviver[0], "")
	}
	
	return result, nil
}

// ParseTo parses JSON string into specific type
func ParseTo[T any](jsonString string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(jsonString), &result)
	return result, err
}

// processWithReplacer processes object with replacer function
func processWithReplacer(value interface{}, replacer func(string, interface{}) interface{}, key string) interface{} {
	processed := replacer(key, value)
	
	switch v := processed.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, val := range v {
			result[k] = processWithReplacer(val, replacer, k)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = processWithReplacer(val, replacer, strconv.Itoa(i))
		}
		return result
	default:
		return processed
	}
}

// processWithReviver processes object with reviver function
func processWithReviver(value interface{}, reviver func(string, interface{}) interface{}, key string) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, val := range v {
			result[k] = processWithReviver(val, reviver, k)
		}
		return reviver(key, result)
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = processWithReviver(val, reviver, strconv.Itoa(i))
		}
		return reviver(key, result)
	default:
		return reviver(key, value)
	}
}

// Object provides TypeScript-like Object functionality
type ObjectUtils struct{}

// Global Object utilities instance
var Object = ObjectUtils{}

// Keys returns object keys (like Object.keys())
func (ObjectUtils) Keys(obj interface{}) []string {
	val := reflect.ValueOf(obj)
	
	switch val.Kind() {
	case reflect.Map:
		keys := make([]string, 0, val.Len())
		for _, key := range val.MapKeys() {
			keys = append(keys, fmt.Sprintf("%v", key.Interface()))
		}
		return keys
	case reflect.Struct:
		typ := val.Type()
		keys := make([]string, 0, val.NumField())
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			if field.IsExported() {
				// Check for json tag
				if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
					// Handle json:",omitempty" tags
					if parts := strings.Split(tag, ","); len(parts) > 0 && parts[0] != "" {
						keys = append(keys, parts[0])
					} else {
						keys = append(keys, field.Name)
					}
				} else {
					keys = append(keys, field.Name)
				}
			}
		}
		return keys
	case reflect.Ptr:
		if !val.IsNil() {
			return Object.Keys(val.Elem().Interface())
		}
	}
	
	return []string{}
}

// Values returns object values (like Object.values())
func (ObjectUtils) Values(obj interface{}) []interface{} {
	val := reflect.ValueOf(obj)
	
	switch val.Kind() {
	case reflect.Map:
		values := make([]interface{}, 0, val.Len())
		for _, key := range val.MapKeys() {
			values = append(values, val.MapIndex(key).Interface())
		}
		return values
	case reflect.Struct:
		values := make([]interface{}, 0, val.NumField())
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			if field.IsExported() {
				values = append(values, val.Field(i).Interface())
			}
		}
		return values
	case reflect.Ptr:
		if !val.IsNil() {
			return Object.Values(val.Elem().Interface())
		}
	}
	
	return []interface{}{}
}

// Entries returns key-value pairs (like Object.entries())
func (ObjectUtils) Entries(obj interface{}) []types.Tuple2[string, interface{}] {
	val := reflect.ValueOf(obj)
	
	switch val.Kind() {
	case reflect.Map:
		entries := make([]types.Tuple2[string, interface{}], 0, val.Len())
		for _, key := range val.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			value := val.MapIndex(key).Interface()
			entries = append(entries, types.NewTuple2(keyStr, value))
		}
		return entries
	case reflect.Struct:
		typ := val.Type()
		entries := make([]types.Tuple2[string, interface{}], 0, val.NumField())
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			if field.IsExported() {
				keyName := field.Name
				if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
					if parts := strings.Split(tag, ","); len(parts) > 0 && parts[0] != "" {
						keyName = parts[0]
					}
				}
				value := val.Field(i).Interface()
				entries = append(entries, types.NewTuple2(keyName, value))
			}
		}
		return entries
	case reflect.Ptr:
		if !val.IsNil() {
			return Object.Entries(val.Elem().Interface())
		}
	}
	
	return []types.Tuple2[string, interface{}]{}
}

// FromEntries creates object from key-value pairs (like Object.fromEntries())
func (ObjectUtils) FromEntries(entries []types.Tuple2[string, interface{}]) map[string]interface{} {
	result := make(map[string]interface{})
	for _, entry := range entries {
		result[entry.First] = entry.Second
	}
	return result
}

// Assign merges objects (like Object.assign())
func (ObjectUtils) Assign(target map[string]interface{}, sources ...map[string]interface{}) map[string]interface{} {
	if target == nil {
		target = make(map[string]interface{})
	}
	
	for _, source := range sources {
		for key, value := range source {
			target[key] = value
		}
	}
	
	return target
}

// HasOwnProperty checks if object has property (like obj.hasOwnProperty())
func (ObjectUtils) HasOwnProperty(obj interface{}, prop string) bool {
	val := reflect.ValueOf(obj)
	
	switch val.Kind() {
	case reflect.Map:
		return val.MapIndex(reflect.ValueOf(prop)).IsValid()
	case reflect.Struct:
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			if field.IsExported() {
				if field.Name == prop {
					return true
				}
				if tag := field.Tag.Get("json"); tag != "" {
					if parts := strings.Split(tag, ","); len(parts) > 0 && parts[0] == prop {
						return true
					}
				}
			}
		}
		return false
	case reflect.Ptr:
		if !val.IsNil() {
			return Object.HasOwnProperty(val.Elem().Interface(), prop)
		}
	}
	
	return false
}

// GetProperty gets property value from object
func (ObjectUtils) GetProperty(obj interface{}, prop string) types.Optional[interface{}] {
	val := reflect.ValueOf(obj)
	
	switch val.Kind() {
	case reflect.Map:
		value := val.MapIndex(reflect.ValueOf(prop))
		if value.IsValid() {
			return types.Some[interface{}](value.Interface())
		}
	case reflect.Struct:
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			if field.IsExported() {
				match := false
				if field.Name == prop {
					match = true
				} else if tag := field.Tag.Get("json"); tag != "" {
					if parts := strings.Split(tag, ","); len(parts) > 0 && parts[0] == prop {
						match = true
					}
				}
				
				if match {
					return types.Some[interface{}](val.Field(i).Interface())
				}
			}
		}
	case reflect.Ptr:
		if !val.IsNil() {
			return Object.GetProperty(val.Elem().Interface(), prop)
		}
	}
	
	return types.None[interface{}]()
}

// SetProperty sets property value on object (for maps)
func (ObjectUtils) SetProperty(obj map[string]interface{}, prop string, value interface{}) map[string]interface{} {
	if obj == nil {
		obj = make(map[string]interface{})
	}
	obj[prop] = value
	return obj
}

// DeleteProperty deletes property from object (for maps)
func (ObjectUtils) DeleteProperty(obj map[string]interface{}, prop string) map[string]interface{} {
	if obj != nil {
		delete(obj, prop)
	}
	return obj
}

// DeepClone creates a deep copy of an object
func (ObjectUtils) DeepClone(obj interface{}) (interface{}, error) {
	// Use JSON round-trip for deep cloning
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	
	var result interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, err
	}
	
	return result, nil
}

// DeepEqual checks if two objects are deeply equal
func (ObjectUtils) DeepEqual(obj1, obj2 interface{}) bool {
	return reflect.DeepEqual(obj1, obj2)
}

// Freeze creates a "frozen" copy (read-only) - simulated with a wrapper
type FrozenObject struct {
	data map[string]interface{}
}

// Get retrieves a value from frozen object
func (fo *FrozenObject) Get(key string) types.Optional[interface{}] {
	if value, exists := fo.data[key]; exists {
		return types.Some[interface{}](value)
	}
	return types.None[interface{}]()
}

// Keys returns keys of frozen object
func (fo *FrozenObject) Keys() []string {
	keys := make([]string, 0, len(fo.data))
	for key := range fo.data {
		keys = append(keys, key)
	}
	return keys
}

// Values returns values of frozen object
func (fo *FrozenObject) Values() []interface{} {
	values := make([]interface{}, 0, len(fo.data))
	for _, value := range fo.data {
		values = append(values, value)
	}
	return values
}

// Freeze creates a frozen (immutable) object
func (ObjectUtils) Freeze(obj map[string]interface{}) *FrozenObject {
	data := make(map[string]interface{})
	for key, value := range obj {
		data[key] = value
	}
	return &FrozenObject{data: data}
}

// Merge deeply merges objects
func (ObjectUtils) Merge(objects ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	
	for _, obj := range objects {
		for key, value := range obj {
			if existing, exists := result[key]; exists {
				// If both are maps, merge recursively
				if existingMap, ok := existing.(map[string]interface{}); ok {
					if valueMap, ok := value.(map[string]interface{}); ok {
						result[key] = Object.Merge(existingMap, valueMap)
						continue
					}
				}
			}
			result[key] = value
		}
	}
	
	return result
}

// Pick creates object with only specified keys
func (ObjectUtils) Pick(obj map[string]interface{}, keys ...string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, key := range keys {
		if value, exists := obj[key]; exists {
			result[key] = value
		}
	}
	return result
}

// Omit creates object without specified keys
func (ObjectUtils) Omit(obj map[string]interface{}, keys ...string) map[string]interface{} {
	result := make(map[string]interface{})
	keySet := make(map[string]bool)
	for _, key := range keys {
		keySet[key] = true
	}
	
	for key, value := range obj {
		if !keySet[key] {
			result[key] = value
		}
	}
	return result
}

// MapValues transforms all values in object
func (ObjectUtils) MapValues(obj map[string]interface{}, fn func(interface{}) interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range obj {
		result[key] = fn(value)
	}
	return result
}

// MapKeys transforms all keys in object
func (ObjectUtils) MapKeys(obj map[string]interface{}, fn func(string) string) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range obj {
		newKey := fn(key)
		result[newKey] = value
	}
	return result
}

// FilterKeys filters object by keys
func (ObjectUtils) FilterKeys(obj map[string]interface{}, predicate func(string) bool) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range obj {
		if predicate(key) {
			result[key] = value
		}
	}
	return result
}

// FilterValues filters object by values
func (ObjectUtils) FilterValues(obj map[string]interface{}, predicate func(interface{}) bool) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range obj {
		if predicate(value) {
			result[key] = value
		}
	}
	return result
}

// IsEmpty checks if object is empty
func (ObjectUtils) IsEmpty(obj map[string]interface{}) bool {
	return len(obj) == 0
}

// Size returns number of properties in object
func (ObjectUtils) Size(obj map[string]interface{}) int {
	return len(obj)
}

// ToQueryString converts object to URL query string
func (ObjectUtils) ToQueryString(obj map[string]interface{}) string {
	var parts []string
	for key, value := range obj {
		parts = append(parts, fmt.Sprintf("%s=%v", key, value))
	}
	return strings.Join(parts, "&")
}