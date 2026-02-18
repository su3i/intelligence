package etl

type EtlType string

const (
	EtlTypeRedis   EtlType = "airbyte"
)

// Minimal ETL operations
type ETL interface {
	CreateSourceConnection(name string, configuration map[string]interface{}) (*string, error)
	DeleteSourceConnection(sourceId string) error
	TestSourceConnection(sourceId string) error
	RetrieveSourceSchemas(sourceId string) ([]SourceSchema, error)
}

type AirbyteSourceStream struct {
	StreamName              string          `json:"streamName"`
	StreamNamespace         string          `json:"streamnamespace"`
	DefaultCursorField      []string        `json:"defaultCursorField"`
	SourceDefinedCursorField bool           `json:"sourceDefinedCursorField"`
	SourceDefinedPrimaryKey [][]string      `json:"sourceDefinedPrimaryKey"`
	PropertyFields          [][]string      `json:"propertyFields"`
}

type AirbyteSourceStreamsResponse struct {
	Streams []AirbyteSourceStream `json:"streams"`
}

type SourceSchema struct {
	Name              string          `json:"name"`
	Namespace         string          `json:"namespace"`
	PrimaryKeys 		[][]string      `json:"primaryKeys"`
	Fields          [][]string      `json:"fields"`
}
