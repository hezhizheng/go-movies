package utils

type SpiderTask interface {
	Start()
	PageDetail(id string)
	DoRecentUpdate()
}
