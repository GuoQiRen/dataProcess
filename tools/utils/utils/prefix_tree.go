package utils

import (
	"dataProcess/constants/plat"
	"errors"
	"strings"
)

var Root *PrefixTree

/**
children: 表示孩子节点有多少
text: 表示节点的上下文路径字符串
pass: 表示是否经过该节点
end: 表示字符串是否以该节点结尾
*/
type PrefixTree struct {
	children map[string]*PrefixTree
	text     string
	pass     int
	end      int
}

/**
初始化前缀树的根节点root
*/
func CreatePrefixTreeRoot(children map[string]*PrefixTree, text string, pass, end int) {
	Root = &PrefixTree{children: children, text: text, pass: pass, end: end}
}

/**
插入节点至前缀树
*/
func (r *PrefixTree) InsertNodeToPrefixTree(jinnPaths []string) (err error) {
	if r == nil {
		err = errors.New("please init prefixTree root node")
		return
	}

	node := r

	for _, jinnPath := range jinnPaths {
		node = r

		jinnArray := strings.Split(jinnPath, plat.LinuxSpiltRex)

		for _, text := range jinnArray {
			if len(text) == 0 {
				continue
			}
			if node.children == nil { // 表示不在该tree中，直接插入
				node.children = make(map[string]*PrefixTree)
			}
			if _, ok := node.children[text]; !ok {
				node.children[text] = &PrefixTree{}
			}
			node.children[text].pass++
			node = node.children[text]
		}
		node.end++
	}

	return
}

/**
找到prefix中pass最多且end=1的nodePath
*/
func (r *PrefixTree) FindMaxPassNodeWithEnd(nodePaths *[]string, initPath []string) (err error) {
	if r == nil {
		return
	}

	if r.end == 1 { // 当前节点是要找的节点
		path := plat.LinuxSpiltRex + strings.Join(initPath, plat.LinuxSpiltRex)
		*nodePaths = append(*nodePaths, path)
		return
	}

	for key, node := range r.children {
		initPath = append(initPath, key)
		err = node.FindMaxPassNodeWithEnd(nodePaths, initPath)
		if err != nil {
			return
		}
		initPath = initPath[:len(initPath)-1]
	}

	return
}
