package dict

import (
	"errors"
	"fmt"
	"maps"
	"time"
)

var (
	ErrKeyNotFound           = errors.New("key not found")
	ErrValueTypeMismatch     = errors.New("value type mismatch")
	ErrCannotCoerceValueType = errors.New("cannot coerce value type")
)

type Dict map[string]any

func GetPtr[T any](d Dict, key string, cfn CoercePtrFunc[T]) (*T, error) {
	if v, ok := d[key]; !ok {
		return nil, fmt.Errorf("%w: %q", ErrKeyNotFound, key)
	} else if as, err := cfn(v); err != nil {
		return nil, fmt.Errorf("%w: expected %T, saw %T", ErrValueTypeMismatch, (*T)(nil), v)
	} else {
		return as, nil
	}
}

func MustGet[T any](d Dict, key string, cfn func(v any) (T, error)) T {
	if as, err := Coerce(d, key, cfn); err != nil {
		panic(err.Error())
	} else {
		return as
	}
}

func MustGetPtr[T any](d Dict, key string, cfn CoercePtrFunc[T]) *T {
	if as, err := GetPtr(d, key, cfn); err != nil {
		panic(err.Error())
	} else {
		return as
	}
}

func GetExact[T any](v any) (T, error) {
	var zero T
	if vt, ok := v.(T); !ok {
		return zero, fmt.Errorf("%w: %T to %T", ErrCannotCoerceValueType, v, 0)
	} else {
		return vt, nil
	}
}

func GetExactPtr[T any](v any) (*T, error) {
	if vt, ok := v.(T); !ok {
		if vt, ok := v.(*T); !ok {
			return nil, fmt.Errorf("%w: %T to %T", ErrCannotCoerceValueType, v, 0)
		} else {
			return vt, nil
		}
	} else {
		return &vt, nil
	}
}

// ShallowMerge clones the local Dict, overwriting any / all top-level keys with values from "other".  This does _not_
// recursively merge nested structures.
func (d Dict) ShallowMerge(other Dict) Dict {
	out := maps.Clone(d)
	for k, v := range other {
		out[k] = v
	}
	return out
}

func (d Dict) GetInt(key string) (int, error) {
	return Coerce(d, key, coerceInt)
}

func (d Dict) GetIntOr(key string, def int) int {
	return CoerceNonZeroComparableOr(d, key, coerceInt, def)
}

func (d Dict) MustGetInt(key string) int {
	return MustGet(d, key, coerceInt)
}

func (d Dict) GetString(key string) (string, error) {
	return Coerce(d, key, coerceString)
}

func (d Dict) GetStringOr(key string, def string) string {
	return CoerceNonZeroComparableOr(d, key, coerceString, def)
}

func (d Dict) MustGetString(key string) string {
	return MustGet(d, key, coerceString)
}

func (d Dict) GetBool(key string) (bool, error) {
	return Coerce(d, key, coerceBool)
}

func (d Dict) GetBoolOr(key string, def bool) bool {
	return CoerceOr(d, key, coerceBool, def)
}

func (d Dict) MustGetBool(key string) bool {
	return MustGet(d, key, coerceBool)
}

func (d Dict) GetDuration(key string) (time.Duration, error) {
	return Coerce(d, key, coerceDuration)
}

func (d Dict) GetDurationOr(key string, def time.Duration) time.Duration {
	return CoerceNonZeroComparableOr(d, key, coerceDuration, def)
}

func (d Dict) MustGetDuration(key string) time.Duration {
	return MustGet(d, key, coerceDuration)
}

func (d Dict) GetDict(key string) (Dict, error) {
	return Coerce(d, key, coerceDict)
}

func (d Dict) GetDictOr(key string, def Dict) Dict {
	return CoerceOr(d, key, coerceDict, def)
}

func (d Dict) MustGetDict(key string) Dict {
	return MustGet(d, key, coerceDict)
}

func (d Dict) GetStrings(key string) ([]string, error) {
	return Coerce(d, key, CoerceValueSlice(coerceString))
}

func (d Dict) MustGetStrings(key string) []string {
	return MustGet(d, key, CoerceValueSlice(coerceString))
}

func (d Dict) GetStringsOr(key string, def []string) []string {
	return CoerceOr(d, key, CoerceValueSlice(coerceString), def)
}

func (d Dict) GetInts(key string) ([]int, error) {
	return Coerce(d, key, CoerceValueSlice(coerceInt))
}

func (d Dict) MustGetInts(key string) []int {
	return MustGet(d, key, CoerceValueSlice(coerceInt))
}

func (d Dict) GetIntsOr(key string, def []int) []int {
	return CoerceOr(d, key, CoerceValueSlice(coerceInt), def)
}
