package authority

type RolePermissionRequest struct {
	RoleID        uint64   `json:"role_id"`
	PermissionIDs []uint64 `json:"permission_ids"`
}

type Request struct {
	Includes []string `json:"includes" query:"includes" param:"includes"`
	Sorts    []string `json:"sorts" query:"sorts" param:"sorts"`
	Excludes []string `json:"excludes" query:"excludes" param:"excludes"`
}

type RequestPagination struct {
	PerPage int `json:"per_page" query:"per_page" param:"per_page"`
	Page    int `json:"page" query:"page" param:"page"`
	Offset  int `json:"offset" query:"offset" param:"offset"`
}

type RequestFilters struct {
	Search string `json:"search" query:"search" param:"search"`
}

type RequestData struct {
	Request
	RequestFilters
	RequestPagination
}

func (r *RequestData) SetDefault() {
	if r.PerPage <= 0 {
		r.PerPage = 10
	}

	if r.Page <= 0 {
		r.Page = 1
	}

	if r.Offset <= 0 {
		r.Offset = 0
	}

	if r.Page > 1 {
		r.Offset = (r.Page - 1) * r.PerPage
	}

	if len(r.Sorts) <= 0 {
		r.Sorts = []string{"created_at:desc"}
	}
}
