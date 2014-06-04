package gaudi_test

import (
	"code.google.com/p/gomock/gomock"
	. "launchpad.net/gocheck"
	"os"
	"testing"

	"github.com/marmelab/gaudi/docker" // mock
	"github.com/marmelab/gaudi/gaudi"
	"github.com/marmelab/gaudi/util" // mock
)

func Test(t *testing.T) { TestingT(t) }

type GaudiTestSuite struct{}

var _ = Suite(&GaudiTestSuite{})

func disableLog() {
	util.MOCK().DisableMock("LogError")
	util.MOCK().DisableMock("PrintGreen")
	util.MOCK().DisableMock("PrintWithColor")
	util.MOCK().DisableMock("BuildReflectArguments")
}

func enableLog() {
	util.MOCK().EnableMock("LogError")
	util.MOCK().EnableMock("PrintGreen")
	util.MOCK().EnableMock("PrintWithColor")
	util.MOCK().EnableMock("BuildReflectArguments")
}

func (s *GaudiTestSuite) TestInitShouldTrowAndErrorOnMalformedYmlContent(c *C) {
	g := gaudi.Gaudi{}

	// Disable the util package mock
	util.MOCK().DisableMock("LogError")

	c.Assert(func() {
		g.Init(`
		applications:
			tabulated:
				type: varnish
`)
	}, PanicMatches, "YAML error: line 1: found character that cannot start any token")
}

func (s *GaudiTestSuite) TestInitShouldTrowAndErrorOnWrongContent(c *C) {
	g := gaudi.Gaudi{}

	// Disable the util package mock
	util.MOCK().DisableMock("LogError")

	c.Assert(func() { g.Init("<oldFormat>Skrew you, i'm not yml</oldFormat>") }, Panics, "No application or binary to start. Are you missing a 'applications' or 'binaries' field in your configuration ?")
}

func (s *GaudiTestSuite) TestInitShouldCreateApplications(c *C) {
	os.RemoveAll("/var/tmp/gaudi/")

	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()
	docker.MOCK().SetController(ctrl)

	// Setup the util mock package
	util.MOCK().SetController(ctrl)

	// Disable the util package mock
	util.MOCK().DisableMock("IsDir")
	util.MOCK().DisableMock("IsFile")

	util.EXPECT().PrintGreen("Retrieving templates ...")

	docker.EXPECT().ImageExists(gomock.Any()).Return(true).Times(1)
	docker.EXPECT().HasDocker().Return(true).Times(1)
	docker.EXPECT().Inspect(gomock.Any()).Return([]byte("[{\"ID\": \"123\", \"State\":{\"Running\": false}, \"NetworkSettings\": {\"IPAddress\": \"\"}}]"), nil)

	g := gaudi.Gaudi{}
	g.Init(`
applications:
    app:
        type: php-fpm
        links: [db]
    db:
        type: mysql
        ports:
            3306: 9000
`)

	c.Assert(len(g.Applications), Equals, 2)
	c.Assert(g.GetApplication("app").Name, Equals, "app")
	c.Assert(g.GetApplication("app").Type, Equals, "php-fpm")
	c.Assert(g.GetApplication("app").Dependencies[0].Name, Equals, "db")
	c.Assert(g.GetApplication("db").GetFirstPort(), Equals, "3306")
	c.Assert(g.GetApplication("db").IsRunning(), Equals, false)
}

