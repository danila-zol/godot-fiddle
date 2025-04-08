package config

type DatabaseConfig struct {
	MigrateDatabse  bool
	Migrations      any
	MigrationsDir   string // TODO: Relative to where?
	ExpectedVersion int
	VersionTable    string
}
