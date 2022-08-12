package async

func Run[T any](
	fn func() (T, error),
	empty T,
) *Future[T] {
	f := newFuture(empty)
	go func() {
		f.fulfill(fn())
	}()
	return f
}
