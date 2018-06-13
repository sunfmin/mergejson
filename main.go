package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fileSource := os.Args[1]
	fileTarget := os.Args[2]
	fs, err := os.Open(fileSource)
	if err != nil {
		fmt.Println("source:", fileSource)
		panic(err)
	}
	ft, err := os.Open(fileTarget)
	if err != nil {
		fmt.Println("target:", fileTarget)
		panic(err)
	}
	fsb, err := ioutil.ReadAll(fs)
	if err != nil {
		panic(err)
	}
	ftb, err := ioutil.ReadAll(ft)
	if err != nil {
		panic(err)
	}

	var mapSource map[string]interface{}

	err = json.Unmarshal(fsb, &mapSource)
	if err != nil {
		panic(err)
	}

	var mapTarget map[string]interface{}
	err = json.Unmarshal(ftb, &mapTarget)

	var merged = make(map[string]interface{})
	mergeOrMark(mapSource, mapTarget, merged)

	out, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}

func conflictedKey(k string) string {
	return fmt.Sprintf("!!!conflicted!!!%s", k)
}

func sourceDeleted(k string) string {
	return fmt.Sprintf("!!!source_deleted!!!%s", k)

}

func mergeOrMark(source, target, merged map[string]interface{}) {
	for k, sv := range source {
		if mtv, yes := sv.(map[string]interface{}); yes {
			newm := make(map[string]interface{})
			merged[k] = newm

			var newtarget = make(map[string]interface{})
			if _, texists := target[k]; texists {
				if nt, ok := target[k].(map[string]interface{}); ok {
					newtarget = nt
				} else {
					merged[conflictedKey(k)] = target[k]
				}
			}

			mergeOrMark(mtv, newtarget, newm)
			continue
		}

		if tv, ok := target[k]; ok {
			merged[k] = tv
		} else {
			merged[k] = sv
		}
	}

	for tk, tv := range target {
		if _, ok := source[tk]; !ok {
			merged[sourceDeleted(tk)] = tv
		}
	}
}
