package packet

import (
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	t.Log("testing....")

	//{"data":{"mission_type_list":[3]},"type":8307}
	//gqR0eXBlzSBzpGRhdGGBsW1pc3Npb25fdHlwZV9saXN0kQM=
	//gqR0eXBl0gAAIHOkZGF0YYGxbWlzc2lvbl90eXBlX2xpc3SR0gAAAAM=

	encode := Encode(`{"type":8307,"data":{"mission_type_list":[3]}}`)

	fmt.Println(string(encode))
}
