package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"slices"
)

type Instances []Instance

func (i Instances) FilterOut(flowLogs FlowLogs) Instances {
	var out Instances
	names := flowLogs.NamesSet()
	for _, v := range i {
		name := fmt.Sprintf("instance-%s", v.Name)
		if _, ok := names[name]; !ok {
			out = append(out, v)
		}
	}
	return out
}

func (i Instances) GetById(id string) (Instance, bool) {
	for _, instance := range i {
		if instance.Id == id {
			return instance, true
		}
	}
	return Instance{}, false
}

// GetByNames returns sorted names (keys) and map where key is name and value instances
func (i Instances) GetByNames() ([]string, map[string]Instances) {
	var names []string
	out := make(map[string]Instances)
	for _, instance := range i {
		if _, ok := out[instance.Name]; !ok {
			out[instance.Name] = Instances{}
			names = append(names, instance.Name)
		}
		out[instance.Name] = append(out[instance.Name], instance)
	}
	slices.Sort(names)
	return names, out
}

type Instance struct {
	VpcId               string
	SubnetId            string
	Id                  string
	Name                string
	NetworkInterfaceIds []string
	tags                map[string]string
}

func (i Instance) Tags() map[string]string {
	out := make(map[string]string)
	for k, v := range i.tags {
		out[k] = v
	}
	return out
}

func (i Instance) String() string {
	name := i.Name
	if name == "" {
		name = "-"
	}
	return fmt.Sprintf("%s [%s]", name, i.Id)
}

func toInstances(in []types.Instance) Instances {
	var out Instances
	for _, v := range in {
		out = append(out, toInstance(v))
	}
	return out
}

func toInstance(in types.Instance) Instance {
	tags := toTags(in.Tags)
	var networkInterfaceIds []string
	for _, v := range in.NetworkInterfaces {
		if v.Status == types.NetworkInterfaceStatusInUse {
			networkInterfaceIds = append(networkInterfaceIds, aws.ToString(v.NetworkInterfaceId))
		}
	}

	return Instance{
		VpcId:               aws.ToString(in.VpcId),
		SubnetId:            aws.ToString(in.SubnetId),
		Id:                  aws.ToString(in.InstanceId),
		Name:                tags["Name"],
		NetworkInterfaceIds: networkInterfaceIds,
		tags:                tags,
	}
}
