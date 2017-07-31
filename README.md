# envtoconf

[![Go Report Card](https://goreportcard.com/badge/github.com/benoahriz/envtoconf)](https://goreportcard.com/report/github.com/benoahriz/envtoconf)

Striving to be a simple env program to process go templates primarily used for docker containers to read environment variables at runtime. The funcMap() for the template processor comes from [sprig](https://github.com/Masterminds/sprig) which gives you all kinds of functions if you wish to use them.

Below is some exmaples of using environment variables in a template.

``` bash
home is :{{ env "HOME" }} endhome.
path is :{{ expandenv "Your path is set to $PATH" }} endpath.

```


TODO:

  Create tests for malformed template
  Write article about going from an idea to production in golang.
  Create option for required vars strict mode.
