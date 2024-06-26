package dict

import (
	"fmt"
	"strconv"
	"time"
)

type (
	CoerceValueFunc[T any] func(any) (T, error)

	CoercePtrFunc[T any] func(any) (*T, error)
)

type (
	durer interface {
		Duration() time.Duration
	}
	durerErr interface {
		Duration() (time.Duration, error)
	}
	stringer interface {
		String() string
	}
	toStringer interface {
		ToString() string
	}
)

func coerceString(v any) (string, error) {
	switch v.(type) {
	case string:
		return v.(string), nil

	case int, int64, int32, float64, float32, uint, uint64, uint32:
		return fmt.Sprintf("%d", v), nil

	case bool:
		return strconv.FormatBool(v.(bool)), nil

	case time.Duration:
		return v.(time.Duration).String(), nil

	case stringer:
		return v.(stringer).String(), nil

	case toStringer:
		return v.(toStringer).ToString(), nil

	default:
		return "", fmt.Errorf("%w: %T to %T", ErrCannotCoerceValueType, v, "")
	}
}

func coerceInt(v any) (int, error) {
	switch v.(type) {
	case int:
		return v.(int), nil
	case int64:
		return int(v.(int64)), nil
	case int32:
		return int(v.(int32)), nil
	case float64:
		return int(v.(float64)), nil
	case float32:
		return int(v.(float32)), nil
	case uint:
		return int(v.(uint)), nil
	case uint64:
		return int(v.(uint64)), nil
	case uint32:
		return int(v.(uint32)), nil

	case bool:
		if v.(bool) {
			return 1, nil
		}
		return 0, nil

	case string:
		return strconv.Atoi(v.(string))

	case time.Duration:
		return int(v.(time.Duration).Nanoseconds()), nil

	case durer:
		return int(v.(durer).Duration()), nil

	case durerErr:
		d, err := v.(durerErr).Duration()
		return int(d), err

	default:
		return 0, fmt.Errorf("%w: %T to %T", ErrCannotCoerceValueType, v, 0)
	}
}

func coerceBool(v any) (bool, error) {
	switch v.(type) {
	case bool:
		return v.(bool), nil

	case string:
		return strconv.ParseBool(v.(string))

	case int:
		return v.(int) > 0, nil
	case int64:
		return v.(int64) > 0, nil
	case int32:
		return v.(int32) > 0, nil
	case float64:
		return v.(float64) > 0, nil
	case float32:
		return v.(float32) > 0, nil
	case uint:
		return v.(uint) > 0, nil
	case uint64:
		return v.(uint64) > 0, nil
	case uint32:
		return v.(uint32) > 0, nil

	default:
		return false, fmt.Errorf("%w: %T to %T", ErrCannotCoerceValueType, v, 0)
	}
}

func coerceDuration(v any) (time.Duration, error) {
	switch v.(type) {
	case time.Duration:
		return v.(time.Duration), nil

	case int:
		return time.Duration(v.(int)), nil
	case int64:
		return time.Duration(v.(int64)), nil
	case int32:
		return time.Duration(v.(int32)), nil
	case float64:
		return time.Duration(v.(float64)), nil
	case float32:
		return time.Duration(v.(float32)), nil
	case uint:
		return time.Duration(v.(uint)), nil
	case uint64:
		return time.Duration(v.(uint64)), nil
	case uint32:
		return time.Duration(v.(uint32)), nil

	case string:
		return time.ParseDuration(v.(string))

	case durer:
		return v.(durer).Duration(), nil

	case durerErr:
		return v.(durerErr).Duration()

	default:
		return 0, fmt.Errorf("%w: %T to %T", ErrCannotCoerceValueType, v, 0)
	}
}

func coerceDict(v any) (Dict, error) {
	switch v.(type) {
	case Dict:
		return v.(Dict), nil
	case map[string]any:
		return v.(map[string]any), nil

	default:
		return nil, fmt.Errorf("%w: %T to %T", ErrCannotCoerceValueType, v, 0)
	}
}

func CoerceExact[T any](v any) (T, error) {
	if vt, ok := v.(T); !ok {
		var zero T
		return zero, fmt.Errorf("%w: %T to %T", ErrCannotCoerceValueType, v, 0)
	} else {
		return vt, nil
	}
}

func CoerceExactPtr[T any](v any) (*T, error) {
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
