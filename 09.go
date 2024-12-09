package main

import (
	"bytes"
	"slices"
)

func init() {
	addSolutions(9, problem9)
}

func problem9(ctx *problemContext) {
	input := bytes.TrimSpace(ctx.readAll())
	ctx.reportLoad()

	d1 := disk1FromMap(bytes.Clone(input))
	d1.compact()
	ctx.reportPart1(d1.checksum())

	d2 := disk2FromMap(bytes.Clone(input))
	d2.compact()
	ctx.reportPart2(d2.checksum())
}

type disk1 []int

func disk1FromMap(b []byte) disk1 {
	var d disk1
	inFree := false
	var id int
	for _, c := range b {
		n := int(c - '0')
		val := id
		if inFree {
			val = -1
			id++
		}
		for i := 0; i < n; i++ {
			d = append(d, val)
		}
		inFree = !inFree
	}
	if !inFree {
		panic("bad")
	}
	return d
}

func (d disk1) compact() {
	i := 0
	j := len(d) - 1
	for i < j {
		if d[i] >= 0 {
			i++
			continue
		}
		if d[j] < 0 {
			j--
			continue
		}
		d[i], d[j] = d[j], d[i]
		i++
		j--
	}
}

func (d disk1) checksum() int64 {
	var sum int64
	for i, v := range d {
		if v >= 0 {
			sum += int64(i * v)
		}
	}
	return sum
}

type diskFile struct {
	id     int // not set for free blocks
	offset int
	size   int
}

type disk2 struct {
	files []*diskFile
	free  []*diskFile
}

func disk2FromMap(b []byte) *disk2 {
	var d disk2
	inFree := false
	var id int
	var offset int
	for _, c := range b {
		f := &diskFile{
			offset: offset,
			size:   int(c - '0'),
		}
		offset += f.size
		if inFree {
			if f.size > 0 {
				d.free = append(d.free, f)
			}
			id++
		} else {
			f.id = id
			d.files = append(d.files, f)
		}
		inFree = !inFree
	}
	if !inFree {
		panic("bad")
	}
	return &d
}

func (d *disk2) compact() {
outer:
	for _, f := range slices.Backward(d.files) {
		for i, free := range d.free {
			if free.offset >= f.offset {
				d.free = d.free[:i]
				break
			}
			if free.size >= f.size {
				f.offset = free.offset
				if free.size > f.size {
					free.offset += f.size
					free.size -= f.size
				} else {
					d.free = slices.Delete(d.free, i, i+1)
				}
				continue outer
			}
		}
	}
}

func (d *disk2) checksum() int64 {
	var sum int64
	for _, f := range d.files {
		for i := 0; i < f.size; i++ {
			sum += int64(f.offset+i) * int64(f.id)
		}
	}
	return sum
}
