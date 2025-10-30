package shared

// SetupQueryPagination sets up pagination parameters for a query
func SetupQueryPagination(page *int, limit *int) *Pagination {
	var (
		currentPage int
		perPage     int
		offset      int
	)

	if page != nil {
		currentPage = *page
	} else {
		currentPage = 1
	}

	if limit != nil {
		perPage = *limit
	} else {
		perPage = 100
	}

	offset = (currentPage - 1) * perPage

	return &Pagination{
		Limit:  perPage,
		Offset: offset,
	}
}

// SetupQuerySorting sets up sorting parameters for a query
func SetupQuerySorting(orderBy *string, asc *bool) *SortMethod {
	if orderBy != nil {
		isAsc := true
		if asc != nil {
			isAsc = *asc
		}
		return &SortMethod{
			Asc:     isAsc,
			OrderBy: *orderBy,
		}
	}

	return &SortMethod{
		Asc:     true,
		OrderBy: "id",
	}
}
