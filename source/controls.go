package source

import "net/url"

func ControlFormInputs(u string) {
	isVul, p := CD.FormTest(u, Payloads)
	if isVul {
		if _, contains := VulnerableURLs[u]; !contains {
			VulnerableURLs[u] = Empty{Payload: p}
		}
	}
}

func ControlQueryParameters(u string) {
	// convert string to url.URL
	originalURL, _ := url.Parse(u)
	modifiedURL := originalURL

	// testing for DOM XSS

	// browser does not allow , encodes url
	// TODO: find a way to bypass
	// chrome-headless disable-xss-auditor ?

	if modifiedURL.Fragment != "" {
		for _, payload := range Payloads {

			modifiedURL.Fragment = payload

			if CD.PageTest(modifiedURL.String()) {
				if _, contains := VulnerableURLs[modifiedURL.String()]; !contains {
					VulnerableURLs[modifiedURL.String()] = Empty{}
				}
				break
			}
		}
	}

	// get query parameters and values as map[string][]string
	q := originalURL.Query()

	// checking query parameters Ex: username=xx&email=yy
	// only one parameter is changed at once # fix ?
	if len(q) > 0 {
		for _, payload := range Payloads {
			for parameter := range q {


				modifiedURL = originalURL
				q.Set(parameter, payload)

			}
			modifiedURL.RawQuery = q.Encode()

			if CD.PageTest(modifiedURL.String()) {
				if _, contains := VulnerableURLs[modifiedURL.String()]; !contains {
					VulnerableURLs[modifiedURL.String()] = Empty{}
				}
				break
			}
		}
	}
}

