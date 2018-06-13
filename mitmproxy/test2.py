from mitmproxy import http
import urllib.parse as urlparse
from urllib.parse import urlencode

def request(flow: http.HTTPFlow) -> None:
	path = flow.request.path
	# if we have sentinel in the url
	if '$' in flow.request.path:
		arr = path.split('$')
		payload = arr[1]

		# set the path as normal
		flow.request.path = arr[0]
		
		# set headers
		flow.request.headers["Referer"] = payload
		flow.request.headers["User-Agent"] = payload

		print(flow.request.url)
