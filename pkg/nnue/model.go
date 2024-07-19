package nnue

import (
	"embed"
	"encoding/json"
)

type model struct {
	L0 linearLayer `json:"l0"`
	L1 linearLayer `json:"l1"`
	L2 linearLayer `json:"l2"`
}

const (
	_NUM_FEATURES = 81920
	_M            = 4
	_N            = 8
	_K            = 1
)

//go:embed model.json
var fs embed.FS

var m model

func init() {
	f, err := fs.Open("model.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	d := json.NewDecoder(f)
	err = d.Decode(&m)
	if err != nil {
		panic(err)
	}
}
