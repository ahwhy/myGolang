package trie_test

import (
	"fmt"
	"testing"

	"github.com/ahwhy/myGolang/data_structure/trie"
)

func TestTrie(t *testing.T) {
	tree := new(trie.TrieTree)
	tree.AddTerm("分散")
	tree.AddTerm("分散精力")
	tree.AddTerm("分散投资")
	tree.AddTerm("分布式")
	tree.AddTerm("工程")
	tree.AddTerm("工程师")

	terms := tree.Retrieve("分散")
	fmt.Println(terms)
	terms = tree.Retrieve("人工")
	fmt.Println(terms)
	terms = tree.Retrieve("工程")
	fmt.Println(terms)
}
