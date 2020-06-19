package strfmt

import (
	"bytes"
	"fmt"
	"net/textproto"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/go-courier/reflectx"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Resource(name = cpu, targetAverageUtilization = 70)
func ParseMetrics(metricsStr string) (Metrics, error) {
	if metricsStr == "" {
		return nil, nil
	}

	s := &scanner.Scanner{}
	s.Init(bytes.NewReader([]byte(metricsStr)))

	ms := make([]metric, 0)

	m := metric{values: map[string]string{}}

	lastKey := ""
	lastTextToken := ""

	setValue := func() {
		if v, err := strconv.Unquote(lastTextToken); err != nil {
			m.values[lastKey] = lastTextToken
		} else {
			m.values[lastKey] = v
		}
		lastTextToken = ""
	}

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		switch tok {
		case '(':
			m.tpe = autoscalingv2beta1.MetricSourceType(textproto.CanonicalMIMEHeaderKey(lastTextToken))
			lastTextToken = ""
		case ')':
			setValue()

			ms = append(ms, m)
			m = metric{values: map[string]string{}}
		case '=':
			lastKey = lastTextToken
			lastTextToken = ""
		case ',':
			setValue()
		default:
			lastTextToken += s.TokenText()
		}
	}

	metrics := Metrics{}

	for i := range ms {
		m := ms[i]

		metric := autoscalingv2beta1.MetricSpec{Type: m.tpe}

		switch m.tpe {
		case autoscalingv2beta1.ResourceMetricSourceType:
			metric.Resource = &autoscalingv2beta1.ResourceMetricSource{}
			if err := unmarshalValues(m.values, metric.Resource); err != nil {
				return nil, err
			}
		case autoscalingv2beta1.PodsMetricSourceType:
			metric.Pods = &autoscalingv2beta1.PodsMetricSource{}
			if err := unmarshalValues(m.values, metric.Object); err != nil {
				return nil, err
			}
		case autoscalingv2beta1.ObjectMetricSourceType:
			metric.Object = &autoscalingv2beta1.ObjectMetricSource{}
			if err := unmarshalValues(m.values, metric.Object); err != nil {
				return nil, err
			}
		case autoscalingv2beta1.ExternalMetricSourceType:
			metric.External = &autoscalingv2beta1.ExternalMetricSource{}
			if err := unmarshalValues(m.values, metric.Object); err != nil {
				return nil, err
			}
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

type metric struct {
	tpe    autoscalingv2beta1.MetricSourceType
	values map[string]string
}

type Metrics []autoscalingv2beta1.MetricSpec

func (Metrics) OpenAPISchemaType() []string { return []string{"string"} }
func (Metrics) OpenAPISchemaFormat() string { return "metrics" }

func (metrics Metrics) String() string {
	b := bytes.NewBuffer(nil)

	writeValues := func(values map[string]string) {
		keys := make([]string, 0)

		for k := range values {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for i, key := range keys {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(key)
			b.WriteString(" = ")
			b.WriteString(strconv.Quote(values[key]))
		}
	}

	for i := range metrics {
		if i != 0 {
			b.WriteRune(' ')
		}

		m := metrics[i]

		b.WriteString(string(m.Type))
		b.WriteRune('(')

		switch m.Type {
		case autoscalingv2beta1.ResourceMetricSourceType:
			values, err := marshalValues(m.Resource)
			if err == nil {
				writeValues(values)
			}
		case autoscalingv2beta1.PodsMetricSourceType:
			values, err := marshalValues(m.Pods)
			if err == nil {
				writeValues(values)
			}
		case autoscalingv2beta1.ObjectMetricSourceType:
			values, err := marshalValues(m.Object)
			if err == nil {
				writeValues(values)
			}
		case autoscalingv2beta1.ExternalMetricSourceType:
			values, err := marshalValues(m.External)
			if err == nil {
				writeValues(values)
			}
		}

		b.WriteRune(')')
	}

	return b.String()
}

func marshalValues(v interface{}) (map[string]string, error) {
	values := map[string]string{}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	st := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		fieldRv := rv.Field(i)
		fieldType := st.Field(i)

		v, ok := fieldType.Tag.Lookup("json")
		if ok {
			name := strings.Split(v, ",")[0]
			if name == "" {
				name = fieldType.Name
			}

			switch fieldValue := fieldRv.Interface().(type) {
			case autoscalingv2beta1.CrossVersionObjectReference:
				values[name] = fmt.Sprintf("%s.%s#%s", fieldValue.APIVersion, fieldValue.Kind, fieldValue.Name)

			case *metav1.LabelSelector:
				if fieldValue != nil {
					values[name] = metav1.FormatLabelSelector(fieldValue)
				}
			case *resource.Quantity:
				if fieldValue != nil {
					values[name] = fieldValue.String()
				}
			case resource.Quantity:
				values[name] = fieldValue.String()
			default:
				b, err := reflectx.MarshalText(fieldRv)
				if err != nil {
					return nil, err
				}
				if len(b) > 0 {
					values[name] = string(b)
				}
			}
		}
	}

	return values, nil
}

func unmarshalValues(values map[string]string, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if !rv.CanAddr() {
		return fmt.Errorf("invalid %v", v)
	}

	st := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		fieldRv := rv.Field(i)
		fieldType := st.Field(i)

		v, ok := fieldType.Tag.Lookup("json")
		if ok {
			name := strings.Split(v, ",")[0]
			if name == "" {
				name = fieldType.Name
			}

			if v, ok := values[name]; ok {
				switch fieldValue := fieldRv.Interface().(type) {
				case autoscalingv2beta1.CrossVersionObjectReference:
					parts := strings.Split(v, "#")
					if len(parts) == 2 {
						fieldValue.Name = parts[1]

						parts = strings.Split(parts[0], ".")

						if len(parts) == 2 {
							fieldValue.APIVersion = parts[0]
							fieldValue.Kind = parts[1]
						}
					}

					fieldRv.Set(reflect.ValueOf(fieldValue))
				case *metav1.LabelSelector:
					s, err := metav1.ParseToLabelSelector(v)
					if err != nil {
						return err
					}
					fieldValue = s

					fieldRv.Set(reflect.ValueOf(fieldValue))
				case *resource.Quantity:
					s, err := resource.ParseQuantity(v)
					if err != nil {
						return err
					}
					fieldValue = &s

					fieldRv.Set(reflect.ValueOf(fieldValue))
				case resource.Quantity:
					s, err := resource.ParseQuantity(v)
					if err != nil {
						return err
					}
					fieldValue = s
					fieldRv.Set(reflect.ValueOf(fieldValue))
				default:
					if err := reflectx.UnmarshalText(fieldRv, []byte(v)); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
