package linker

import (
	"fmt"
	"github.com/obgnail/LinkGameCheater/config"
	image2 "github.com/obgnail/LinkGameCheater/image"
	"github.com/obgnail/LinkGameCheater/utils"
	"image"
	"log"
	"strings"
)

var table *Table

type Table struct {
	RowLen  int
	LineLen int
	Table   [][]*Point
}

func newTable(linkGameTable [][]int) *Table {
	rowLen := len(linkGameTable)
	lineLen := len(linkGameTable[0])
	t := &Table{
		RowLen:  rowLen,
		LineLen: lineLen,
		Table:   make([][]*Point, rowLen),
	}

	for rowIdx := 0; rowIdx < t.RowLen; rowIdx++ {
		t.Table[rowIdx] = make([]*Point, lineLen)
		for lineIdx := 0; lineIdx < t.LineLen; lineIdx++ {
			typeCode := linkGameTable[rowIdx][lineIdx]
			point := NewPoint(rowIdx, lineIdx, typeCode)
			t.Table[rowIdx][lineIdx] = point
		}
	}
	return t
}

func (t *Table) String() string {
	var rows []string
	for rowIdx := 0; rowIdx < t.RowLen; rowIdx++ {
		var line []string
		for lineIdx := 0; lineIdx < t.LineLen; lineIdx++ {
			point := t.Table[rowIdx][lineIdx]
			s := fmt.Sprintf("%d", point.TypeCode)
			line = append(line, s)
		}
		rows = append(rows, strings.Join(line, "\t"))
	}
	return strings.Join(rows, "\n") + "\n"
}

func NewTableFromArr(tableArr [][]int) *Table {
	return newTable(tableArr)
}

func NewTableFromRandom(typeCodeCount, rowLen, lineLen int) *Table {
	total := lineLen * rowLen
	TableList, err := utils.GenRandomTableList(typeCodeCount, total)
	if err != nil {
		log.Fatal("[ERROR] Gen TableList", err)
	}
	TableArr, err := utils.GenTableArr(TableList, lineLen, rowLen)
	if err != nil {
		log.Fatal("[ERROR] Gen TableArr", err)
	}
	table := newTable(TableArr)
	return table
}

func NewTableFromImageByCount(
	imagePath string,
	fixRectangleMinPointX, fixRectangleMinPointY, fixRectangleMaxPointX, fixRectangleMaxPointY int,
	rowLen, lineLen int,
	emptyIndies []*image2.Idx,
) *Table {
	img, err := image2.NewImage(imagePath, fixRectangleMinPointX, fixRectangleMinPointY, fixRectangleMaxPointX, fixRectangleMaxPointY)
	if err != nil {
		log.Fatal(err)
	}
	subImages, err := img.GetSubImagesByCount(rowLen, lineLen)
	if err != nil {
		log.Fatal(err)
	}

	table, err := NewTableByImageArr(subImages, emptyIndies)
	if err != nil {
		log.Fatal(err)
	}
	return table
}

func NewTableFromImageByPixel(
	imagePath string,
	fixRectangleMinPointX, fixRectangleMinPointY, fixRectangleMaxPointX, fixRectangleMaxPointY int,
	subImgDW, subImgDH int,
	emptyIndies []*image2.Idx,
) *Table {
	img, err := image2.NewImage(imagePath, fixRectangleMinPointX, fixRectangleMinPointY, fixRectangleMaxPointX, fixRectangleMaxPointY)
	if err != nil {
		log.Fatal(err)
	}
	subImages, err := img.GetSubImagesByPixel(subImgDW, subImgDH)
	if err != nil {
		log.Fatal(err)
	}

	table, err := NewTableByImageArr(subImages, emptyIndies)
	if err != nil {
		log.Fatal(err)
	}
	return table
}

func NewTableByImageArr(imageArr [][]*image.NRGBA, emptyIndies []*image2.Idx) (*Table, error) {
	linkGameTable, err := image2.GenTableArrByImages(imageArr, emptyIndies)
	if err != nil {
		return nil, err
	}
	withEmpty := utils.AddOutEmptyPoint(linkGameTable)
	table := newTable(withEmpty)
	return table, nil
}

func (t *Table) GetPoint(rowIdx, lineIdx int) (*Point, error) {
	if 0 > rowIdx || rowIdx >= t.RowLen || 0 > lineIdx || lineIdx >= t.LineLen {
		return nil, fmt.Errorf("point(%d, %d) is out of boundary(%d, %d)", rowIdx, lineIdx, t.RowLen, t.LineLen)
	}
	return t.Table[rowIdx][lineIdx], nil
}

func (t *Table) SetEmpty(rowIdx, lineIdx int) error {
	p, err := t.GetPoint(rowIdx, lineIdx)
	if err != nil {
		return err
	}
	p.TypeCode = config.PointTypeCodeEmpty
	return nil
}

func InitTable(method string) {
	switch method {
	case "FromRandom":
		table = NewTableFromRandom(config.TypeCodeCount, config.LineLen, config.RowLen)
	case "FromArr":
		table = NewTableFromArr(config.Table)
	case "FromImageByCount":
		emptyIndies := image2.NewIndies(config.EmptySubImageIndies)
		table = NewTableFromImageByCount(
			config.ImagePath,
			config.FixImageRectangleMinPointX,
			config.FixImageRectangleMinPointY,
			config.FixImageRectangleMaxPointX,
			config.FixImageRectangleMaxPointY,
			config.ImageRowCount,
			config.ImageLineCount,
			emptyIndies,
		)
	case "FromImageByPixel":
		emptyIndies := image2.NewIndies(config.EmptySubImageIndies)
		table = NewTableFromImageByPixel(
			config.ImagePath,
			config.FixImageRectangleMinPointX,
			config.FixImageRectangleMinPointY,
			config.FixImageRectangleMaxPointX,
			config.FixImageRectangleMaxPointY,
			config.EachSubImageRowPixel,
			config.EachSubImageLinePixel,
			emptyIndies,
		)
	default:
		log.Fatal("ERROR Method")
	}
}

func GetTable() *Table { return table }
