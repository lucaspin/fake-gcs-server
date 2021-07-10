package fakestorage

import (
	"fmt"
	"reflect"
	"testing"
)

type FieldTestCase struct {
	testCase           string
	fields             string
	expectedFields     []string
	expectedItemFields []string
	expectedError      error
}

func getAllFieldsTestCases() []FieldTestCase {
	return []FieldTestCase{
		{
			"kind, prefixes and full items",
			"kind,prefixes,items",
			[]string{"kind", "prefixes", "items"},
			[]string{},
			nil,
		},
		{
			"kind, prefixes and full items with spaces",
			"kind    ,    prefixes,  items  ",
			[]string{"kind", "prefixes", "items"},
			[]string{},
			nil,
		},
		{
			"kind and specific items field",
			"kind,items(name)",
			[]string{"kind"},
			[]string{"name"},
			nil,
		},
		{
			"empty items",
			"items()",
			[]string{},
			[]string{},
			nil,
		},
		{
			"only specific items fields",
			"items(name),items(bucket)",
			[]string{},
			[]string{"name", "bucket"},
			nil,
		},
		{
			"only specific items fields with spaces",
			"items(  name  ),items(  bucket  )",
			[]string{},
			[]string{"name", "bucket"},
			nil,
		},
		{
			"only specific items fields in items(<field1>,<field2>) format",
			"items(name,bucket)",
			[]string{},
			[]string{"name", "bucket"},
			nil,
		},
		{
			"only specific items fields in both formats",
			"items(name,bucket),items(size)",
			[]string{},
			[]string{"name", "bucket", "size"},
			nil,
		},
		{
			"only invalid field",
			"invalid-field",
			[]string{},
			[]string{},
			fmt.Errorf("invalid-field is invalid"),
		},
		{
			"valid field and invalid field",
			"kind,invalid-field",
			[]string{"kind"},
			[]string{},
			fmt.Errorf("invalid-field is invalid"),
		},
		{
			"only invalid item field",
			"items(invalid-field)",
			[]string{},
			[]string{},
			fmt.Errorf("invalid-field is invalid"),
		},
		{
			"valid and invalid item field",
			"items(name,invalid-field)",
			[]string{},
			[]string{},
			fmt.Errorf("invalid-field is invalid"),
		},
	}
}

func TestFieldsParsing(t *testing.T) {
	testCases := getAllFieldsTestCases()
	for _, testCase := range testCases {
		t.Run(testCase.testCase, func(t *testing.T) {
			result, err := ParseFields(testCase.fields)
			if testCase.expectedError == nil && err != nil {
				t.Errorf("expected no error\ngot  %#v", err)
			} else if testCase.expectedError != nil && err == nil {
				t.Errorf("expected error %#v\ngot %#v", testCase.expectedError, err)
			}

			if !reflect.DeepEqual(result.Fields, testCase.expectedFields) {
				t.Errorf("wrong fields returned\nwant %#v\ngot  %#v", testCase.expectedFields, result.Fields)
			}

			if !reflect.DeepEqual(result.ItemFields, testCase.expectedItemFields) {
				t.Errorf("wrong item fields returned\nwant %#v\ngot  %#v", testCase.expectedItemFields, result.ItemFields)
			}
		})
	}
}
