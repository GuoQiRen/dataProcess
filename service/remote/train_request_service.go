package remote

import (
	"net/url"
)

type TrainReq interface {
	RequestRemote(query url.Values, heads map[string]string) (resp interface{}, err error)
}
