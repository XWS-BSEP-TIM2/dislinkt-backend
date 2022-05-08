package base

type ErrorMapper interface {
	ToStatusError(domainErr error) (statusErr error)
}
