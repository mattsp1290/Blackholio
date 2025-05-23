package tables

import (
	"github.com/clockworklabs/SpacetimeDB/crates/bindings-go/pkg/spacetimedb/schema"
)

// Schema Migration Example: Converting Blackholio tables to use the new schema framework
// This shows how the manually defined table metadata can be replaced with the schema framework

// Example: Converting the existing TableDefinitions map to use schema.TableInfo

// CreateBlackholioSchema creates all table definitions using the new schema framework
func CreateBlackholioSchema() []*schema.TableInfo {
	tables := []*schema.TableInfo{}

	// Config table
	configTable := schema.NewTableInfo("config")
	configTable.PublicRead = true
	configTable.Columns = []schema.Column{
		schema.NewPrimaryKeyColumn("id", schema.TypeU32),
		schema.NewColumn("world_size", schema.TypeU64),
	}
	tables = append(tables, configTable)

	// Entity table
	entityTable := schema.NewTableInfo("entity")
	entityTable.PublicRead = true
	entityTable.Columns = []schema.Column{
		schema.NewAutoIncColumn("entity_id", schema.TypeU32),
		schema.NewColumn("position", "DbVector2"), // Custom type
		schema.NewColumn("mass", schema.TypeU32),
	}
	tables = append(tables, entityTable)

	// Circle table with index
	circleTable := schema.NewTableInfo("circle")
	circleTable.PublicRead = true
	circleTable.Columns = []schema.Column{
		schema.NewPrimaryKeyColumn("entity_id", schema.TypeU32),
		schema.NewColumn("player_id", schema.TypeU32),
		schema.NewColumn("direction", "DbVector2"), // Custom type
		schema.NewColumn("speed", schema.TypeF32),
		schema.NewColumn("last_split_time", schema.TypeTimestamp),
	}
	circleTable.Indexes = []schema.Index{
		schema.NewBTreeIndex("idx_player_id", []string{"player_id"}),
	}
	tables = append(tables, circleTable)

	// Player table
	playerTable := schema.NewTableInfo("player")
	playerTable.PublicRead = true
	playerTable.Columns = []schema.Column{
		schema.NewPrimaryKeyColumn("identity", schema.TypeIdentity),
		{
			Name:    "player_id",
			Type:    schema.TypeU32,
			AutoInc: true,
			Unique:  true,
			NotNull: true,
		},
		schema.NewColumn("name", schema.TypeString),
	}
	tables = append(tables, playerTable)

	// Logged out player table (private)
	loggedOutPlayerTable := schema.NewTableInfo("logged_out_player")
	loggedOutPlayerTable.PublicRead = false // Private table
	loggedOutPlayerTable.Columns = []schema.Column{
		schema.NewPrimaryKeyColumn("identity", schema.TypeIdentity),
		{
			Name:    "player_id",
			Type:    schema.TypeU32,
			AutoInc: true,
			Unique:  true,
			NotNull: true,
		},
		schema.NewColumn("name", schema.TypeString),
	}
	tables = append(tables, loggedOutPlayerTable)

	// Food table
	foodTable := schema.NewTableInfo("food")
	foodTable.PublicRead = true
	foodTable.Columns = []schema.Column{
		schema.NewPrimaryKeyColumn("entity_id", schema.TypeU32),
	}
	tables = append(tables, foodTable)

	// Timer tables - all have similar structure
	timerTables := []string{
		"move_all_players_timer",
		"spawn_food_timer",
		"circle_decay_timer",
	}

	for _, tableName := range timerTables {
		timerTable := schema.NewTableInfo(tableName)
		timerTable.Columns = []schema.Column{
			schema.NewAutoIncColumn("scheduled_id", schema.TypeU64),
			schema.NewColumn("scheduled_at", schema.TypeScheduleAt),
		}
		tables = append(tables, timerTable)
	}

	// Circle recombine timer (has extra player_id field)
	circleRecombineTable := schema.NewTableInfo("circle_recombine_timer")
	circleRecombineTable.Columns = []schema.Column{
		schema.NewAutoIncColumn("scheduled_id", schema.TypeU64),
		schema.NewColumn("scheduled_at", schema.TypeScheduleAt),
		schema.NewColumn("player_id", schema.TypeU32),
	}
	tables = append(tables, circleRecombineTable)

	// Consume entity timer (has extra entity ID fields)
	consumeEntityTable := schema.NewTableInfo("consume_entity_timer")
	consumeEntityTable.Columns = []schema.Column{
		schema.NewAutoIncColumn("scheduled_id", schema.TypeU64),
		schema.NewColumn("scheduled_at", schema.TypeScheduleAt),
		schema.NewColumn("consumed_entity_id", schema.TypeU32),
		schema.NewColumn("consumer_entity_id", schema.TypeU32),
	}
	tables = append(tables, consumeEntityTable)

	return tables
}

// RegisterBlackholioSchema registers all Blackholio tables with the global schema registry
func RegisterBlackholioSchema() error {
	tables := CreateBlackholioSchema()

	// Register all tables with validation enabled
	options := schema.RegistrationOptions{
		ValidateSchema: true,
		AllowOverwrite: false,
		AssignIDs:      true,
	}

	return schema.GlobalRegisterAll(tables, options)
}

// GetBlackholioTableStats returns statistics about the registered Blackholio tables
func GetBlackholioTableStats() schema.RegistryStats {
	return schema.GlobalGetStats()
}

// ValidateBlackholioSchema validates all registered Blackholio table schemas
func ValidateBlackholioSchema() error {
	return schema.GlobalValidateAll()
}

// Example usage functions showing the benefits of the schema framework

// FindTableByName demonstrates safe table lookup
func FindTableByName(name string) (*schema.TableInfo, bool) {
	return schema.GlobalGetTable(name)
}

// GetAllTableNames returns all registered table names in sorted order
func GetAllTableNames() []string {
	return schema.GlobalGetTableNames()
}

// PrintTableInfo prints detailed information about a table
func PrintTableInfo(tableName string) {
	table, exists := schema.GlobalGetTable(tableName)
	if !exists {
		println("Table not found:", tableName)
		return
	}

	println("Table:", table.String())
	println("  Columns:")
	for i, col := range table.Columns {
		println("   ", i+1, ":", col.String())
	}

	if len(table.Indexes) > 0 {
		println("  Indexes:")
		for i, idx := range table.Indexes {
			println("   ", i+1, ":", idx.String())
		}
	}
}

// Benefits of using the schema framework:
//
// 1. **Type Safety**: Strongly typed column definitions with validation
// 2. **Consistency**: All tables use the same schema structure
// 3. **Validation**: Built-in validation for schema correctness
// 4. **Thread Safety**: Safe concurrent access to table definitions
// 5. **Introspection**: Easy lookup and enumeration of tables/columns
// 6. **Future Extensibility**: Framework can be extended with new features
// 7. **JSON Support**: Built-in JSON serialization for schema inspection
// 8. **Performance**: Optimized lookups and minimal memory overhead
//
// Migration benefits:
// - Replaces manual TableDefinitions map with validated schema registry
// - Eliminates custom Column/Index structs in favor of framework types
// - Provides better error handling and validation
// - Enables schema introspection and documentation generation
// - Prepares for future SpacetimeDB Go tooling integration
