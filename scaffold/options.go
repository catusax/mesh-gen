package scaffold

// Options represents the options for the generator.
type Options struct {
	// Service is the name of the service the generator will generate files
	// for.
	Service string
	// Vendor is the service vendor.
	Vendor string
	// Directory is the directory where the files will be generated to.
	Directory string

	//ContainerTag is the container's tag
	ContainerTag string

	// Version is the version of the container
	Version string

	// Port of your service
	Port string

	// Namespace in your k8s cluster
	Namespace string

	//RegistryPrefix eg: gcr.io/username/repo
	RegistryPrefix string

	//Mesh service mesh name eg:istio
	Mesh string

	//Replica number of ReplicaSet
	Replica int
}

// Option manipulates the Options passed.
type Option func(o *Options)

// Service sets the service name.
func Service(s string) Option {
	return func(o *Options) {
		o.Service = s
	}
}

// Vendor sets the service vendor.
func Vendor(v string) Option {
	return func(o *Options) {
		o.Vendor = v
	}
}

// Directory sets the directory in which files are generated.
func Directory(d string) Option {
	return func(o *Options) {
		o.Directory = d
	}
}

func ContainerTag(r string) Option {
	return func(o *Options) {
		o.ContainerTag = r
	}
}

func ContainerVersion(r string) Option {
	return func(o *Options) {
		o.Version = r
	}
}

func Port(r string) Option {
	return func(o *Options) {
		o.Port = r
	}
}

func Namespace(r string) Option {
	return func(o *Options) {
		o.Namespace = r
	}
}

func RegistryPrefix(r string) Option {
	return func(o *Options) {
		o.RegistryPrefix = r
	}
}

func Mesh(r string) Option {
	return func(o *Options) {
		o.Mesh = r
	}
}

func Replica(r int) Option {
	return func(o *Options) {
		o.Replica = r
	}
}
