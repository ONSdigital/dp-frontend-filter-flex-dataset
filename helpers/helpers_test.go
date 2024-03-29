package helpers

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mocks"
	"github.com/ONSdigital/dp-renderer/v2/helper"
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

func TestPluralise(t *testing.T) {
	helper.InitialiseLocalisationsHelper(mocks.MockAssetFunction)
	req := httptest.NewRequest("GET", "http://localhost:20100", nil)

	Convey("Given a valid key with lookup prefix", t, func() {
		input := "Country"
		expectedOutput := "Countries"
		lookupPrefix := "AreaType"
		sut := Pluralise(req, input, "en", lookupPrefix, 4)

		Convey("Then pluralised value is returned", func() {
			So(sut, ShouldEqual, expectedOutput)
		})
	})

	Convey("Given a valid key without lookup prefix", t, func() {
		input := "Test"
		expectedOutput := "Tests"
		sut := Pluralise(req, input, "en", "", 4)

		Convey("Then pluralised value is returned", func() {
			So(sut, ShouldEqual, expectedOutput)
		})
	})

	Convey("Given a valid key with spaces and mixed case, without lookup prefix", t, func() {
		input := "aRea tYPE Country"
		expectedOutput := "Countries"
		sut := Pluralise(req, input, "en", "", 4)

		Convey("Then pluralised value is returned", func() {
			So(sut, ShouldEqual, expectedOutput)
		})
	})

	Convey("Given a valid key without lookup prefix in Welsh", t, func() {
		input := "Test"
		expectedOutput := "Tests (cy)"
		sut := Pluralise(req, input, "cy", "", 4)

		Convey("Then pluralised value is returned", func() {
			So(sut, ShouldEqual, expectedOutput)
		})
	})

	Convey("Given an invalid key", t, func() {
		input := "Blah blah"
		expectedOutput := ""
		sut := Pluralise(req, input, "", "", 1)

		Convey("Then empty string is returned", func() {
			So(sut, ShouldEqual, expectedOutput)
		})
	})
}

func TestTrimCategoryValue(t *testing.T) {
	Convey("Given a string with the category substring", t, func() {
		tc := []struct{ given, expected string }{
			{
				given:    "abc_1a",
				expected: "abc",
			},
			{
				given:    "abc_1",
				expected: "abc_1",
			},
			{
				given:    "abc_123",
				expected: "abc",
			},
			{
				given:    "abc_123a",
				expected: "abc",
			},
			{
				given:    "abc",
				expected: "abc",
			},
		}

		Convey("When the string is passed to the TrimCategoryValue function", func() {
			for _, t := range tc {
				Convey(fmt.Sprintf("Then the given string '%s' returns '%s'", t.given, t.expected), func() {
					sut := TrimCategoryValue(t.given)
					So(sut, ShouldEqual, t.expected)
				})
			}
		})
	})
}

func TestIsBoolPtr(t *testing.T) {
	Convey("When the value is nil", t, func() {
		Convey("Then the returned value is false", func() {
			So(IsBoolPtr(nil), ShouldBeFalse)
		})
	})
	Convey("When the value is a false pointer", t, func() {
		ptr := *new(bool)
		Convey("Then the returned value is false", func() {
			So(IsBoolPtr(&ptr), ShouldBeFalse)
		})
	})
	Convey("When the value is a true pointer", t, func() {
		ptr := bool(true)
		Convey("Then the returned value is true", func() {
			So(IsBoolPtr(&ptr), ShouldBeTrue)
		})
	})
}
