package pargraph

import "fmt"

type SourceRow struct {
	ID            int64
	HasParent     bool
	ParentID      int64
	SomeDataField string
}

type TreeNode struct {
	ID            int64
	SomeDataField string
	Children      []*TreeNode
}

type nodeMapItem struct {
	index   int
	visited bool
	node    *TreeNode
}

type couplingType uint8

const (
	None   couplingType = 0
	Child  couplingType = 1
	Parent couplingType = 2
)

type couplingMatrix map[int64]map[int64]couplingType

func MakeReachableTree(originID int64, rows []*SourceRow) (tree *TreeNode, err error) {
	nm, err := prepareNodeMap(rows)
	if err != nil {
		// добавить текстовку всплытия ошибки
		return
	}

	item := nm[originID]
	if item == nil {
		// среди исходного набора данных нет элемента с id = originID
		// возвращаем nil, поскольку дерево не из чего строить
		return
	}

	item.node = NewNode(rows[item.index])
	item.visited = true

	cm := prepareCouplingMatrix(rows)

	// набор элементов, которые надо обойти на текущей итерации
	edge := []int64{
		originID,
	}

	for len(edge) > 0 {
		// очередь обхода на следующей итерации
		nextEdge := []int64{}

		for _, id := range edge {
			current := nm[id].node

			// пока дерево строится попутно ищем корневой элемент
			if !rows[nm[id].index].HasParent {
				if tree != nil {
					// мы уже находили корневой элемент ранее, то есть в графе
					// их оказалось несколько, следовательно это не дерево;
					// возвращаем ошибку
					tree = nil
					err = fmt.Errorf("найден второй корневой элемент при построении дерева")
					return
				}
				tree = current
			}

			// смотрим соседей текущего элемента
			neighbors := cm[id]
			// nil map is safe for reading
			for neighborID, ct := range neighbors {
				neighbor := nm[neighborID]
				if !neighbor.visited {
					// создаем новый элемент дерева
					node := NewNode(rows[neighbor.index])
					neighbor.node = node
					neighbor.visited = true

					// подсоединяем новый элемент к старым элементам дерева
					if ct == Child {
						current.Children = append(current.Children, node)
					} else if ct == Parent {
						node.Children = append(node.Children, current)
					}

					// добавляем непосещенных соседей в очередь на обход
					nextEdge = append(nextEdge, neighborID)
				}
			}

		}
		edge = nextEdge
	}

	// дополнительные проверки правильности получившегося дерева
	if tree == nil {
		// в графе нет корневого элемента, значит это не дерево;
		// возвращаем ошибку
		err = fmt.Errorf("не обнаружен корневой элемент")
		return
	}
	if HasLoop(tree) {
		tree = nil
		err = fmt.Errorf("в построенном дереве обнаружен цикл")
		return
	}

	return
}

func NewNode(row *SourceRow) (node *TreeNode) {
	return &TreeNode{
		ID:            row.ID,
		SomeDataField: row.SomeDataField,
	}
}

func HasLoop(tree *TreeNode) bool {
	visitedMap := make(map[int64]bool)
	return hasLoop(tree, visitedMap)
}

func hasLoop(tree *TreeNode, visitedMap map[int64]bool) bool {
	if visitedMap[tree.ID] {
		return true
	}
	visitedMap[tree.ID] = true

	// nil slice is safe for range iteration
	for _, child := range tree.Children {
		if hasLoop(child, visitedMap) {
			return true
		}
	}
	return false
}

func prepareNodeMap(rows []*SourceRow) (nodeMap map[int64]*nodeMapItem, err error) {
	nodeMap = make(map[int64]*nodeMapItem)
	for i, row := range rows {
		nodeMap[row.ID] = &nodeMapItem{
			index: i,
		}
	}

	// проверка наличия родителей среди исходных элементов
	for _, row := range rows {
		if row.HasParent {
			_, ok := nodeMap[row.ParentID]
			if !ok {
				// родитель элемента отстутствует в списке исходных элементов;
				// возвращаем ошибку
				nodeMap = nil
				err = fmt.Errorf("у элемента %d не найден родитель с указанным id = %d", row.ID, row.ParentID)
				return
			}
		}
	}
	return
}

func prepareCouplingMatrix(rows []*SourceRow) (matrix couplingMatrix) {
	matrix = make(map[int64]map[int64]couplingType)
	for _, row := range rows {
		if row.HasParent {
			if matrix[row.ID] == nil {
				matrix[row.ID] = make(map[int64]couplingType)
			}
			matrix[row.ID][row.ParentID] = Parent

			if matrix[row.ParentID] == nil {
				matrix[row.ParentID] = make(map[int64]couplingType)
			}
			matrix[row.ParentID][row.ID] = Child
		}
	}
	return
}
