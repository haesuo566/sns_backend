package impls

type Topic interface {
	ExecuteEvent(string, string, []byte) error
}
