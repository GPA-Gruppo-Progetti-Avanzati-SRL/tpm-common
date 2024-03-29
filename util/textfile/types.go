package textfile

type CSVFieldInfo struct {
	Id         string `yaml:"id,omitempty" mapstructure:"id,omitempty" json:"id,omitempty"`
	Name       string `yaml:"name,omitempty" mapstructure:"name,omitempty" json:"name,omitempty"`
	Validation string `yaml:"validation,omitempty" mapstructure:"validation,omitempty" json:"validation,omitempty"`
	Help       string `yaml:"help,omitempty" mapstructure:"help,omitempty" json:"help,omitempty"`
	Index      int    `yaml:"index,omitempty" mapstructure:"index,omitempty" json:"index,omitempty"`
}
