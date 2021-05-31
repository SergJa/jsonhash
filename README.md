# Json elective hash 

Calculates hash **sha256** of JSON objects with ignoring some fields option 
and array order unimportance. It can work with common golang objects, 
but correctness is not guaranteed

## Usage

The most simple variant

```go
import "github.com/SergJa/jsonhash"

func main() {
	jsonContent, err := ioutil.ReadFile("file.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	hash, err := jsonhash.CalculateJsonHash(jsonContent, []string{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(hex.EncodeToString(hash[:]))
}
```

You may use already unmarshalled object

```go
hash, err := jsonhash.CalculateHash(unmarshalledObject, []string{})
```

If you need to ignore some fields, you can fill the second parameter. 
There you can describe fields using dot (".") as separator and leader of each string.
Iportant! Every string _must_ begin with dot. Example: `[]string{".uuid"}` - 
field with name _uuid_ will be ignored while hash calculaion.

Example 2:
```go 
[]string{".subobj.subfield"}
```

```json
{
 "field1": "data1",
 "subobj": {
  "subfield": "will be ignored",
  "else_field": "will take part in hash calculation"
 }
}
```
Arrays do not consider level. Example 3:

```go 
[]string{".subobj.subfield"}
```

```json
{
 "subarr": [
  {
   "subfield": "will be ignored",
   "else_field": "will take part in hash calculation"
  },
  {
   "subfield": "will be ignored",
   "else_field2": "will take part in hash calculation"
  }
 ]
}
```

## Object hashCalculator

If you need to count plenty of hashes with same ignoring fields list, 
you can use `hashCalculator`. Creating: 

```calc := jsonhash.NewHashCalculator([]string{".field.subfield"})```

This `calc` has just the same methods as described below, 
excluding second parameter whith ignoring fields list. This list will be taken
from `calc`. Example:

```go
calc := jsonhash.NewHashCalculator([]string{".ignoreme", ".ignoreme_root.ignore"})

hash, err := calc.CalculateJsonHash([]byte(`{"json":content", "ignoreme":"anything"}`))

hash2, err := calc.CalculateJsonHash([]byte(`{"ignoreme_root":{"ignore": "anything"}`))
```