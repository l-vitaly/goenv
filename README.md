# Go Env

### Example

[See](https://github.com/l-vitaly/goenv/example/README.md)

#### Generate from env file

Install:

go get -u github.com/l-vitaly/goenv/cmd/envgen

To generate a config based on the environment variable file, run the following command:

`envgen <env_file> <prefix>`

-   env_file - environment variable file
-   prefix - this value will be removed from the field names of the structure and constants

You can also group the fields in the structure by separating them with `__` characters. For example, if the name of the environment variable is `DB__CONN_SRT`, then the result will be like this:

```golang
type Config {
  Db {
    ConnSrt string
  }
}
```
