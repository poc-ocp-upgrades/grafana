package dashdiffs

import (
	"testing"
	"github.com/grafana/grafana/pkg/components/simplejson"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDiff(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	const (
		leftJSON	= `{
			"key": "value",
			"object": {
				"key": "value",
				"anotherObject": {
					"same": "this field is the same in rightJSON",
					"change": "this field should change in rightJSON",
					"delete": "this field doesn't appear in rightJSON"
				}
			},
			"array": [
				"same",
				"change",
				"delete"
			],
			"embeddedArray": {
				"array": [
					"same",
					"change",
					"delete"
				]
			}
		}`
		rightJSON	= `{
			"key": "differentValue",
			"object": {
				"key": "value",
				"newKey": "value",
				"anotherObject": {
					"same": "this field is the same in rightJSON",
					"change": "this field should change in rightJSON",
					"add": "this field is added"
				}
			},
			"array": [
				"same",
				"changed!",
				"add"
			],
			"embeddedArray": {
				"array": [
					"same",
					"changed!",
					"add"
				]
			}
		}`
	)
	Convey("Testing dashboard diffs", t, func() {
		baseData, err := simplejson.NewJson([]byte(leftJSON))
		So(err, ShouldBeNil)
		newData, err := simplejson.NewJson([]byte(rightJSON))
		So(err, ShouldBeNil)
		left, jsonDiff, err := getDiff(baseData, newData)
		So(err, ShouldBeNil)
		Convey("The JSONFormatter should produce the expected JSON tokens", func() {
			f := NewJSONFormatter(left)
			_, err := f.Format(jsonDiff)
			So(err, ShouldBeNil)
			changeCounts := make(map[ChangeType]int)
			for _, line := range f.Lines {
				changeCounts[line.Change]++
			}
			expectedChangeCounts := map[ChangeType]int{ChangeNil: 12, ChangeAdded: 2, ChangeDeleted: 1, ChangeOld: 5, ChangeNew: 5, ChangeUnchanged: 5}
			So(changeCounts, ShouldResemble, expectedChangeCounts)
		})
		Convey("The BasicFormatter should produce the expected BasicBlocks", func() {
			f := NewBasicFormatter(left)
			_, err := f.Format(jsonDiff)
			So(err, ShouldBeNil)
			bd := &BasicDiff{}
			blocks := bd.Basic(f.jsonDiff.Lines)
			changeCounts := make(map[ChangeType]int)
			for _, block := range blocks {
				for _, change := range block.Changes {
					changeCounts[change.Change]++
				}
				for _, summary := range block.Summaries {
					changeCounts[summary.Change]++
				}
				changeCounts[block.Change]++
			}
			expectedChangeCounts := map[ChangeType]int{ChangeNil: 3, ChangeAdded: 2, ChangeDeleted: 1, ChangeOld: 3}
			So(changeCounts, ShouldResemble, expectedChangeCounts)
		})
	})
}
