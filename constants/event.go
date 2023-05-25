package constants

import (
	"sync"
)

type eventItem struct {
	eventid int
	params  []interface{}
}
type eventCB struct {
	eventid int
	cb      func(evtid int, args ...interface{})
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////
type IEventManager interface {
	AddEventLister(evtid int, cb func(evtid int, args ...interface{})) *eventCB
	RemoveEventLister(c *eventCB)

	DispatchEvent(evtid int, params ...interface{})
}

type eventManager struct {
	eventsLis sync.Map
	eventBuf  chan *eventItem
}

func (m *eventManager) AddEventLister(evtid int, cb func(evtid int, args ...interface{})) *eventCB {
	var es1 []*eventCB
	es, ok := m.eventsLis.Load(evtid)
	if !ok {
		es1 = []*eventCB{}
	} else {
		es1 = es.([]*eventCB)
	}
	evtcb := &eventCB{
		eventid: evtid,
		cb:      cb,
	}
	es = append(es1, evtcb)
	m.eventsLis.Store(evtid, es)

	return evtcb
}
func (m *eventManager) RemoveEventLister(c *eventCB) {
	if c == nil {
		return
	}

	es, ok := m.eventsLis.Load(c.eventid)
	if ok {
		es1 := es.([]*eventCB)
		idx := -1
		for n, cb := range es1 {
			if cb == c {
				idx = n
				break
			}
		}
		if idx >= 0 {
			es = append(es1[:idx], es1[idx+1:]...)
			m.eventsLis.Store(c.eventid, es)
		}
	}
}
func (m *eventManager) DispatchEvent(evtid int, params ...interface{}) {
	t := &eventItem{
		eventid: evtid,
		params:  append([]interface{}{}, params...),
	}

	m.eventBuf <- t
}

func (m *eventManager) run() {
	for evt := range m.eventBuf {
		if evt != nil {
			es, ok := m.eventsLis.Load(evt.eventid)
			if ok {
				es1 := es.([]*eventCB)
				for _, cb := range es1 {
					go cb.cb(evt.eventid, evt.params...)
				}
			}
		}
	}
}
func NewEventManager() IEventManager {
	t := &eventManager{
		eventsLis: sync.Map{},
		eventBuf:  make(chan *eventItem, 1024),
	}
	go t.run()
	return t
}
