# Chanson
Stream Json Objects in Go

# Why
Let's say you want to json-encode a large amount of data from your database and send it through the network.
So you pack the data into a structure and then you call `json.NewEncoder(w).Encode(theData)`, but this has some implications:
_a)_ If the data is to big, you will end up using a massive amount of memory and
_b)_ your client will need to wait all that time before receiving a single bit.

# Usage (TODO)
