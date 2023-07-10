package cloudtasks

import (
	"context"
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

}


func (c *Client) get(name string) (*cloudtasks.Queue, error) {

}

func (c *Client) create(queue Queue) error {
	
}

func (c *Client) update(queue Queue) error {
	
}
