package container_test

import (
	"testing"
	"github.com/marmelab/arch-o-matic/container"

	"code.google.com/p/gomock/gomock"
	"github.com/marmelab/arch-o-matic/docker" // mock
)

func TestStartedContainerShouldRetrieveItsIp(t *testing.T) {
	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup the docker mock package
	docker.MOCK().SetController(ctrl)
	docker.EXPECT().Start(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("123")
	docker.EXPECT().Inspect(gomock.Any()).Return([]byte("[{\"ID\": \"123\", \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"))

	// @TODO : find a way to mock time.Sleep

	container := container.Container{Name: "Test"}
	container.Start()

	if !container.IsRunning() {
		t.Error("Started container should be marked as running")
	}

	if container.Ip != "172.17.0.10" {
		t.Error("Started container IP not retrieved correctly")
	}
}

func TestGetFirstPortShouldReturnTheFirstDeclaredPort (t *testing.T) {
	container := container.Container{Name: "Test"}

	if container.GetFirstPort() != "" {
		t.Error("Container without port should return empty value when calling GetFirstPort")
	}

	container.Ports = make(map[string]string)
	container.Ports["80"] = "8080"
	container.Ports["9000"] = "9000"

	if container.GetFirstPort() != "80" {
		t.Error("Container with at least a port should returns the first when calling GetFirstPort")
	}
}

func TestCallCleanShouldStopTheContainer(t *testing.T) {
	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup the docker mock package
	docker.MOCK().SetController(ctrl)
	docker.EXPECT().Clean(gomock.Any()).Return()

	done := make(chan bool, 1)
	container := container.Container{Name: "Test"}
	container.Running = true
	container.Clean(done)
	<-done

	if container.IsRunning() {
		t.Error("Cleaned container should not be marked as running")
	}
}

func TestContainerWithReadyDependenciesShouldBeReady(t *testing.T) {
	dep1 := container.Container{Name: "abc"}
	dep1.Running = false

	dep2 := container.Container{Name: "abc"}
	dep2.Running = true

	c := container.Container{Name: "Test"}
	c.Dependencies = make([]*container.Container, 0)

	c.AddDependency(&dep1)
	c.AddDependency(&dep2)

	if c.IsReady() {
		t.Error("Container with non running dependencies should not be ready")
	}

	dep1.Running = true
	if !c.IsReady() {
		t.Error("Container with running dependencies should be ready")
	}
}
