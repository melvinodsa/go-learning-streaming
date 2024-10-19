# Streaming exercise

Trying to play around streaming with Go.

## Sequential Streaming

This program streams a file and stores it in buffer. Then writes the content to another file.
This implementation can be viewed in `sequential-streaming` branch

We can play around by change the knobs on delay and buffer size

## Concurrent Streaming

This program streams a file and stores it in queue. Concurrently the writer consumes from the queue.
This implementation can be viewed in `concurrent-streaming` branch

We can play around by change the knobs on delay
