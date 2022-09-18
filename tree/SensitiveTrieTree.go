package tree

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/pkg/errors"
	"go-utils/converter"
	"regexp"
	"strings"
)

// SensitiveTrieTree - Trie tree that used to store the sensitive words and find the sensitive words in text string.
type SensitiveTrieTree struct {
	replaceChar rune
	TrieTree
}

func NewSensitiveTrieTree() *SensitiveTrieTree {
	return &SensitiveTrieTree{
		replaceChar: '*',
		TrieTree: TrieTree{root: &TrieNode{
			childNodes: nil,
			Data:       "",
			End:        false,
		}},
	}
}

func (tree *SensitiveTrieTree) AddWord(sensitiveWord string) {
	sensitiveWord = filterChars(sensitiveWord)
	sensitiveRunes := []rune(sensitiveWord)
	node := tree.root
	for _, sensitiveRune := range sensitiveRunes {
		node = node.AddChild(sensitiveRune)
	}
	node.End = true
}

func (tree *SensitiveTrieTree) AddWords(sensitiveWords []string) {
	for _, sensitiveWord := range sensitiveWords {
		tree.AddWord(sensitiveWord)
	}
}

func (tree *SensitiveTrieTree) AddChineseWords(sensitiveWords []string) {
	for _, sensitiveWord := range sensitiveWords {
		tree.AddWord(sensitiveWord)

		if pinyin, err := converter.WordToPinyin(sensitiveWord); err == nil {
			tree.AddWord(pinyin)
		}
	}
}

// replaceRune replace the runes start at the 'begin' pos, end at the 'end' pos, attention: includes the end pos.
func (tree *SensitiveTrieTree) replaceRune(chars []rune, begin int, end int) {
	for i := begin; i <= end; i++ {
		chars[i] = tree.replaceChar
	}
}

func filterChars(text string) string {
	otherCharReg := regexp.MustCompile("[^\u4e00-\u9fa5a-zA-Z\\d]|\\s")
	text = otherCharReg.ReplaceAllString(text, "")
	text = strings.ToLower(text)
	return text
}

func (tree *SensitiveTrieTree) Match(text string) (sensitiveWords []string, replaceText string, err error) {
	if tree.root == nil {
		return nil, text, errors.New("Sensitive trie tree not built yet ...")
	}
	originText := text
	text = filterChars(text)
	// create a set that used store the sensitive words occur in text.
	sensitiveSet := mapset.NewSet()
	textChars := []rune(text)
	textCharsCopy := make([]rune, len(textChars))
	copy(textCharsCopy, textChars)
	for i, textLen := 0, len(textChars); i < textLen; i++ {
		trieNode := tree.root.FindChild(textChars[i], false)
		if trieNode == nil {
			continue
		}
		// if only one rune matches.
		// if all the sensitive words matches, then starts to replace the sensitive words.
		if trieNode.End {
			if sensitiveSet.Add(trieNode.Data) {
				sensitiveWords = append(sensitiveWords, trieNode.Data)
			}
			tree.replaceRune(textCharsCopy, i, i)
			continue
		}
		// more than 1 rune match.
		j := i + 1
		for ; j < textLen && trieNode != nil; j++ {
			trieNode = trieNode.FindChild(textChars[j], false)
			if trieNode == nil {
				break
			}
			// if all the sensitive words matches, then starts to replace the sensitive words.
			if trieNode.End {
				if sensitiveSet.Add(trieNode.Data) {
					sensitiveWords = append(sensitiveWords, trieNode.Data)
				}
				tree.replaceRune(textCharsCopy, i, j)
				i = j
				break
			}
		}
	}

	if len(sensitiveWords) > 0 {
		replaceText = string(textCharsCopy)
	} else {
		replaceText = originText
	}
	return sensitiveWords, replaceText, nil
}
