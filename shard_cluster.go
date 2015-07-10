package wizard

import (
	"github.com/evalphobia/wizard/errors"
)

// ShardCluster is struct for sharded database cluster
type ShardCluster struct {
	List     []*ShardSet // sharded database clusters
	slotsize int64
}

// Master is dummy method for interface
func (c ShardCluster) Master() *Node {
	return nil
}

// Masters returns all db masters from the sharded clusters
func (c ShardCluster) Masters() []*Node {
	var result []*Node
	for _, s := range c.List {
		if s.set == nil {
			continue
		}
		result = append(result, s.set.Master())
	}
	return result
}

// Slave is dummy method for interface
func (c ShardCluster) Slave() *Node {
	return nil
}

// SelectBySlot returns sharded cluster by hash slot id
func (c ShardCluster) SelectBySlot(i int64) *StandardCluster {
	mod := i % c.slotsize
	for _, shard := range c.List {
		if shard.InRange(mod) {
			return shard.set
		}
	}
	return nil
}

// RegisterShard adds cluster with hash slot range(min and max) 
func (c *ShardCluster) RegisterShard(min, max int64, s *StandardCluster) error {
	err := c.checkOverlapped(min, max)
	if err != nil {
		return err
	}

	ss := &ShardSet{
		min: min,
		max: max,
		set: s,
	}
	err = ss.checkSlotSize(c.slotsize)
	if err != nil {
		return err
	}

	c.List = append(c.List, ss)
	return nil
}

// checkOverlapped checks the hash slot range is not overlapped among the shards
func (c *ShardCluster) checkOverlapped(min, max int64) error {
	for _, ss := range c.List {
		switch {
		case ss.InRange(min):
			return errors.NewErrSlotMinOverlapped(min)
		case ss.InRange(max):
			return errors.NewErrSlotMaxOverlapped(max)
		}
	}
	return nil
}

// ShardSet is struct of sharded cluster
type ShardSet struct {
	min int64
	max int64
	set *StandardCluster
}

// InRange checks given number is in range of this shard
func (ss ShardSet) InRange(v int64) bool {
	return ss.min <= v && v <= ss.max
}

// checkSlotSize checks given number is not minus and within slotsize
func (ss ShardSet) checkSlotSize(slot int64) error {
	switch {
	case !ss.isMinAboveZero():
		return errors.NewErrSlotSizeMin(ss.min)
	case !ss.isMaxInSlotSize(slot):
		return errors.NewErrSlotSizeMax(ss.max, slot)
	}
	return nil
}

// isMinAboveZero checks given number is not minus
func (ss ShardSet) isMinAboveZero() bool {
	return ss.min >= 0
}

// isMaxInSlotSize checks given number is within slotsize
func (ss ShardSet) isMaxInSlotSize(slot int64) bool {
	return ss.max < slot
}
