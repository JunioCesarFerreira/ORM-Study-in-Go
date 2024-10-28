package columnfieldmap

type ColumnFieldPair struct {
	ColumnName string
	Field      interface{}
}

type Mapped interface {
	Mapped() []ColumnFieldPair
}

func ColumnsNames(entity Mapped) []string {
	mapped := entity.Mapped()
	columnNames := make([]string, 0, len(mapped))

	for _, pair := range mapped {
		columnNames = append(columnNames, pair.ColumnName)
	}

	return columnNames
}

func Fields(entity Mapped) []interface{} {
	mapped := entity.Mapped()
	fields := make([]interface{}, 0, len(mapped))

	for _, pair := range mapped {
		fields = append(fields, pair.Field)
	}

	return fields
}
