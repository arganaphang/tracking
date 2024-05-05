package pagination

func ToLimitOffset(page, perPage uint) (limit, offset uint) {
	limit = perPage
	offset = (page - 1) * perPage
	return
}
