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
