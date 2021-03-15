package tree

type Trie struct {
	isWord   bool
	children map[rune]*Trie
}

func NewTrie() *Trie {
	return &Trie{
		isWord:   false,
		children: make(map[rune]*Trie),
	}
}

func (this *Trie) Insert(word string) {
	parent := this
	for _, ch := range word {
		if child, ok := parent.children[ch]; ok {
			parent = child
		} else {
			newChild := &Trie{
				children: make(map[rune]*Trie),
			}
			parent.children[ch] = newChild
			parent = newChild

		}
	}
	parent.isWord = true
}

func (this *Trie) Search(word string) bool {
	parent := this
	for _, ch := range word {
		if child, ok := parent.children[ch]; ok {
			parent = child
			continue
		}
		return false
	}
	return parent.isWord
}

func (this *Trie) StartsWith(prefix string) bool {
	parent := this
	for _, ch := range prefix {
		if child, ok := parent.children[ch]; ok {
			parent = child
			continue
		}
		return false
	}
	return true
}
