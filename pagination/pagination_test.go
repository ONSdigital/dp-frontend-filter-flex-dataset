package pagination

import (
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ONSdigital/dp-renderer/v2/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetTotalPages(t *testing.T) {
	t.Parallel()
	Convey("Given the total count and page size of 50", t, func() {
		testcases := []struct{ totalCount, expectedPages int }{
			{totalCount: 1, expectedPages: 1},
			{totalCount: 50, expectedPages: 1},
			{totalCount: 51, expectedPages: 2},
			{totalCount: 200, expectedPages: 4},
			{totalCount: 1000, expectedPages: 20},
			{totalCount: -1, expectedPages: 0},
		}
		Convey("When the 'GetTotalPages' function is called", func() {
			for _, tc := range testcases {
				sut := GetTotalPages(tc.totalCount, 50)
				Convey(fmt.Sprintf("Then %s pages are returned for a count of %s", strconv.Itoa(tc.expectedPages), strconv.Itoa(tc.totalCount)), func() {
					So(sut, ShouldEqual, tc.expectedPages)
				})
			}
		})
	})
}

func TestGetOffset(t *testing.T) {
	t.Parallel()
	Convey("Given the current page number and page size of 50", t, func() {
		testcases := []struct{ pageNo, expectedOffset int }{
			{pageNo: 1, expectedOffset: 0},
			{pageNo: 2, expectedOffset: 50},
			{pageNo: 3, expectedOffset: 100},
			{pageNo: 10, expectedOffset: 450},
			{pageNo: -1, expectedOffset: 0},
		}
		Convey("When the 'GetOffset' function is called", func() {
			for _, tc := range testcases {
				sut := GetOffset(50, tc.pageNo)
				Convey(fmt.Sprintf("Then an offset of %s is returned for page number %s", strconv.Itoa(tc.expectedOffset), strconv.Itoa(tc.pageNo)), func() {
					So(sut, ShouldEqual, tc.expectedOffset)
				})
			}
		})
	})
}

func TestGetPagesToDisplay(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "/a/page", nil)
	Convey("Given an http request, total pages and current page parameters", t, func() {
		testcases := []struct {
			totalPages    int
			currentPage   int
			expectedModel []model.PageToDisplay
		}{
			{
				totalPages:  2,
				currentPage: 1,
				expectedModel: []model.PageToDisplay{
					{
						PageNumber: 1,
						URL:        "/a/page?page=1",
					},
					{
						PageNumber: 2,
						URL:        "/a/page?page=2",
					},
				},
			},
			{
				totalPages:  3,
				currentPage: 1,
				expectedModel: []model.PageToDisplay{
					{
						PageNumber: 1,
						URL:        "/a/page?page=1",
					},
					{
						PageNumber: 2,
						URL:        "/a/page?page=2",
					},
					{
						PageNumber: 3,
						URL:        "/a/page?page=3",
					},
				},
			},
			{
				totalPages:  5,
				currentPage: 5,
				expectedModel: []model.PageToDisplay{
					{
						PageNumber: 1,
						URL:        "/a/page?page=1",
					},
					{
						PageNumber: 2,
						URL:        "/a/page?page=2",
					},
					{
						PageNumber: 3,
						URL:        "/a/page?page=3",
					},
					{
						PageNumber: 4,
						URL:        "/a/page?page=4",
					},
					{
						PageNumber: 5,
						URL:        "/a/page?page=5",
					},
				},
			},
			{
				totalPages:  10,
				currentPage: 5,
				expectedModel: []model.PageToDisplay{
					{
						PageNumber: 4,
						URL:        "/a/page?page=4",
					},
					{
						PageNumber: 5,
						URL:        "/a/page?page=5",
					},
					{
						PageNumber: 6,
						URL:        "/a/page?page=6",
					},
				},
			},
		}
		Convey("When the 'GetPagesToDisplay' function is called", func() {
			for _, tc := range testcases {
				sut := GetPagesToDisplay(tc.currentPage, tc.totalPages, req)
				Convey(fmt.Sprintf("Then with 'total pages' of %s and the 'current page' is page %s the model should resemble the expected model", strconv.Itoa(tc.totalPages), strconv.Itoa(tc.currentPage)), func() {
					So(sut, ShouldResemble, tc.expectedModel)
				})
			}
		})
	})
}

