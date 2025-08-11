package mongodb

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"time"
)

type TimeCodec struct {
	//useLocalTimeZone bool
}

var tTime = reflect.TypeOf(time.Time{})
var emptyValue = reflect.Value{}

const timeFormatString = "2006-01-02T15:04:05.999Z07:00"

func (tc *TimeCodec) DecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tTime {
		return bsoncodec.ValueDecoderError{Name: "TimeDecodeValue", Types: []reflect.Type{tTime}, Received: val}
	}

	elem, err := tc.decodeType(dc, vr, tTime)
	if err != nil {
		return err
	}

	val.Set(elem)
	return nil
}

func (tc *TimeCodec) decodeType(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) {
	if t != tTime {
		return emptyValue, bsoncodec.ValueDecoderError{
			Name:     "TimeDecodeValue",
			Types:    []reflect.Type{tTime},
			Received: reflect.Zero(t),
		}
	}

	var timeVal time.Time
	switch vrType := vr.Type(); vrType {
	case bsontype.DateTime:
		dt, err := vr.ReadDateTime()
		if err != nil {
			return emptyValue, err
		}
		timeVal = time.Unix(dt/1000, dt%1000*1000000)
	case bsontype.String:
		// assume strings are in the isoTimeFormat
		timeStr, err := vr.ReadString()
		if err != nil {
			return emptyValue, err
		}
		timeVal, err = time.Parse(timeFormatString, timeStr)
		if err != nil {
			return emptyValue, err
		}
	case bsontype.Int64:
		i64, err := vr.ReadInt64()
		if err != nil {
			return emptyValue, err
		}
		timeVal = time.Unix(i64/1000, i64%1000*1000000)
	case bsontype.Timestamp:
		t, _, err := vr.ReadTimestamp()
		if err != nil {
			return emptyValue, err
		}
		timeVal = time.Unix(int64(t), 0)
	case bsontype.Null:
		if err := vr.ReadNull(); err != nil {
			return emptyValue, err
		}
	case bsontype.Undefined:
		if err := vr.ReadUndefined(); err != nil {
			return emptyValue, err
		}
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a time.Time", vrType)
	}
	tt := timeVal.UTC()
	timeVal = tt.Add(time.Duration(8) * time.Hour)
	return reflect.ValueOf(timeVal), nil
}

func (tc *TimeCodec) EncodeValue(context bsoncodec.EncodeContext, writer bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tTime {
		return bsoncodec.ValueEncoderError{Name: "TimeEncodeValue", Types: []reflect.Type{tTime}, Received: val}
	}
	tt := val.Interface().(time.Time)
	dt := primitive.NewDateTimeFromTime(tt)
	return writer.WriteDateTime(int64(dt))
}
