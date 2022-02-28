package chrome

import (
	"context"
	"strings"

	"github.com/chromedp/chromedp"
)

var defaultOpts = []chromedp.ExecAllocatorOption{
	chromedp.NoFirstRun,
	chromedp.NoDefaultBrowserCheck,
	chromedp.IgnoreCertErrors,
	chromedp.Flag("test-type ", true), // This prevents the `ignore-certificate-errors` warning.
}

type Config struct {
	ProxyServer      string
	ProxyBypassHosts []string
}

// NewExecAllocator returns a new context setup with a chromedp.ExecAllocator.
// Its `context.Context` return value can be used to create subsequent contexts for interacting
// with an allocated Chrome browser.
func NewExecAllocator(ctx context.Context, cfg Config) (context.Context, context.CancelFunc) {
	proxyBypass := strings.Join(append([]string{"<-loopback"}, cfg.ProxyBypassHosts...), ";")
	//nolint:gocritic
	opts := append(defaultOpts,
		chromedp.ProxyServer(cfg.ProxyServer),
		chromedp.Flag("proxy-bypass-list", proxyBypass),
	)

	return chromedp.NewExecAllocator(ctx, opts...)
}
