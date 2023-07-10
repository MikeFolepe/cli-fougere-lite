// ©Copyright 2023 Metrio
package cloudtasks

import (
	"context"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/api/option"
	"google.golang.org/api/cloudtasks/v2"
	"metrio.net/fougere-lite/internal/utils"
)

// Helper method to create a client
func getMockedClient(serverURL string) *Client {
	client, err := NewClient(context.Background(), option.WithoutAuthentication(), option.WithEndpoint(serverURL))
	if err != nil {
		Fail(err.Error())
	}
	return client
}

var _ = Describe("Cloud Tasks client", func() {
	var queue Queue

	BeforeEach(func() {
		queue = Queue{
			Name:       "projects/project-123/locations/northamerica-northeast1/queues/queue-1",
			Region:     "northamerica-northeast1",
			ProjectId:  "project-123",
			MinBackoff: "1s",
			MaxBackoff: "10s",
			MaxConcurrentDispatches: 100,
			MaxDispatchesPerSecond:  1000,
			clientName: "banane",
		}
	})
	
	Describe("Create queue", func() {
		It("successfully creates the queue", func() {
			mockServerCalls := make(chan utils.MockServerCall, 1)
			mockServerCalls <- utils.MockServerCall{
				UrlMatchFunc: func(url string) bool {
					return strings.Contains(url, "projects/project-123/locations/northamerica-northeast1/queues")
				},
				Method: "post",

				ResponseBody: &cloudtasks.Queue{
					Name: "projects/project-123/locations/northamerica-northeast1/queues/queue-1",
					RateLimits: &cloudtasks.RateLimits{
						MaxDispatchesPerSecond: 1000,
						MaxConcurrentDispatches: 100,
					},
					RetryConfig: &cloudtasks.RetryConfig{
						MinBackoff: "1s",
						MaxBackoff: "10s",
					},
				}, 
				ResponseCode: 200,
				
			}
			mockServer := utils.NewMockServer(mockServerCalls)
			defer mockServer.Close()
	
			client := getMockedClient(mockServer.URL)
	
			err := client.create(queue)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Update queue", func() {
		It("successfully updates the queue", func() {
			mockServerCalls := make(chan utils.MockServerCall, 1)
			mockServerCalls <- utils.MockServerCall{
				UrlMatchFunc: func(url string) bool {
					return strings.Contains(url, "projects/project-123/locations/northamerica-northeast1/queues/queue-1")
				},
				Method: "patch",

				ResponseBody: &cloudtasks.Queue{
					Name: "projects/project-123/locations/northamerica-northeast1/queues/queue-1",
					RateLimits: &cloudtasks.RateLimits{
						MaxDispatchesPerSecond: 1000,
						MaxConcurrentDispatches: 100,
					},
					RetryConfig: &cloudtasks.RetryConfig{
						MinBackoff: "1s",
						MaxBackoff: "10s",
					},
				}, 
				ResponseCode: 200,
			}
			mockServer := utils.NewMockServer(mockServerCalls)
			defer mockServer.Close()
	
			client := getMockedClient(mockServer.URL)
	
			err := client.update(queue)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	
	
	
})
