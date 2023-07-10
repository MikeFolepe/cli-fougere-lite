package cloudtasks

type Config struct {
	Queues map[string]Queue `mapstructure:"cloudTasks" validate:"dive"`
}

type Queue struct {
	Name       string `json:"name" validate:"required"`
	Region     string `json:"region" validate:"required"`
	ProjectId  string `json:"projectId" validate:"required"`
}