package cmf

/*
CheckErr is used to stop execution if an error occurs, without having to put
`if err != nil` everywhere. I mostly use it for quick scripts.
*/
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

/*
Must is used to turn a value and an error into a value. Panics on error similarly
to regexp.MustCompile. I mostly use it for quick scripts, or things that need to happen
in `init` functions when program execution should be terminated on failure
*/
func Must[T any](t T, err error) T {
	CheckErr(err)
	return t
}
