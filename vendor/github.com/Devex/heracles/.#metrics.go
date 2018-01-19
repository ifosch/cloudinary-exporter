package heracles

import (
	"fmt"
	"log"
	"reflect"
)

const tagName = "metric"

func PrintMetric(metric interface{}, ns string) {
	v := reflect.ValueOf(metric)
	p := reflect.Indirect(v)
	t := p.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)
		fieldValue, _ := t.FieldByName(field.Name)
		if fieldValue.Type.Kind() == reflect.Struct {
			log.Printf("%v is struct", field.Name)
			PrintMetric(field, fmt.Sprintf("%s_%s", ns, tag))
		} else {
			if tag != "" {
				log.Printf("%s_%s %s", ns, tag, field.Name)
			}
		}
	}
}
