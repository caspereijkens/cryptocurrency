package transaction

type Transaction struct {
	Version  string
	Inputs   []string
	Outputs  []string
	Locktime string
}
