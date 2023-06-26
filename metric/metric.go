package metric

import "time"

type Tag struct {
	Key   string
	Value string
}

type Field struct {
	Key   string
	Value interface{}
}

type Metric struct {
	Name   string
	Tags   []Tag
	Fields []Field
	Time   time.Time
}

func (mt *Metric) Copy() Metric {
	tags := make([]Tag, len(mt.Tags))
	copy(tags, mt.Tags)

	fields := make([]Field, len(mt.Fields))
	copy(fields, mt.Fields)

	return Metric{
		Name:   mt.Name,
		Tags:   tags,
		Fields: fields,
		Time:   mt.Time,
	}
}
