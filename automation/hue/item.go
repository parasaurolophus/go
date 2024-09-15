// Copyright 2024 Kirk Rader

package hue

// Alias for map[string]any used as the basic data model for the Hue Bridge API
// V2.
type Item map[string]any

// A few predefined keys for Item maps. There are countless other keys present
// in Hue Bridge payloads used for device-specific types of items and for the
// maps they contain.
const (

	// Keys common to all Items.
	Id   = "id"
	IdV1 = "id_v1"
	Type = "type"

	// A "metadata" field in some types of Items. When present, it contains a
	// "name" field.
	Metadata = "metadata"

	// Metadata's "name" field.
	Name = "name"

	// An "owner" field is present in some types of Items. It is used for
	// cross-referencing items in the massive and bizarrely designed data
	// structure returned by the /resource endpoint.
	Owner = "owner"

	// Keys for the map that is the value of "owner" when present.
	Rid   = "rid"
	Rtype = "rtype"
)
