package helpers

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHasStringInSlice(t *testing.T) {
	Convey("no string given and empty array returns false", t, func() {
		So(HasStringInSlice("", []string{}), ShouldBeFalse)
	})
	Convey("no string given valid string array returns false", t, func() {
		So(HasStringInSlice("", []string{"hello", "world"}), ShouldBeFalse)
	})
	Convey("valid string given empty array returns false", t, func() {
		So(HasStringInSlice("hello", []string{}), ShouldBeFalse)
	})
	Convey("valid string given with valid string array returns true", t, func() {
		So(HasStringInSlice("hello", []string{"hello", "world"}), ShouldBeTrue)
	})
}

func TestCheckForExistingParams(t *testing.T) {
	Convey("persist existing query string values and ignore given value", t, func() {
		queryStrValues := []string{"Value 1", "Value 2"}
		ignoreValue := "Value 1"
		key := "testKey"
		q := url.Values{}

		PersistExistingParams(queryStrValues, key, ignoreValue, q)
		So(q[key], ShouldResemble, []string{"Value 2"})
	})

	Convey("persist existing query string values no value to ignore", t, func() {
		queryStrValues := []string{"Value 1", "Value 2"}
		existingValue := ""
		key := "testKey"
		q := url.Values{}

		PersistExistingParams(queryStrValues, key, existingValue, q)
		So(q[key], ShouldResemble, queryStrValues)
	})

	Convey("invalid key given no values persisted", t, func() {
		queryStrValues := []string{"Value 1", "Value 2"}
		existingValue := ""
		key := "testKey"
		q := url.Values{}

		PersistExistingParams(queryStrValues, key, existingValue, q)
		So(q["another key"], ShouldBeNil)
		So(q[key], ShouldResemble, queryStrValues)
	})
}

func TestToBoolPtr(t *testing.T) {
	Convey("Given a bool, a pointer is returned", t, func() {
		So(ToBoolPtr(false), ShouldResemble, new(bool))
		So(ToBoolPtr(false), ShouldNotBeNil)
		truePtr := func(b bool) *bool { return &b }(true)
		So(ToBoolPtr(true), ShouldResemble, truePtr)
	})
}
