package main

import "net/url"

func controlFormInputs(u string) {
	isVul, p := cd.formTest(u, payloads)
	if isVul {
		if _, contains := vulnerableURLs[u]; !contains {
			vulnerableURLs[u] = empty{payload: p}
		}
	}
}

func controlQueryParameters(u string) {
	// convert string to url.URL
	originalURL, _ := url.Parse(u)
	modifiedURL := originalURL

	// testing for DOM XSS

	// browser does not allow , encodes url
	// TODO: find a way to bypass
	// chrome-headless disable-xss-auditor ?

	if modifiedURL.Fragment != "" {
		for _, payload := range payloads {

			modifiedURL.Fragment = payload

			if cd.pageTest(modifiedURL.String()) {
				if _, contains := vulnerableURLs[modifiedURL.String()]; !contains {
					vulnerableURLs[modifiedURL.String()] = empty{}
				}
				break
			}
		}
	}

	// get query parameters and values as map[string][]string
	q := originalURL.Query()

	// checking query parameters Ex: username=xx&email=yy
	// only one parameter is changed at once
	if len(q) > 0 {
		for parameter := range q {
			for _, payload := range payloads {

				modifiedURL = originalURL
				q.Set(parameter, payload)
				modifiedURL.RawQuery = q.Encode()

				if cd.pageTest(modifiedURL.String()) {
					if _, contains := vulnerableURLs[modifiedURL.String()]; !contains {
						vulnerableURLs[modifiedURL.String()] = empty{}
					}
					break
				}
			}
		}
	}
}

