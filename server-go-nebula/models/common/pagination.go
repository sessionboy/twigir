package models

// 分页
type UserPagination struct {
	Limit uint8      `json:"limit"` // 返回n条数据
	Page  uint8      `json:"page"`  // 跳过m条数据
	Total string     `json:"total"`
	Data  []ListUser `json:"data"`
}
