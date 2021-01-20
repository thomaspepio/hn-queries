# Home Assignment

### Project manipulation & layout

#### Manipulation

All commands are to be run at the root of the project.

##### Running the tests
   - with code coverage analysis : `go test --coverprofile=coverage.out ./... && go tool cover -func=coverage.out` 
   - without code coverage : `go test ./...`

##### Building the project
`go get && go build`

##### Running the app
1. copy the `hn_logs.tsv` file at the root of the project
2. launch the binary : `./hn-queries`

A server should start on `localhost:8080`.

#### Layout
- _avltree_ : almost complete implementation of an AVL tree (the delete operation is not supported)
- _constant_ : stores values used across multiple packages
- _endpoint_ : API endpoints configuration and http parameters management
- _index_ : main indexing structure
- _parser_ : typed representation of a log line and its parser
- _query_ : queries the API supports, the unique call point for endpoints
- _util_ : utility functions used across multiple packages

### Analysis 
#### 1. Looking at the data
A quick glance at `hn_logs.tsv` shows that each line of the file is structured as such : `<YYYY-MM-DD HH:mm:SS><tab><url> `

We choose to not take extra care of any line that would not respect this structure : it will be discarded by the parser.

#### 2. Design

- GET /1/queries/count/<DATE_PREFIX> :
   - INPUT  : year | year-month | year-month-day | year-month-day hour:minute
   - OUTPUT : number of requests

- GET /1/queries/popular/<DATE_PREFIX>?size=<SIZE>
   - INPUTS : year | year-month | year-month-day | year-month-day hour:minute, size
   - OUTPUT : list of queries

We want both APIs responses time to be fast, whether we search targeting a specific minute or a whole year :
   - Reading https://www.bigocheatsheet.com/, it's tempting to go for a hashmap be cause it has _O(1)_ average search time. But our APIs supports range searches, which binary search trees are better at.
   - We choose to go for an AVLTree : because it's a self balancing BST, it offers _O(log n)_ for all scenarios.
   - Six keys are extracted from a date. For instance, given the date _2015-08-01 00:03:50_, we extract six keys : 
  
      | key            | reference                  |
      | -------------- | -------------------------- |
      | 20150000000000 | year 2015                  |
      | 20150800000000 | month 2015-08              |
      | 20150801000000 | day 2015-08-01             |
      | 20150801010000 | hour 2015-08-01 00         |
      | 20150801010400 | minute 2015-08-01 00:03    |
      | 20150801010451 | second 2015-08-01 00:03:50 |

      This design maps each request to a single node in the tree, which should lead to fast response time.

      _(note that we add 1 to the hour/minute/second to avoid conflicting keys. This preserves the order that exists between dates)_

   - This design regarding the key does not invalidate the choice for an AVL Tree : searches are sped up for _"whole"_ intervals (e.g. : the whole 2015 year, a whole month, a whole day, a whole minute...), and look just like hashmap lookups, but ranged searches are still required for _"overlapping"_ intervals<sup>1</sup> (e.g. between 2021-01-01 00:01:30 and 2021-01-03 00:00:00).

   - Since we expect URLs to be duplicated in the log file, our data structure will maintain an index or URLs (a map URL -> ID)

#### 3. Concerns
At the eve of returning this home assignment, I'm concerned that the choice of extracting six keys for each date poses a huge memory problem.

The other design I could have opted for would have been to simply use the dates as keys, but this would have increased the search response time (since each node no longer holds the whole information required to answer). And I believe the fastest response times are of paramount importance at Algolia.

I've done a bit of calculation to measure the order of magnitudes at which both of these designs operate memory wise, extrapolating whole years of indexed data from the log file, and could not find a deciding factor for one or the other.

##### Footnotes
<sup>1</sup> : granted, this falls under the category of over-engineering and is not required to complete the assignment.