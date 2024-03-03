# Day 09

<h3 id="ex00">Exercise 00: Sleepsort</h3>

So, let's write a toy algorithm as an example. It is pretty much useless for production, but it helps to grasp the concept.

You need to write a function called `sleepSort` that accepts an unsorted slice of integers and returns an integer channel, where these numbers will be written one by one in sorted order. To test it, in main goroutine you should read and print output values from a returned channel and gracefully stop the application in the end.

The idea of Sleepsort (what makes it a "toy") is that we're spawning a number of goroutines equal to the size of an input slice. Then, each goroutine sleeps for amount of seconds equal to the received number. After that it wakes up and this number to the output channel. It's easy to notice that this way numbers will be returned in a sorted order.

<h3 id="ex01">Exercise 01: Spider-Sense</h3>

You probably remember how Peter Parker realised he now has superpowers when he woke up in the morning. Well, let's write our own spider (or crawler) for parsing the web. You need to implement a function `crawlWeb` which will accept an input channel (for sending in URLs) and return another channel for crawling results (pointers to web page body as a string). Also, at any moment in time there shouldn't be more than 8 goroutines querying pages in parallel.

But we want to be quick and flexible, so another requirement is to be able to stop the crawling process at any time by pressing Ctrl+C (and your code should perform a graceful shutdown). For that you may add more input arguments to `crawlWeb` function, which should be either `context.Context` for cancellation or `done` channel. If not interrupted, the program should gracefully stop after all given URLs are processed.

<h3 id="ex02">Exercise 02: Dr. Octopus</h3>

Okay, so now we have to slain the villain! The main problem with Dr.Octopus is that he has a lot of tech tentacles, and it's hard to keep track of them all. Let's tie them together!

For this exercise, you need to write a function called `multiplex` which should be *variadic* (accepts a variable number of arguments). It should accept channels (`chan interface{}`) for arguments and return a single channel of the same type. The idea is to redirect any incoming messages from these input channels into just one single output channel, effectively implementing "fan-in" pattern.

As a proof of work, you should write a test on sample data which will explicitly show that all values randomly sent to any input channels are received further on the same output channel.

And...you've just defeated a villain!