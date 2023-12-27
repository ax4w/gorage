package Gorage

import (
	"fmt"
	"hash/fnv"
	"os"
)

func gprint(a, s string) {
	fmt.Printf("[Gorage](%s) - %s\n", a, s)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func computeHash(data []interface{}) uint32 {
	var h string
	for _, v := range data {
		switch v.(type) {
		case string:
			h += v.(string)
			break
		case int:
			h += fmt.Sprintf("%d", v)
			break
		case float64:
			h += fmt.Sprintf("%d", int(v.(float64)))
			break
		case float32:
			h += fmt.Sprintf("%d", int(v.(float32)))
			break
		case bool:
			if v.(bool) {
				h += "True"
			} else {
				h += "False"
			}
		}
	}
	ha := fnv.New32a()
	ha.Write([]byte(h))
	return ha.Sum32()
}

func compareRows(a, b []interface{}) bool {
	return computeHash(a) == computeHash(b)
}

func validateDatatype(is interface{}, c Column) bool {
	switch is.(type) {
	case int:
		if c.Datatype != INT {
			return false
		}
	case string:
		if c.Datatype != STRING {
			return false
		}
	case bool:
		if c.Datatype != BOOLEAN {
			return false
		}
	case float64:
		if c.Datatype != FLOAT {
			return false
		}
	case float32:
		if c.Datatype != FLOAT {
			return false
		}
	}
	return true

}
