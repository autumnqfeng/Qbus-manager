package cluster

type AddClusterVo struct {
	Name          string `json:"name"`
	ZookeeperList string `json:"zookeeper_list"`
}
