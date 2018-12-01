package layout_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vcraescu/go-xrandr"
	"github.com/vcraescu/rescreen/config"
	"github.com/vcraescu/rescreen/layout"
	"testing"
)

func TestLayoutNodeByID(t *testing.T) {
	nodes := layout.Nodes{
		{
			Monitor: xrandr.Monitor{
				ID: "Node1",
			},
		},
		{
			Monitor: xrandr.Monitor{
				ID: "Node2",
			},
		},
		{
			Monitor: xrandr.Monitor{
				ID: "Node3",
			},
		},
		{
			Monitor: xrandr.Monitor{
				ID: "Node4",
			},
		},
	}

	node := layout.FindLayoutNodeById("Node3", nodes)

	assert.NotNil(t, node)
	assert.NotNil(t, "Node3", node.ID())
}

func TestCalculateRescaledNodeResolution(t *testing.T) {
	node := layout.Node{
		Monitor: xrandr.Monitor{
			Modes: []xrandr.Mode{
				{
					Resolution: xrandr.Size{
						Width:  1000,
						Height: 1000,
					},
					RefreshRates: []xrandr.RefreshRate{
						{
							Current: true,
						},
					},
				},
			},
		},
	}

	size, err := layout.CalculateRescaledNodeResolution(node, 1.3)
	assert.NoError(t, err)
	assert.Equal(t, float32(1300), size.Width)
	assert.Equal(t, float32(1300), size.Height)
}

func TestCalculateNodeResolution(t *testing.T) {
	node := layout.Node{
		Monitor: xrandr.Monitor{
			ID: "Monitor1",
			Modes: []xrandr.Mode{
				{
					Resolution: xrandr.Size{
						Width:  1000,
						Height: 1000,
					},
					RefreshRates: []xrandr.RefreshRate{
						{
							Current: true,
						},
					},
				},
			},
		},
	}

	cfg := config.Config{
		Monitors: config.MonitorsConfig{
			"Monitor1": {
				Scale: 1.3,
			},
		},
	}

	size, err := layout.CalculateNodeResolution(node, cfg)
	assert.NoError(t, err)
	assert.Equal(t, float32(1300), size.Width)
	assert.Equal(t, float32(1300), size.Height)

	size, err = layout.CalculateNodeResolution(node, config.Config{})
	assert.NoError(t, err)
	assert.Equal(t, float32(1000), size.Width)
	assert.Equal(t, float32(1000), size.Height)
}

func TestCalculateNodePosition(t *testing.T) {
	node := layout.Node{}

	pos := layout.CalculateNodePosition(node)
	assert.Equal(t, 0, pos.X)
	assert.Equal(t, 0, pos.Y)

	node.Left = &layout.Node{
		Resolution: xrandr.Size{
			Width:  1000,
			Height: 1000,
		},
	}

	pos = layout.CalculateNodePosition(node)
	assert.Equal(t, 1000, pos.X)
	assert.Equal(t, 0, pos.Y)

	node.Left = nil
	node.Top = &layout.Node{
		Resolution: xrandr.Size{
			Width:  1000,
			Height: 1000,
		},
	}
	pos = layout.CalculateNodePosition(node)
	assert.Equal(t, 0, pos.X)
	assert.Equal(t, 1000, pos.Y)
}

func TestCalculateNodePosition4x4(t *testing.T) {
	node := layout.Node{}
	topLeft := &layout.Node{
		Resolution: xrandr.Size{
			Width:  300,
			Height: 400,
		},
	}
	bottomLeft := &layout.Node{
		Resolution: xrandr.Size{
			Width:  200,
			Height: 600,
		},
		Top: topLeft,
	}

	topRight := &layout.Node{
		Resolution: xrandr.Size{
			Width:  300,
			Height: 500,
		},
		Left:   topLeft,
		Bottom: &node,
	}

	node.Top = topRight
	node.Left = bottomLeft

	pos := layout.CalculateNodePosition(node)
	assert.Equal(t, 300, pos.X)
	assert.Equal(t, 500, pos.Y)
}

func TestCalculateNodePosition3x3(t *testing.T) {
	// *-----------------------------*
	// | 0       | 1       | 2       |
	// *-----------------------------*-*
	// | 300x400 | 400x500 | 200x600 |0|
	// | 500x600 | 200x550 | 300x200 |1|
	// | 400x200 | 400x700 | x       |2|
	// *-----------------------------*-*

	node := &layout.Node{}

	node00 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  300,
			Height: 400,
		},
	}
	node01 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  400,
			Height: 500,
		},
		Left: node00,
	}
	node00.Right = node01
	node02 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  200,
			Height: 600,
		},
		Left: node01,
	}
	node01.Right = node02

	node10 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  500,
			Height: 600,
		},
		Top: node00,
	}
	node00.Bottom = node10
	node11 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  200,
			Height: 550,
		},
		Top:  node01,
		Left: node10,
	}
	node10.Right = node11
	node01.Bottom = node11
	node12 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  300,
			Height: 200,
		},
		Top:    node02,
		Left:   node11,
		Bottom: node,
	}
	node11.Right = node12
	node02.Bottom = node12

	node20 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  400,
			Height: 200,
		},
		Top: node10,
	}
	node10.Bottom = node20
	node21 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  400,
			Height: 700,
		},
		Top:   node11,
		Left:  node20,
		Right: node,
	}
	node11.Bottom = node21
	node20.Right = node21

	node.Left = node21
	node.Top = node12

	pos := layout.CalculateNodePosition(*node)
	assert.Equal(t, 800, pos.X)
	assert.Equal(t, 1200, pos.Y)
}

