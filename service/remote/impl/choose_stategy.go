package impl

import (
	"dataProcess/service/remote"
	"net/url"
)

type RemoteReq struct {
	TrainReq remote.TrainReq
}

func (d *RemoteReq) RequestRemote(query url.Values, heads map[string]string) (res interface{}, err error) {
	res, err = d.TrainReq.RequestRemote(query, heads)
	return
}

func CreateRemoteReq(trainReq remote.TrainReq) *RemoteReq {
	return &RemoteReq{
		TrainReq: trainReq}
}
