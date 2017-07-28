# envtoconf

Striving to be a simple env program to process go templates primarily used for docker containers to read environment variables at runtime. The funcMap() for the template processor comes from [sprig](https://github.com/Masterminds/sprig) which gives you all kinds of functions if you wish to use them.

Below is some exmaples of using environment variables in a template.

``` bash
home is :{{ env "HOME" }} endhome.
path is :{{ expandenv "Your path is set to $PATH" }} endpath.

```
