package utils

type Page struct {
	PageNo     int         `json:"page_no"`
	PageSize   int         `json:"page_size"`
	TotalPage  int         `json:"total_page"`
	TotalCount int         `json:"total_count"`
	FirstPage  bool        `json:"first_page"`
	LastPage   bool        `json:"last_page"`
	List       interface{} `json:"list"`
}

// 分页工具类
func PageUtil(count int, pageNo int, pageSize int, list interface{}) Page {
	tp := count / pageSize
	if count%pageSize > 0 {
		tp = count/pageSize + 1
	}
	return Page{PageNo: pageNo, PageSize: pageSize, TotalPage: tp, TotalCount: count, FirstPage: pageNo == 1, LastPage: pageNo == tp, List: list}
}
