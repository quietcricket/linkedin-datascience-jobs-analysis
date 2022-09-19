package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

func flatternJson(input map[string]interface{}, output *[]string, keys *map[string]int, parentKey string) {
	for k, v := range input {
		var newKey string
		if len(parentKey) > 0 {
			newKey = parentKey + "." + k
		} else {
			newKey = k
		}
		if reflect.ValueOf(v).Kind() == reflect.Map {
			flatternJson(v.(map[string]interface{}), output, keys, newKey)
		} else if v != nil {
			(*output)[(*keys)[newKey]] = fmt.Sprintf("%v", v)
		}
	}
}

func ProcessJson() {
	dataDir := "./data/json"
	keys := extractKeys(dataDir, 10)

	csvOutput, _ := os.OpenFile("./data/jobs.csv", os.O_WRONLY|os.O_CREATE, 0644)
	writer := csv.NewWriter(csvOutput)
	writer.Write(sortKeys(&keys))
	defer csvOutput.Close()

	counter := 0
	filepath.Walk(dataDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		s, _ := os.ReadFile(path)
		var data map[string]interface{}
		json.Unmarshal(s, &data)
		row := make([]string, len(keys))
		flatternJson(data, &row, &keys, "")
		writer.Write(row)

		counter++
		if counter > 10 {
			// return io.EOF
		}
		return nil
	})
	writer.Flush()

}

// This function is very similar to flattern
// It's possible to create a more powerful version of flattern function
// But I think it will be too confusing. Thus, just duplicate it and make
// minor changes here
func extractKeyRecurse(inputJson map[string]interface{}, keys *map[string]int, parentKey string) {
	for k, v := range inputJson {
		var newKey string
		if len(parentKey) > 0 {
			newKey = parentKey + "." + k
		} else {
			newKey = k
		}
		if reflect.ValueOf(v).Kind() == reflect.Map {
			extractKeyRecurse(v.(map[string]interface{}), keys, newKey)
		} else {
			_, ok := (*keys)[newKey]
			if !ok {
				(*keys)[newKey] = len(*keys)
			}
		}
	}
}

/**
Find all possible keys
The full structure of the json is unknown and it may chance over the time
Some entries may have more data than others such as salary range.
Some of the entries from UK and CN contains salaray range.
This function scan through a percentage of documents (hopefully it's enough) and find all possible keys
*/

func extractKeys(folder string, skip int) map[string]int {
	keys := make(map[string]int)
	counter := -1
	filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		counter++
		if counter%skip > 0 {
			return nil
		}
		s, _ := os.ReadFile(path)
		var data map[string]interface{}
		json.Unmarshal(s, &data)
		extractKeyRecurse(data, &keys, "")
		return nil
	})
	return keys
}

func sortKeys(keys *map[string]int) []string {
	type kv struct {
		k string
		v int
	}
	kvPairs := make([]kv, 0)
	for k, v := range *keys {
		kvPairs = append(kvPairs, kv{k, v})
	}

	sort.Slice(kvPairs, func(i, j int) bool {
		return kvPairs[i].v < kvPairs[j].v
	})
	fmt.Println(kvPairs)
	sortedKeys := make([]string, len(kvPairs))
	for i, kv := range kvPairs {
		sortedKeys[i] = strings.Replace(kv.k, "@", "", -1)
	}

	return sortedKeys
}
