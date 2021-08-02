package main

import "fmt"

type TrieNode struct {
	Word     rune               //当前节点存储的字符。byte只能表示英文字符，rune可以表示任意字符
	Children map[rune]*TrieNode //孩子节点，用一个map存储
	Term     string
}

type TrieTree struct {
	root *TrieNode
}

//add 把words[beginIndex:]插入到Trie树中
func (node *TrieNode) add(words []rune, term string, beginIndex int) {
	if beginIndex >= len(words) { //words已经遍历完了
		node.Term = term
		return
	}

	if node.Children == nil {
		node.Children = make(map[rune]*TrieNode)
	}

	word := words[beginIndex] //把这个word放到node的子节点中
	if child, exists := node.Children[word]; !exists {
		newNode := &TrieNode{Word: word}
		node.Children[word] = newNode
		newNode.add(words, term, beginIndex+1) //递归
	} else {
		child.add(words, term, beginIndex+1) //递归
	}
}

//walk words[0]就是当前节点上存储的字符，按照words的指引顺着树往下走，最终返回words最后一个字符对应的节点
func (node *TrieNode) walk(words []rune, beginIndex int) *TrieNode {
	if beginIndex == len(words)-1 {
		return node
	}
	beginIndex += 1
	word := words[beginIndex]
	if child, exists := node.Children[word]; exists {
		return child.walk(words, beginIndex)
	} else {
		return nil
	}
}

//traverseTerms 遍历一个node下面所有的term。注意要传数组的指针，才能真正修改这个数组
func (node *TrieNode) traverseTerms(terms *[]string) {
	if len(node.Term) > 0 {
		*terms = append(*terms, node.Term)
	}
	for _, child := range node.Children {
		child.traverseTerms(terms)
	}
}

func (tree *TrieTree) AddTerm(term string) {
	if len(term) <= 1 {
		return
	}
	words := []rune(term)

	if tree.root == nil {
		tree.root = new(TrieNode)
	}

	tree.root.add(words, term, 0)
}

func (tree *TrieTree) Retrieve(prefix string) []string {
	if tree.root == nil || len(tree.root.Children) == 0 {
		return nil
	}
	words := []rune(prefix)
	firstWord := words[0]
	if child, exists := tree.root.Children[firstWord]; exists {
		end := child.walk(words, 0)
		if end == nil {
			return nil
		} else {
			terms := make([]string, 0, 100)
			end.traverseTerms(&terms)
			return terms
		}
	} else {
		return nil
	}
}

func main7() {
	tree := new(TrieTree)
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
}

//go run data_structure/trie.go
