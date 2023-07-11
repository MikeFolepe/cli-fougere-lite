// ©Copyright 2023 Metrio
package cloudtasks

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var validQueueConfig = []byte(`
cloudTasks:
  metrio-task:
    region: us-central1
    projectId: some-project
    minBackoff: 100ms
    maxBackoff: 1000ms
    maxConcurrentDispatches: 10
    maxDispatchesPerSecond: 100.0`)

var invalidConfig = []byte(`
cloudTasks:
  some-task:
    region:
      - should_not_be_an_array`)

var _ = Describe("config", func() {
	BeforeEach(func() {
		viper.Reset()
		viper.SetConfigType("yaml")
	})

	Describe("GetCloudTasksConfig", func() {
		It("should successfully parse a cloud task config", func() {
			err := viper.ReadConfig(bytes.NewBuffer(validQueueConfig))
			Expect(err).ToNot(HaveOccurred())
			cloudTasksConfig, err := GetCloudTasksConfig(viper.GetViper(), "metrio-client")
			Expect(err).To(BeNil())
			Expect(len(cloudTasksConfig.Queues)).To(Equal(1))
			queue := cloudTasksConfig.Queues["metrio-task"]
			Expect(queue.Region).To(Equal("us-central1"))
			Expect(queue.ProjectId).To(Equal("some-project"))
			Expect(queue.MinBackoff).To(Equal("100ms"))
			Expect(queue.MaxBackoff).To(Equal("1000ms"))
			Expect(queue.MaxConcurrentDispatches).To(Equal(10))
			Expect(queue.MaxDispatchesPerSecond).To(Equal(100.0))
		})

		It("returns an error if cannot parse the config", func() {
			err := viper.ReadConfig(bytes.NewBuffer(invalidConfig))
			Expect(err).ToNot(HaveOccurred())
			_, err = GetCloudTasksConfig(viper.GetViper(), "metrio-client")
			Expect(err).NotTo(BeNil())
		})
	})

	Context("validates cloud task queues", func() {
		It("should not detect error", func() {
			config := &Config{
				Queues: map[string]Queue{
					"foooo": {
						Name:       "foooo",
						Region:     "us-central1",
						ProjectId:  "mock-project",
						MinBackoff: "100ms",
						MaxBackoff: "1000ms",
						MaxConcurrentDispatches: 10,
						MaxDispatchesPerSecond:  100.0,
					},
				},
			}
			err := ValidateConfig(config)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should detect empty name cloud tasks queues", func() {
			config := &Config{
				Queues: map[string]Queue{
					"foooo": {
						Region:    "us-central1",
						ProjectId: "mock-project",
						MinBackoff: "100ms",
						MaxBackoff: "1000ms",
						MaxConcurrentDispatches: 10,
						MaxDispatchesPerSecond:  100.0,
					},
				},
			}
			err := ValidateConfig(config)
			Expect(err).Should(MatchError(ContainSubstring("validate failed on the required rule")))
		})

		It("should detect an empty region", func() {
			config := &Config{
				Queues: map[string]Queue{
					"foooo": {
						Name:       "foooo",
						ProjectId:  "mock-project",
						MinBackoff: "100ms",
						MaxBackoff: "1000ms",
						MaxConcurrentDispatches: 10,
						MaxDispatchesPerSecond:  100.0,
					},
				},
			}
			err := ValidateConfig(config)
			Expect(err).Should(MatchError(ContainSubstring("validate failed on the required rule")))
		})


		It("should detect an empty project id", func() {
			config := &Config{
				Queues: map[string]Queue{
					"foooo": {
						Name:       "foooo",
						Region:    "us-central1",
						MinBackoff: "100ms",
						MaxBackoff: "1000ms",
						MaxConcurrentDispatches: 10,
						MaxDispatchesPerSecond:  100.0,
					},
				},
			}
			err := ValidateConfig(config)
			Expect(err).Should(MatchError(ContainSubstring("validate failed on the required rule")))
		})
	})
})
