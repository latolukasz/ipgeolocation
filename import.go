package ipgeolocation

import (
	"context"

	"github.com/latolukasz/beeorm"
)

func Import(ctx context.Context) error {

	var countryEntity *CountryEntity

	registry := beeorm.NewRegistry()
	registry.RegisterMySQLPool("root:root@tcp(localhost:3315)/ipgeolocation")
	registry.RegisterRedis("localhost:6383", 0)
	registry.RegisterEntity(countryEntity)

	validatedRegistry, err := registry.Validate(ctx)
	if err != nil {
		return err
	}
	engine := validatedRegistry.CreateEngine(ctx)

	engine.GetRegistry().GetTableSchemaForEntity(countryEntity).UpdateSchemaAndTruncateTable(engine)
	return nil
}
