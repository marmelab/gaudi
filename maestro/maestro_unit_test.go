package maestro_test

import (
	"testing"
	"code.google.com/p/gomock/gomock"
	. "launchpad.net/gocheck"

	"github.com/marmelab/arch-o-matic/docker" // mock
	"github.com/marmelab/arch-o-matic/maestro"
)

func Test(t *testing.T) { TestingT(t) }

type MaestroTestSuite struct{}
var _ = Suite(&MaestroTestSuite{})

func (s *MaestroTestSuite) TestInitFromStringShouldTrowAndErrorOnMalformedYmlContent(c *C) {
	m := maestro.Maestro{}

	c.Assert(func() { m.InitFromString(`
		containers:
			tabulated:
				type: varnish
`, "") }, PanicMatches, "YAML error: line 1: found character that cannot start any token")
}

func (s *MaestroTestSuite) TestInitFromStringShouldTrowAndErrorOnWrongContent(c *C) {
	m := maestro.Maestro{}

	c.Assert(func() { m.InitFromString("<oldFormat>Skrew you, i'm not yml</oldFormat>", "") }, PanicMatches, "No container to start")
}

func (s *MaestroTestSuite) TestInitFromStringShouldCreateAMaestro (c *C) {
	m := maestro.Maestro{}
	m.InitFromString(`
containers:
    app:
        type: php-fpm
        links: [db]
    db:
        type: mysql
        ports:
            3306: 9000
`, "")

	c.Assert(len(m.Containers), Equals, 2)
	c.Assert(m.GetContainer("app").Name, Equals, "app")
	c.Assert(m.GetContainer("app").Type, Equals, "php-fpm")
	c.Assert(m.GetContainer("app").Dependencies[0].Name, Equals, "db")
	c.Assert(m.GetContainer("db").GetFirstPort(), Equals, "3306")
	c.Assert(m.GetContainer("db").IsRunning(), Equals, false)
}

func (s *MaestroTestSuite) TestStartContainerShouldCleanAndBuildThem (c *C) {
	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	// Setup the docker mock package
	docker.MOCK().SetController(ctrl)
	docker.EXPECT().Clean(gomock.Any()).Return().Times(2)
	docker.EXPECT().Build(gomock.Any()).Return().Times(2)
	docker.EXPECT().Start(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("123").Times(2)
	docker.EXPECT().Inspect(gomock.Any()).Return([]byte("[{\"ID\": \"123\", \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]")).Times(2)

	m := maestro.Maestro{}
	m.InitFromString(`
containers:
    app:
        type: php-fpm
        links: [db]
    db:
        type: mysql
        ports:
            3306: 9000
`, "")

	c.Assert(len(m.Containers), Equals, 2)

	m.Start()
	c.Assert(m.GetContainer("db").IsRunning(), Equals, true)
	c.Assert(m.GetContainer("app").IsRunning(), Equals, true)
}

func (s *MaestroTestSuite) TestStartContainerShouldStartThemByOrderOfDependencies(c *C) {
	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	// Setup the docker mock package
	docker.MOCK().SetController(ctrl)

	m := maestro.Maestro{}
	m.InitFromString(`
containers:
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

	c.Assert(len(m.Containers), Equals, 5)

	docker.EXPECT().Clean(gomock.Any()).Return().Times(5)
	docker.EXPECT().Build(gomock.Any()).Return().Times(5)

	gomock.InOrder(
		docker.EXPECT().Start("db", gomock.Any(), gomock.Any(), gomock.Any()).Return("123"),
		docker.EXPECT().Start("app", gomock.Any(), gomock.Any(), gomock.Any()).Return("123"),
		docker.EXPECT().Start("front1", gomock.Any(), gomock.Any(), gomock.Any()).Return("123"),
		docker.EXPECT().Start("front2", gomock.Any(), gomock.Any(), gomock.Any()).Return("123"),
		docker.EXPECT().Start("lb", gomock.Any(), gomock.Any(), gomock.Any()).Return("123"),
	)

	docker.EXPECT().Inspect(gomock.Any()).Return([]byte("[{\"ID\": \"123\", \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]")).Times(5)

	m.Start()
}
