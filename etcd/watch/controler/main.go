package main

import "github.com/ahwhy/myGolang/etcd/watch"

func main() {
	watch.WatchConfig("/registry/configs")
}
