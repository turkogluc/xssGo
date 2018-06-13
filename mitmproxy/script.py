from mitmproxy import http

def request(flow: http.HTTPFlow) -> None:

    h = flow.request.headers
    print(h)
    flow.request.headers["Referer"] = "cemal"

    form = flow.request.urlencoded_form
    print(form)
    #print(form)
    #for k,v in form.items():
    #    print(k + " = " + v)
    #    flow.request.urlencoded_form[k] = "haydaa"
    
