package mime

const charsetUTF8 = "charset=utf-8"
const (
	// ApplicationJSON MIME type
	ApplicationJSON = `application/json`

	// ApplicationJSONCharsetUTF8 MIME type
	ApplicationJSONCharsetUTF8 = ApplicationJSON + `; ` + charsetUTF8

	// ApplicationJavaScript MIME type
	ApplicationJavaScript = `application/javascript`

	// ApplicationJavaScriptCharsetUTF8 MIME type
	ApplicationJavaScriptCharsetUTF8 = ApplicationJavaScript + `; ` + charsetUTF8

	// ApplicationXML MIME type
	ApplicationXML = `application/xml`

	// ApplicationXMLCharsetUTF8 MIME type
	ApplicationXMLCharsetUTF8 = ApplicationXML + `; ` + charsetUTF8

	// ApplicationForm MIME type
	ApplicationForm = `application/x-www-form-urlencoded`

	// ApplicationProtobuf MIME type
	ApplicationProtobuf = `application/protobuf`

	// ApplicationMsgpack MIME type
	ApplicationMsgpack = `application/msgpack`

	// TextHTML MIME type
	TextHTML = `text/html`

	// TextHTMLCharsetUTF8 MIME type
	TextHTMLCharsetUTF8 = TextHTML + `; ` + charsetUTF8

	// TextCSS MIME type
	TextCSS = `text/css`

	// TextCSSCharsetUTF8 MIME type
	TextCSSCharsetUTF8 = TextCSS + `; ` + charsetUTF8

	// TextPlain MIME type
	TextPlain = `text/plain`

	// TextPlainCharsetUTF8 MIME type
	TextPlainCharsetUTF8 = TextPlain + `; ` + charsetUTF8

	// TextXML MIME type
	TextXML = `text/xml`

	// TextXMLCharsetUTF8 MIME type
	TextXMLCharsetUTF8 = TextXML + `; ` + charsetUTF8

	// TextJavascript MIME type
	TextJavascript = `text/javascript`

	// MultipartForm MIME type
	MultipartForm = `multipart/form-data`

	// OctetStream MIME type
	OctetStream = `application/octet-stream`

	// FavIcon MIME type
	FavIcon = `image/vnd.microsoft.icon`

	// ImagePNG MIME type
	ImagePNG = `image/png`

	// ImageSVG MIME type
	ImageSVG = `image/svg+xml`

	// ImageICO Official registered MIME type of ico files
	ImageICO = `image/vnd.microsoft.icon`

	// ImageXICON MIME type of ico files, which is used by Microsoft by stupidity
	ImageXICON = `image/x-icon`

	// AutoReplied auto-replied header of RFC3464
	AutoReplied = `auto-replied`

	// MessageDeliveryStatus MIME type of RFC3464
	MessageDeliveryStatus = `message/delivery-status`

	// MessageRfc822 MIME type of RFC822
	MessageRfc822 = `message/rfc822`

	// TextRfc822Headers MIME type of RFC822
	TextRfc822Headers = `text/rfc822-headers`
)
