package dict

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"time"
)

var (
	ErrKeyNotFound           = errors.New("key not found")
	ErrValueTypeMismatch     = errors.New("value type mismatch")
	ErrCannotCoerceValueType = errors.New("cannot coerce value type")
)

type Dict map[string]any

func Get[T any](d Dict, key string, cfn CoerceValueFunc[T]) (T, error) {
	var zero T
	if v, ok := d[key]; !ok {
		return zero, fmt.Errorf("%w: %q", ErrKeyNotFound, key)
	} else if as, err := cfn(v); err != nil {
		return zero, fmt.Errorf("%w: expected %T, saw %T", ErrValueTypeMismatch, zero, v)
	} else {
		return as, nil
	}
}

func MustGet[T any](d Dict, key string, cfn CoerceValueFunc[T]) T {
	if as, err := Get(d, key, cfn); err != nil {
		panic(err.Error())
	} else {
		return as
	}
}

func GetOr[T any](d Dict, key string, cfn CoerceValueFunc[T], or T) T {
	if v, err := Get[T](d, key, cfn); err != nil {
		return or
	} else {
		return v
	}
}

// GetNonZeroComparableOr is significantly different from GetOr in that it considers a comparable's zero-val to be
// "empty", thus returning the value provided to "or".
//
// This is an important distinction.
func GetNonZeroComparableOr[T comparable](d Dict, key string, cfn CoerceValueFunc[T], or T) T {
	var zero T
	if v, err := Get[T](d, key, cfn); err != nil || v == zero {
		return or
	} else {
		return v
	}
}

func GetExact[T any](d Dict, key string) (T, error) {
	return Get(d, key, CoerceExact[T])
}

func MustGetExact[T any](d Dict, key string) T {
	return MustGet(d, key, CoerceExact[T])
}

func GetExactPtr[T any](d Dict, key string) (*T, error) {
	return GetPtr(d, key, CoerceExactPtr[T])
}

func MustGetExactPtr[T any](d Dict, key string) *T {
	return MustGetPtr(d, key, CoerceExactPtr[T])
}

func GetPtr[T any](d Dict, key string, cfn CoercePtrFunc[T]) (*T, error) {
	if v, ok := d[key]; !ok {
		return nil, fmt.Errorf("%w: %q", ErrKeyNotFound, key)
	} else if as, err := cfn(v); err != nil {
		return nil, fmt.Errorf("%w: expected %T, saw %T", ErrValueTypeMismatch, (*T)(nil), v)
	} else {
		return as, nil
	}
}

func MustGetPtr[T any](d Dict, key string, cfn CoercePtrFunc[T]) *T {
	if as, err := GetPtr(d, key, cfn); err != nil {
		panic(err.Error())
	} else {
		return as
	}
}

func GetSlice[T any](cfn CoerceValueFunc[T]) func(any) ([]T, error) {
	return func(v any) ([]T, error) {
		switch v.(type) {
		case []T:
			return slices.Clone(v.([]T)), nil

		case []any:
			var (
				out  []T
				errs []error
			)
			for i, vv := range v.([]any) {
				if tv, err := cfn(vv); err != nil {
					errs = append(errs, fmt.Errorf("index %d: %w", i, err))
				} else {
					out = append(out, tv)
				}
			}
			if len(errs) > 0 {
				return nil, errors.Join(errs...)
			}
			return out, nil

		default:
			var tzero []T
			return nil, fmt.Errorf("%w: %T to %T", ErrCannotCoerceValueType, v, tzero)
		}
	}
}

func GetSlicePtr[T any](cfn CoercePtrFunc[T]) func(any) ([]*T, error) {
	return func(v any) ([]*T, error) {
		switch v.(type) {
		case []*T:
			return slices.Clone(v.([]*T)), nil

		case []any:
			var (
				out  []*T
				errs []error
			)
			for i, vv := range v.([]any) {
				if tv, err := cfn(vv); err != nil {
					errs = append(errs, fmt.Errorf("index %d: %w", i, err))
				} else {
					out = append(out, tv)
				}
			}
			if len(errs) > 0 {
				return nil, errors.Join(errs...)
			}
			return out, nil

		default:
			var tzero []*T
			return nil, fmt.Errorf("%w: %T to %T", ErrCannotCoerceValueType, v, tzero)
		}
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
	return Get(d, key, coerceInt)
}

func (d Dict) GetIntOr(key string, def int) int {
	return GetNonZeroComparableOr(d, key, coerceInt, def)
}

func (d Dict) MustGetInt(key string) int {
	return MustGet(d, key, coerceInt)
}

func (d Dict) GetString(key string) (string, error) {
	return Get(d, key, coerceString)
}

func (d Dict) GetStringOr(key string, def string) string {
	return GetNonZeroComparableOr(d, key, coerceString, def)
}

func (d Dict) MustGetString(key string) string {
	return MustGet(d, key, coerceString)
}

func (d Dict) GetBool(key string) (bool, error) {
	return Get(d, key, coerceBool)
}

func (d Dict) GetBoolOr(key string, def bool) bool {
	return GetOr(d, key, coerceBool, def)
}

func (d Dict) MustGetBool(key string) bool {
	return MustGet(d, key, coerceBool)
}

func (d Dict) GetDuration(key string) (time.Duration, error) {
	return Get(d, key, coerceDuration)
}

func (d Dict) GetDurationOr(key string, def time.Duration) time.Duration {
	return GetNonZeroComparableOr(d, key, coerceDuration, def)
}

func (d Dict) MustGetDuration(key string) time.Duration {
	return MustGet(d, key, coerceDuration)
}

func (d Dict) GetDict(key string) (Dict, error) {
	return Get(d, key, coerceDict)
}

func (d Dict) GetDictOr(key string, def Dict) Dict {
	return GetOr(d, key, coerceDict, def)
}

func (d Dict) MustGetDict(key string) Dict {
	return MustGet(d, key, coerceDict)
}

func (d Dict) GetStrings(key string) ([]string, error) {
	return Get(d, key, GetSlice(coerceString))
}

func (d Dict) MustGetStrings(key string) []string {
	return MustGet(d, key, GetSlice(coerceString))
}

func (d Dict) GetStringsOr(key string, def []string) []string {
	return GetOr(d, key, GetSlice(coerceString), def)
}

func (d Dict) GetInts(key string) ([]int, error) {
	return Get(d, key, GetSlice(coerceInt))
}

func (d Dict) MustGetInts(key string) []int {
	return MustGet(d, key, GetSlice(coerceInt))
}

func (d Dict) GetIntsOr(key string, def []int) []int {
	return GetOr(d, key, GetSlice(coerceInt), def)
}
