package response

type PaginationCursorResponseMeta struct {
	TotalData       int64       `json:"totalData,omitempty"`
	TotalDataOnPage int64       `json:"totalDataOnPage"`
	PrevCursor      interface{} `json:"prevCursor"`
	NextCursor      interface{} `json:"nextCursor"`
}
