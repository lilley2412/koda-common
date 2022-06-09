package util

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func JSONMustPrettyPrint(in interface{}) {
	d, err := json.Marshal(in)
	if err != nil {
		fmt.Printf("error marshaling json: %s\n", err)
		return
	}

	var buf bytes.Buffer
	err = json.Indent(&buf, d, "", "    ")
	if err != nil {
		fmt.Printf("error indenting json: %s\n", err)
		return
	}

	fmt.Println(buf.String())
}
