package config

import "github.com/spf13/viper"

type TrainConf struct {
	Head                   string
	Host                   string
	ImageName              string
	ImageTag               string
	UserGroup              string
	ResGroupName           string
	ResGroupLabel          string
	ResourceGroupCloudId   string
	ResourceGroupServerCpu int32
	ResourceCpuOverRation  float64
	ResourceGroupType      string
	IsLoop                 bool
	FrameType              string
	MountContext           string
	MountTrainContext      string
	JobContext             string
	RuntimeContext         string
	TaskOperateContext     string
	TrainLogContext        string
	IsMulti                string
	NodeToUse              int32
	Priority               int32
	GpuType                int32
	HardwareType           int32
	Tag                    string
	ResourceModel          int32
}

func InitTrainConf(cfg *viper.Viper) *TrainConf {
	return &TrainConf{
		Head:                   cfg.GetString("head"),
		Host:                   cfg.GetString("host"),
		ImageName:              cfg.GetString("imageName"),
		ResGroupName:           cfg.GetString("resGroupName"),
		ResGroupLabel:          cfg.GetString("resGroupLabel"),
		ResourceGroupCloudId:   cfg.GetString("resourceGroupCloudId"),
		ResourceGroupServerCpu: cfg.GetInt32("resourceGroupServerCpu"),
		ResourceCpuOverRation:  cfg.GetFloat64("resourceCpuOverRation"),
		ResourceGroupType:      cfg.GetString("resourceGroupType"),
		ImageTag:               cfg.GetString("imageTag"),
		IsMulti:                cfg.GetString("isMulti"),
		NodeToUse:              cfg.GetInt32("nodeToUse"),
		Priority:               cfg.GetInt32("priority"),
		GpuType:                cfg.GetInt32("gpuType"),
		HardwareType:           cfg.GetInt32("hardwareType"),
		Tag:                    cfg.GetString("tag"),
		ResourceModel:          cfg.GetInt32("resourceModel"),
		MountContext:           cfg.GetString("mountContext"),
		MountTrainContext:      cfg.GetString("mountTrainContext"),
		UserGroup:              cfg.GetString("userGroup"),
		IsLoop:                 cfg.GetBool("isLoop"),
		FrameType:              cfg.GetString("frameType"),
		JobContext:             cfg.GetString("jobContext"),
		RuntimeContext:         cfg.GetString("runtimeContext"),
		TrainLogContext:        cfg.GetString("trainLogContext"),
		TaskOperateContext:     cfg.GetString("taskOperateContext"),
	}
}
