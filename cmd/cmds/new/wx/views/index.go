package views

var IndexViewTmpl = `<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="utf-8">
    <base href="/">
    <meta name="viewport" content="width=device-width, initial-scale=1, minimum-scale=1.0, maximum-scale=1.0, user-scalable=0">
    <link rel="dns-prefetch" href="//cdn.static.hztx18.com">
</head>
	hello {@pShortName}
      <input type="hidden" id="{{.xsrf_key}}" value="{{.xsrf_token}}">
</body>`
