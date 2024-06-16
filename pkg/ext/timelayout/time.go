package timelayout

// On feature: https://github.com/golang/go/issues/21990 - proposal: encoding/json: support struct tag for time.Format in JSON Marshaller/Unmarshaller
// Ref: https://developpaper.com/how-to-customize-the-time-format-for-json-serialization-of-golang-structure/
// Src: https://github.com/liamylian/jsontime
// Why need this ext? Lack of format time, cause format time is messy each time
//
// Usage in entity:
// Override MarshalJSON
// 		Default format YYYY-mm-dd HH:mm:ss
//
// time_format:"date" or time_format:"time"  .. date only, or time only in metadata
//
// In the

import (
    "time"
    "unsafe"

    "store-api/pkg/data/constant"

    "github.com/sirupsen/logrus"

    jsoniter "github.com/json-iterator/go"
)

// time format alias
const (
    Date = "date"
    Time = "time"

    ANSIC       = "ANSIC"
    UnixDate    = "UnixDate"
    RubyDate    = "RubyDate"
    RFC822      = "RFC822"
    RFC822Z     = "RFC822Z"
    RFC850      = "RFC850"
    RFC1123     = "RFC1123"
    RFC1123Z    = "RFC1123Z"
    RFC3339     = "RFC3339"
    RFC3339Nano = "RFC3339Nano"
    Kitchen     = "Kitchen"
    Stamp       = "Stamp"
    StampMilli  = "StampMilli"
    StampMicro  = "StampMicro"
    StampNano   = "StampNano"
)

// time zone alias
const (
    Local = "Local"
    UTC   = "UTC"
)

const (
    tagNameTimeFormat   = "time_format"
    tagNameTimeLocation = "time_location"
)

var TimeFormat = jsoniter.ConfigCompatibleWithStandardLibrary

var formatAlias = map[string]string{
    Date:        constant.DefaultDateLayout,
    Time:        constant.DefaultTimeLayout,
    ANSIC:       time.ANSIC,
    UnixDate:    time.UnixDate,
    RubyDate:    time.RubyDate,
    RFC822:      time.RFC822,
    RFC822Z:     time.RFC822Z,
    RFC850:      time.RFC850,
    RFC1123:     time.RFC1123,
    RFC1123Z:    time.RFC1123Z,
    RFC3339:     time.RFC3339,
    RFC3339Nano: time.RFC3339Nano,
    Kitchen:     time.Kitchen,
    Stamp:       time.Stamp,
    StampMilli:  time.StampMilli,
    StampMicro:  time.StampMicro,
    StampNano:   time.StampNano,
}

var localeAlias = map[string]*time.Location{
    Local: time.Local,
    UTC:   time.UTC,
}

var (
    defaultFormat = time.RFC3339
    defaultLocale = time.Local
)

func Init() {
    logrus.Debugln("Register Time Format Extension")
    TimeFormat.RegisterExtension(&CustomTimeExtension{})
}

func AddTimeFormatAlias(alias, format string) {
    formatAlias[alias] = format
}

func AddLocaleAlias(alias string, locale *time.Location) {
    localeAlias[alias] = locale
}

func SetDefaultTimeFormat(timeFormat string, timeLocation *time.Location) {
    defaultFormat = timeFormat
    defaultLocale = timeLocation
}

// Create instance for register
type CustomTimeExtension struct {
    jsoniter.DummyExtension
}

// Handler override here time.Time, *time.Time
func (extension *CustomTimeExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
    for _, binding := range structDescriptor.Fields {
        var typeErr error
        var isPtr bool
        typeName := binding.Field.Type().String()
        if typeName == "time.Time" {
            isPtr = false
        } else if typeName == "*time.Time" {
            isPtr = true
        } else {
            continue
        }

        timeFormat := defaultFormat
        formatTag := binding.Field.Tag().Get(tagNameTimeFormat)
        if format, ok := formatAlias[formatTag]; ok {
            timeFormat = format
        } else if formatTag != "" {
            timeFormat = formatTag
        }
        locale := defaultLocale
        if localeTag := binding.Field.Tag().Get(tagNameTimeLocation); localeTag != "" {
            if loc, ok := localeAlias[localeTag]; ok {
                locale = loc
            } else {
                loc, err := time.LoadLocation(localeTag)
                if err != nil {
                    typeErr = err
                } else {
                    locale = loc
                }
            }
        }

        binding.Encoder = &funcEncoder{fun: func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
            if typeErr != nil {
                stream.Error = typeErr
                return
            }

            var tp *time.Time
            if isPtr {
                tpp := (**time.Time)(ptr)
                tp = *(tpp)
            } else {
                tp = (*time.Time)(ptr)
            }

            if tp != nil {
                lt := tp.In(locale)
                str := lt.Format(timeFormat)
                stream.WriteString(str)
            } else {
                stream.Write([]byte("null"))
            }
        }}

        binding.Decoder = &funcDecoder{fun: func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
            if typeErr != nil {
                iter.Error = typeErr
                return
            }

            str := iter.ReadString()
            var t *time.Time
            if str != "" {
                var err error
                tmp, err := time.ParseInLocation(timeFormat, str, locale)
                if err != nil {
                    iter.Error = err
                    return
                }
                t = &tmp
            } else {
                t = nil
            }

            if isPtr {
                tpp := (**time.Time)(ptr)
                *tpp = t
            } else {
                tp := (*time.Time)(ptr)
                if tp != nil && t != nil {
                    *tp = *t
                }
            }
        }}
    }
}

type funcDecoder struct {
    fun jsoniter.DecoderFunc
}

func (decoder *funcDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
    decoder.fun(ptr, iter)
}

type funcEncoder struct {
    fun         jsoniter.EncoderFunc
    isEmptyFunc func(ptr unsafe.Pointer) bool
}

func (encoder *funcEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
    encoder.fun(ptr, stream)
}

func (encoder *funcEncoder) IsEmpty(ptr unsafe.Pointer) bool {
    if encoder.isEmptyFunc == nil {
        return false
    }
    return encoder.isEmptyFunc(ptr)
}
