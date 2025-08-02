package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
	"typescript-golang/types"
)

// StringUtils provides TypeScript-like string manipulation methods
type StringUtils struct{}

// Global string utilities instance
var Strings = StringUtils{}

// Length returns the length of string (like string.length in TypeScript)
func (StringUtils) Length(s string) int {
	return utf8.RuneCountInString(s)
}

// CharAt returns character at specified index (like string.charAt())
func (StringUtils) CharAt(s string, index int) string {
	runes := []rune(s)
	if index < 0 || index >= len(runes) {
		return ""
	}
	return string(runes[index])
}

// CharCodeAt returns Unicode code of character at index (like string.charCodeAt())
func (StringUtils) CharCodeAt(s string, index int) types.Optional[int] {
	runes := []rune(s)
	if index < 0 || index >= len(runes) {
		return types.None[int]()
	}
	return types.Some(int(runes[index]))
}

// Concat concatenates strings (like string.concat())
func (StringUtils) Concat(s string, others ...string) string {
	return s + strings.Join(others, "")
}

// IndexOf returns first index of substring (like string.indexOf())
func (StringUtils) IndexOf(s, substr string) int {
	index := strings.Index(s, substr)
	if index == -1 {
		return -1
	}
	// Convert byte index to rune index
	runes := []rune(s)
	byteIndex := 0
	for i, r := range runes {
		if byteIndex == index {
			return i
		}
		byteIndex += len(string(r))
	}
	return -1
}

// LastIndexOf returns last index of substring (like string.lastIndexOf())
func (StringUtils) LastIndexOf(s, substr string) int {
	index := strings.LastIndex(s, substr)
	if index == -1 {
		return -1
	}
	// Convert byte index to rune index
	runes := []rune(s)
	byteIndex := 0
	for i, r := range runes {
		if byteIndex == index {
			return i
		}
		byteIndex += len(string(r))
	}
	return -1
}

// Substring returns substring (like string.substring())
func (StringUtils) Substring(s string, start int, end ...int) string {
	runes := []rune(s)
	length := len(runes)
	
	// Handle negative start
	if start < 0 {
		start = 0
	}
	if start > length {
		start = length
	}
	
	// Handle end parameter
	endPos := length
	if len(end) > 0 {
		endPos = end[0]
		if endPos < 0 {
			endPos = 0
		}
		if endPos > length {
			endPos = length
		}
	}
	
	// Ensure start <= end
	if start > endPos {
		start, endPos = endPos, start
	}
	
	return string(runes[start:endPos])
}

// Substr returns substring with length (like string.substr())
func (StringUtils) Substr(s string, start int, length ...int) string {
	runes := []rune(s)
	runeLength := len(runes)
	
	// Handle negative start
	if start < 0 {
		start = runeLength + start
		if start < 0 {
			start = 0
		}
	}
	
	if start >= runeLength {
		return ""
	}
	
	// Handle length parameter
	extractLength := runeLength - start
	if len(length) > 0 && length[0] >= 0 {
		if length[0] < extractLength {
			extractLength = length[0]
		}
	}
	
	end := start + extractLength
	if end > runeLength {
		end = runeLength
	}
	
	return string(runes[start:end])
}

// Slice returns slice of string (like string.slice())
func (StringUtils) Slice(s string, start int, end ...int) string {
	runes := []rune(s)
	length := len(runes)
	
	// Handle negative start
	if start < 0 {
		start = length + start
		if start < 0 {
			start = 0
		}
	}
	if start > length {
		start = length
	}
	
	// Handle end parameter
	endPos := length
	if len(end) > 0 {
		endPos = end[0]
		if endPos < 0 {
			endPos = length + endPos
			if endPos < 0 {
				endPos = 0
			}
		}
		if endPos > length {
			endPos = length
		}
	}
	
	if start >= endPos {
		return ""
	}
	
	return string(runes[start:endPos])
}

// ToUpperCase converts to uppercase (like string.toUpperCase())
func (StringUtils) ToUpperCase(s string) string {
	return strings.ToUpper(s)
}

// ToLowerCase converts to lowercase (like string.toLowerCase())
func (StringUtils) ToLowerCase(s string) string {
	return strings.ToLower(s)
}

