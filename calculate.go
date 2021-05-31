package jsonhash

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

var TRUE = sha256.Sum256([]byte("true"))
var FALSE = sha256.Sum256([]byte("false"))
var NULL = sha256.Sum256([]byte("null"))
var BLANK = sha256.Sum256([]byte(""))

// вариант со структурой, чтобы не пересоздавать мапу каждый раз при многократном использовании
type HashCalculator struct {
	mapExclude map[string]bool
}

func NewHashCalculator(exclude []string) (calc HashCalculator) {
	calc.mapExclude = make(map[string]bool, len(exclude))
	for i := range exclude {
		calc.mapExclude[exclude[i]] = true
	}
	return
}

func (h *HashCalculator) CalculateHash(in interface{}) [32]byte {
	return calculateHash(in, &(h.mapExclude), "")
}

func (h *HashCalculator) CalculateJsonHash(in []byte) ([32]byte, error) {
	var interfaceData interface{}
	err := json.Unmarshal(in, &interfaceData)
	return h.CalculateHash(interfaceData), err
}

// Подсчёт хеша json-объекта без учёта порядка в слайсах и мапах
// hash([1,2,3]) == hash([3,1,2])
func CalculateJsonHash(in []byte, exclude []string) ([32]byte, error) {
	var interfaceData interface{}
	err := json.Unmarshal(in, &interfaceData)
	return CalculateHash(interfaceData, exclude), err
}

func CalculateHash(in interface{}, exclude []string) [32]byte {
	excludeNames := make(map[string]bool, len(exclude))
	for i := range exclude {
		excludeNames[exclude[i]] = true
	}
	return calculateHash(in, &excludeNames, "")
}

func calculateHash(in interface{}, exclude *map[string]bool, fullFieldName string) [32]byte {
	if (*exclude)[fullFieldName] {
		return BLANK
	}
	if in == nil {
		return NULL
	}
	valIn := reflect.Indirect(reflect.ValueOf(in))
	switch valIn.Kind() {
	case reflect.String:
		return sha256.Sum256([]byte(valIn.String()))
	case reflect.Int:
		return sha256.Sum256([]byte(strconv.Itoa(int(valIn.Int()))))
	case reflect.Bool:
		if valIn.Bool() {
			return TRUE
		} else {
			return FALSE
		}
	case reflect.Float64:
		return sha256.Sum256(float64ToByte(valIn.Float()))
	case reflect.Slice:
		resStrings := make([]string, valIn.Len())
		for i := 0; i < valIn.Len(); i++ {
			hash := calculateHash(valIn.Index(i).Interface(), exclude, fullFieldName)
			resStrings[i] = hex.EncodeToString(hash[:])
		}
		sort.Strings(resStrings)
		return calculateHash(strings.Join(resStrings, "|"), exclude, fullFieldName)
	case reflect.Map:
		var resStrings []string
		iter := valIn.MapRange()
		for iter.Next() {
			newFieldName := fullFieldName + "." + fmt.Sprint(iter.Key().Interface())
			if (*exclude)[newFieldName] { // мелкая оптимизация
				continue
			}
			keyHash := calculateHash(iter.Key().Interface(), exclude, newFieldName)
			valueHash := calculateHash(iter.Value().Interface(), exclude, newFieldName)
			resStrings = append(resStrings, strings.Join([]string{
				hex.EncodeToString(keyHash[:]),
				hex.EncodeToString(valueHash[:]),
			}, ":"))
		}
		// panic(...)  этот код не должен выполниться при подаче на вход json
		sort.Strings(resStrings)
		return calculateHash(strings.Join(resStrings, "|"), exclude, fullFieldName)
	}

	jsonBytes, _ := json.Marshal(valIn)
	return sha256.Sum256(jsonBytes)
}

func float64ToByte(f float64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}
