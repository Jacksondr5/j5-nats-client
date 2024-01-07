TOOD

- wake on lan to start the nas back up when the power comes back on
- Instead of just shutting down the pis when the last one acks, use the acks to determine how many pis are responsive. Then poll the responsive pis with https://pkg.go.dev/github.com/tatsushid/go-fastping#section-readme to see if they're still alive. If they're not, then shut down the rest of the pis.
  - Set a timeout on the ping of 10 seconds. That's around what's needed to let them finish shutting down after they stop responding.
  - This will probably need lots of goroutines, see if its necessary to make things multithreaded. Go is singlethreaded by default.
