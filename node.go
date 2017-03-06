package gomigration

type Method struct {
	Id       string
	Callback MigrationCallBack
	Next     *Node
}

type Node struct {
	Parent  *Node
	Methods []Method
}

func (n *Node) Add(id string, callback MigrationCallBack) *Node {
	newNode := &Node{Parent: n, Methods: []Method{}}
	n.Methods = append(n.Methods, Method{Id: id, Callback: callback, Next: newNode})
	return newNode
}
