package backend

import "github.com/webx-top/echo/subdomains"

func URLFor(purl string, relative ...bool) string {
	if len(relative) > 0 && relative[0] {
		return subdomains.Default.RelativeURL(purl, `backend`)
	}
	return subdomains.Default.URL(purl, `backend`)
}

func FrontendURLFor(purl string, relative ...bool) string {
	if len(relative) > 0 && relative[0] {
		return subdomains.Default.RelativeURL(purl, `frontend`)
	}
	return subdomains.Default.URL(purl, `frontend`)
}
