# `hello`

```console
$ go build -o hello pkg/command/testdata/examples/hello/main.go
```

```console
$ ./hello 
Usage: hello [--version] [--help] <command> [<args>]

Available commands are:
    say    Say hello
```

```console
$ ./hello say
missing positional argument(s): expected 1, given 0

Usage:
         say <who> [options]
Description:
        Say hello
Flags:
        --cowboy   say hello like a cowboy does (default: false) 
        --help     Print this help menu (default: false)
```

```console
$ ./hello say kent
👋 hello kent
```

```console
$ ./hello say kent --cowboy
🤠 howdy kent
```