package gaudi_functional_test

import (
	. "launchpad.net/gocheck"
	"testing"

	"github.com/marmelab/gaudi/gaudi"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func Test(t *testing.T) { TestingT(t) }

type GaudiTestSuite struct{}

var _ = Suite(&GaudiTestSuite{})

// Apache
func (s *GaudiTestSuite) TestStartApacheShouldStartedItCorrectly(c *C) {
	g := gaudi.Gaudi{}
	g.Init(`
applications:
    front:
        type: apache
        ports:
            80: 80
`)

	c.Assert(len(g.Applications), Equals, 1)
	g.StartApplications()

	// Test apache is running
	resp, err := http.Get("http://" + g.GetApplication("front").Ip)
	defer resp.Body.Close()

	c.Check(err, Equals, nil)
	c.Check(resp.StatusCode, Equals, 200)
}

// Apache + php-fpm
func (s *GaudiTestSuite) TestStartPhpAndApacheShouldStartedThemCorrectly(c *C) {
	err := os.MkdirAll("/tmp/php", 0775)
	ioutil.WriteFile("/tmp/php/ok.php", []byte("<?php echo 'ok';"), 0775)

	g := gaudi.Gaudi{}
	g.Init(`
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
`)

	c.Assert(len(g.Applications), Equals, 2)
	g.StartApplications()
	time.Sleep(2 * time.Second)

	// Test apache is running
	resp, err := http.Get("http://" + g.GetApplication("front").Ip + "/ok.php")
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)

	c.Check(err, Equals, nil)
	c.Check(resp.StatusCode, Equals, 200)
	c.Check(string(content), Equals, "ok")
}
