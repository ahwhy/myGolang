package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// 0Start()
	// BinToChar()
	fmt.Println(HumanBytesLoaded(1024 * 1024 * 1024))
}

/*
	请将这段二进制翻译成中文(unicode编码)
	1000 1000 0110 0011
	0110 0111 0000 1101
	0101 0101 1001 1100
	0110 1011 0010 0010
	0111 1010 0111 1111  】
	0100 1110 0010 1101
	0101 0110 1111 1101
	0111 1110 1010 0010
*/
func BinToChar() {
	var binary []int = []int{
		0b1000100001100011,
		0b0110011100001101,
		0b0101010110011100,
		0b0110101100100010,
		0b0111101001111111,
		0b0100111000101101,
		0b0101011011111101,
		0b0111111010100010,
	}
	for _, b := range binary {
		fmt.Printf("二进制: %b, unicode编码: %U 代表字符: %c\n", b, b, b)
	}
	fmt.Println()
}

const (
	bu = 1 << 10
	kb = 1 << 20
	mb = 1 << 30
	gb = 1 << 40
	tb = 1 << 50
	eb = 1 << 60
)

// HumanBytesLoaded 单位转换
func HumanBytesLoaded(bytesLength int64) string {
	if bytesLength < bu {
		return fmt.Sprintf("%dB", bytesLength)
	} else if bytesLength < kb {
		return fmt.Sprintf("%.2fKB", float64(bytesLength)/float64(bu))
	} else if bytesLength < mb {
		return fmt.Sprintf("%.2fMB", float64(bytesLength)/float64(kb))
	} else if bytesLength < gb {
		return fmt.Sprintf("%.2fGB", float64(bytesLength)/float64(mb))
	} else if bytesLength < tb {
		return fmt.Sprintf("%.2fTB", float64(bytesLength)/float64(gb))
	} else {
		return fmt.Sprintf("%.2fEB", float64(bytesLength)/float64(tb))
	}
}

func Start() {
	rand.Seed(time.Now().Unix())
	anser := -1
	for {
		if anser != -1 {
			continue
		}

		m, n := rand.Intn(10), rand.Intn(10)
		fmt.Printf("请回答如下问题: %d + %d = ", m, n)
		fmt.Scanf("%d\n", &anser)

		if anser != m+n {
			fmt.Println("洗洗睡了吧，你已经不行了!")
			return
		}

		fmt.Println("恭喜你回答正确,即将进入抽奖环节....")
		if duang() {
			return
		}
		fmt.Println("就差亿点点，再试试吧！")

		anser = -1
	}
}

var (
	rollCount  = 100
	totalCount = 0
	lucky      = rand.Intn(rollCount)
)

func duang() bool {
	count := rand.Intn(10) + 1
	fmt.Printf("你当前获取[%d]次抽奖机会\n", count)
	
	for i := 1; i <= count; i++ {
		ticket := rand.Intn(rollCount)
		totalCount++
		if ticket == lucky {
			fmt.Printf("当前是你第%d次抽奖: 抽奖结果 中奖, 幸运数: %d\n, 你当前号码: %d", totalCount, lucky, ticket)
			return true
		}

		fmt.Printf("当前是你第%d次抽奖: 抽奖结果 未中奖, 幸运数: %d, 你当前号码: %d\n", totalCount, lucky, ticket)
	}

	return false
}
