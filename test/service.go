package test

type TestService struct {
}

type Args struct {
	First  int
	Second int
}

type Args2 struct {
	First1  string
	Second1 int
}

func (t TestService) Add(args Args) int {
	ret := args.First + args.Second
	return ret
}
