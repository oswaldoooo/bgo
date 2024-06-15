package bgo

import (
	"encoding/json"

	"github.com/oswaldoooo/bgo/types"
)

func Debug(src types.Packages) []byte {
	content, _ := json.MarshalIndent(src, "", "    ")
	return content
}
