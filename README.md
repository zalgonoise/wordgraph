# wordgraph
An implementation of a data graph, as part of [Kata19](http://codekata.com/kata/kata19-word-chains/)



___________

### Targets

- Building a chain of words with relationships between them, to allow querying for approximate, valid words
- With this word chain, creating logic to generate a chain where an origin word becomes the target, with single incremental steps, using all valid words. 
- Deliver the shortest chain (slice of strings) possible from the origin to the target

### Example

The word `cat` can become `dog` in four steps:

1. `cat`
1. `cot`
1. `cog`
1. `dog`

### Implementation

Breaking down the different modules, this implementation is based on a graph data structure that is non-cyclical, uni-directional, unweighted and map-based.

This approach will allow loading all the different words in the (defined) dictionary into memory, for more efficient queries.

Non-cyclical because words will not loop around eachother (e.g. the last character will not point to the first or middle character, suddenly); uni-directional because it will follow one direction (from beginning to the end, and never in reverse); unweighted because each step is a single unit (traversing through each letter of the word). Lastly, this implementation uses a map instead of a list of pointers (for instance) for faster access to the initialized keys. Although this would mean a bigger memory consumption, the list used (containing over 450k words) is loaded pretty quickly even so -- but this means that the root node contains a map for each character, that points to a new node containing a map to each character, etc.

As for traversing the graph, this implementation chooses a linear-direction approach (traverssing each word until its end, and exploring neighbours progressively in the same way). Although in retrospective I think it would be interesting to change this approach to a breadth-first search and comparing how it differs from one to the other.

To make up for the inevitable (potentially never-ending) recursion, this graph adds a few approaches that will reduce the query and output time:

1. Added a relevance system
1. Added goroutines and channels to distribute the load, on `FindRoute()`'s logic

##### Relevance system

This is a pretty simple approach to create a system to evaluate whether a word sits closest to the target or not. This is done by setting a weight unit size (100 divided by the number of characters in the word), where each matching character from the origin will award it one weight point. The word is seen as a match when the minimum accuracy value (of 98) is met. This makes up for situations like 3-letter words which would get a weight unit size of 33, and even matchin all 3 characters would never sum up to 100.

Apart from weight, there is also a potential value, which is the number of siblings (real words similar to the origin, with one changed character). The more this word can morph, the bigger the probability of finding a heavier-weight word. This metric isn't as relevant as the weight.

All retrieved results are passed through a quicksort function that will order them by weight, and then by potential, to ensure that the most relevant keywords are explored first.


##### `FindRoute()` approach

This section is mostly focused on creating a communication framework for all spawned goroutines. As the many iterations could potentially consume a long time, they were set to be called as goroutines.

As channels are inexpensive in Go, the spawned goroutines will collect each valid route and send it to a filter method (`findBestRoute()`). This method will listen to a done channel and the responses channel, as it accumulates routes from the running goroutines.

To stop the routines, either:

1. The overall query timer runs out; where the shortest existing route is returned
1. The done signal is tripped; where the shortest existing route is returned
1. If it takes more than the maximum no-response time limit, between receiving routes; where it trips the done signal
1. If the routes list exceeds the maximum limit of routes to evaluate; where it trips the done signal

This makes queries overall consistent, correct and reliable. I noticed however that it still varies greatly in time (for the same query, may be sometimes faster, sometimes slower). Some notes on this topic in the last section, below.



##### Final thoughts

I really enjoyed working on this particular project, it was very fun. As a last note I noticed that I missed a simple technique to quickly find a heavier-weight match, which is by swapping individual letters (from the same index) between the origin and target, and checking if that word exists. It should provide a quicker set of results while having the weight / potential system as an _exploration_ approach.