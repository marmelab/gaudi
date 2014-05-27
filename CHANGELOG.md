# Version 0.1

- gaudi can now be installed via apt-get
- [Ambassadors](http://marmelab.com/blog/2014-05-12-gaudi-news-gaudi-io-apt-get-install)
- New [phpMyAdmin container](http://marmelab.com/blog/2014-05-12-gaudi-news-gaudi-io-apt-get-install)
- Redirect `stdin` & `stdin` to console during docker build
- Force rebuild when configuration file changes
- Fixes #56, #48, #47

# Version 0.2

## gaudi

- Allows to use custom templates.
- New components [`python`](https://github.com/marmelab/gaudi/tree/master/example/python), [`django`](https://github.com/marmelab/gaudi/tree/master/example/django) & [`golang`](https://github.com/marmelab/gaudi/tree/master/example/golag)
- Add an `empty-cmd` option to debug containers
- Use color in gaudi output
- Updated [examples](https://github.com/marmelab/gaudi/tree/master/example) with use cases
- Add a check of gaudi version to retrieve new templates
- Fix #47, #59, #58

## Builder

- Use browserify to include Javascript files
- Add unit tests
