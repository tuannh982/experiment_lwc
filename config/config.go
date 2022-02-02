package config

type Configuration struct {
	ID        string             `yaml:"id"`
	RootFs    string             `yaml:"root-fs"`
	Chdir     string             `yaml:"chdir"`
	UID       int                `yaml:"uid"`
	GID       int                `yaml:"gid"`
	Hostname  string             `yaml:"hostname"`
	Resources Resources          `yaml:"resources"`
	Mounts    []Mount            `yaml:"mounts"`
	Networks  map[string]Network `yaml:"networks"`
}

type Resources struct {
	Cgroup string  `yaml:"cgroup"`
	CPU    float64 `yaml:"cpu"`
	Memory int     `yaml:"memory"`
}

type Mount struct {
	Source string `yaml:"source"`
	Target string `yaml:"Target"`
	Fs     string `yaml:"Fs"`
	Flags  int    `yaml:"flags"`
	Data   string `yaml:"data"`
}

type Network struct {
	Address string `yaml:"address"`
	CIDR    string `yaml:"cidr"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Mounts:   make([]Mount, 0),
		Networks: make(map[string]Network),
	}
}
