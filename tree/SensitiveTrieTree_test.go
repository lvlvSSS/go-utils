package tree

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSensitiveTrieTree_filterChars(t *testing.T) {
	str := "这是一个 测\t 12试， \n 哈哈,"
	targetStr := "这是一个测12试哈哈"
	chars := filterChars(str)
	t.Log(len([]rune(chars)))
	t.Log(chars)
	t.Logf("%v", []byte(chars))
	t.Logf("%v", []byte(targetStr))

	if chars != targetStr {
		t.Fail()
	}
}

func trieDemo(sensitiveWords []string, matchContents []string) {

	// 汉字转拼音
	trie := NewSensitiveTrieTree()
	trie.AddChineseWords(sensitiveWords)

	//trie.AddWords(pinyinContents)
	//for _, content := range contents {
	//	trie.AddWord(content)
	//}

	for _, srcText := range matchContents {
		matchSensitiveWords, replaceText, _ := trie.Match(srcText)
		fmt.Println("srcText        -> ", srcText)
		fmt.Println("replaceText    -> ", replaceText)
		fmt.Println("sensitiveWords -> ", matchSensitiveWords)
		fmt.Println()
	}

	// 动态添加
	trie.AddWord("牛大大")
	content := "今天，牛大大挑战灰大大"
	matchSensitiveWords, replaceText, _ := trie.Match(content)
	fmt.Println("srcText        -> ", content)
	fmt.Println("replaceText    -> ", replaceText)
	fmt.Println("sensitiveWords -> ", matchSensitiveWords)
}

func TestNewSensitiveTrieTree(t *testing.T) {

	sensitiveWords := []string{
		"傻逼",
		"傻叉",
		"垃圾",
		"妈的",
		"sb",
	}

	matchContents := []string{
		"你是一个大傻&逼，大傻 叉",
		"你是傻☺叉",
		"shabi东西",
		"他made东西",
		"什么垃圾打野，傻逼一样，叫你来开龙不来，SB",
		"正常的内容☺",
	}

	fmt.Println("\n--------- 前缀树匹配敏感词 ---------")
	trieDemo(sensitiveWords, matchContents)
}

func TestSensitiveTrieTree_Match(t *testing.T) {
	type fields struct {
		replaceChar rune
		TrieTree    TrieTree
	}
	tests := []struct {
		name               string
		fields             fields
		text               string
		wantSensitiveWords []string
		wantReplaceText    string
		wantErr            bool
	}{
		{
			name: "test1",
			fields: fields{
				replaceChar: '*',
				TrieTree: TrieTree{root: &TrieNode{
					childNodes: nil,
					Data:       "",
					End:        false,
				}},
			},
			text: "你是一个大傻&逼，大傻 叉",
			wantSensitiveWords: []string{
				"傻逼",
				"傻叉",
			},
			wantReplaceText: "你是一个大**大**",
			wantErr:         false,
		},
		{
			name: "test2",
			fields: fields{
				replaceChar: '*',
				TrieTree: TrieTree{root: &TrieNode{
					childNodes: nil,
					Data:       "",
					End:        false,
				}},
			},
			text: "你是傻☺叉",
			wantSensitiveWords: []string{
				"傻叉",
			},
			wantReplaceText: "你是**",
			wantErr:         false,
		},
		{
			name: "test3",
			fields: fields{
				replaceChar: '*',
				TrieTree: TrieTree{root: &TrieNode{
					childNodes: nil,
					Data:       "",
					End:        false,
				}},
			},
			text: "shabi东西",
			wantSensitiveWords: []string{
				"傻逼",
				"傻叉",
				"垃圾",
				"妈的",
				"sb",
			},
			wantReplaceText: "*****东西",
			wantErr:         false,
		},
		{
			name: "test4",
			fields: fields{
				replaceChar: '*',
				TrieTree: TrieTree{root: &TrieNode{
					childNodes: nil,
					Data:       "",
					End:        false,
				}},
			},
			text: "他made东西",
			wantSensitiveWords: []string{
				"傻逼",
				"傻叉",
				"垃圾",
				"妈的",
				"sb",
			},
			wantReplaceText: "他****东西",
			wantErr:         false,
		},
		{
			name: "test5",
			fields: fields{
				replaceChar: '*',
				TrieTree: TrieTree{root: &TrieNode{
					childNodes: nil,
					Data:       "",
					End:        false,
				}},
			},
			text: "什么垃圾打野，傻逼一样，叫你来开龙不来，SB",
			wantSensitiveWords: []string{
				"傻逼",
				"傻叉",
				"垃圾",
				"妈的",
				"sb",
			},
			wantReplaceText: "什么**打野**一样叫你来开龙不来**",
			wantErr:         false,
		},
		{
			name: "test6",
			fields: fields{
				replaceChar: '*',
				TrieTree: TrieTree{root: &TrieNode{
					childNodes: nil,
					Data:       "",
					End:        false,
				}},
			},
			text: "正常的内容☺",
			wantSensitiveWords: []string{
				"傻逼",
				"傻叉",
				"垃圾",
				"妈的",
				"sb",
			},
			wantReplaceText: "正常的内容☺",
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := &SensitiveTrieTree{
				replaceChar: tt.fields.replaceChar,
				TrieTree:    tt.fields.TrieTree,
			}
			tree.AddChineseWords(tt.wantSensitiveWords)
			gotSensitiveWords, gotReplaceText, err := tree.Match(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSensitiveWords, tt.wantSensitiveWords) {
				t.Errorf("Match() gotSensitiveWords = %v, want %v", gotSensitiveWords, tt.wantSensitiveWords)
			}
			if gotReplaceText != tt.wantReplaceText {
				t.Errorf("Match() gotReplaceText = %v, want %v", gotReplaceText, tt.wantReplaceText)
			}
		})
	}
}
