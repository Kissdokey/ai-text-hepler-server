package utils
import (
	"golang.org/x/net/html"
	"strings"
)

// extractText 从HTML中提取纯文本
func ExtractText(htmlStr string) string {
	var buf strings.Builder
	node, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return ""
	}
	visitNodes(node, &buf)
	return buf.String()
}

// visitNodes 递归访问HTML节点，并提取文本
func visitNodes(n *html.Node, buf *strings.Builder) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		visitNodes(c, buf)
	}
}