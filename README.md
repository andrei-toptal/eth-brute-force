# yet another ethereum brute force tool

A tiny command line app to find the private key matching the given ETH address.

## usage

`go run main.go 0x1234567890123456789012345678901234567890`

where `0x1234567890123456789012345678901234567890` is the ethereum address.

## why?

Every other blockchain engineer implements this exercise, for fun.

If you are really serious about using this tool, then bear in mind the probability of encountering a private key 
that corresponds to someone else’s Ethereum address is around 1 in 2^256. To cover just 1% of that key space, 
even if we used computing resources that would allow us to generate 100 trillion keys per second, it would take us roughly years.
Eventually, your electricity bills and the cost of hardware you used will exceed the benefit of "hacking" someone's address.
You are warned.
