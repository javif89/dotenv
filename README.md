# dotenv

Programatically manipulate .env files:
- Add/remove keys
- Change values
- Check the difference between two .env files (ex your .env and .env.example)
- Merge two .env files
- Format a .env file (uppercase keys with _, quote multi word values, etc)

You can use this as a go module or use the cli tool.

# Use case

I made this as part of a larger project I'm working on and thought it would be easier to maintain as a separate package and it might be useful to someone else.

# Installation

### Go package

```bash
go get github.com/javif89/dotenv
```

[Go quickstart](#go-quickstart-example)

### CLI Utility

```bash
brew install javif89/javif89/dotenv
```

[CLI quickstart](#cli-quickstart-example)

# Go quickstart example

```go
package main

import (
    "github.com/javif89/dotenv"
)

func main() {
    file := dotenv.New("./mypath/.env")

    file.Add("DB_USERNAME", "root")
    file.Add("DB_PASSWORD", "password")

    file.Save()
}
```

This would create a file in ./mypath/.env with the contents:

```dotenv
DB_USERNAME=root
DB_PASSWORD=password
```

### Editing existing files

Modifying existing files is just as easy. If we have a file that looks like:

```dotenv
DB_USERNAME=root
DB_PASSWORD=password
```

We can do:

```go
package main

import (
    "github.com/javif89/dotenv"
)

func main() {
    file := dotenv.Load("./.env")

    file.Add("DB_HOST", "127.0.0.1")

    file.Save()
}
```

Now the file will look like this:

```dotenv
DB_USERNAME=root
DB_PASSWORD=password
DB_HOST="127.0.0.1"
```

### Modifying values

```go
package main

import (
    "github.com/javif89/dotenv"
)

func main() {
    file := dotenv.Load("./.env")

    file.Set("DB_HOST", "localhost")

    file.Save()
}
```

Result:

```dotenv
DB_USERNAME=root
DB_PASSWORD=password
DB_HOST=localhost
```

# CLI

**Help text**

```
NAME:
   dotenv - Create and manipulate .env files in your system

USAGE:
   dotenv [path to file] command [command options]

AUTHOR:
   Javier Feliz

COMMANDS:
   set, s   Set an environment variable
   get, g   Get the value of a key in a file
   keys, k  List all the keys in a file
   fmt      Format the env file and fix any issues
   diff     Show the difference between two files
   merge    Merge file 2 into file 1
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

# CLI Quickstart Example

```shell
dotenv set -f .env -k DB_USERNAME -v root
dotenv set -f .env -k DB_PASSWORD -v password
```

This would create a file in ./.env with the contents:

```dotenv
DB_USERNAME=root
DB_PASSWORD=password
```

### Modifying values

```shell
dotenv set -f .env -k DB_PASSWORD -v newpassword
```

Result:

```dotenv
DB_USERNAME=root
DB_PASSWORD=newpassword
```