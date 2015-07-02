package wizard

type ShardCluster struct {
	table    interface{}
	List     []*ShardSet
	slotsize int64
}

func (c ShardCluster) Table() interface{} {
	return c.table
}

func (c ShardCluster) Master() *Node {
	return c.SelectBySlot(0).Master()
}

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

func (c ShardCluster) Slave() *Node {
	return c.SelectBySlot(0).Slave()
}

func (c ShardCluster) SelectBySlot(i int64) *StandardCluster {
	mod := i % c.slotsize
	for _, shard := range c.List {
		if shard.InRange(mod) {
			return shard.set
		}
	}
	return nil
}

func (c *ShardCluster) RegisterShard(min, max int64, s *StandardCluster) {
	ss := &ShardSet{
		min: min,
		max: max,
		set: s,
	}
	c.List = append(c.List, ss)
}

type ShardSet struct {
	min int64
	max int64
	set *StandardCluster
}

func (ss ShardSet) InRange(v int64) bool {
	return ss.min <= v && v <= ss.max
}
