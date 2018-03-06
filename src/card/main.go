package main

import (
	"bufio"
	"fmt"
	"hash"
	"os"

	"github.com/clarkduvall/hyperloglog"
	"github.com/spaolacci/murmur3"
	cli "gopkg.in/urfave/cli.v1"
)

func hashIt(value []byte) hash.Hash64 {
	h := murmur3.New64()
	h.Write(value)
	return h
}

func execute(c *cli.Context) error {
	p := uint8(c.Uint("precision"))
	if p < 4 || p > 16 {
		return cli.NewExitError("p must be 4 >= x <= 16", 22)
	}

	hll, _ := hyperloglog.NewPlus(p)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		value := []byte(scanner.Text())
		hll.Add(hashIt(value))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	fmt.Fprintln(os.Stdout, hll.Count())
	return nil
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.UintFlag{
			Name:  "precision, p",
			Value: 14,
			Usage: "Precision to use (must be 4 >= x <= 16)",
		},
	}
	app.Name = "card"
	app.Usage = "Estimates cardinality of stdin using HyperLogLog++ algorithm"
	app.Action = execute
	app.Run(os.Args)
}
