all: clean
	gb build all

clean:
	rm -rf bin

deps:
	gvt fetch -tag v1.19.1 github.com/urfave/cli
	gvt fetch -revision a0107a5d80040c9b5efb012203b9a08d6817d196 github.com/clarkduvall/hyperloglog
	gvt fetch -revision f09979ecbc725b9e6d41a297405f65e7e8804acc github.com/spaolacci/murmur3
	mv vendor vendor-tmp
	mkdir -p vendor/src
	mv vendor-tmp/* vendor/src/
	rmdir vendor-tmp