func (s *GaudiTestSuite) TestStartApplicationShouldCleanAndBuildThem(c *C) {
	os.RemoveAll("/var/tmp/gaudi/templates/")

	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	// Setup the util mock package
	util.MOCK().SetController(ctrl)

	util.MOCK().DisableMock("IsFile")
	util.MOCK().DisableMock("IsDir")

	// Retrieving templates (1)
	util.EXPECT().PrintGreen(gomock.Any()).Times(1)
	// Killing, Clearing, Building, Starting (3*2)
	util.EXPECT().PrintGreen(gomock.Any(), gomock.Any(), gomock.Any()).Times(6)
	// Started (1*2)
	util.EXPECT().PrintGreen(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(2)

	// Setup the docker mock package
	docker.MOCK().SetController(ctrl)
	docker.EXPECT().ImageExists(gomock.Any()).Return(true).Times(1)
	docker.EXPECT().HasDocker().Return(true).Times(1)
	docker.EXPECT().Kill(gomock.Any()).Return().Times(2)
	docker.EXPECT().Remove(gomock.Any()).Return().Times(2)
	docker.EXPECT().Build(gomock.Any(), gomock.Any()).Return().Times(2)
	docker.EXPECT().Start(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("123").Times(2)
	docker.EXPECT().Inspect(gomock.Any()).Return([]byte("[{\"ID\": \"123\", \"State\":{\"Running\": false}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil).Times(2)

	g := gaudi.Gaudi{}
	g.Init(`
applications:
    app:
        type: php-fpm
        links: [db]
    db:
        type: mysql
        ports:
            3306: 9000
`)

	c.Assert(len(g.Applications), Equals, 2)

	g.StartApplications(true)
	c.Assert(g.GetApplication("db").IsRunning(), Equals, true)
	c.Assert(g.GetApplication("app").IsRunning(), Equals, true)
}

func (s *GaudiTestSuite) TestStartApplicationShouldStartThemByOrderOfDependencies(c *C) {
	os.RemoveAll("/var/tmp/gaudi/templates/")

	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	// Setup the docker mock package
	docker.MOCK().SetController(ctrl)

	// Setup the util mock package
	util.MOCK().SetController(ctrl)

	// Disable the util package mock
	util.MOCK().DisableMock("IsDir")
	util.MOCK().DisableMock("IsFile")
	util.MOCK().EnableMock("PrintGreen")

	// Retrieving templates (1)
	util.EXPECT().PrintGreen(gomock.Any()).Times(1)
	// Killing, Clearing, Building, Starting (3*5)
	util.EXPECT().PrintGreen(gomock.Any(), gomock.Any(), gomock.Any()).Times(15)
	// Started (1*5)
	util.EXPECT().PrintGreen(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(5)

	docker.EXPECT().ImageExists(gomock.Any()).Return(true).Times(1)
	docker.EXPECT().HasDocker().Return(true).Times(1)
	docker.EXPECT().Kill(gomock.Any()).Return().Times(5)
	docker.EXPECT().Remove(gomock.Any()).Return().Times(5)
	docker.EXPECT().Build(gomock.Any(), gomock.Any()).Return().Times(5)

	gomock.InOrder(
		docker.EXPECT().Start("db", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("100"),
		docker.EXPECT().Inspect("100").Return([]byte("[{\"ID\": \"100\", \"State\":{\"Running\": true}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil),

		docker.EXPECT().Start("app", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("101"),
		docker.EXPECT().Inspect("101").Return([]byte("[{\"ID\": \"101\", \"State\":{\"Running\": true}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil),

		docker.EXPECT().Start("front1", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("102"),

		docker.EXPECT().Start("front2", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("103"),

		docker.EXPECT().Inspect("102").Return([]byte("[{\"ID\": \"102\", \"State\":{\"Running\": true}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil),
		docker.EXPECT().Inspect("103").Return([]byte("[{\"ID\": \"103\", \"State\":{\"Running\": true}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil),

		docker.EXPECT().Start("lb", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("104"),
		docker.EXPECT().Inspect("104").Return([]byte("[{\"ID\": \"104\", \"State\":{\"Running\": true}, \"NetworkSettings\": {\"IPAddress\": \"172.17.0.10\"}}]"), nil),
	)

	g := gaudi.Gaudi{}
	g.Init(`
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
`)

	g.StartApplications(true)
	c.Assert(len(g.Applications), Equals, 5)
}

func (s *GaudiTestSuite) TestCheckRunningContainerShouldUseDockerPs(c *C) {
	os.RemoveAll("/var/tmp/gaudi/templates/")

	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	// Setup the docker mock package
	docker.MOCK().SetController(ctrl)

	// Setup the util mock package
	util.MOCK().SetController(ctrl)

	// Disable the util package mock
	util.MOCK().DisableMock("IsDir")
	util.MOCK().DisableMock("IsFile")

	psResult := make(map[string]string)
	psResult["gaudi/lb"] = "123"
	psResult["gaudi/front1"] = "124"
	psResult["gaudi/db"] = "125"

	util.EXPECT().PrintGreen("Retrieving templates ...").Times(1)

	docker.EXPECT().ImageExists(gomock.Any()).Return(true).Times(1)
	docker.EXPECT().HasDocker().Return(true).Times(1)
	docker.EXPECT().SnapshotProcesses().Return(psResult, nil)

	docker.EXPECT().Inspect("123").Return([]byte("[{\"ID\": \"123\", \"State\":{\"Running\": true}, \"NetworkSettings\": {\"IPAddress\": \"123.124.125.126\"}}]"), nil)
	docker.EXPECT().Inspect("124").Return([]byte("[{\"ID\": \"123\", \"State\":{\"Running\": true}, \"NetworkSettings\": {\"IPAddress\": \"123.124.125.127\"}}]"), nil)
	docker.EXPECT().Inspect("125").Return([]byte("[{\"ID\": \"123\", \"State\":{\"Running\": true}, \"NetworkSettings\": {\"IPAddress\": \"123.124.125.128\"}}]"), nil)

	util.EXPECT().PrintOrange("Application", "lb", "is running", "(123.124.125.126:)")
	util.EXPECT().PrintOrange("Application", "front1", "is running", "(123.124.125.127:)")
	util.EXPECT().PrintOrange("Application", "db", "is running", "(123.124.125.128:3306)")

	g := gaudi.Gaudi{}
	g.Init(`
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
`)

	g.Check()
}

func (s *GaudiTestSuite) TestStartBinariesShouldCleanAndBuildThem(c *C) {
	os.RemoveAll("/var/tmp/gaudi/")

	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	// Setup the docker mock package
	docker.MOCK().SetController(ctrl)

	// Setup the util mock package
	util.MOCK().SetController(ctrl)

	// Disable the util package mock
	util.MOCK().DisableMock("IsDir")
	util.MOCK().DisableMock("IsFile")

	util.EXPECT().PrintGreen("Retrieving templates ...")

	docker.EXPECT().ImageExists(gomock.Any()).Return(true).Times(1)
	docker.EXPECT().HasDocker().Return(true).Times(1)

	util.EXPECT().PrintGreen("Building", "gaudi/npm", "...")
	docker.EXPECT().Build(gomock.Any(), gomock.Any()).Times(1)

	util.EXPECT().PrintGreen("Running", "npm", "update", "...")
	docker.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return().Times(1)

	g := gaudi.Gaudi{}
	g.Init(`
binaries:
    npm:
        type: npm
`)

	c.Assert(len(g.Applications), Equals, 0)
	c.Assert(len(g.Binaries), Equals, 1)

	g.Run("npm", []string{"update"})
}

func (s *GaudiTestSuite) TestUseCustomTemplateShouldUseIt(c *C) {
	os.RemoveAll("/var/tmp/gaudi/")

	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()
	docker.MOCK().SetController(ctrl)

	// Setup the util mock package
	util.MOCK().SetController(ctrl)

	util.MOCK().EnableMock("IsDir")
	util.MOCK().EnableMock("IsFile")
	util.MOCK().EnableMock("LogError")

	g := gaudi.Gaudi{}
	g.ApplicationDir = "/vagrant"

	docker.EXPECT().HasDocker().Return(true).Times(1)
	docker.EXPECT().ImageExists(gomock.Any()).Return(true).Times(1)

	util.EXPECT().IsDir("/var/tmp/gaudi/templates/").Return(false)
	util.EXPECT().IsFile("/vagrant/.gaudi/version.txt").Return(false)
	util.EXPECT().PrintGreen("Retrieving templates ...")

	util.EXPECT().IsFile("/vagrant/front/Dockerfile").Return(true).Times(3)
	util.EXPECT().LogError("Application 'custom' is not supported. Check http://gaudi.io/components.html for a list of supported applications.").Times(3)

	g.Init(`
applications:
    app:
        type: custom
        template: ./front/Dockerfile

    app2:
        type: custom
        template: front/Dockerfile

    app3:
        type: custom
        template: /vagrant/front/Dockerfile
`)
}

func (s *GaudiTestSuite) TestExtendsShouldCopyElements(c *C) {
	os.RemoveAll("/var/tmp/gaudi/")

	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()
	docker.MOCK().SetController(ctrl)

	util.MOCK().DisableMock("IsDir")
	util.MOCK().DisableMock("IsFile")

	disableLog()

	docker.EXPECT().ImageExists(gomock.Any()).Return(true).Times(1)
	docker.EXPECT().HasDocker().Return(true).Times(1)

	g := gaudi.Gaudi{}

	g.Init(`
applications:
    a:
        type: apache
        before_script: echo hello
    b:
        extends: a
    c:
        extends: a
        before_script: echo ok
    d:
        extends: c
        type: mysql
`)

	c.Check(g.Applications["a"].Type, Equals, "apache")
	c.Check(g.Applications["b"].BeforeScript, Equals, "echo hello")

	c.Check(g.Applications["b"].Type, Equals, "apache")
	c.Check(g.Applications["b"].BeforeScript, Equals, "echo hello")

	c.Check(g.Applications["c"].Type, Equals, "apache")
	c.Check(g.Applications["c"].BeforeScript, Equals, "echo ok")

	c.Check(g.Applications["d"].Type, Equals, "mysql")
	c.Check(g.Applications["d"].BeforeScript, Equals, "echo ok")

	enableLog()
}

func (s *GaudiTestSuite) TestExtendsShouldThrowAnErrorWhenTheElementDoesNotExists(c *C) {
	os.RemoveAll("/var/tmp/gaudi/")

	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()
	docker.MOCK().SetController(ctrl)

	util.MOCK().DisableMock("IsDir")
	util.MOCK().DisableMock("IsFile")
	util.MOCK().DisableMock("LogError")

	disableLog()

	g := gaudi.Gaudi{}

	c.Assert(func() {
		g.Init(`
applications:
    a:
        type: apache
    b:
        extends: c
`)
	}, PanicMatches, "b extends a non existing application : c")

	enableLog()
}

func (s *GaudiTestSuite) TestExtendsShouldCopyElementsOfNonOrderedComponent(c *C) {
	os.RemoveAll("/var/tmp/gaudi/")

	// Create a gomock controller, and arrange for it's finish to be called
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()
	docker.MOCK().SetController(ctrl)

	util.MOCK().DisableMock("IsDir")
	util.MOCK().DisableMock("IsFile")

	disableLog()

	docker.EXPECT().ImageExists(gomock.Any()).Return(true).Times(1)
	docker.EXPECT().HasDocker().Return(true).Times(1)

	g := gaudi.Gaudi{}

	g.Init(`
applications:
    c:
        extends: a
        before_script: echo ok
    a:
        extends: b
    b:
        type: apache
        before_script: echo hello
`)

	c.Check(g.Applications["a"].Type, Equals, "apache")
	c.Check(g.Applications["b"].BeforeScript, Equals, "echo hello")

	c.Check(g.Applications["b"].Type, Equals, "apache")
	c.Check(g.Applications["b"].BeforeScript, Equals, "echo hello")

	c.Check(g.Applications["c"].Type, Equals, "apache")
	c.Check(g.Applications["c"].BeforeScript, Equals, "echo ok")

	enableLog()
}