// Trim removes whitespace from both ends (like string.trim())
func (StringUtils) Trim(s string) string {
	return strings.TrimSpace(s)
}

// TrimStart removes whitespace from start (like string.trimStart())
func (StringUtils) TrimStart(s string) string {
	return strings.TrimLeftFunc(s, unicode.IsSpace)
}

// TrimEnd removes whitespace from end (like string.trimEnd())
func (StringUtils) TrimEnd(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// PadStart pads string at start (like string.padStart())
func (StringUtils) PadStart(s string, targetLength int, padString ...string) string {
	runes := []rune(s)
	currentLength := len(runes)
	
	if currentLength >= targetLength {
		return s
	}
	
	pad := " "
	if len(padString) > 0 && padString[0] != "" {
		pad = padString[0]
	}
	
	padLength := targetLength - currentLength
	padRunes := []rune(pad)
	
	if len(padRunes) == 0 {
		return s
	}
	
	var result []rune
	for len(result) < padLength {
		remaining := padLength - len(result)
		if remaining >= len(padRunes) {
			result = append(result, padRunes...)
		} else {
			result = append(result, padRunes[:remaining]...)
		}
	}
	
	return string(result) + s
}

// PadEnd pads string at end (like string.padEnd())
func (StringUtils) PadEnd(s string, targetLength int, padString ...string) string {
	runes := []rune(s)
	currentLength := len(runes)
	
	if currentLength >= targetLength {
		return s
	}
	
	pad := " "
	if len(padString) > 0 && padString[0] != "" {
		pad = padString[0]
	}
	
	padLength := targetLength - currentLength
	padRunes := []rune(pad)
	
	if len(padRunes) == 0 {
		return s
	}
	
	var result []rune
	for len(result) < padLength {
		remaining := padLength - len(result)
		if remaining >= len(padRunes) {
			result = append(result, padRunes...)
		} else {
			result = append(result, padRunes[:remaining]...)
		}
	}
	
	return s + string(result)
}

// Replace replaces first occurrence (like string.replace() with string)
func (StringUtils) Replace(s, old, new string) string {
	return strings.Replace(s, old, new, 1)
}

// ReplaceAll replaces all occurrences (like string.replaceAll())
func (StringUtils) ReplaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// ReplaceRegex replaces using regex (like string.replace() with RegExp)
func (StringUtils) ReplaceRegex(s, pattern, replacement string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return s, err
	}
	return re.ReplaceAllString(s, replacement), nil
}

// Split splits string by separator (like string.split())
func (StringUtils) Split(s, separator string, limit ...int) []string {
	if separator == "" {
		// Split into individual characters
		runes := []rune(s)
		result := make([]string, len(runes))
		for i, r := range runes {
			result[i] = string(r)
		}
		return result
	}
	
	parts := strings.Split(s, separator)
	
	// Apply limit if specified
	if len(limit) > 0 && limit[0] > 0 && limit[0] < len(parts) {
		// Join remaining parts
		remaining := strings.Join(parts[limit[0]-1:], separator)
		result := make([]string, limit[0])
		copy(result, parts[:limit[0]-1])
		result[limit[0]-1] = remaining
		return result
	}
	
	return parts
}

// StartsWith checks if string starts with prefix (like string.startsWith())
func (StringUtils) StartsWith(s, prefix string, position ...int) bool {
	start := 0
	if len(position) > 0 {
		start = position[0]
	}
	
	if start < 0 {
		start = 0
	}
	
	runes := []rune(s)
	if start >= len(runes) {
		return prefix == ""
	}
	
	return strings.HasPrefix(string(runes[start:]), prefix)
}

// EndsWith checks if string ends with suffix (like string.endsWith())
func (StringUtils) EndsWith(s, suffix string, length ...int) bool {
	endLength := utf8.RuneCountInString(s)
	if len(length) > 0 {
		endLength = length[0]
	}
	
	if endLength < 0 {
		endLength = 0
	}
	
	runes := []rune(s)
	if endLength > len(runes) {
		endLength = len(runes)
	}
	
	return strings.HasSuffix(string(runes[:endLength]), suffix)
}

// Includes checks if string contains substring (like string.includes())
func (StringUtils) Includes(s, substr string, position ...int) bool {
	start := 0
	if len(position) > 0 {
		start = position[0]
	}
	
	if start < 0 {
		start = 0
	}
	
	runes := []rune(s)
	if start >= len(runes) {
		return substr == ""
	}
	
	return strings.Contains(string(runes[start:]), substr)
}

