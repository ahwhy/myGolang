package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

var readWg = sync.WaitGroup{}
var dealWg = sync.WaitGroup{}

var textChan = make(chan string, 10000)
var numChan = make(chan int, 10000)
var writeFinishChan = make(chan struct{})

func readFile(infile string) {
	defer readWg.Done()
	fin, err := os.Open(infile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fin.Close()
	reader := bufio.NewReader(fin)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err)
			}
		}
		// n := calculate(string(line))
		// numChan <- n
		textChan <- string(line)
	}
}

func dealLine() {
	defer dealWg.Done()
	for {
		line, ok := <-textChan
		if !ok {
			break
		} else {
			n := calculate(string(line))
			numChan <- n
		}
	}
}

func writeFile(outfile string) {
	fout, err := os.OpenFile(outfile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fout.Close()
	writer := bufio.NewWriter(fout)
	for {
		n, ok := <-numChan
		if !ok {
			break
		} else {
			writer.WriteString(strconv.Itoa(n))
			writer.WriteString("\n")
		}
	}
	writer.Flush()
	writeFinishChan <- struct{}{}
}

func calculate(line string) int {
	sum := 0
	for _, c := range line {
		sum += int(c)
	}
	return sum
}

func main12() {
	//第1阶段，io密集型，并行执行提高速度
	readWg.Add(2)
	go func() {
		readFile("data/rsa_private_key.pem")
	}()
	go func() {
		readFile("data/rsa_public_key.pem")
	}()

	//第2阶段，cpu密集型，多分配几个内核线程
	dealWg.Add(4)
	for i := 0; i < 4; i++ {
		go dealLine()
	}

	//第3阶段，汇总，写一个文件
	go writeFile("data/digit.txt")

	//第1阶段结束后，关闭管道textChan
	readWg.Wait()
	close(textChan)

	//第2阶段结束后，关闭管道numChan
	dealWg.Wait()
	close(numChan)

	//第3阶段结束后，往writeFinishChan里send一下
	<-writeFinishChan
}
