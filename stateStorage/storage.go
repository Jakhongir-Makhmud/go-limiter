package stateStorage

import (
	"context"
	"time"

	"github.com/xaionaro-go/atomicmap"
)

func New() *stateStorage {
	aMap := atomicmap.New()
	aMap.SetThreadSafety(true)
	aMap.SetForbidGrowing(false)
	return &stateStorage{
		state: aMap,
	}
}

type stateStorage struct {
	state atomicmap.Map
}

type ValueAndDur struct {
	value    uint64
	saveTime time.Time
	dur      time.Duration
}

func (s *stateStorage) SetState(ctx context.Context, key string, val uint64, dur time.Duration) error {
	valDur := ValueAndDur{value: val, dur: dur, saveTime: time.Now()}
	return s.state.Set(key, valDur)
}

func (s *stateStorage) GetState(ctx context.Context, key string) (uint64, error) {
	valInterf, err := s.state.Get(key)
	if err != nil {
		return 0, err
	}

	val, ok := valInterf.(ValueAndDur)
	if ok {
		if val.saveTime.Add(val.dur).Before(time.Now()) {
			return val.value, nil
		}
		s.state.Unset(key)
	}

	return 0, nil
}
