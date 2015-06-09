package libs

import (
	"github.com/sogko/slumber/domain"
)

// MergeACLMap Returns a new map
func MergeACLMap(to *domain.ACLMap, from *domain.ACLMap) domain.ACLMap {
	res := domain.ACLMap{}
	for k, v := range *to {
		res[k] = v
	}
	for k, v := range *from {
		res[k] = v
	}
	return res
}

// MergeRoutes Returns a new Routes
func MergeRoutes(to *domain.Routes, from *domain.Routes) domain.Routes {
	res := domain.Routes{}
	for _, v := range *to {
		res = append(res, v)
	}
	for _, v := range *from {
		res = append(res, v)
	}
	return res
}
