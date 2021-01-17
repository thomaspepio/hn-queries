# Thought process

### Looking at the data
A quick glance at `hn_logs.tsv` shows that each line of the file is structured as such : `<YYYY-MM-DD HH:mm:SS><tab><url> `

We choose to not take extra care of any line that would not respect this structure : it will be discarded by any parser we will write.

### Designing the data structure that will support the API (round 1) : 

   - GET /1/queries/count/<DATE_PREFIX> :
      - INPUT  : year | year-month | year-month-day
      - OUTPUT : number of requests
   
   - GET /1/queries/popular/<DATE_PREFIX>?size=<SIZE>
      - INPUTS : year | year-month | year-month-day, size
      - OUTPUT : list of queries

   - Design :
     - for the date parameter, three are search patterns : `YYYY`, `YYYY-MM` and `YYYY-MM-DD`. Response time of the APIs should not vary too much whether the search targets a specific day or a whole year.
     - we know that balanced binary trees offer _O(log n)_ search. Our problem is that we have three different kind of keys. A naïve approach would see us use three kidns of binary trees, one per search pattern, nesting related tree into one another.
     - this naïve approach seems complicated to implement and to maintain. We can solve this issue by making the three patterns comparable by padding _zeroes_ when necessary (`2020` becomes `20200000`, `202001` becomes `20200100`, `20200101` does not change).
     - this less naïve approach yields 1095 keys, for each year we want to index (indexing every HN search since it's birth would require 14k keys).
     - URLs can be deduplicated and referenced by a _primary key_, tree nodes can reference those keys to avoid duplication.
     - retrieving the Nth most popular querries can be managed by updating query frequency at the node level.

Given the following URL log :
```
2021-01-14  Foo
2021-02-01  Bar
2021-03-01  Foo
2020-10-28  Baz
2020-06-01  Foo
2019-02-10  Bar
```

We extract a deduplicated table of URLs :
```
Foo  1
Bar  2
Baz  3
```

And the following search indexes (examples covering 2021 only ):
```
20210000    // the whole 20201 year
20210100    // january 20201
20210200    // february 2021
20210300    // march 2021
20210114    // jan 14th 2021
20210201    // etc.
...
```

And use them as keys in the following binary tree :
```
                                             (20210300, [(1,1)])
                                             /                  \
                          (20210100, [(1, 1)])                   (20210114, [(1, 1)])
                          /                   \                                      \
(20210000, [(1, 2), (2, 1)]                    (20210200, [(2, 1)]                    (20210201, [(2, 1)])
```

In this tree, each value is a list or pairs `(KEY_ID, COUNT)`, which should be maintained ordered by COUNT to meet the performance requirements of the `/1/queries/popular/` API.

### Designing : round 2
Hashtables support search, insert and delete operations with O(1) time, which beats O(log n).