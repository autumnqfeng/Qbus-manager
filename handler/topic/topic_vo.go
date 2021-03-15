package topic

type CreateTopicVo struct {
	Cluster      string `json:"cluster"`
	Topic        string `json:"topic"`
	Partitions   int    `json:"partitions"`
	Replications int    `json:"replications"`
	Retention    int    `json:"retention"`
	MaxMessage   int    `json:"max_message"`
}

type UpdateTopicVo struct {
	Cluster    string `json:"cluster"`
	Topic      string `json:"topic"`
	Retention  int    `json:"retention"`
	MaxMessage int    `json:"max_message"`
}
