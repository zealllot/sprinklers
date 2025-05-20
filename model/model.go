package model

// Point 表示一个二维点
type Point struct {
	X, Y float64
}

// Room 表示一个房间
type Room struct {
	Walls                   []Point // 房间的墙体点，按顺时针或逆时针顺序排列
	SprinklerCoverageRadius float64 // 喷头覆盖半径，用于计算最大间距（最大间距 = 2 * 覆盖半径）
	MinSprinklerDistance    float64 // 喷头之间的最小距离
	MinWallDistance         float64 // 喷头到墙的最小距离
}

// Sprinkler 表示一个喷头
type Sprinkler struct {
	Position Point   // 喷头位置
	Coverage float64 // 覆盖半径
}

// BoundingBox 表示一个边界框
type BoundingBox struct {
	MinX, MinY, MaxX, MaxY float64
}

// Grid 表示网格
type Grid struct {
	Cells    [][]bool    // 网格单元，true 表示该单元被占用
	CellSize float64     // 网格单元大小
	BBox     BoundingBox // 网格的边界框
}
