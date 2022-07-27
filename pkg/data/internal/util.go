package internal

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

func EnsureDownload(
	ctx context.Context,
	client *http.Client,
	url string,
	dst string,
) error {
	if _, err := os.Stat(dst); err == nil {
		return nil
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if s := res.StatusCode; s != http.StatusOK {
		return fmt.Errorf("http status %d", s)
	}

	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()

	if _, err := io.Copy(w, res.Body); err != nil {
		return err
	}

	return nil
}
