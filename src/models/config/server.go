package config

type ServerConfig struct {
	Database   MySQL   `yaml:"database" mapstructure:"database"`
	Zap        Zap     `yaml:"zap" mapstructure:"zap"`
	WFGroups   []int64 `yaml:"wf_groups" mapstructure:"wf_groups"`
	XDGroups   []int64 `yaml:"xd_groups" mapstructure:"xd_groups"`
	CQHttpAddr string  `yaml:"cq_http_addr" mapstructure:"cq_http_addr"`
}

type MySQL struct {
	Path         string `mapstructure:"path" json:"path" yaml:"path"`
	Config       string `mapstructure:"config_models" json:"config-models" yaml:"config_models"`
	Dbname       string `mapstructure:"db_name" json:"dbname" yaml:"db_name"`
	Username     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	MaxIdleConns int    `mapstructure:"max_idle_conns" json:"max-idle-conns" yaml:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns" json:"max-open-conns" yaml:"max_open_conns"`
	LogMode      bool   `mapstructure:"log_mode" json:"log-mode" yaml:"log_mode"`
	LogZap       string `mapstructure:"log_zap" json:"log-zap" yaml:"log_zap"`
}

func (m *MySQL) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
}

type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`
	Format        string `mapstructure:"format" json:"format" yaml:"format"`
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`
	LinkName      string `mapstructure:"link_name" json:"link-name" yaml:"link_name"`
	ShowLine      bool   `mapstructure:"show-line" json:"show-line" yaml:"show_line"`
	EncodeLevel   string `mapstructure:"encode-level" json:"encode-level" yaml:"encode_level"`
	StacktraceKey string `mapstructure:"stacktrace_key" json:"stacktrace-key" yaml:"stacktrace_key"`
	LogInConsole  bool   `mapstructure:"log_in_console" json:"log-in-console" yaml:"log_in_console"`
}
