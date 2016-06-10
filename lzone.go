package lzoned

/*
	LZone (LazyZone) is a data-structure that manages an arbitrary number of
	user-defined 'zones'.  Each zone can be in either empty, clean or dirty
	states. Each zone can have a bound function for handling flushing and
	fetching. This is useful for things like models where certain parts of model
	needed to be loaded from different data sources depending on actions taken.
*/

// ---------------------------------------------------------------------------
// State for how fresh the data is in a single zone.
// ---------------------------------------------------------------------------
const (
	LZEmpty = iota // Data in this zone is unloaded and needs to be fetched.
	LZClean        // Data in this zone is clean
	LZDirty        // Data in this zone is dirty and needs to be flushed.
)

// An arena holds a number of zones
type LZArena struct {
	ops []LZOps
}

// Include this in the structure you want to support LazyZone.
type LZoned struct {
	LZStates []int
	LZArena  *LZArena
	LZObj    interface{}
}

// Operations for each zone for things like fetching and flushing
// changes.
type _LZFetch func(obj interface{})
type _LZFlush func(obj interface{})
type LZOps struct {
	Fetch _LZFetch
	Flush _LZFlush
}

func NewLZArena() *LZArena {
	return &LZArena{}
}

func (a *LZArena) AddZone(ops LZOps) int {
	a.ops = append(a.ops, ops)

	return len(a.ops) - 1
}

// Initiatialize an instance
func (lz *LZoned) Init(arena *LZArena, obj interface{}) {
	for i := 0; i < len(arena.ops); i++ {
		lz.LZStates = append(lz.LZStates, 0)
	}

	lz.LZArena = arena
	lz.LZObj = obj
}

func (lz *LZoned) GetState(zone int) int {
	return lz.LZStates[zone]
}

func (lz *LZoned) SetDirty(zone int) {
	lz.LZStates[zone] = LZDirty
}

func (lz *LZoned) SetClean(zone int) {
	lz.LZStates[zone] = LZClean
}

func (lz *LZoned) Fetch(zone int) {
	if lz.GetState(zone) == LZEmpty {
		zoneOps := lz.LZArena.ops[zone]
		zoneOps.Fetch(lz.LZObj)
	}

	lz.SetDirty(zone)
}

func (lz *LZoned) _flush(zone int) {
	if lz.GetState(zone) == LZDirty {
		zoneOps := lz.LZArena.ops[zone]
		zoneOps.Flush(lz.LZObj)
	}

	lz.SetClean(zone)
}

func (lz *LZoned) Flush() {
	for i := 0; i < len(lz.LZStates); i++ {
		lz._flush(i)
	}
}
