package fakestorage

import (
	"testing"
)

// items
// items()
// items(name)
// items(name,bucket)
// items(name),items(bucket)
// kind,items(name,bucket)

var objects = []Object{
	{BucketName: "bucket1", Name: "file1.txt"},
}

func Test__AllFields(t *testing.T) {
	result, _ := ProcessFields("kind,prefixes,items")
	if len(result.Fields) != 3 {
		t.Errorf("Wrong result")
	}

	if !isInSlice(result.Fields, "kind") {
		t.Errorf("Wrong result")
	}

	if !isInSlice(result.Fields, "prefixes") {
		t.Errorf("Wrong result")
	}

	if !isInSlice(result.Fields, "items") {
		t.Errorf("Wrong result")
	}

	result.GenerateResponse([]string{}, objects)
}

func Test__SpecificItemsFields(t *testing.T) {
	result, _ := ProcessFields("items(name),items(bucket)")
	if len(result.Fields) != 0 {
		t.Errorf("Wrong result")
	}

	if len(result.ItemFields) != 2 {
		t.Errorf("Wrong result")
	}

	if !isInSlice(result.ItemFields, "name") {
		t.Errorf("Wrong result")
	}

	if !isInSlice(result.ItemFields, "bucket") {
		t.Errorf("Wrong result")
	}

	result.GenerateResponse([]string{}, objects)
}

func Test__ItemsAloneOverridesSpecificItemsFields(t *testing.T) {
	result, _ := ProcessFields("items,items(name),items(bucket)")
	if len(result.Fields) != 1 {
		t.Errorf("Wrong result")
	}

	if len(result.ItemFields) != 0 {
		t.Errorf("Wrong result: %v", result.ItemFields)
	}

	if !isInSlice(result.Fields, "items") {
		t.Errorf("Wrong result")
	}
}
