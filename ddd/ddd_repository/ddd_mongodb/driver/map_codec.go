package driver

import (
	"encoding"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"reflect"
)

type MapCodec struct {
	base *bsoncodec.MapCodec
}

var defaultMapCodec = NewMapCodec()

// NewMapCodec returns a MapCodec with options opts.
//
// Deprecated: Use [go.mongodb.org/mongo-driver/bson.NewRegistry] to get a registry with the
// MapCodec registered.
func NewMapCodec(opts ...*bsonoptions.MapCodecOptions) *MapCodec {
	base := bsoncodec.NewMapCodec(opts...)
	return &MapCodec{base: base}
}

// DecodeValue is the ValueDecoder for map[string/decimal]* types.
func (mc *MapCodec) DecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if val.Kind() != reflect.Map || (!val.CanSet() && val.IsNil()) {
		return bsoncodec.ValueDecoderError{Name: "MapDecodeValue", Kinds: []reflect.Kind{reflect.Map}, Received: val}
	}

	switch vrType := vr.Type(); vrType {
	case bsontype.Type(0), bsontype.EmbeddedDocument:
	case bsontype.Null:
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadNull()
	case bsontype.Undefined:
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into a %s", vrType, val.Type())
	}

	dr, err := vr.ReadDocument()
	if err != nil {
		return err
	}

	if val.IsNil() {
		val.Set(reflect.MakeMap(val.Type()))
	}

	if val.Len() > 0 && (mc.base.DecodeZerosMap || dc.zeroMaps) {
		clearMap(val)
	}

	eType := val.Type().Elem()
	decoder, err := dc.LookupDecoder(eType)
	if err != nil {
		return err
	}
	eTypeDecoder, _ := decoder.(typeDecoder)

	if eType == tEmpty {
		dc.Ancestor = val.Type()
	}

	keyType := val.Type().Key()

	for {
		key, vr, err := dr.ReadElement()
		if err == bsonrw.ErrEOD {
			break
		}
		if err != nil {
			return err
		}

		k, err := mc.decodeKey(key, keyType)
		if err != nil {
			return err
		}

		elem, err := decodeTypeOrValueWithInfo(decoder, eTypeDecoder, dc, vr, eType, true)
		if err != nil {
			return newDecodeError(key, err)
		}

		val.SetMapIndex(k, elem)
	}
	return nil
}

func clearMap(m reflect.Value) {
	var none reflect.Value
	for _, k := range m.MapKeys() {
		m.SetMapIndex(k, none)
	}
}

func (mc *MapCodec) decodeKey(key string, keyType reflect.Type) (reflect.Value, error) {
	keyVal := reflect.ValueOf(key)
	var err error
	switch {
	// First, if EncodeKeysWithStringer is not enabled, try to decode withKeyUnmarshaler
	case !mc.base.EncodeKeysWithStringer && reflect.PtrTo(keyType).Implements(keyUnmarshalerType):
		keyVal = reflect.New(keyType)
		v := keyVal.Interface().(KeyUnmarshaler)
		err = v.UnmarshalKey(key)
		keyVal = keyVal.Elem()
	// Try to decode encoding.TextUnmarshalers.
	case reflect.PtrTo(keyType).Implements(textUnmarshalerType):
		keyVal = reflect.New(keyType)
		v := keyVal.Interface().(encoding.TextUnmarshaler)
		err = v.UnmarshalText([]byte(key))
		keyVal = keyVal.Elem()
	// Otherwise, go to type specific behavior
	default:
		switch keyType.Kind() {
		case reflect.String:
			keyVal = reflect.ValueOf(key).Convert(keyType)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, parseErr := strconv.ParseInt(key, 10, 64)
			if parseErr != nil || reflect.Zero(keyType).OverflowInt(n) {
				err = fmt.Errorf("failed to unmarshal number key %v", key)
			}
			keyVal = reflect.ValueOf(n).Convert(keyType)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n, parseErr := strconv.ParseUint(key, 10, 64)
			if parseErr != nil || reflect.Zero(keyType).OverflowUint(n) {
				err = fmt.Errorf("failed to unmarshal number key %v", key)
				break
			}
			keyVal = reflect.ValueOf(n).Convert(keyType)
		case reflect.Float32, reflect.Float64:
			if mc.EncodeKeysWithStringer {
				parsed, err := strconv.ParseFloat(key, 64)
				if err != nil {
					return keyVal, fmt.Errorf("Map key is defined to be a decimal type (%v) but got error %v", keyType.Kind(), err)
				}
				keyVal = reflect.ValueOf(parsed)
				break
			}
			fallthrough
		default:
			return keyVal, fmt.Errorf("unsupported key type: %v", keyType)
		}
	}
	return keyVal, err
}
