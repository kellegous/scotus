package async

type Future[T any] struct {
	empty T
	rch   chan T
	ech   chan error
}

func newFuture[T any](
	empty T,
) *Future[T] {
	return &Future[T]{
		empty: empty,
		rch:   make(chan T, 1),
		ech:   make(chan error, 1),
	}
}

func (f *Future[T]) fulfill(r T, err error) {
	if err != nil {
		f.ech <- err
	} else {
		f.rch <- r
	}
}

func (f *Future[T]) Resolve() (T, error) {
	select {
	case r := <-f.rch:
		return r, nil
	case err := <-f.ech:
		return f.empty, err
	}
}
