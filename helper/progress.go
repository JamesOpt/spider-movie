package helper

import "fmt"

type Progress struct {
	percent int64 // 百分比
	cur int64 // 当前进度
	total int64 // 总进度
	rate string // 进度条
	graph string // 显示符号
	remark string
}

func NewProgress(start, total int64, remark string) *Progress {
	bar := Progress{
		cur:    start,
		total:  total,
		remark: remark,
		graph: "█",
	}
	
	bar.percent = bar.getPresent()

	for i := 0; i < int(bar.percent); i+=2 {
		bar.rate += bar.graph
	}

	return &bar
}

func (bar *Progress) getPresent() int64 {
	return int64(float32(bar.cur) / float32(bar.total) * 100)
}

func (bar *Progress) Play (cur int64)  {
	bar.cur = cur
	last := bar.percent
	bar.percent = bar.getPresent()
	if bar.percent != last && bar.percent%2 == 0 {
		bar.rate += bar.graph
	}

	fmt.Printf("\r[%-50s]%3d%%  %8d/%d %s", bar.rate, bar.percent, bar.cur, bar.total, bar.remark)
}