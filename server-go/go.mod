module github.com/clockworklabs/Blackholio/server-go

go 1.21

require github.com/clockworklabs/SpacetimeDB/crates/bindings-go v0.2.0

// Replace directive to use the local fork
replace github.com/clockworklabs/SpacetimeDB/crates/bindings-go => ../../SpacetimeDB/crates/bindings-go
