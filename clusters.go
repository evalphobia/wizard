package wizard

var defaultSlot int64 = 9973

type Cluster interface {
	Table() interface{}
	SelectBySlot(i int64) *StandardCluster
	Master() *Node
	Masters() []*Node
	Slave() *Node
}

type Wizard struct {
	clusters map[interface{}]Cluster
	slot     int64
}

func NewWizard() *Wizard {
	return &Wizard{
		clusters: make(map[interface{}]Cluster),
		slot:     defaultSlot,
	}
}

func (w *Wizard) Slot(v int64) {
	w.slot = v
}

func (w *Wizard) CreateCluster(obj interface{}, db interface{}) *StandardCluster {
	c := NewCluster(db)
	w.clusters[NormalizeValue(obj)] = c
	return c
}

func (w *Wizard) CreateShardCluster(obj interface{}, slot int64) *ShardCluster {
	c := &ShardCluster{
		slotsize: slot,
	}
	w.clusters[NormalizeValue(obj)] = c
	return c
}

func (w *Wizard) Select(obj interface{}) *StandardCluster {
	c, ok := w.clusters[NormalizeValue(obj)]
	if !ok {
		return nil
	}
	switch v := c.(type) {
	case *StandardCluster:
		return v
	case *ShardCluster:
		return v.SelectBySlot(getID(obj))
	default:
		return nil
	}
}

func (w *Wizard) SelectBySlot(obj interface{}, id interface{}) *StandardCluster {
	c, ok := w.clusters[NormalizeValue(obj)]
	if !ok {
		return nil
	}
	i := getInt64(id)
	return c.SelectBySlot(i)
}

func Slot(v int64) {
	defaultSlot = v
}
