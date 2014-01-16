package container_test

import (
	"testing"
	"github.com/marmelab/arch-o-matic/container"

	"code.google.com/p/gomock/gomock"
	"github.com/marmelab/arch-o-matic/docker" // mock
)

func TestStartedContainerShouldBeMarkedAsRunning (t *testing.T) {
	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup the ext mock package
	docker.MOCK().SetController(ctrl)
	docker.EXPECT().Start(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("123")
	docker.EXPECT().Inspect(gomock.Any()).Return([]byte("[{\"ID\": \"123\", \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"))

	container := container.Container{Name: "Test"}
	container.Start()

	if !container.IsRunning() {
		t.Error("Container not running")
	}

	if container.Ip != "172.17.0.10" {
		t.Error("Container IP not retrieved correctly")
	}
}

