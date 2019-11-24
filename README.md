DEPRECATED
==========
The novendor check was originally developed to ensure the consistency of the vendor directory for projects. This was
necessary when the vendor directory was managed manually or using basic tools such as [govendor](https://github.com/kardianos/govendor).
However, tools such as [dep](https://github.com/golang/dep) and the support for [modules in Go](https://blog.golang.org/using-go-modules) 
have made the functionality provided by this tool unnecessary. As such, active development on this project has ended.

godel-okgo-asset-novendor
=========================
godel-okgo-asset-novendor is an asset for the gödel [okgo plugin](https://github.com/palantir/okgo). It provides the
functionality of the [go-novendor](https://github.com/palantir/go-novendor) check.

This check verifies that unused packages are not vendored.
