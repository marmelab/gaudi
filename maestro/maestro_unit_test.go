package maestro_test

import (
	"code.google.com/p/gomock/gomock"
	. "launchpad.net/gocheck"
	"testing"

	mockfmt "fmt"                      // mock
	"github.com/marmelab/gaudi/docker" // mock
	"github.com/marmelab/gaudi/maestro"
)

func Test(t *testing.T) { TestingT(t) }

type MaestroTestSuite struct{}

var _ = Suite(&MaestroTestSuite{})

func (s *MaestroTestSuite) TestInitFromStringShouldTrowAndErrorOnMalformedYmlContent(c *C) {
	m := maestro.Maestro{}

	c.Assert(func() {
		m.InitFromString(`
		applications:
			tabulated:
				type: varnish
`, "")
	}, PanicMatches, "YAML error: line 1: found character that cannot start any token")
}

func (s *MaestroTestSuite) TestInitFromStringShouldTrowAndErrorOnWrongContent(c *C) {
	m := maestro.Maestro{}

	c.Assert(func() { m.InitFromString("<oldFormat>Skrew you, i'm not yml</oldFormat>", "") }, PanicMatches, "No application to start")
}

func (s *MaestroTestSuite) TestInitFromStringShouldCreateAMaestro(c *C) {
	m := maestro.Maestro{}
	m.InitFromString(`
applications:
    app:
        type: php-fpm
        links: [db]
    db:
        type: mysql
        ports:
            3306: 9000
`, "")

	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()
	docker.MOCK().SetController(ctrl)
	docker.EXPECT().Inspect(gomock.Any()).Return([]byte("[{\"ID\": \"123\", \"State\":{\"Running\": false}, \"NetworkSettings\": {\"IPAddress\": \"\"}}]"), nil)

	c.Assert(len(m.Applications), Equals, 2)
	c.Assert(m.GetContainer("app").Name, Equals, "app")
	c.Assert(m.GetContainer("app").Type, Equals, "php-fpm")
	c.Assert(m.GetContainer("app").Dependencies[0].Name, Equals, "db")
	c.Assert(m.GetContainer("db").GetFirstPort(), Equals, "3306")
	c.Assert(m.GetContainer("db").IsRunning(), Equals, false)
}

