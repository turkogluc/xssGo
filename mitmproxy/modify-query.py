from mitmproxy import http
import urllib.parse as urlparse
from urllib.parse import urlencode

def request(flow: http.HTTPFlow) -> None:
    parsed = urlparse.urlparse(flow.request.url)
    parsedList = list(parsed)
    print(parsedList)
    query = dict(urlparse.parse_qsl(parsedList[4]))
    for k,v in query.items():
        flow.request.query[k] = "<script>alert(999);</script>"
        print(flow.request.query)
    
