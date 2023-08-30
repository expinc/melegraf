package conveyor

import (
	"testing"

	"github.com/expinc/melegraf/metric"
)

func TestConveyor(t *testing.T) {
	conveyor := NewConveyor("test", 3)

	// Test Put
	err := conveyor.Put(metric.Metric{Name: "metric1"})
	if nil != err {
		t.Fatal(err)
	}
	err = conveyor.Put(metric.Metric{Name: "metric2"})
	if nil != err {
		t.Fatal(err)
	}
	err = conveyor.Put(metric.Metric{Name: "metric3"})
	if nil != err {
		t.Fatal(err)
	}
	err = conveyor.Put(metric.Metric{Name: "metric4"})
	if nil == err {
		t.Fatal("metric4 should fail to put in the full conveyor")
	}

	// Test getting metrics from channel
	mt := <-conveyor.GetChannel()
	if mt.Name != "metric1" {
		t.Fatal("metric1 should be the first metric")
	}
	mt = <-conveyor.GetChannel()
	if mt.Name != "metric2" {
		t.Fatal("metric2 should be the second metric")
	}
	mt = <-conveyor.GetChannel()
	if mt.Name != "metric3" {
		t.Fatal("metric3 should be the third metric")
	}

	conveyor.Dispose()
	// Test Dispose is idempotent
	conveyor.Dispose()

	// Test Put after Dispose
	err = conveyor.Put(metric.Metric{Name: "metric5"})
	if nil == err {
		t.Fatal("metric5 should fail to put in the disposed conveyor")
	}

	// Test GetChannel after Dispose
	_, ok := <-conveyor.GetChannel()
	if ok == true {
		t.Fatal("GetChannel should return a closed channel")
	}
}
