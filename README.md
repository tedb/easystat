Easystat for Golang
===

Very easily record statistics for your program, in a concurrency-safe way

Copyright Licensing
---

See enclosed LICENSE.txt file.

Example Usage
---

The following will print sums of foo = 5 and bar = 10 to stderr:

	stats := easystat.NewWriter(os.Stderr, 500*time.Millisecond)
	stats.Add("foo", 1)
	stats.Add("bar", 2)
	stats.Add("foo", 4)
	stats.Add("bar", 8)
	stats.Stop()

