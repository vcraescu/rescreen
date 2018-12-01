package layout

import (
	"fmt"
	"github.com/vcraescu/go-xrandr"
	"github.com/vcraescu/rescreen/config"
)

// Layout stores all layout nodes
type Layout struct {
	Nodes      Nodes
	Resolution xrandr.Size
	DPI        int
}

// Node individual output from layout
type Node struct {
	Monitor    xrandr.Monitor
	Resolution xrandr.Size
	Position   xrandr.Position
	Primary    bool
	Scale      float32
	Left       *Node
	Right      *Node
	Top        *Node
	Bottom     *Node
}

// Nodes list
type Nodes []*Node

// New creates a new layout based on the unserialized config file and current screen configuration
func New(cfg config.Config, screens xrandr.Screens) (*Layout, error) {
	l := &Layout{}
	l.Nodes = createLayoutNodes(cfg, screens)
	l.Resolution = calculateLayoutResolution(l.Nodes)
	dpi, err := calculateDPI(l.Nodes)
	if err != nil {
		return nil, err
	}

	l.DPI = dpi

	return l, nil
}

// ID returns the node ID which is actually the monitor ID
func (n Node) ID() string {
	return n.Monitor.ID
}

func findLayoutNodeByID(id string, nodes Nodes) *Node {
	for _, node := range nodes {
		if node.Monitor.ID == id {
			return node
		}
	}

	return nil
}

func calculateRescaledNodeResolution(node Node, scale float32) (xrandr.Size, error) {
	res := xrandr.Size{}
	cm, ok := node.Monitor.CurrentMode()
	if !ok {
		return res, fmt.Errorf(`monitor "%s" current mode not found`, node.Monitor.ID)
	}

	return cm.Resolution.Rescale(scale), nil
}

func calculateNodeResolution(node Node, cfg config.Config) (xrandr.Size, error) {
	var scale float32 = 1
	mCfg, ok := cfg.Monitors[node.Monitor.ID]
	if ok {
		scale = mCfg.Scaling()
	}

	return calculateRescaledNodeResolution(node, scale)
}

func calculateNodePosition(node Node) xrandr.Position {
	pos := xrandr.Position{}

	nl := node.Left
	for nl != nil {
		max := nl.Resolution.Width
		nt := node.Top
		for nt != nil {
			if nt.Resolution.Width > max {
				max = nt.Resolution.Width
			}
			nt = nt.Top
		}
		pos.X += int(max)
		nl = nl.Left
	}

	nt := node.Top
	for nt != nil {
		max := nt.Resolution.Height
		nl = nt.Left
		for nl != nil {
			if nl.Resolution.Height > max {
				max = nl.Resolution.Height
			}
			nl = nl.Left
		}
		pos.Y += int(max)
		nt = nt.Top
	}

	return pos
}

func createLayoutNodes(cfg config.Config, screens xrandr.Screens) Nodes {
	nodes := make(Nodes, 0)
	matrix := cfg.Layout.Matrix()
	for i, row := range matrix {
		for j, id := range row {
			if id == "" {
				continue
			}

			monitor, found := screens.MonitorByID(id)
			if !found {
				continue
			}

			node := &Node{Monitor: monitor}
			nodes = append(nodes, node)

			if j > 0 {
				leftNode := findLayoutNodeByID(matrix[i][j-1], nodes)
				if leftNode != nil {
					node.Left = leftNode
					leftNode.Right = node
				}
			}

			if i > 0 {
				topNode := findLayoutNodeByID(matrix[i-1][j], nodes)
				if topNode != nil {
					node.Top = topNode
					topNode.Bottom = node
				}
			}

			node.Primary = cfg.Monitors[id].Primary
			node.Scale = 1
			if cfg.Monitors[id].Scale > 0 {
				node.Scale = cfg.Monitors[id].Scale
			}

			res, err := calculateNodeResolution(*node, cfg)
			if err != nil {
				panic(err)
			}
			node.Resolution = res
		}
	}

	for _, node := range nodes {
		node.Position = calculateNodePosition(*node)
	}

	return nodes
}

func calculateLayoutResolution(nodes Nodes) xrandr.Size {
	res := xrandr.Size{}
	if len(nodes) == 0 {
		return res
	}

	node00 := nodes[0]
	for node00.Left != nil {
		node00 = node00.Left
	}

	for node00.Top != nil {
		node00 = node00.Top
	}

	node := node00
	for node != nil {
		max, _ := calculateMaxWidthHeightPerColumn(node)
		res.Width += max
		node = node.Right
	}

	node = node00
	for node != nil {
		_, max := calculateMaxWidthHeightPerRow(node)
		res.Height += max
		node = node.Bottom
	}

	return res
}

func calculateDPI(nodes Nodes) (int, error) {
	var dpi float32 = 999999
	for _, node := range nodes {
		mDPI, err := node.Monitor.DPI()
		if err != nil {
			return 0, err
		}

		if dpi > mDPI {
			dpi = mDPI
		}
	}

	return int(float64(dpi)), nil
}

func calculateMaxWidthHeightPerColumn(node *Node) (float32, float32) {
	for node.Top != nil {
		node = node.Top
	}

	var maxW, maxH float32
	for node != nil {
		if node.Resolution.Width > maxW {
			maxW = node.Resolution.Width
		}

		if node.Resolution.Height > maxH {
			maxH = node.Resolution.Height
		}

		node = node.Bottom
	}

	return maxW, maxH
}

func calculateMaxWidthHeightPerRow(node *Node) (float32, float32) {
	for node.Left != nil {
		node = node.Left
	}

	var maxW, maxH float32
	for node != nil {
		if node.Resolution.Width > maxW {
			maxW = node.Resolution.Width
		}

		if node.Resolution.Height > maxH {
			maxH = node.Resolution.Height
		}

		node = node.Right
	}

	return maxW, maxH
}
