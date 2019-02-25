package etc

import (
	"fmt"
	"io"
	"reflect"
)

type FixedUint16 uint16
type FixedUint32 uint32
type FixedUint64 uint64
type FixedInt16 int16
type FixedInt32 int32
type FixedInt64 int64
type FixedUint16BE uint16
type FixedUint32BE uint32
type FixedUint64BE uint64
type FixedInt16BE int16
type FixedInt32BE int32
type FixedInt64BE int64

type Encoder struct {
	*Buffer
}

func Marshal(v interface{}) ([]byte, error) {
	buf := NewBuffer()
	enc := NewEncoder(buf)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Unmarshal(b []byte, v interface{}) error {
	buf := FromBytes(b)
	dec := &Decoder{buf}
	return dec.Decode(v)
}

func NewEncoder(out io.Writer) *Encoder {
	return &Encoder{
		DummyWriter(out),
	}
}

func intClass(v interface{}) (underlying interface{}, isInt bool, fixed bool, big bool, bits int, signed bool) {
	if reflect.ValueOf(v).Kind() == reflect.Ptr {
		v = reflect.ValueOf(v).Elem().Interface()
	}

	isInt = true
	big = false
	switch m := v.(type) {
	case FixedInt16:
		bits = 16
		signed = true
		fixed = true
		underlying = int16(m)
	case FixedInt32:
		bits = 32
		signed = true
		fixed = true
		underlying = int32(m)
	case FixedInt64:
		bits = 64
		signed = true
		fixed = true
		underlying = int64(m)
	case FixedUint16:
		bits = 16
		signed = false
		fixed = true
		underlying = uint16(m)
	case FixedUint32:
		bits = 32
		signed = false
		fixed = true
		underlying = uint32(m)
	case FixedUint64:
		bits = 64
		signed = false
		fixed = true
		underlying = uint64(m)
	case FixedInt16BE:
		big = true
		bits = 16
		signed = true
		fixed = true
		underlying = int16(m)
	case FixedInt32BE:
		big = true
		bits = 32
		signed = true
		fixed = true
		underlying = int32(m)
	case FixedInt64BE:
		big = true
		bits = 64
		signed = true
		fixed = true
		underlying = int64(m)
	case FixedUint16BE:
		big = true
		bits = 16
		signed = false
		fixed = true
		underlying = uint16(m)
	case FixedUint32BE:
		big = true
		bits = 32
		signed = false
		fixed = true
		underlying = uint32(m)
	case FixedUint64BE:
		big = true
		bits = 64
		signed = false
		fixed = true
		underlying = uint64(m)
	case uint64:
		big = false
		bits = 64
		signed = false
		fixed = false
		underlying = m
	case int64:
		big = false
		bits = 64
		signed = true
		fixed = false
		underlying = m
	default:
		isInt = false
	}
	return
}

func (e *Encoder) Encode(v interface{}) error {
	object := reflect.ValueOf(v)

	realType, isInt, fixed, bigEndian, bits, signed := intClass(v)

	if isInt {
		e.WriteTypedInt(fixed, bits, signed, bigEndian, &realType)
		return nil
	}

	// Check non-int types
	switch object.Kind() {
	case reflect.Ptr:
		return e.Encode(object.Elem().Interface())
	case reflect.Struct:
		for i := 0; i < object.NumField(); i++ {
			field := object.Field(i)
			exported := !reflect.TypeOf(v).Field(i).Anonymous
			if exported {
				if err := e.Encode(field.Interface()); err != nil {
					return err
				}
			}
		}
	case reflect.Slice:
		// bytes
		if reflect.TypeOf(v).Elem().Kind() == reflect.Uint8 {
			sli := v.([]byte)
			if err := e.Encode(uint64(len(sli))); err != nil {
				return err
			}
			if _, err := e.Write(sli); err != nil {
				return err
			}
			return nil
		}

		sz := object.Len()
		if err := e.Encode(uint64(sz)); err != nil {
			return err
		}
		for i := 0; i < sz; i++ {
			if err := e.Encode(object.Index(i).Interface()); err != nil {
				return err
			}
		}
	case reflect.Array:
		sz := object.Len()
		for i := 0; i < sz; i++ {
			e.Encode(object.Index(i).Interface())
		}
	case reflect.String:
		e.WriteUString(v.(string))
	case reflect.Float32:
		e.WriteFloat32(v.(float32))
	case reflect.Float64:
		e.WriteFloat64(v.(float64))
	case reflect.Bool:
		e.WriteBool(v.(bool))
	case reflect.Map:
		sz := object.Len()
		if err := e.Encode(uint64(sz)); err != nil {
			return err
		}

		keys := object.MapKeys()
		for _, key := range keys {
			if err := e.Encode(key.Interface()); err != nil {
				return err
			}
		}

		for _, key := range keys {
			value := object.MapIndex(key)
			if err := e.Encode(value.Interface()); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("etc: unable to encode type %s", object.String())
	}

	return nil
}

type Decoder struct {
	*Buffer
}

func NewDecoder(i io.Reader, size int64) *Decoder {
	return &Decoder{
		DummyReader(i, size),
	}
}

func (d *Decoder) Decode(v interface{}) error {
	object := reflect.ValueOf(v)

	realType, isInt, fixed, bigEndian, bits, signed := intClass(v)
	if isInt {
		d.ReadTypedInt(fixed, bits, signed, bigEndian, &realType)
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(realType))
		return nil
	}

	switch v.(type) {
	case UUID:
		v = d.ReadUUID()
		return nil
	}

	if object.Kind() == reflect.Slice {
		object = object.Addr()
	}

	if object.Kind() == reflect.Ptr {
		realType, isInt, fixed, bigEndian, bits, signed = intClass(v)
		if isInt {
			d.ReadTypedInt(fixed, bits, signed, bigEndian, &realType)
			reflect.ValueOf(v).Elem().Set(reflect.ValueOf(realType))
			return nil
		}

		switch object.Elem().Kind() {
		case reflect.Ptr:
			iface := object.Elem()

			if iface.IsNil() {
				nw := reflect.New(object.Type().Elem().Elem())
				err := d.Decode(nw.Interface())
				object.Elem().Set(nw)
				if err != nil {
					return err
				}
				return nil
			}

			return d.Decode(iface.Interface())
		case reflect.Struct:
			for i := 0; i < object.Elem().NumField(); i++ {
				constructor := object.Elem().Type().Field(i)
				exported := !object.Elem().Type().Field(i).Anonymous
				if exported {
					out := reflect.New(constructor.Type)
					err := d.Decode(out.Interface())
					if err != nil {
						return err
					}
					object.Elem().Field(i).Set(out.Elem())
				}
			}
			return nil
		case reflect.String:
			value := d.ReadUString()
			object.Elem().SetString(value)
		case reflect.Float32:
			value := d.ReadFloat32()
			object.Elem().SetFloat(float64(value))
		case reflect.Array:
			for i := 0; i < object.Elem().Len(); i++ {
				err := d.Decode(object.Elem().Index(i).Addr().Interface())
				if err != nil {
					return err
				}
			}
		case reflect.Slice:
			sz := int(d.ReadUint())
			slice := reflect.MakeSlice(object.Elem().Type(), sz, sz)
			object.Elem().Set(slice)
			constructor := slice.Index(0).Type()
			for i := 0; i < int(sz); i++ {
				out := reflect.New(constructor)
				err := d.Decode(out.Interface())
				if err != nil {
					return err
				}
				object.Elem().Index(i).Set(out.Elem())
			}
			return nil
		case reflect.Map:
			sz := int(d.ReadUint())
			if object.Elem().IsNil() {
				newMap := reflect.MakeMap(object.Elem().Type())
				fmt.Println("Type of new map", newMap.Type())
				fmt.Println("Type of new map", newMap.Type().Elem())
				object.Elem().Set(newMap)
				fmt.Println(object.Elem().Type())
				fmt.Println(object.Elem().IsNil())
			}

			kt := object.Elem().Type().Key()
			vt := object.Elem().Type().Elem()

			var keys []reflect.Value

			for i := 0; i < sz; i++ {
				idxKey := reflect.New(kt)
				err := d.Decode(idxKey.Interface())
				if err != nil {
					return err
				}
				keys = append(keys, idxKey)
			}

			for i := 0; i < sz; i++ {
				idxVal := reflect.New(vt)
				err := d.Decode(idxVal.Interface())
				if err != nil {
					return err
				}
				object.Elem().SetMapIndex(keys[i].Elem(), idxVal.Elem())
			}
			return nil
		default:
			return fmt.Errorf("etc: unknown type: %s", object.Elem().Type().String())
		}
		return nil
	}

	return fmt.Errorf("etc: unknown type: %s", object.String())
}
