package main

import (
	"bufio"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/clarkduvall/hyperloglog"
	"github.com/spaolacci/murmur3"
	cli "gopkg.in/urfave/cli.v1"
)

const PRECISION_STR = "4 >= p <= 16"
const EINVAL_INVALID_ARG = 22
const ERR = 255

func execute(c *cli.Context) error {
	precision := uint8(c.Uint("precision"))
	if !isValidPrecision(precision) {
		return cli.NewExitError(fmt.Sprintf("precision must be %s", PRECISION_STR), EINVAL_INVALID_ARG)
	}

	file, err := openFile(c.Args().Get(0))
	if err != nil {
		return cli.NewExitError(err, ERR)
	}
	defer file.Close()

	estimate, err := estimateCardinality(precision, file)
	if err != nil {
		return cli.NewExitError(err, ERR)
	}

	fmt.Println(estimate)
	return nil
}

func isValidPrecision(p uint8) bool {
	if p < 4 || p > 16 {
		return false
	}

	return true
}

func openFile(filename string) (*os.File, error) {
	if filename == "" || filename == "-" {
		return os.Stdin, nil
	}

	return os.Open(filename)
}

func estimateCardinality(precision uint8, reader io.Reader) (uint64, error) {
	hll, err := hyperloglog.NewPlus(precision)

	if err != nil {
		return 0.0, err
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		value := []byte(scanner.Text())
		hll.Add(hashValue(value))
	}

	if err := scanner.Err(); err != nil {
		return 0.0, err
	}

	return hll.Count(), nil
}

func hashValue(value []byte) hash.Hash64 {
	h := murmur3.New64()
	h.Write(value)
	return h
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.UintFlag{
			Name:  "precision, p",
			Value: 16,
			Usage: fmt.Sprintf("precision to use (must be %s)", PRECISION_STR),
		},
	}
	app.Name = "card"
	app.Usage = "Estimates cardinality (count-distinct) of inputs using HyperLogLog++ algorithm"
	app.ArgsUsage = "file to read from (optional - defaults to stdin)"
	app.UsageText = "card [-p] <file>\n\n\t if <file> is not provided, stdin will be used instead"
	app.Action = execute
	app.Run(os.Args)
}
