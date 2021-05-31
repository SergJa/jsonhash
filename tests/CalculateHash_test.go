package hash_tests

import (
	"encoding/hex"
	"fmt"
	"github.com/SergJa/jsonhash"
	"testing"
)

func TestCalculateHash(t *testing.T) {
	json1Hash := [32]byte{200, 169, 8, 84, 249, 249, 95, 40, 29, 19, 135, 67, 144, 100, 156, 112, 217, 179, 55, 91, 152, 212, 224, 32, 28, 106, 174, 223, 119, 46, 180, 131}
	jsonArrHash := [32]byte{215, 81, 102, 200, 184, 151, 26, 164, 100, 225, 117, 235, 127, 132, 24, 11, 70, 16, 186, 162, 244, 28, 8, 131, 125, 75, 67, 162, 108, 56, 18, 224}

	json1 := `
{
 "digit":1,
 "string":"one",
 "object": 
   {
     "sub":"two"
   },
 "array": [1, "two", 3.14151926, null],
 "nil":null
}
`

	json1_1 := `
{
 "string":"one",
 "digit":1,
 "object": 
   {
     "sub":"two"
   },
 "array": [1, "two", null, 3.14151926],
 "nil":null
}
`

	jsonArr := `
[
{
 "string":"one",
 "digit":1,
 "object":  {
     "sub":"two"
   },
 "array": [1, "two", null, 3.14151926],
 "nil":null
},
10,
"SomeString",
[null]
]
`

	jsonIgnoreField := `
{
 "field1":"data",
 "ignoreme": "ignoring data",
 "ignoreme_subfield": {
   "matters":"data",
   "ignore":"ignoring data"
  },
 "ignore_subarr": [
  "matters",
  {
   "matters_subobject": "data",
   "ignore_subobject": "ignoring data"
  }
 ]
}
`

	jsonIgnoreFieldEq := `
{
 "field1":"data",
 "ignoreme": "ignoring data another",
 "ignoreme_subfield": {
   "matters":"data",
   "ignore":"ignoring data another"
  },
 "ignore_subarr": [
  "matters",
  {
   "matters_subobject": "data",
   "ignore_subobject": "ignoring data another"
  }
 ]
}
`

	jsonIgnoreFieldNEq1 := `
{
 "field1":"data",
 "ignoreme": "ignoring data another",
 "ignoreme_subfield": {
   "matters":"data not eq",
   "ignore":"ignoring data another"
  },
 "ignore_subarr": [
  "matters",
  {
   "matters_subobject": "data",
   "ignore_subobject": "ignoring data another"
  }
 ]
}
`

	jsonIgnoreFieldNEq2 := `
{
 "field1":"data",
 "ignoreme": "ignoring data another",
 "ignoreme_subfield": {
   "matters":"data",
   "ignore":"ignoring data another"
  },
 "ignore_subarr": [
  "matters",
  {
   "matters_subobject": "data not eq",
   "ignore_subobject": "ignoring data another"
  }
 ]
}
`
	hasherCommon := jsonhash.NewHashCalculator([]string{})

	hash1, err := hasherCommon.CalculateJsonHash([]byte(json1))
	if err != nil {
		t.Error(err)
	}
	hash1_1, err := hasherCommon.CalculateJsonHash([]byte(json1_1))
	if err != nil {
		t.Error(err)
	}
	hashArr, err := hasherCommon.CalculateJsonHash([]byte(jsonArr))
	if err != nil {
		t.Error(err)
	}

	hasherWithFields := jsonhash.NewHashCalculator([]string{".ignoreme", ".ignoreme_subfield.ignore",
		".ignore_subarr.ignore_subobject"})

	hashIgnore, err := hasherWithFields.CalculateJsonHash([]byte(jsonIgnoreField))
	if err != nil {
		t.Error(err)
	}
	hashIgnoreEq, err := hasherWithFields.CalculateJsonHash([]byte(jsonIgnoreFieldEq))
	if err != nil {
		t.Error(err)
	}
	hashIgnoreNEq1, err := hasherWithFields.CalculateJsonHash([]byte(jsonIgnoreFieldNEq1))
	if err != nil {
		t.Error(err)
	}
	hashIgnoreNEq2, err := hasherWithFields.CalculateJsonHash([]byte(jsonIgnoreFieldNEq2))
	if err != nil {
		t.Error(err)
	}

	if hash1 != json1Hash {
		t.Error(fmt.Sprintf("First is %s, but must be hash must be %s",
			hex.EncodeToString(json1Hash[:]), hex.EncodeToString(hash1[:])))
	}
	if hash1 != hash1_1 {
		t.Error("hash1 and hash1_1 are not equivalent")
	}
	if hashArr != jsonArrHash {
		t.Error("hashArr and jsonHashArr are not equivalent")
	}

	if hashIgnore != hashIgnoreEq {
		t.Error("hashIgnore != hashIgnoreEq")
	}
	if hashIgnore == hashIgnoreNEq1 {
		t.Error("hashIgnore == hashIgnoreNEq1")
	}
	if hashIgnore == hashIgnoreNEq2 {
		t.Error("hashIgnore == hashIgnoreNEq2")
	}
}
