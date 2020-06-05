package main

import (
	"encoding/json"
	"os"
)

func F00() {
	enc := json.NewEncoder(os.Stdout)
	d := map[string]int{"apple": 5, "lettuce": 7}
	enc.Encode(d)
}
