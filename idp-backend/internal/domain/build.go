package domain

type FrameworkType string

const (
	FrameworkStatic FrameworkType = "STATIC"
	FrameworkNextSSR FrameworkType = "NEXT_SSR"
	FrameworkUnknown FrameworkType = "UNKNOWN"
)

type BuildResult struct {
	Framework FrameworkType
	ImageName string
	Port      int
}
