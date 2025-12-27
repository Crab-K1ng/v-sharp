package cli

func Run(args []string) {
	if len(args) < 2 {
		return
	}
	run(args[1])
}
