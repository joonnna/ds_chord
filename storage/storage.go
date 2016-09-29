package storage

import (
	"github.com/joonnna/ds_chrod/logger"
	"github.com/joonnna/ds_chrod/util"
	"strings"
)



func SplitStorage(id string, newId string, prevId string vals map[string]string) map[string]string {
	ret = make(map[string]string)

	for key, val := range vals {
		if utils.InKeySpace(newId, key, prevId) {
			ret[key] = val
			delete(vals, key)
		}
	}

	return ret
}
