# card

Command line tool that estimates the cardinality (count-distinct) of inputs using HyperLogLog++ algorithm.

See https://djhworld.github.io/hyperloglog for more information.


## Installation

```bash
go get -u github.com/djhworld/card
```

## Usage

Using stdin:

```bash
cat items.txt | card
```

Using file:

```bash
card items.txt
```

## Performance

The file `items.txt` has 2,000,000 items in it, and 1,000,000 of them are unique.

```bash
djhworld/card|master ✗ ▶ wc -l items.txt
2000000 items.txt
```

Using `sort | uniq`, this takes roughly 25 seconds to complete on my machine. 

```bash
djhworld/card|master ✗ ▶ time sort items.txt | uniq | wc -l
 1000000
sort items.txt  25.88s user 0.33s system 95% cpu 27.336 total
uniq  1.56s user 0.03s system 5% cpu 27.335 total
wc -l  0.01s user 0.01s system 0% cpu 27.334 total
```

Using `card` to get an approximate count, this takes about 0.5 seconds

```bash
djhworld/card|master ✗ ▶ time card items.txt
999149
card items.txt  0.51s user 0.04s system 98% cpu 0.556 total
```

Most of the `sort | uniq` time is spent during the `sort` phase, so it's not a fully fair comparison, you could use something like my other tool [count](https://github.com/djhworld/count) that does one pass over the input. 

However, this could be memory intensive if the cardinality of your input set is high, `card` has a very minimal memory footprint.
