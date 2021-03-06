---
title: Reference
weight: 4
any: false
---
# Python Authorization Library

The Python version of oso is available on [PyPI](https://pypi.org/project/oso/) and can be installed using
`pip`:

```
$ pip install oso==0.9.0
```

To install Python framework integrations, see:


* [Flask](./flask)
* [Django](./django)
* [SQLAlchemy](./sqlalchemy)


**Requirements**


* Python version 3.6 or greater
* Supported platforms:
   * Linux
   * OS X
   * Windows

The Python version is known to work on glibc-based distributions but not on musl-based ones
(like Alpine Linux).  Wheels built against musl that you can use on
Alpine Linux can be downloaded from [the releases page on GitHub](https://github.com/osohq/oso/releases/latest).

