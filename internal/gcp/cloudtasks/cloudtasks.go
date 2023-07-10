package cloudtasks

import (
	"context"
	"strings"
	"google.golang.org/api/cloudtasks/v2"
	"google.golang.org/api/option"
	"metrio.net/fougere-lite/internal/common"
	"metrio.net/fougere-lite/internal/utils"
)

type Client struct {
	cloudTasksService *cloudtasks.Service
}

func NewClient(ctx context.Context, opts ...option.ClientOption) (*Client, error) {
	cloudTasksService, err := cloudtasks.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{
		cloudTasksService: cloudTasksService,
	}, nil
}

func (c *Client) Create(config *Config) error {
	createChannel := make(chan common.Response, len(config.Queues))
	for _, queue := range config.Queues {
		go func(resp chan common.Response, queue Queue) {
			_, err := c.get(queue.Name)
			if err != nil {
				if strings.Contains(err.Error(), "NotFound") {
					utils.Logger.Debug("[%s] queue not found", queue.Name)
					if err := c.create(queue); err != nil {
						resp <- common.Response{Err: err}
						return
					}
				} else {
					utils.Logger.Errorf("[%s] error getting queue: %s", queue.Name, err)
					resp <- common.Response{Err: err}
					return
				}
			} else {
				if err := c.update(queue); err != nil {
					resp <- common.Response{Err: err}
					return
				}
			}
			resp <- common.Response{}
		}(createChannel, queue)
	}
	for range config.Queues {
		resp := <-createChannel
		if resp.Err != nil {
			return resp.Err
		}
	}
	return nil
}

func (c *Client) get(name string) (*cloudtasks.Queue, error) {
	utils.Logger.Debug("[%s] getting queue", name)
	queue, err := c.cloudTasksService.Projects.Locations.Queues.Get(name).Do()
	if err != nil {
		return nil, err
	}
	return queue, nil
}

func (c *Client) create(queue Queue) error {
	utils.Logger.Infof("[%s] creating queue", queue.Name)
	spec := c.createQueueSpec(queue)
	_, err := c.cloudTasksService.Projects.Locations.Queues.Create(queue.ProjectId, spec).Do()
	if err != nil {
		utils.Logger.Errorf("[%s] error creating queue: %s", spec.Name, err)
		return err
	}

	return nil
}

func (c *Client) update(queue Queue) error {
	spec := c.createQueueSpec(queue)
	utils.Logger.Infof("[%s] updating queue", spec.Name)
	_, err := c.cloudTasksService.Projects.Locations.Queues.Patch(spec.Name, spec).Do()
	if err != nil {
		utils.Logger.Errorf("[%s] error updating queue: %s", spec.Name, err)
		return err
	}
	return nil
	
}

func (c *Client) createQueueSpec(queue Queue) *cloudtasks.Queue {
	// Fully qualified queue name
	return &cloudtasks.Queue{
		Name: queue.Name,
		RateLimits: &cloudtasks.RateLimits{
			MaxDispatchesPerSecond: queue.MaxDispatchesPerSecond,
			MaxConcurrentDispatches: int64(queue.MaxConcurrentDispatches),
		},
		RetryConfig: &cloudtasks.RetryConfig{
			MinBackoff: queue.MinBackoff,
			MaxBackoff: queue.MaxBackoff,
		},
	}
}

