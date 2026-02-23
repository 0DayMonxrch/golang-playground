## A Minimal Logging tool written in Golang.

### Features:
- Writes logs as JSON (structured logging)

- Runs file writes in a dedicated goroutine

- Prevents request blocking

- Supports backpressure handling

- Gracefully drains logs on shutdown

- Uses WaitGroup to prevent premature exit

### Future Additions:
- Log rotation (size/time-based)

- File locking (if multi-process)

- Structured fields

- Context integration (request ID)

- Sync on panic

- Protection against log flooding/Disk Exhaustion Attacks

