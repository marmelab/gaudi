package maestro_funcitonal_test

import (
	. "launchpad.net/gocheck"
	"testing"

	"github.com/marmelab/gaudi/maestro"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func Test(t *testing.T) { TestingT(t) }

type MaestroTestSuite struct{}

var _ = Suite(&MaestroTestSuite{})

// Apache
func (s *MaestroTestSuite) TestStartApacheShouldStartedItCorrectly(c *C) {
	m := maestro.Maestro{}
	m.InitFromString(`
applications:
    front:
        type: apache
        ports:
            80: 80
`, "")

	c.Assert(len(m.Applications), Equals, 1)
	m.Start(true)

	// Test apache is running
	resp, err := http.Get("http://" + m.GetContainer("front").Ip)
	defer resp.Body.Close()

	c.Check(err, Equals, nil)
	c.Check(resp.StatusCode, Equals, 200)
}

// Apache + php-fpm
func (s *MaestroTestSuite) TestStartPhpAndApacheShouldStartedThemCorrectly(c *C) {
	err := os.MkdirAll("/tmp/php", 0775)
	ioutil.WriteFile("/tmp/php/ok.php", []byte("<?php echo 'ok';"), 0775)

	m := maestro.Maestro{}
	m.InitFromString(`
applications:
    front:
        type: apache
        links: [app]
        ports:
            80: 80
        volumes:
            /tmp/php: /var/www
        custom:
            fastCgi: app

    app:
        type: php-fpm
        ports:
            9000: 9000
        volumes:
            /tmp/php: /var/www
`, "")

	c.Assert(len(m.Applications), Equals, 2)
	m.Start(true)
	time.Sleep(2 * time.Second)

	// Test apache is running
	resp, err := http.Get("http://" + m.GetContainer("front").Ip + "/ok.php")
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)

	c.Check(err, Equals, nil)
	c.Check(resp.StatusCode, Equals, 200)
	c.Check(string(content), Equals, "ok")
}
