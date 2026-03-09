package generate

import (
	"context"
	"fmt"
	"log/slog"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

func makeColumnReadOnly(tag field.GormTag) field.GormTag {
	slog.Info("Making column readonly")
	return tag.Append("->")
}

var globalGenColumnOptions = []gen.ModelOpt{
	gen.FieldJSONTagWithNS(addOmitEmpty),
	gen.FieldGORMTagReg(".*_at$", makeColumnReadOnly),
}

func ORMCode(ctx context.Context, generator *gen.Generator) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			return
		}
	}()

	generator.WithImportPkgPath("github.com/jghiloni/coredns-pg/common/resolve/records")
	generator.ApplyInterface(func(DNSZoneQueries) {}, generator.GenerateModel("zones", globalGenColumnOptions...))

	generator.ApplyInterface(func(DNSRecordQueries) {},
		generator.GenerateModel("records",
			append(globalGenColumnOptions,
				gen.FieldGORMTag("id", makeColumnReadOnly),
				gen.FieldType("content", "records.DNSRecordContent"),
				gen.FieldType("record_type", "records.RecordType"),
				gen.FieldType("ttl", "uint32"),
				gen.FieldGORMTag("content", func(tag field.GormTag) field.GormTag {
					return tag.Set("type", "jsonb")
				}),
				gen.FieldGORMTag("record_type", func(tag field.GormTag) field.GormTag {
					return tag.Set("type", "record_type")
				}),
			)...,
		),
	)

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Execute doesn't return an error, it panics, hence the defer
	generator.Execute()

	return nil
}

func addOmitEmpty(field string) string {
	return fmt.Sprintf("%s,omitempty", field)
}
