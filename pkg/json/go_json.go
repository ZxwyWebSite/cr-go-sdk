//go:build !stdjson

package json

import json "github.com/goccy/go-json"

var Marshal = json.MarshalNoEscape

var Unmarshal = json.UnmarshalNoEscape

var MarshalIndent = json.MarshalIndent

var NewDecoder = json.NewDecoder

var NewEncoder = json.NewEncoder
