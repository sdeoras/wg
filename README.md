# wg
wg stands for wait-group. It is a shell command that executes other shell commands asynchronously.

# installation
```bash
go get github.com/sdeoras/wg/wg
```

# example
```bash
# initialize wg and send to background
wg init &

# run some commands
wg run -- sleep 1
wg run -- sleep 1
wg run -- sleep 1

# wait for command completion
wg wait

# note that three sleep commands ran asynchronously and total wait time was only 1 second
```

# server and a client
`wg init` starts a server and `wg run` starts a client. It is, therefore, necessary to run `wg init`
in background and to call `wg wait` to not only wait for command completion but also to kill the server.

# limitations
* server is not designed to run in multiple instances. `wg init` will fail if a server is already running.