package kubelet

import (
	"time"
	"net"
	"fmt"
	"strconv"
	"encoding/json"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"golang.org/x/net/context"

	cadvisorapiv2 "github.com/google/cadvisor/info/v2"
	infraapi "github.com/sjpotter/infranetes/pkg/common"
)

func (kl *Kubelet) infranetesInfo(count int) (map[string]cadvisorapiv2.ContainerInfo, error) {
	//1. get local info
	options := cadvisorapiv2.RequestOptions{
		IdType:    cadvisorapiv2.TypeName,
		Count:     count,
		Recursive: false,
	}

	local, err := kl.cadvisor.ContainerInfoV2("/", options)
	if err != nil {
		glog.Infof("infranetesInfo: Couldn't get local node stats: %v", err)
		return nil, fmt.Errorf("Couldn't get local node stats: %v", err)
	}

	//1. dial /tmp/infra
	conn, err := dial("/tmp/infra")
	if err != nil {
		glog.Infof("infranetesInfo: Couldn't dial infranetes: %v", err)
		return nil, fmt.Errorf("Couldn't dial infranetes: %v", err)
	}
	defer conn.Close()

	//2. make client
	client := infraapi.NewMetricsClient(conn)

	//3. call GetMetrics()
	resp, err := client.GetMetrics(context.Background(), &infraapi.GetMetricsRequest{Count: int32(count)})
	if err != nil {
		glog.Infof("infranetesInfo: GetMetrics failed: %v", err)
		return nil, fmt.Errorf("infranetesInfo: GetMetrics failed: %v", err)
	}
	infos := collectResults(resp)

	//4. combine with local info
	infos["/"] = local["/"]

	//5. return result
	return infos, nil
}

func dial(file string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(file, grpc.WithInsecure(), grpc.WithTimeout(5*time.Second),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}))

	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	return conn, nil
}

func collectResults(results *infraapi.GetMetricsResponse) map[string]cadvisorapiv2.ContainerInfo {
	ret := make(map[string]cadvisorapiv2.ContainerInfo)

	for i, jsonResult := range results.JsonMetricResponses {
		k := strconv.Itoa(i)
		var info cadvisorapiv2.ContainerInfo
		err := json.Unmarshal(jsonResult, &info)
		if err != nil {
			glog.Warningf("collectResults: couldn't unmarshall json: %v", err)
			continue
		}
		ret[k] = info
	}

	return ret
}
