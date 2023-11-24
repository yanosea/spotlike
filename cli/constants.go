package cli

type root struct {
	// rootName is the name of root command
	rootName string
	// rootShort is the short description of root command
	rootShort string
	// rootLong is the long description of root command
	rootLong string
}

// variables as constants
var (
	// ROOT is the constants of the root command
	ROOT = root{
		rootName:  "spotlike",
		rootShort: "'spotlike' is the CLI tool to LIKE contents in Spotify.",
		rootLong: `'spotlike' is the CLI tool to LIKE contents in Spotify.

You can get the ID of some contents in Spotify.
You can LIKE the contents in Spotify by ID.`,
	}
)

// RootName returns the name of the root command
func (v *root) RootName() string {
	return v.rootName
}

// RootShort returns the short description of the root command
func (v *root) RootShort() string {
	return v.rootShort
}

// RootLong returns the long description of the root command
func (v *root) RootLong() string {
	return v.rootLong
}
