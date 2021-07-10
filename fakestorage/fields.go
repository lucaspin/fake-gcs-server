package fakestorage

import (
	"fmt"
	"regexp"
	"strings"
)

var jsonFieldNameToStructFieldName = map[string]string{
	"name":   "Name",
	"bucket": "Bucket",
}

var itemFieldNamesAllowed = getMapKeys(jsonFieldNameToStructFieldName)
var itemsPattern = regexp.MustCompile(`items\((.*)\)`)

func getMapKeys(myMap map[string]string) []string {
	keys := make([]string, 0, len(myMap))
	for k := range myMap {
		keys = append(keys, k)
	}
	return keys
}

func ProcessFields(input string) (*fieldsParsingResult, error) {
	fields := strings.Split(input, ",")
	result := fieldsParsingResult{Fields: []string{}, ItemFields: []string{}}

	if len(fields) == 0 {
		return &result, nil
	}

	for _, field := range fields {
		if field == "nextPageToken" {
			// ignore
		} else if field == "kind" || field == "prefixes" {
			result.Fields = append(result.Fields, field)
		} else if field == "items" {
			result.Fields = append(result.Fields, field)
			result.ItemFields = []string{}
		} else {
			if itemsPattern.MatchString(field) {
				if !isInSlice(result.Fields, "items") {
					processItems(&result, field)
				}
			} else {
				return nil, fmt.Errorf("%s is invalid", field)
			}
		}
	}

	return &result, nil
}

func processItems(result *fieldsParsingResult, field string) error {
	match := itemsPattern.FindStringSubmatch(field)
	itemsFields := match[1]
	if itemsFields != "" {
		wantedItemFields := strings.Split(itemsFields, ",")
		for _, wantedItemField := range wantedItemFields {
			if isInSlice(itemFieldNamesAllowed, wantedItemField) {
				result.ItemFields = append(result.ItemFields, wantedItemField)
			} else {
				return fmt.Errorf("%s is invalid", field)
			}
		}
	}

	return nil
}

type fieldsParsingResult struct {
	Fields     []string
	ItemFields []string
}

func (r *fieldsParsingResult) GenerateResponse(prefixes []string, objects []Object) interface{} {
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
		response["items"] = objects
	} else {
		items := r.generateItems(objects)
		if len(items) > 0 {
			response["items"] = items
		}
	}

	return response
}

func (r *fieldsParsingResult) generateItems(objects []Object) []map[string]interface{} {
	items := []map[string]interface{}{}
	for _, object := range objects {
		item := map[string]interface{}{}
		for _, itemField := range r.ItemFields {
			structFieldName := jsonFieldNameToStructFieldName[itemField]
			item[itemField] = object.getValueByFieldName(structFieldName)
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
