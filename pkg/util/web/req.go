package web

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	// 默认单页数据量
	_defaultPageSize = 20
	// 默认当前页
	_defaultCurrentPage = 1
)

// Request request object
type Request struct {
	//  页码
	Page int `query:"page"`
	// 每页数量
	Limit int `query:"limit"`
	// 项目名称
	Project string `query:"project"`
	// 集群名称
	Cluster string `query:"cluster" param:"cluster"`
	// 命名空间
	Namespace string `query:"namespace" param:"namespace"`
	// 资源名称
	Name string `query:"name" param:"name"`
	// Workload名称
	Workload string `query:"workload"`
	// 资源owner
	Owner string `query:"owner"`
	// 资源owner类型
	OwnerKind string `query:"ownerKind"`
	// pod 可能会用到的参数
	Log       bool   `query:"log"`
	Event     bool   `query:"event"`
	Container string `query:"container"`
	Follow    bool   `query:"follow"`
	Tail      int    `query:"tail"`
	Previous  bool   `query:"previous"`
	SinceTime string `query:"sinceTime"`
	Describe  bool   `query:"describe"`
}

// Query k8s resource query structure
type Query struct {
	// cluster to query
	Cluster string
	// namespace to query
	Namespace string
	// current page
	CurrentPage int
	// page size, the limit of one page
	PageSize int
	// list option for k8s resource
	ListOption *metav1.ListOptions
	// get option for k8s resource
	GetOption *metav1.GetOptions
	// update option for k8s resource
	UpdateOption *metav1.UpdateOptions
	// raw label selector
	Selector labels.Selector
}

// HandleQueryParam 处理分页参数，返回总页数
func (q *Query) HandleQueryParam(total int) int {
	totalPages := 1
	// pagesize不合理的数据可用来表示获取所有
	if q.PageSize <= 0 {
		q.CurrentPage = 1
		return totalPages
	}
	if q.CurrentPage <= 0 {
		q.CurrentPage = _defaultCurrentPage
	}
	if q.PageSize <= 0 {
		q.PageSize = _defaultPageSize
	}
	if total > q.PageSize {
		totalPages = total / q.PageSize
		if total%q.PageSize > 0 {
			totalPages = totalPages + 1
		}
	}
	if q.CurrentPage > totalPages {
		q.CurrentPage = totalPages
	}
	return totalPages
}
