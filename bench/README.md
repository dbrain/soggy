# Soggy Benchmarks
These were run on a 13" 2011 MacBook Air with an i7 processor. The server and ab are run on the same laptop, so figures aren't a great example of actual performance.. more just comparisons between previous versions. basicserver.bench is a hello world for Go's defualt http server, compare soggy numbers to these.

I figure running these and committing will give me some kind of historical information so I can see if new changes have slowed soggy down significantly. In future I'll also have benchmark tests for the code paths.

Have a look at benchmark.sh and benchmark.go to see what maps to what .bench file.
