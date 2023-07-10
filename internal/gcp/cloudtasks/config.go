package cloudtasks
import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)
type Config struct {
	Queues map[string]Queue `mapstructure:"cloudTasks" validate:"dive"`
}

type Queue struct {
	Name       string `json:"name" validate:"required"`
	Region     string `json:"region" validate:"required"`
	ProjectId  string `json:"projectId" validate:"required"`
	MinBackoff              string `json:"minBackoff" validate:"required"`
	MaxBackoff              string `json:"maxBackoff" validate:"required"`
	MaxConcurrentDispatches int    `json:"maxConcurrentDispatches" validate:"required"`
	MaxDispatchesPerSecond  float64 `json:"maxDispatchesPerSecond" validate:"required"`
	clientName string
}


func GetCloudTasksConfig(viperConfig *viper.Viper, clientName string) (*Config, error) {
	if viperConfig == nil {
		return nil, nil
	}

	var cloudTasksConfig Config
	err := viperConfig.Unmarshal(&cloudTasksConfig)
	if err != nil {
		return nil, err
	}

	for name, queue := range cloudTasksConfig.Queues {
		queue.Name = strings.Join([]string{clientName, name, queue.ProjectId}, "-")

		cloudTasksConfig.Queues[name] = queue
	}
	return &cloudTasksConfig, nil
}

func ValidateConfig(config *Config) error {
	v := validator.New()
	if err := v.Struct(config); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("%s validate failed on the %s rule", err.Namespace(), err.Tag())
		}
	}
	return nil
}

func (q *Queue) Parent() string {
	return fmt.Sprintf("projects/%s/locations/%s", q.ProjectId, q.Region)
}