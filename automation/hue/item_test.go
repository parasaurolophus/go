// Copyright 2024 Kirk Rader

package hue_test

import (
	"encoding/json"
	"parasaurolophus/automation/hue"
	"testing"
)

func TestId(t *testing.T) {

	// happy path
	s := `{
		"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
		"id_v1": "/sensors/40",
		"light": {
			"light_level": 7457,
			"light_level_report": {
			"changed": "2024-09-13T19:20:52.121Z",
			"light_level": 7457
			},
			"light_level_valid": true
		},
		"owner": {
			"rid": "ce065303-2711-4689-b488-6ef835afada4",
			"rtype": "device"
		},
		"type": "light_level"
		}`

	var item1 hue.Item

	if err := json.Unmarshal([]byte(s), &item1); err != nil {
		t.Error(err.Error())
	} else if id, err := item1.Id(); err != nil {
		t.Error(err.Error())
	} else if id != "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5" {
		t.Errorf(`expected "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5" but got "%s"`, id)
	}

	// missing id
	s = `{
		"id_v1": "/sensors/40",
		"light": {
			"light_level": 7457,
			"light_level_report": {
			"changed": "2024-09-13T19:20:52.121Z",
			"light_level": 7457
			},
			"light_level_valid": true
		},
		"owner": {
			"rid": "ce065303-2711-4689-b488-6ef835afada4",
			"rtype": "device"
		},
		"type": "light_level"
		}`

	var item2 hue.Item

	if err := json.Unmarshal([]byte(s), &item2); err != nil {
		t.Error(err.Error())
	} else if _, err := item2.Id(); err == nil {
		t.Error("expected missing id error")
	}

	// invalid id
	s = `{
		"id": 42,
		"id_v1": "/sensors/40",
		"light": {
			"light_level": 7457,
			"light_level_report": {
			"changed": "2024-09-13T19:20:52.121Z",
			"light_level": 7457
			},
			"light_level_valid": true
		},
		"owner": {
			"rid": "ce065303-2711-4689-b488-6ef835afada4",
			"rtype": "device"
		},
		"type": "light_level"
		}`

	var item3 hue.Item

	if err := json.Unmarshal([]byte(s), &item3); err != nil {
		t.Error(err.Error())
	} else if _, err := item3.Id(); err == nil {
		t.Error("expected invalid id error")
	}
}

func TestIdV1(t *testing.T) {

	// happy path
	s := `{
		"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
		"id_v1": "/sensors/40",
		"light": {
			"light_level": 7457,
			"light_level_report": {
			"changed": "2024-09-13T19:20:52.121Z",
			"light_level": 7457
			},
			"light_level_valid": true
		},
		"owner": {
			"rid": "ce065303-2711-4689-b488-6ef835afada4",
			"rtype": "device"
		},
		"type": "light_level"
		}`

	var item1 hue.Item

	if err := json.Unmarshal([]byte(s), &item1); err != nil {
		t.Error(err.Error())
	} else if idV1, err := item1.IdV1(); err != nil {
		t.Error(err.Error())
	} else if idV1 != "/sensors/40" {
		t.Errorf(`expected "/sensors/40" but got "%s"`, idV1)
	}

	// missing id
	s = `{
		"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
		"light": {
			"light_level": 7457,
			"light_level_report": {
			"changed": "2024-09-13T19:20:52.121Z",
			"light_level": 7457
			},
			"light_level_valid": true
		},
		"owner": {
			"rid": "ce065303-2711-4689-b488-6ef835afada4",
			"rtype": "device"
		},
		"type": "light_level"
		}`

	var item2 hue.Item

	if err := json.Unmarshal([]byte(s), &item2); err != nil {
		t.Error(err.Error())
	} else if _, err := item2.IdV1(); err == nil {
		t.Error("expected missing id_v2 error")
	}

	// invalid id
	s = `{
		"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
		"id_v1": 42,
		"light": {
			"light_level": 7457,
			"light_level_report": {
			"changed": "2024-09-13T19:20:52.121Z",
			"light_level": 7457
			},
			"light_level_valid": true
		},
		"owner": {
			"rid": "ce065303-2711-4689-b488-6ef835afada4",
			"rtype": "device"
		},
		"type": "light_level"
		}`

	var item3 hue.Item

	if err := json.Unmarshal([]byte(s), &item3); err != nil {
		t.Error(err.Error())
	} else if _, err := item3.IdV1(); err == nil {
		t.Error("expected invalid id error")
	}
}

