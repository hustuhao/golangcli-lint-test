# golangcli-lint-test


Linter sqlclosecheck produces different result in vet and sqlclosecheck：

While I run such two command on the root of my project seperately, I get different results.

`go vet -vettool=$(which sqlclosecheck)  ./... `

`golangci-lint run --disable-all --enable=sqlclosecheck ./... `

The result seems that vet scans more directories and files than golangci-lint does , so I get more lines in vet.