func TestGetFirstAndLastPages(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "/a/page", nil)
	Convey("Given an http request and the total pages parameter", t, func() {
		testcases := []struct {
			totalPages    int
			expectedModel []model.PageToDisplay
		}{
			{
				totalPages: 1,
				expectedModel: []model.PageToDisplay{
					{
						PageNumber: 1,
						URL:        "/a/page?page=1",
					},
					{
						PageNumber: 1,
						URL:        "/a/page?page=1",
					},
				},
			},
			{
				totalPages: 2,
				expectedModel: []model.PageToDisplay{
					{
						PageNumber: 1,
						URL:        "/a/page?page=1",
					},
					{
						PageNumber: 2,
						URL:        "/a/page?page=2",
					},
				},
			},
			{
				totalPages: 5,
				expectedModel: []model.PageToDisplay{
					{
						PageNumber: 1,
						URL:        "/a/page?page=1",
					},
					{
						PageNumber: 5,
						URL:        "/a/page?page=5",
					},
				},
			},
		}
		Convey("When the 'GetFirstAndLastPages' function is called", func() {
			for _, tc := range testcases {
				sut := GetFirstAndLastPages(req, tc.totalPages)
				Convey(fmt.Sprintf("Then with total pages of %s the model should resemble the expected model", strconv.Itoa(tc.totalPages)), func() {
					So(sut, ShouldResemble, tc.expectedModel)
				})
			}
		})
	})
}

func TestGetStartEndPage(t *testing.T) {
	t.Parallel()
	Convey("Given a set of parameters expressing: the 'current page number', out of a 'total number of pages', and the 'window size'", t, func() {
		testcases := []struct{ current, total, window, exStart, exEnd int }{
			{current: 1, total: 1, window: 1, exStart: 1, exEnd: 1},

			{current: 1, total: 2, window: 1, exStart: 2, exEnd: 2},
			{current: 2, total: 2, window: 1, exStart: 1, exEnd: 1},

			{current: 1, total: 3, window: 2, exStart: 1, exEnd: 2},
			{current: 2, total: 3, window: 2, exStart: 2, exEnd: 3},
			{current: 3, total: 3, window: 2, exStart: 2, exEnd: 3},

			{current: 1, total: 3, window: 3, exStart: 1, exEnd: 3},
			{current: 2, total: 3, window: 3, exStart: 1, exEnd: 3},
			{current: 3, total: 3, window: 3, exStart: 1, exEnd: 3},

			{current: 3, total: 4, window: 3, exStart: 2, exEnd: 4},
			{current: 3, total: 4, window: 5, exStart: 1, exEnd: 4},

			{current: 28, total: 32, window: 5, exStart: 26, exEnd: 30},
			{current: 31, total: 32, window: 5, exStart: 28, exEnd: 32},
		}
		Convey("check the generated start and end page numbers are correct", func() {
			for _, tc := range testcases {
				sp, ep := getWindowStartEndPage(tc.current, tc.total, tc.window)
				So(sp, ShouldEqual, tc.exStart)
				So(ep, ShouldEqual, tc.exEnd)
			}
		})
	})
}

func TestGetPageRange(t *testing.T) {
	t.Parallel()
	Convey("Given the parameter total pages", t, func() {
		testcases := []struct{ totalPages, expectedRange int }{
			{totalPages: 1, expectedRange: 5},
			{totalPages: 3, expectedRange: 5},
			{totalPages: 10, expectedRange: 3},
			{totalPages: 20, expectedRange: 3},
		}
		Convey("When the 'getPageRange' function is called", func() {
			for _, tc := range testcases {
				sut := getPageRange(tc.totalPages)
				Convey(fmt.Sprintf("Then a range of %s is returned for total pages of %s", strconv.Itoa(tc.expectedRange), strconv.Itoa(tc.totalPages)), func() {
					So(sut, ShouldEqual, tc.expectedRange)
				})
			}
		})
	})
}
