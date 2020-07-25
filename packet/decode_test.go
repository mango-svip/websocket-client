package packet

import (
	"fmt"
	"testing"
)

func TestDecode(t *testing.T) {

	t.Log("testing")

	decodeString := "gqR0eXBlzSBzpGRhdGGBsW1pc3Npb25fdHlwZV9saXN0kQM="

	decode := DecodeString(decodeString)

	fmt.Println(string(decode))

}
