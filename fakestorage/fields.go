package fakestorage

import (
	"fmt"
	"regexp"
	"strings"
)

var valueByFieldMap = map[string]interface{}{
	"kind":            func(o Object) interface{} { return "storage#object" },
	"id":              func(o Object) interface{} { return o.id() },
	"name":            func(o Object) interface{} { return o.Name },
	"bucket":          func(o Object) interface{} { return o.BucketName },
	"size":            func(o Object) interface{} { return fmt.Sprintf("%d", len(o.Content)) },
	"contentType":     func(o Object) interface{} { return o.ContentType },
	"contentEncoding": func(o Object) interface{} { return o.ContentEncoding },
	"acl":             func(o Object) interface{} { return getAccessControlsListFromObject(o) },
	"crc32c":          func(o Object) interface{} { return o.Crc32c },
	"md5Hash":         func(o Object) interface{} { return o.Md5Hash },
	"timeCreated":     func(o Object) interface{} { return o.Created.Format(timestampFormat) },
	"timeDeleted":     func(o Object) interface{} { return o.Deleted.Format(timestampFormat) },
	"updated":         func(o Object) interface{} { return o.Updated.Format(timestampFormat) },
	"generation":      func(o Object) interface{} { return fmt.Sprintf("%d", o.Generation) },
	"metadata":        func(o Object) interface{} { return o.Metadata },
}

// there's not yet an implementation for nextPageToken
var ignoreableFields = []string{"", "nextPageToken"}
var itemFieldNamesAllowed = getMapKeys(valueByFieldMap)
var itemsPattern = regexp.MustCompile(`items\(([^)]*)\)`)

func ParseFields(input string) (*fieldsResult, error) {
	result := fieldsResult{Fields: []string{}, ItemFields: []string{}}
	input, itemFields, err := findItemFields(input)
	if err != nil {
		return &result, err
	}

	result.ItemFields = itemFields

	if input == "" {
		return &result, nil
	}

	fields := strings.Split(input, ",")
	for _, field := range fields {
		trimmedField := strings.Trim(field, " ")
		if isInSlice(ignoreableFields, trimmedField) {
			// just ignore it
		} else if trimmedField == "kind" || trimmedField == "prefixes" {
			result.Fields = append(result.Fields, trimmedField)
		} else if trimmedField == "items" {
			result.Fields = append(result.Fields, trimmedField)
			// if there's a "items" field, we ignore all the other specific "items(<field>)" fields
			result.ItemFields = []string{}
		} else {
			return &result, fmt.Errorf("%s is invalid", trimmedField)
		}
	}

	return &result, nil
}

func findItemFields(input string) (string, []string, error) {
	if itemsPattern.MatchString(input) {
		itemFields, err := processItemFields(input)
		if err != nil {
			return input, []string{}, err
		}

		fieldsWithNoItems := itemsPattern.ReplaceAll([]byte(input), []byte(""))
		return string(fieldsWithNoItems), itemFields, nil
	} else {
		return input, []string{}, nil
	}
}

func processItemFields(input string) ([]string, error) {
	itemFields := []string{}
	matches := itemsPattern.FindAllStringSubmatch(input, -1)
	if matches == nil {
		return []string{}, nil
	}

	for _, match := range matches {
		group := match[1]
		if group == "" {
			continue
		}

		wantedItemFields := strings.Split(group, ",")
		for _, wantedItemField := range wantedItemFields {
			trimmedField := strings.Trim(wantedItemField, " ")
			if isInSlice(itemFieldNamesAllowed, trimmedField) {
				itemFields = append(itemFields, trimmedField)
			} else {
				return nil, fmt.Errorf("%s is invalid", trimmedField)
			}
		}
	}

	return itemFields, nil
}

func (r *fieldsResult) GenerateResponse(prefixes []string, objects []Object) interface{} {
	if len(r.Fields) == 0 && len(r.ItemFields) == 0 {
		return newListObjectsResponse(objects, prefixes)
	}

	var response = map[string]interface{}{}
	if isInSlice(r.Fields, "kind") {
		response["kind"] = "storage#objects"
	}

	if isInSlice(r.Fields, "prefixes") {
		response["prefixes"] = prefixes
	}

	if isInSlice(r.Fields, "items") {
		items := make([]interface{}, len(objects))
		for i, obj := range objects {
			items[i] = newObjectResponse(obj)
		}
		response["items"] = items
	} else if len(r.ItemFields) > 0 {
		items := r.generateItems(objects)
		if len(items) > 0 {
			response["items"] = items
		}
	}

	return response
}

type fieldsResult struct {
	Fields     []string
	ItemFields []string
}

func (r *fieldsResult) generateItems(objects []Object) []map[string]interface{} {
	items := []map[string]interface{}{}
	for _, object := range objects {
		item := map[string]interface{}{}
		for _, itemField := range r.ItemFields {
			getValue := valueByFieldMap[itemField]
			item[itemField] = getValue.(func(Object) interface{})(object)
		}
		items = append(items, item)
	}

	return items
}

func isInSlice(slice []string, value string) bool {
	for _, element := range slice {
		if element == value {
			return true
		}
	}

	return false
}

func getMapKeys(myMap map[string]interface{}) []string {
	keys := make([]string, 0, len(myMap))
	for k := range myMap {
		keys = append(keys, k)
	}
	return keys
}
