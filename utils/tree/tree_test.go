package tree

import (
	"fmt"
	"testing"
)

func TestTrie_Search(t *testing.T) {
	trie := NewTrie()
	trie.Insert("word")
	trie.Insert("wore")
	fmt.Println(trie.Search("word"))
	fmt.Println(trie.Search("worf"))
	fmt.Println(trie.StartsWith("wor"))
}
