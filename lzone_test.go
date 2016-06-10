package lzoned

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLZone(t *testing.T) {
	Convey("Can use Lazy-Zones properly", t, func() {
		type Foo struct {
			LZ LZoned

			x int
		}
		foo := Foo{x: 4}

		// Create a new arena that will store
		// operations.
		arena := NewLZArena()

		var fetched bool
		var flushed bool
		var fetchedObject interface{}
		var flushedObject interface{}
		var flushedTags []string
		// Add zones to our arena
		zoneA := arena.AddZone(LZOps{
			Fetch: func(obj interface{}) {
				fetched = true
				fetchedObject = obj
			},
			Flush: func(obj interface{}, tags []string) {
				flushed = true
				flushedObject = obj
				for _, tag := range tags {
					flushedTags = append(flushedTags, tag)
				}
			},
		})

		//Before fetching we should be clean
		foo.LZ.Init(arena, foo)
		So(foo.LZ.GetState(zoneA), ShouldEqual, LZEmpty)

		// Explicitly set dirty
		foo.LZ.SetDirty(zoneA)
		So(foo.LZ.GetState(zoneA), ShouldEqual, LZDirty)

		// Explicitly set clean
		foo.LZ.SetClean(zoneA)
		So(foo.LZ.GetState(zoneA), ShouldEqual, LZClean)

		// Clear
		foo = Foo{}
		foo.LZ.Init(arena, foo)

		// Does fetch a zone
		So(fetched, ShouldEqual, false)
		foo.LZ.Fetch(zoneA)
		So(fetched, ShouldEqual, true)
		So(fetchedObject.(Foo).x, ShouldEqual, foo.x)

		// Doesn't fetch twice if not empty
		fetched = false
		foo.LZ.Fetch(zoneA)
		So(fetched, ShouldEqual, false)

		// Doesn't flush if not dirty
		foo.LZ.SetClean(zoneA)
		foo.LZ.Flush()
		So(flushed, ShouldEqual, false)

		// Does flush if dirty and can set keys
		foo.LZ.SetDirty(zoneA, "a")
		foo.LZ.SetDirty(zoneA, "b")
		foo.LZ.Flush()
		So(flushed, ShouldEqual, true)
		So(flushedObject.(Foo).x, ShouldEqual, foo.x)
		So(len(flushedTags), ShouldEqual, 2)

		foo.LZ.SetDirty(zoneA)
		flushedTags = []string{}
		foo.LZ.Flush()
		So(len(flushedTags), ShouldEqual, 0)
	})
}
