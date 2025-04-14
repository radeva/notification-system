package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"notification-system/pkg/config"
	"notification-system/pkg/model"
	"notification-system/pkg/storage"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Slack Integration Test", func() {
	var (
		db            *storage.Database
		apiURL        string
		notification  model.Notification
	)

	ginkgo.BeforeEach(func() {
		// Load test configuration
		cfg, err := config.LoadConfigFromFile("../.env.test")
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Initialize database
		db, err = storage.NewDatabase(cfg.Database)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Set up test notification
		notification = model.Notification{
			Channel:   model.ChannelSlack,
			Recipient: "C08NKAKQ4N4",
			Message:   "Test Slack message",
		}

		// Set API URL
		apiURL = fmt.Sprintf("http://localhost:%s", cfg.Server.Port)
	})

	ginkgo.AfterEach(func() {
		if db != nil {
			db.Close()
		}
	})

	ginkgo.It("should create and send a Slack notification via API", func() {
		// Create notification via API
		notificationJSON, err := json.Marshal(notification)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		resp, err := http.Post(
			fmt.Sprintf("%s/notifications", apiURL),
			"application/json",
			bytes.NewBuffer(notificationJSON),
		)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusAccepted))

		// Parse response to get notification ID
		var response struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(response.Status).To(gomega.Equal(string(model.StatusPending)))

		// Wait for the worker to process the message
		time.Sleep(5 * time.Second)

		// Check notification status via API
		resp, err = http.Get(fmt.Sprintf("%s/notifications/%s/status", apiURL, response.ID))
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK))

		var statusResponse model.Notification
		err = json.NewDecoder(resp.Body).Decode(&statusResponse)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(statusResponse.Status).To(gomega.Equal(model.StatusSent))
		gomega.Expect(statusResponse.Attempts).To(gomega.Equal(1))
	})
})