// Repeat repeats string n times (like string.repeat())
func (StringUtils) Repeat(s string, count int) string {
	if count < 0 {
		panic("repeat count must be non-negative")
	}
	return strings.Repeat(s, count)
}

// Match matches against regex (like string.match())
func (StringUtils) Match(s, pattern string) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return []string{}, nil
	}
	
	return matches, nil
}

// MatchAll finds all matches (like string.matchAll())
func (StringUtils) MatchAll(s, pattern string) ([][]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	
	return re.FindAllStringSubmatch(s, -1), nil
}

// Search searches for regex match (like string.search())
func (StringUtils) Search(s, pattern string) int {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return -1
	}
	
	loc := re.FindStringIndex(s)
	if loc == nil {
		return -1
	}
	
	// Convert byte index to rune index
	runes := []rune(s)
	byteIndex := 0
	for i, r := range runes {
		if byteIndex == loc[0] {
			return i
		}
		byteIndex += len(string(r))
	}
	
	return -1
}

// LocaleCompare compares strings (simplified version of string.localeCompare())
func (StringUtils) LocaleCompare(s1, s2 string) int {
	if s1 < s2 {
		return -1
	}
	if s1 > s2 {
		return 1
	}
	return 0
}

// ToString converts value to string (like String() constructor)
func (StringUtils) ToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

// FromCharCode creates string from Unicode values (like String.fromCharCode())
func (StringUtils) FromCharCode(codes ...int) string {
	runes := make([]rune, len(codes))
	for i, code := range codes {
		runes[i] = rune(code)
	}
	return string(runes)
}

// IsEmpty checks if string is empty
func (StringUtils) IsEmpty(s string) bool {
	return s == ""
}

// IsBlank checks if string is empty or only whitespace
func (StringUtils) IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

// Capitalize capitalizes first letter
func (StringUtils) Capitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Uncapitalize uncapitalizes first letter
func (StringUtils) Uncapitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// Reverse reverses string
func (StringUtils) Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ToCamelCase converts string to camelCase
func (StringUtils) ToCamelCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	
	if len(words) == 0 {
		return ""
	}
	
	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		if len(words[i]) > 0 {
			result += strings.ToUpper(string(words[i][0])) + strings.ToLower(words[i][1:])
		}
	}
	
	return result
}

// ToPascalCase converts string to PascalCase
func (StringUtils) ToPascalCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	
	var result string
	for _, word := range words {
		if len(word) > 0 {
			result += strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	
	return result
}

// ToKebabCase converts string to kebab-case
func (StringUtils) ToKebabCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	
	var result []string
	for _, word := range words {
		if len(word) > 0 {
			result = append(result, strings.ToLower(word))
		}
	}
	
	return strings.Join(result, "-")
}

// ToSnakeCase converts string to snake_case
func (StringUtils) ToSnakeCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	
	var result []string
	for _, word := range words {
		if len(word) > 0 {
			result = append(result, strings.ToLower(word))
		}
	}
	
	return strings.Join(result, "_")
}

// ParseInt parses string to integer (like parseInt())
func (StringUtils) ParseInt(s string, base ...int) (int, error) {
	b := 10
	if len(base) > 0 {
		b = base[0]
	}
	result, err := strconv.ParseInt(s, b, 64)
	return int(result), err
}

// ParseFloat parses string to float (like parseFloat())
func (StringUtils) ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// IsNumeric checks if string represents a number
func (StringUtils) IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// WordCount counts words in string
func (StringUtils) WordCount(s string) int {
	return len(strings.Fields(s))
}

// Truncate truncates string to specified length with optional suffix
func (StringUtils) Truncate(s string, maxLength int, suffix ...string) string {
	runes := []rune(s)
	if len(runes) <= maxLength {
		return s
	}
	
	suf := "..."
	if len(suffix) > 0 {
		suf = suffix[0]
	}
	
	sufRunes := []rune(suf)
	if len(sufRunes) >= maxLength {
		return string(sufRunes[:maxLength])
	}
	
	truncateAt := maxLength - len(sufRunes)
	return string(runes[:truncateAt]) + suf
}