package mset_test

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/ahwhy/myGolang/map/mset"
)

func TestMset(t *testing.T) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	rand.Seed(time.Now().Unix())
	a := mset.NewSet()
	b := mset.NewSet()

	a.Add("aa", 123, false, 456.789, 1000)
	b.Add("aa", 123, false, 456.789, 1000)

	a.Delete(1000)
	b.Modify("aa", 123, false, 456.789)

	log.Printf("%T %+v ", a, a.Mset)
	log.Printf("%T %#v ", b, b.Mset)

	b.Clear()
	b.Add("aa", 123, false, 456.789, 1000)

	log.Println(a.Contains("aa"), a.Size(), a.Equel(b), a.IsSubset(b))

	c := mset.NewSet()
	for i := 0; i < 10000; i++ {
		v := fmt.Sprintf("KEY_%d", i)
		go func() {
			c.Add(v)
			for i := 0; i < 1000; i++ {
				c.Add(rand.Intn(1000))
			}
		}()

		go func() {
			c.Modify(v)
			for i := 0; i < 1000; i++ {
				c.Modify(rand.Intn(1000))
			}
		}()
	}

	for i := 0; i < 8000; i++ {
		v := fmt.Sprintf("KEY_%d", rand.Intn(1000))
		go func() {
			c.Delete(v)
			for i := 0; i < 1000; i++ {
				c.Delete(rand.Intn(1000))
			}
		}()

		go func() {
			c.Modify(v)
			for i := 0; i < 1000; i++ {
				c.Modify(rand.Intn(1000))
			}
		}()
	}

	log.Println(c.Contains(rand.Intn(1000)), c.Size(), c.Equel(b), a.IsSubset(c))
	log.Println(c.Mset)
}
