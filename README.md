# Go Env

### Example

``` golang
package main

import "github.com/l-vitaly/goenv"

const (
  DBConnEnvName = "DB_CONN"
)

type config struct {
  DBConn string
}

func main() {
  cfg := &config{}

  goenv.StringVar(cfg.DBConn, DBConnEnvName, "")
  goenv.Parse()
}
```