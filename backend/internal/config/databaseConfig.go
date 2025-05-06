package config

type DatabaseConfig struct {
	MigrateDatabse  bool
	Migrations      any
	ExpectedVersion int
	VersionTable    string
}
