package lzoned

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
		var shouldError = false
		// Add zones to our arena
		zoneA := arena.AddZone(LZOps{
			Fetch: func(obj interface{}) {
				fetched = true
				fetchedObject = obj
			},
			Commit: func(obj interface{}, tags []string) error {
				flushed = true
				flushedObject = obj
				for _, tag := range tags {
					flushedTags = append(flushedTags, tag)
				}

				if shouldError {
					return fmt.Errorf("holah")
				}

				return nil
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
		foo.LZ.Commit()
		So(flushed, ShouldEqual, false)

		// Stays empty if flushed
		foo.LZ.LZStates[zoneA].state = 0
		foo.LZ.Commit()
		So(foo.LZ.GetState(zoneA), ShouldEqual, 0)

		// Does flush if dirty and can set keys
		foo.LZ.SetDirty(zoneA, "a")
		foo.LZ.SetDirty(zoneA, "b")
		err := foo.LZ.Commit()
		So(flushed, ShouldEqual, true)
		So(flushedObject.(Foo).x, ShouldEqual, foo.x)
		So(len(flushedTags), ShouldEqual, 2)
		So(err, ShouldEqual, nil)

		// Test zone tag clear & errors
		foo.LZ.SetDirty(zoneA)
		flushedTags = []string{}
		shouldError = true
		err = foo.LZ.Commit()
		So(len(flushedTags), ShouldEqual, 0)
		So(err, ShouldNotEqual, nil)
	})
}
