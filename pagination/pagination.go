package pagination

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	coreModel "github.com/ONSdigital/dp-renderer/model"
)

// GetTotalPages returns the total pages from the total results count and pagesize parameters
func GetTotalPages(totalCount, pageSize int) int {
	if totalCount <= 0 || pageSize <= 0 {
		return 0
	}

	if totalCount <= pageSize {
		return 1
	}

	if totalCount%pageSize != 0 {
		return (totalCount / pageSize) + 1
	}

	return totalCount / pageSize
}

// GetOffset returns the offset (0 based) into a list, given a page number (1 based) and the size of a page.
// A pageSize <= 0 or a pageNo <= 0 will give an offset of 0
func GetOffset(pageSize, pageNo int) int {
	if pageSize <= 0 || pageNo <= 0 {
		return 0
	}

	return (pageSize * pageNo) - pageSize
}

// GetPagesToDisplay returns the pages to be displayed within the first and last pages
func GetPagesToDisplay(currentPage, totalPages int, req *http.Request) []coreModel.PageToDisplay {
	pageRange := getPageRange(totalPages)
	start, end := getWindowStartEndPage(currentPage, totalPages, pageRange)

	var pagesToDisplay []coreModel.PageToDisplay

	for i := start; i <= end; i++ {
		pagesToDisplay = append(pagesToDisplay, coreModel.PageToDisplay{
			PageNumber: i,
			URL:        getPageUrl(req, i),
		})
	}

	return pagesToDisplay
}

// GetFirstAndLastPages returns the first and last pages of a paginated results
func GetFirstAndLastPages(req *http.Request, totalPages int) []coreModel.PageToDisplay {
	return []coreModel.PageToDisplay{
		{
			PageNumber: 1,
			URL:        getPageUrl(req, 1),
		},
		{
			PageNumber: totalPages,
			URL:        getPageUrl(req, totalPages),
		},
	}
}

// getWindowStartEndPage calculates the start and end page of the moving window of size windowSize, over the set of pages
// whose current page is currentPage, and whose size is totalPages
// It is an error to pass a parameter whose value is < 1, or a currentPage > totalPages, and the function will panic in this case
func getWindowStartEndPage(currentPage, totalPages, windowSize int) (int, int) {
	if currentPage < 1 || totalPages < 1 || windowSize < 1 || currentPage > totalPages {
		panic("invalid parameters for getWindowStartEndPage - see documentation")
	}
	switch {
	case windowSize == 1:
		se := (currentPage % totalPages) + 1
		return se, se
	case windowSize >= totalPages:
		return 1, totalPages
	}

	windowOffset := getWindowOffset(windowSize)
	start := currentPage - windowOffset
	switch {
	case start <= 0:
		start = 1
	case start > totalPages-windowSize+1:
		start = totalPages - windowSize + 1
	}

	end := start + windowSize - 1
	if end > totalPages {
		end = totalPages
	}

	return start, end
}

func getWindowOffset(windowSize int) int {
	if windowSize%2 == 0 {
		return (windowSize / 2) - 1
	}

	return windowSize / 2
}

// getPageUrl returns the url for a paginated page that retains existing query string parameters
func getPageUrl(req *http.Request, pg int) string {
	page := strconv.Itoa(pg)
	u, _ := url.Parse(fmt.Sprint(req.URL))

	q := u.Query()
	q.Set("page", page)
	u.RawQuery = q.Encode()

	return fmt.Sprint(u)
}

// getPageRange ensures that the total page numbers displayed is limited to 5 i.e. 1 .. 2, (3), 4 ... 7
func getPageRange(totalPages int) int {
	var pageRange int
	if totalPages <= 5 {
		pageRange = 5
	} else {
		pageRange = 3
	}
	return pageRange
}
