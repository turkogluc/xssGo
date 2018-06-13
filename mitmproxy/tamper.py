from mitmproxy import http
import urllib.parse as urlparse
from urllib.parse import urlencode

def request(flow: http.HTTPFlow) -> None:
    parsed = urlparse.urlparse(flow.request.url)
    parsedList = list(parsed)
    print(flow.request.query)
