package ttlcache

import (
	"fmt"
)

const _EvictionReasonName = "RemovedEvictedSizeExpiredClosed"

var _EvictionReasonIndex = [...]uint8{0, 7, 18, 25, 31}

func (i EvictionReason) String() string {
	if i < 0 || i >= EvictionReason(len(_EvictionReasonIndex)-1) {
		return fmt.Sprintf("EvictionReason(%d)", i)
	}
	return _EvictionReasonName[_EvictionReasonIndex[i]:_EvictionReasonIndex[i+1]]
}

var _EvictionReasonValues = []EvictionReason{0, 1, 2, 3}

var _EvictionReasonNameToValueMap = map[string]EvictionReason{
	_EvictionReasonName[0:7]:   0,
	_EvictionReasonName[7:18]:  1,
	_EvictionReasonName[18:25]: 2,
	_EvictionReasonName[25:31]: 3,
}

// EvictionReasonString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func EvictionReasonString(s string) (EvictionReason, error) {
	if val, ok := _EvictionReasonNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to EvictionReason values", s)
}

// EvictionReasonValues returns all values of the enum
func EvictionReasonValues() []EvictionReason {
	return _EvictionReasonValues
}

// IsAEvictionReason returns "true" if the value is listed in the enum definition. "false" otherwise
func (i EvictionReason) IsAEvictionReason() bool {
	for _, v := range _EvictionReasonValues {
		if i == v {
			return true
		}
	}
	return false
}
