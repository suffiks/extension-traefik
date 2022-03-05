package controllers

import (
	"fmt"
	"strings"
)

type node map[string]node

func (n node) add(parts []string) {
	if len(parts) == 0 {
		return
	}

	if _, ok := n[parts[0]]; !ok {
		n[parts[0]] = make(node)
	}

	n[parts[0]].add(parts[1:])
}

func (n node) contains(parts []string) bool {
	if len(parts) == 0 {
		return true
	}

	if _, ok := n[parts[0]]; !ok {
		if len(parts) == 1 && n["*"] != nil {
			return true
		}
		return false
	}

	return n[parts[0]].contains(parts[1:])
}

func (n node) print(indenter string, indents int) {
	for k, v := range n {
		ind := strings.Repeat(indenter, indents)
		fmt.Println(ind + k)
		v.print(indenter, indents+1)
	}
}

type tree struct {
	node node
}

func (t *tree) Add(host string) {
	parts := strings.Split(host, ".")
	reverse(parts)
	t.node.add(parts)
}

func (t *tree) Contains(host Host) bool {
	if t.node == nil {
		return false
	}

	parts := strings.Split(string(host), ".")
	reverse(parts)
	return t.node.contains(parts)
}

func (t *tree) print() {
	t.node.print("  ", 0)
}

func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