func TestOwner(t *testing.T) {

	// happy path
	s := `{
		"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
		"id_v1": "/sensors/40",
		"light": {
			"light_level": 7457,
			"light_level_report": {
			"changed": "2024-09-13T19:20:52.121Z",
			"light_level": 7457
			},
			"light_level_valid": true
		},
		"owner": {
			"rid": "ce065303-2711-4689-b488-6ef835afada4",
			"rtype": "device"
		},
		"type": "light_level"
		}`

	var item1 hue.Item

	err := json.Unmarshal([]byte(s), &item1)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = item1.Owner()

	if err != nil {
		t.Error(err.Error())
	}

	// missing "owner"
	s = `{
			"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
			"id_v1": "/sensors/40",
			"light": {
				"light_level": 7457,
				"light_level_report": {
				"changed": "2024-09-13T19:20:52.121Z",
				"light_level": 7457
				},
				"light_level_valid": true
			},
			"type": "light_level"
			}`

	var item2 hue.Item

	err = json.Unmarshal([]byte(s), &item2)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = item2.Owner()
	if err == nil {
		t.Error("expected an error")
	}

	// invalid "owner"
	s = `{
			"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
			"id_v1": "/sensors/40",
			"light": {
				"light_level": 7457,
				"light_level_report": {
				"changed": "2024-09-13T19:20:52.121Z",
				"light_level": 7457
				},
				"light_level_valid": true
			},
			"owner": "a string",
			"type": "light_level"
			}`

	var item3 hue.Item

	err = json.Unmarshal([]byte(s), &item3)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = item3.Owner()
	if err == nil {
		t.Error("expected an error")
	}
}

func TestOwnerRid(t *testing.T) {

	// happy path
	s := `{
		"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
		"id_v1": "/sensors/40",
		"light": {
			"light_level": 7457,
			"light_level_report": {
			"changed": "2024-09-13T19:20:52.121Z",
			"light_level": 7457
			},
			"light_level_valid": true
		},
		"owner": {
			"rid": "ce065303-2711-4689-b488-6ef835afada4",
			"rtype": "device"
		},
		"type": "light_level"
		}`

	var item1 hue.Item

	err := json.Unmarshal([]byte(s), &item1)
	if err != nil {
		t.Error(err.Error())
	}

	rid, err := item1.OwnerRid()
	if err != nil {
		t.Error(err.Error())
	}

	if rid != "ce065303-2711-4689-b488-6ef835afada4" {
		t.Errorf(`expected "ce065303-2711-4689-b488-6ef835afada4" but got "%s"`, rid)
	}

	// missing rid
	s = `{
		"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
		"id_v1": "/sensors/40",
		"light": {
			"light_level": 7457,
			"light_level_report": {
			"changed": "2024-09-13T19:20:52.121Z",
			"light_level": 7457
			},
			"light_level_valid": true
		},
		"owner": {
			"rtype": "device"
		},
		"type": "light_level"
		}`

	var item2 hue.Item

	err = json.Unmarshal([]byte(s), &item2)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = item2.OwnerRid()
	if err == nil {
		t.Error("expected an error")
	}

	// invalid rid
	s = `{
			"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
			"id_v1": "/sensors/40",
			"light": {
				"light_level": 7457,
				"light_level_report": {
				"changed": "2024-09-13T19:20:52.121Z",
				"light_level": 7457
				},
				"light_level_valid": true
			},
			"owner": {
				"rid": 1,
				"rtype": "device"
			},
			"type": "light_level"
			}`

	var item3 hue.Item

	err = json.Unmarshal([]byte(s), &item3)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = item3.OwnerRid()
	if err == nil {
		t.Error("expected an error")
	}

	// missing owner
	s = `{
			"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
			"id_v1": "/sensors/40",
			"light": {
				"light_level": 7457,
				"light_level_report": {
				"changed": "2024-09-13T19:20:52.121Z",
				"light_level": 7457
				},
				"light_level_valid": true
			},
			"type": "light_level"
			}`

	var item4 hue.Item

	err = json.Unmarshal([]byte(s), &item4)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = item4.OwnerRid()
	if err == nil {
		t.Error("expected an error")
	}
}

