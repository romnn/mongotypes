## mongotypes

[![Build Status](https://travis-ci.com/romnn/mongotypes.svg?branch=master)](https://travis-ci.com/romnn/mongotypes)
[![GitHub](https://img.shields.io/github/license/romnn/mongotypes)](https://github.com/romnn/mongotypes)
[![GoDoc](https://godoc.org/github.com/romnn/mongotypes?status.svg)](https://godoc.org/github.com/romnn/mongotypes)  [![Test Coverage](https://codecov.io/gh/romnn/mongotypes/branch/master/graph/badge.svg)](https://codecov.io/gh/romnn/mongotypes)
[![Release](https://img.shields.io/github/release/romnn/mongotypes)](https://github.com/romnn/mongotypes/releases/latest)

Provides types for [go.mongodb.org/mongo-driver](https://github.com/mongodb/mongo-go-driver) that can be used to construct and decode mongodb responses.

**Note**: Currently, only types needed for configuration of mongodb replicasets are provided.

**Note**: If you are looking for types for [gopkg.in/mgo.v2](https://github.com/go-mgo/mgo), have a look at [juju/replicaset](https://github.com/juju/replicaset).

#### Example

Lets say you want to run the `replSetGetStatus` command to check the primary of your replicaset. You can use `replicaset.Status` with `Decode` to parse the raw bson result.

```golang
import "github.com/romnn/mongotypes/replicaset"

var statusResult replicaset.Status

// adminDatabase is a connected mongo db admin database
if err := adminDatabase.RunCommand(context.TODO(), bson.D{{"replSetGetStatus", nil}}).Decode(&statusResult); err != nil {
    log.Fatal(err)
}
primary := statusResult.Primary()
if primary == nil {
    log.Fatal("The replicaset has no primary")
}
```

For more examples, see `examples/`.


#### Development

######  Prerequisites

Before you get started, make sure you have installed the following tools::

    $ python3 -m pip install -U cookiecutter>=1.4.0
    $ python3 -m pip install pre-commit bump2version invoke ruamel.yaml halo
    $ go get -u golang.org/x/tools/cmd/goimports
    $ go get -u golang.org/x/lint/golint
    $ go get -u github.com/fzipp/gocyclo
    $ go get -u github.com/mitchellh/gox  # if you want to test building on different architectures

**Remember**: To be able to excecute the tools downloaded with `go get`, 
make sure to include `$GOPATH/bin` in your `$PATH`.
If `echo $GOPATH` does not give you a path make sure to run
(`export GOPATH="$HOME/go"` to set it). In order for your changes to persist, 
do not forget to add these to your shells `.bashrc`.

With the tools in place, it is strongly advised to install the git commit hooks to make sure checks are passing in CI:
```bash
invoke install-hooks
```

You can check if all checks pass at any time:
```bash
invoke pre-commit
```

Note for Maintainers: After merging changes, tag your commits with a new version and push to GitHub to create a release:
```bash
bump2version (major | minor | patch)
git push --follow-tags
```

#### Note

This project is still in the alpha stage and should not be considered production ready.
