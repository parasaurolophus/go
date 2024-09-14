// Copyright 2024 Kirk Rader

package hue

import "fmt"

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

	// An "owner" is present in most, but not quite all, types of Items. It is
	// used for cross-referencing items in the massive and bizarrely designed
	// data structure returned by the /resource endpoint.
	Owner = "owner"

	// Keys for the map that is the value of "owner" when present.
	Rid   = "rid"
	Rtype = "rtype"
)

// Getter for item["id"].
func (item Item) Id() (id string, err error) {

	var (
		a  any
		ok bool
	)

	if a, ok = item[Id]; !ok {
		err = fmt.Errorf("missing id")
	} else if id, ok = a.(string); !ok {
		err = fmt.Errorf("id %v expected to be a string but got %T", a, a)
	}

	return
}

// Getter for item["id_v1"].
func (item Item) IdV1() (idV1 string, err error) {

	var (
		a  any
		ok bool
	)

	if a, ok = item[IdV1]; !ok {
		err = fmt.Errorf("missing id_v1")
	} else if idV1, ok = a.(string); !ok {
		err = fmt.Errorf("id_v1 %v expected to be a string but got %T", a, a)
	}

	return
}

// Getter for item["owner"].
func (item Item) Owner() (owner map[string]any, err error) {

	var (
		o  any
		ok bool
	)

	if o, ok = item[Owner]; !ok {
		err = fmt.Errorf("missing owner")
	} else if owner, ok = o.(map[string]any); !ok {
		err = fmt.Errorf("owner, of type %T, is not a map", o)
	}

	return
}

// Getter for item["owner"]["rid"].
func (msg Item) OwnerRid() (rid string, err error) {

	var (
		owner map[string]any
		i     any
		ok    bool
	)

	if owner, err = msg.Owner(); err == nil {
		if i, ok = owner[Rid]; !ok {
			err = fmt.Errorf("missing rid")
		} else if rid, ok = i.(string); !ok {
			err = fmt.Errorf("rid, of type %T, is not a string", i)
		}
	}

	return
}

// Getter for item["owner"]["rtype"].
func (msg Item) OwnerType() (rtype string, err error) {

	var (
		owner map[string]any
		t     any
		ok    bool
	)

	if owner, err = msg.Owner(); err == nil {
		if t, ok = owner[Rtype]; !ok {
			err = fmt.Errorf("missing rtype")
		} else if rtype, ok = t.(string); !ok {
			err = fmt.Errorf("rtype, of type %T, is not a string", t)
		}
	}

	return
}

// Getter for item["type"].
func (item Item) Type() (typ string, err error) {

	var (
		a  any
		ok bool
	)

	if a, ok = item[Type]; !ok {
		err = fmt.Errorf("missing type")
	} else if typ, ok = a.(string); !ok {
		err = fmt.Errorf("type %v expected to be a string but got %T", a, a)
	}

	return
}
