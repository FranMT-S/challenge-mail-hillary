package cmd

import (
	"strconv"
	"testing"

	"indexer/models"

	"github.com/stretchr/testify/assert"
)

var testCases = []struct {
	name       string
	pageResult models.PageResult
}{
	{"Page 1", models.PageResult{Page: 1, State: models.PageResultStateFinished, Total: 10, Error: ""}},
	{"Page 2", models.PageResult{Page: 2, State: models.PageResultStatePending, Total: 25, Error: "some error"}},
}

func Test_SaveAndLoadStatus(t *testing.T) {

	newCmd := NewCmd(nil, 10, 2)

	newCmd.status.Set("1", testCases[0].pageResult)
	newCmd.status.Set("2", testCases[1].pageResult)
	newCmd.SavePageResults("../data", "test.json", newCmd.status.GetCopy())

	loadedStatus, err := newCmd.LoadPageResults("../data", "test.json")
	if err != nil {
		t.Error(err)
	}

	for i, tc := range testCases {
		key := strconv.Itoa(i + 1)
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.pageResult, loadedStatus[key])
		})
	}
}