func TestCalculateLayoutResolution3x3(t *testing.T) {
	// *-----------------------------*
	// | 0       | 1       | 2       |
	// *-----------------------------*-*
	// | 300x400 | 400x500 | 200x600 |0|
	// | 500x600 | 200x550 | 300x200 |1|
	// | 400x200 | 400x700 | 100x200 |2|
	// *-----------------------------*-*
	nodes := make(layout.Nodes, 9)

	node00 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  300,
			Height: 400,
		},
	}
	nodes[0] = node00
	node01 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  400,
			Height: 500,
		},
		Left: node00,
	}
	nodes[1] = node01
	node00.Right = node01
	node02 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  200,
			Height: 600,
		},
		Left: node01,
	}
	nodes[2] = node02
	node01.Right = node02

	node10 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  500,
			Height: 600,
		},
		Top: node00,
	}
	nodes[3] = node10
	node00.Bottom = node10
	node11 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  200,
			Height: 550,
		},
		Top:  node01,
		Left: node10,
	}
	nodes[4] = node11
	node10.Right = node11
	node01.Bottom = node11
	node12 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  300,
			Height: 200,
		},
		Top:  node02,
		Left: node11,
	}
	nodes[5] = node12
	node11.Right = node12
	node02.Bottom = node12

	node20 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  400,
			Height: 200,
		},
		Top: node10,
	}
	nodes[6] = node20
	node10.Bottom = node20
	node21 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  400,
			Height: 700,
		},
		Top:  node11,
		Left: node20,
	}
	nodes[7] = node21
	node11.Bottom = node21
	node20.Right = node21

	node22 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  100,
			Height: 200,
		},
		Top:  node12,
		Left: node21,
	}
	nodes[8] = node22
	node12.Bottom = node22
	node21.Right = node22

	s := layout.CalculateLayoutResolution(nodes)
	assert.Equal(t, float32(1200), s.Width)
	assert.Equal(t, float32(1900), s.Height)
}

func TestCalculateLayoutResolution3x3Incomplete(t *testing.T) {
	// *-----------------------------*
	// | 0       | 1       | 2       |
	// *-----------------------------*-*
	// | 300x400 | 400x500 | 200x600 |0|
	// | 500x600 | 200x550 |         |1|
	// | 400x200 |         |         |2|
	// *-----------------------------*-*
	nodes := make(layout.Nodes, 6)

	node00 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  300,
			Height: 400,
		},
	}
	nodes[0] = node00
	node01 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  400,
			Height: 500,
		},
		Left: node00,
	}
	nodes[1] = node01
	node00.Right = node01
	node02 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  200,
			Height: 600,
		},
		Left: node01,
	}
	nodes[2] = node02
	node01.Right = node02

	node10 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  500,
			Height: 600,
		},
		Top: node00,
	}
	nodes[3] = node10
	node00.Bottom = node10
	node11 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  200,
			Height: 550,
		},
		Top:  node01,
		Left: node10,
	}
	nodes[4] = node11
	node10.Right = node11
	node01.Bottom = node11

	node20 := &layout.Node{
		Resolution: xrandr.Size{
			Width:  400,
			Height: 200,
		},
		Top: node10,
	}
	nodes[5] = node20
	node10.Bottom = node20

	s := layout.CalculateLayoutResolution(nodes)
	assert.Equal(t, float32(1100), s.Width)
	assert.Equal(t, float32(1400), s.Height)
}

func TestCalculateDPI(t *testing.T) {
	nodes := make(layout.Nodes, 2)

	nodes[0] = &layout.Node{
		Monitor: xrandr.Monitor{
			Modes: []xrandr.Mode{
				{
					Resolution: xrandr.Size{
						Width:  3840,
						Height: 2160,
					},
					RefreshRates: []xrandr.RefreshRate{
						{
							Current: true,
						},
					},
				},
			},
			Size: xrandr.Size{
				Width:  487,
				Height: 247,
			},
		},
	}
	nodes[1] = &layout.Node{
		Monitor: xrandr.Monitor{
			Modes: []xrandr.Mode{
				{
					Resolution: xrandr.Size{
						Width:  3840,
						Height: 2160,
					},
					RefreshRates: []xrandr.RefreshRate{
						{
							Current: true,
						},
					},
				},
			},
			Size: xrandr.Size{
				Width:  597.7,
				Height: 336.2,
			},
		},
	}

	dpi, err := layout.CalculateDPI(nodes)
	assert.Nil(t, err)
	assert.Equal(t, 163, int(dpi))
}
