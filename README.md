# Hashcash

[![GoDoc](https://img.shields.io/badge/reference-007d9c.svg?logo=go&logoColor=white&label=doc)](https://pkg.go.dev/github.com/PoW-HC/hashcash)
[![Go Report Card](https://goreportcard.com/badge/github.com/PoW-HC/hashcash)](https://goreportcard.com/report/github.com/PoW-HC/hashcash)
[![Lint](https://github.com/PoW-HC/hashcash/actions/workflows/lint.yml/badge.svg)](https://github.com/PoW-HC/hashcash/actions/workflows/lint.yml)
[![Tests](https://github.com/PoW-HC/hashcash/actions/workflows/tests.yml/badge.svg)](https://github.com/PoW-HC/hashcash/actions/workflows/tests.yml)
![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/halfi/cce912d48b587ba656b45a0cba34510b/raw/pow-hc-test-coverage.json)
[![Code quality](https://github.com/PoW-HC/hashcash/actions/workflows/codeql.yml/badge.svg)](https://github.com/PoW-HC/hashcash/actions/workflows/codeql.yml)
[![license](https://img.shields.io/github/license/PoW-HC/hashcash)](https://github.com/PoW-HC/hashcash/blob/main/LICENSE)

---

Hashcash is a Go library which implements the hashcash [proof-of-work](https://en.wikipedia.org/wiki/Proof_of_work)
algorithm. Hashcash has been used as a denial-of-service counter measure technique in a number of systems. To learn more
about hashcash visit [official page](http://hashcash.org/).

## Usage

Computing a valid hashcash:

```go
package main

import (
	"context"
	"fmt"

	"github.com/PoW-HC/hashcash/pkg/hash"
	"github.com/PoW-HC/hashcash/pkg/pow"
)

const maxIterations = 1 << 30

func main() {
	hasher, err := hash.NewHasher("sha256")
	if err != nil {
		// handle error
	}

	p := pow.New(hasher)

	hashcash, err := pow.InitHashcash(5, "127.0.0.1", pow.SignExt("secret", hasher))
	if err != nil {
		// handle error
	}

	solution, err := p.Compute(context.Background(), hashcash, maxIterations)
	if err != nil {
		// handle error
	}
	fmt.Println(solution)
}

```

Outputs:

```
1:5:1649257375:127.0.0.1:41965c8500f67f79b35672d7fb7e19fe5af0d51da582cbfcdbedb0f5944198bb:MgjmhCBiRV4=:MTk3YzA0
```

```shell
echo -n "1:5:1649257375:127.0.0.1:41965c8500f67f79b35672d7fb7e19fe5af0d51da582cbfcdbedb0f5944198bb:MgjmhCBiRV4=:MTk3YzA0" | shasum -a 256
0000067c716fa612cee5eb31eab6163e804c002841cfe96cc81f4f8034fd6006  -
```

Verification token:

```go
	err := p.Verify(solution, "127.0.0.1")
	if err != nil {
		// hashcash token failed verification.
	}
```

## Documentation

https://pkg.go.dev/github.com/PoW-HC/hashcash

## License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).
