package restclient

type HarBuilderOption func(hrab *HarBuilder)

type HarBuilder struct {
	creator        string
	creatorVersion string
	comment        string
	entries        []*Entry
}

func WithHarEntry(e *Entry) HarBuilderOption {
	return func(hrab *HarBuilder) {
		hrab.entries = append(hrab.entries, e)
	}
}

func WithHarCreator(creator string, version string) HarBuilderOption {
	return func(hrab *HarBuilder) {
		hrab.creator = creator
		hrab.creatorVersion = version
	}
}

func WithHarCcomment(comment string) HarBuilderOption {
	return func(hrab *HarBuilder) {
		hrab.comment = comment
	}
}

func NewHAR(opts ...HarBuilderOption) *HAR {

	harb := HarBuilder{creator: "rest-client", creatorVersion: "1.0"}
	for _, o := range opts {
		o(&harb)
	}

	har := HAR{
		Log: &Log{
			Version: "1.1",
			Creator: &Creator{
				Name:    harb.creator,
				Version: harb.creatorVersion,
			},
			Comment: harb.comment,
		},
	}

	har.Log.Entries = append(har.Log.Entries, harb.entries...)
	return &har
}