func (s *MaestroTestSuite) TestStartApplicationShouldCleanAndBuildThem(c *C) {
	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	// Setup the docker mock package
	docker.MOCK().SetController(ctrl)
	docker.EXPECT().ImageExists(gomock.Any()).Return(true).Times(1)
	docker.EXPECT().Kill(gomock.Any()).Return().Times(2)
	docker.EXPECT().Remove(gomock.Any()).Return().Times(2)
	docker.EXPECT().Build(gomock.Any(), gomock.Any()).Return().Times(2)
	docker.EXPECT().Start(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("123").Times(2)
	docker.EXPECT().Inspect(gomock.Any()).Return([]byte("[{\"ID\": \"123\", \"State\":{\"Running\": false}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil).Times(4)

	m := maestro.Maestro{}
	m.InitFromString(`
applications:
    app:
        type: php-fpm
        links: [db]
    db:
        type: mysql
        ports:
            3306: 9000
`, "")

	c.Assert(len(m.Applications), Equals, 2)

	m.Start()
	c.Assert(m.GetContainer("db").IsRunning(), Equals, true)
	c.Assert(m.GetContainer("app").IsRunning(), Equals, true)
}

func (s *MaestroTestSuite) TestStartApplicationShouldStartThemByOrderOfDependencies(c *C) {
	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	// Setup the docker mock package
	docker.MOCK().SetController(ctrl)

	m := maestro.Maestro{}
	m.InitFromString(`
applications:
    lb:
        links: [front1, front2]
        type: varnish

    front1:
        links: [app]
        type: apache

    front2:
        links: [app]
        type: apache

    app:
        links: [db]
        type: php-fpm

    db:
      type: mysql
`, "")

	c.Assert(len(m.Applications), Equals, 5)

	docker.EXPECT().ImageExists(gomock.Any()).Return(true).Times(1)
	docker.EXPECT().Kill(gomock.Any()).Return().Times(5)
	docker.EXPECT().Remove(gomock.Any()).Return().Times(5)
	docker.EXPECT().Build(gomock.Any(), gomock.Any()).Return().Times(5)

	gomock.InOrder(
		docker.EXPECT().Inspect("db").Return([]byte("[{\"ID\": \"100\", \"State\":{\"Running\": false}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil),
		docker.EXPECT().Start("db", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("100"),
		docker.EXPECT().Inspect("100").Return([]byte("[{\"ID\": \"100\", \"State\":{\"Running\": true}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil),

		docker.EXPECT().Inspect("app").Return([]byte("[{\"ID\": \"101\", \"State\":{\"Running\": false}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil),
		docker.EXPECT().Start("app", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("101"),
		docker.EXPECT().Inspect("101").Return([]byte("[{\"ID\": \"101\", \"State\":{\"Running\": true}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil),

		docker.EXPECT().Inspect("front1").Return([]byte("[{\"ID\": \"102\", \"State\":{\"Running\": false}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil),
		docker.EXPECT().Start("front1", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("102"),

		docker.EXPECT().Inspect("front2").Return([]byte("[{\"ID\": \"103\", \"State\":{\"Running\": false}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil),
		docker.EXPECT().Start("front2", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("103"),

		docker.EXPECT().Inspect("102").Return([]byte("[{\"ID\": \"102\", \"State\":{\"Running\": true}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil),
		docker.EXPECT().Inspect("103").Return([]byte("[{\"ID\": \"103\", \"State\":{\"Running\": true}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil),

		docker.EXPECT().Inspect("lb").Return([]byte("[{\"ID\": \"104\", \"State\":{\"Running\": false}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil),
		docker.EXPECT().Start("lb", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("104"),
		docker.EXPECT().Inspect("104").Return([]byte("[{\"ID\": \"104\", \"State\":{\"Running\": true}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil),
	)

	m.Start()
}

func (s *MaestroTestSuite) TestCheckRunningContainerShouldUseDockerPs(c *C) {
	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	// Setup the docker mock package
	docker.MOCK().SetController(ctrl)

	// Setup the mockfmt mock package
	mockfmt.MOCK().SetController(ctrl)

	psResult := make(map[string]string)
	psResult["gaudi/lb"] = "123"
	psResult["gaudi/front1"] = "124"
	psResult["gaudi/db"] = "125"

	docker.EXPECT().SnapshotProcesses().Return(psResult, nil)
	docker.EXPECT().Inspect("123").Return([]byte("[{\"ID\": \"123\", \"State\":{\"Running\": true}, \"NetworkSettings\": {\"IPAddress\": \"123.124.125.126\"}}]"), nil)
	docker.EXPECT().Inspect("124").Return([]byte("[{\"ID\": \"123\", \"State\":{\"Running\": true}, \"NetworkSettings\": {\"IPAddress\": \"123.124.125.127\"}}]"), nil)
	docker.EXPECT().Inspect("125").Return([]byte("[{\"ID\": \"123\", \"State\":{\"Running\": true}, \"NetworkSettings\": {\"IPAddress\": \"123.124.125.128\"}}]"), nil)

	mockfmt.EXPECT().Println("Application", "lb", "is running", "(123.124.125.126:)")
	mockfmt.EXPECT().Println("Application", "front1", "is running", "(123.124.125.127:)")
	mockfmt.EXPECT().Println("Application", "db", "is running", "(123.124.125.128:3306)")

	m := maestro.Maestro{}
	m.InitFromString(`
applications:
    lb:
        links: [front1]
        type: varnish

    front1:
        type: apache

    db:
        type: mysql
        ports:
            3306: 9000
`, "")

	m.Check()
}
