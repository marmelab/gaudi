package container_test

import (
	"testing"
	"code.google.com/p/gomock/gomock"
	. "launchpad.net/gocheck"

	"github.com/marmelab/gaudi/container"
	"github.com/marmelab/gaudi/docker" // mock
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ContainerTestSuite struct{}
var _ = Suite(&ContainerTestSuite{})

func (s *ContainerTestSuite) TestStartedContainerShouldRetrieveItsIp(c *C) {
	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	// Setup the docker mock package
	docker.MOCK().SetController(ctrl)
	docker.EXPECT().Start(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("123")
	docker.EXPECT().Inspect(gomock.Any()).Return([]byte("[{\"ID\": \"123\", \"State\":{\"Running\": false}, \"NetworkSettings\": {\"IPAddress\": \"\"}}]"), nil)
	docker.EXPECT().Inspect(gomock.Any()).Return([]byte("[{\"ID\": \"123\", \"State\":{\"Running\": true}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil)

	// @TODO : find a way to mock time.Sleep
	container := container.Container{Name: "Test"}
	container.Start()

	c.Check(container.IsRunning(), Equals, true)
	c.Check(container.Ip, Equals, "172.17.0.10")
}

func (s *ContainerTestSuite) TestGetFirstPortShouldReturnTheFirstDeclaredPort (c *C) {
	container := container.Container{Name: "Test"}

	c.Check(container.GetFirstPort(), Equals, "")

	container.Ports = make(map[string]string)
	container.Ports["80"] = "8080"
	container.Ports["9000"] = "9000"

	c.Check(container.GetFirstPort(), Equals, "80")
}

func (s *ContainerTestSuite) TestCallCleanShouldStopTheContainer(c *C) {
	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	// Setup the docker mock package
	docker.MOCK().SetController(ctrl)
	docker.EXPECT().Inspect(gomock.Any()).Return([]byte("[{\"ID\": \"123\", \"State\":{\"Running\": false}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil)
	docker.EXPECT().Kill(gomock.Any()).Return()
	docker.EXPECT().Remove(gomock.Any()).Return()

	done := make(chan bool, 1)
	container := container.Container{Name: "Test"}
	container.Clean(done)
	<-done

	c.Check(container.IsRunning(), Equals, false)
}

func (s *ContainerTestSuite) TestContainerWithReadyDependenciesShouldBeReady(c *C) {
	dep1 := container.Container{Name: "abc"}
	dep1.Running = false

	dep2 := container.Container{Name: "abc"}
	dep2.Running = true

	mainContainer := container.Container{Name: "Test"}
	mainContainer.Dependencies = make([]*container.Container, 0)

	mainContainer.AddDependency(&dep1)
	mainContainer.AddDependency(&dep2)

	c.Check(mainContainer.IsReady(), Equals, false)

	dep1.Running = true
	c.Check(mainContainer.IsReady(), Equals, true)
}
