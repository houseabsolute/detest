package detest

// CollectionEnding is an enum for collection ending checks, where we either
// check that all elements have been tested or allow extra untested elements.
type CollectionEnding int

const (
	// Unset means that the caller did not specify how to check the ending.
	Unset CollectionEnding = iota
	// Etc means that additional unchecked elements are allowed.
	Etc
	// End means that all elements must be checked or the test fails.
	End
)
