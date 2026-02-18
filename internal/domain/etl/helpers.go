package etl

func MapAirbyteStreamToSourceSchema(s AirbyteSourceStream) SourceSchema {
	return SourceSchema{
		Name:       s.StreamName,
		Namespace:  s.StreamNamespace,
		PrimaryKeys: s.SourceDefinedPrimaryKey,
		Fields:     s.PropertyFields,
	}
}