func TestOwnerRtype(t *testing.T) {

	// happy path
	s := `{
		"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
		"id_v1": "/sensors/40",
		"light": {
			"light_level": 7457,
			"light_level_report": {
			"changed": "2024-09-13T19:20:52.121Z",
			"light_level": 7457
			},
			"light_level_valid": true
		},
		"owner": {
			"rid": "ce065303-2711-4689-b488-6ef835afada4",
			"rtype": "device"
		},
		"type": "light_level"
		}`

	var item1 hue.Item

	err := json.Unmarshal([]byte(s), &item1)
	if err != nil {
		t.Error(err.Error())
	}

	rtype, err := item1.OwnerType()
	if err != nil {
		t.Error(err.Error())
	}

	if rtype != "device" {
		t.Errorf(`expected "device" but got "%s"`, rtype)
	}

	// missing rtype
	s = `{
			"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
			"id_v1": "/sensors/40",
			"light": {
				"light_level": 7457,
				"light_level_report": {
				"changed": "2024-09-13T19:20:52.121Z",
				"light_level": 7457
				},
				"light_level_valid": true
			},
			"owner": {
				"rid": "ce065303-2711-4689-b488-6ef835afada4"
			},
			"type": "light_level"
			}`

	var item2 hue.Item

	err = json.Unmarshal([]byte(s), &item2)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = item2.OwnerType()
	if err == nil {
		t.Error("expected an error")
	}

	// invalid rtype
	s = `{
			"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
			"id_v1": "/sensors/40",
			"light": {
				"light_level": 7457,
				"light_level_report": {
				"changed": "2024-09-13T19:20:52.121Z",
				"light_level": 7457
				},
				"light_level_valid": true
			},
			"owner": {
				"rid": "ce065303-2711-4689-b488-6ef835afada4",
				"rtype": false
			},
			"type": "light_level"
			}`

	var item3 hue.Item

	err = json.Unmarshal([]byte(s), &item3)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = item3.OwnerType()
	if err == nil {
		t.Error("expected an error")
	}

	// missing owner
	s = `{
			"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
			"id_v1": "/sensors/40",
			"light": {
				"light_level": 7457,
				"light_level_report": {
				"changed": "2024-09-13T19:20:52.121Z",
				"light_level": 7457
				},
				"light_level_valid": true
			},
			"type": "light_level"
			}`

	var item4 hue.Item

	err = json.Unmarshal([]byte(s), &item4)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = item4.OwnerType()
	if err == nil {
		t.Error("expected an error")
	}
}

func TestType(t *testing.T) {

	// happy path
	s := `{
		"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
		"id_v1": "/sensors/40",
		"light": {
			"light_level": 7457,
			"light_level_report": {
			"changed": "2024-09-13T19:20:52.121Z",
			"light_level": 7457
			},
			"light_level_valid": true
		},
		"owner": {
			"rid": "ce065303-2711-4689-b488-6ef835afada4",
			"rtype": "device"
		},
		"type": "light_level"
		}`

	var item1 hue.Item

	if err := json.Unmarshal([]byte(s), &item1); err != nil {
		t.Error(err.Error())
	} else if typ, err := item1.Type(); err != nil {
		t.Error(err.Error())
	} else if typ != "light_level" {
		t.Errorf(`expected "light_level" but got "%s"`, typ)
	}

	// missing type
	s = `{
		"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
		"id_v1": "/sensors/40",
		"light": {
			"light_level": 7457,
			"light_level_report": {
			"changed": "2024-09-13T19:20:52.121Z",
			"light_level": 7457
			},
			"light_level_valid": true
		},
		"owner": {
			"rid": "ce065303-2711-4689-b488-6ef835afada4",
			"rtype": "device"
		}
		}`

	var item2 hue.Item

	if err := json.Unmarshal([]byte(s), &item2); err != nil {
		t.Error(err.Error())
	} else if _, err := item2.Type(); err == nil {
		t.Error("expected missing id error")
	}

	// invalid type
	s = `{
		"id": "81dba98d-5e53-4e9c-9ce4-9cbd3e519eb5",
		"id_v1": "/sensors/40",
		"light": {
			"light_level": 7457,
			"light_level_report": {
			"changed": "2024-09-13T19:20:52.121Z",
			"light_level": 7457
			},
			"light_level_valid": true
		},
		"owner": {
			"rid": "ce065303-2711-4689-b488-6ef835afada4",
			"rtype": "device"
		},
		"type": 42
		}`

	var item3 hue.Item

	if err := json.Unmarshal([]byte(s), &item3); err != nil {
		t.Error(err.Error())
	} else if _, err := item3.Type(); err == nil {
		t.Error("expected invalid id error")
	}
}
