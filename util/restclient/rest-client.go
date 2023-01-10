package restclient

import (
	"crypto/tls"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/go-resty/resty/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type LinkedService struct {
	Cfg *Config
}

func NewInstanceWithConfig(cfg *Config) (*LinkedService, error) {
	lks := &LinkedService{Cfg: cfg}
	return lks, nil
}

func (lks LinkedService) NewClient(opts ...Option) (*Client, error) {
	cli := NewClient(lks.Cfg, opts...)
	return cli, nil
}

type Client struct {
	cfg Config

	restClient *resty.Client
	span       opentracing.Span
	spanOwned  bool
}

func NewClient(cfg *Config, opts ...Option) *Client {

	const semLogContext = "restclient.NewClient"

	var clientOptions Config
	if cfg == nil {
		clientOptions = Config{TraceRequestName: "rest-client"}
	} else {
		clientOptions = *cfg
	}

	for _, o := range opts {
		o(&clientOptions)
	}

	s := &Client{
		cfg:  clientOptions,
		span: clientOptions.Span,
	}

	if clientOptions.TraceGroupName != "" {
		s.span = startSpan(clientOptions.Span, clientOptions.TraceGroupName)
		s.spanOwned = true
	}

	s.restClient = resty.New()
	if s.cfg.RestTimeout != 0 {
		s.restClient.SetTimeout(s.cfg.RestTimeout)
		log.Trace().Dur("rest-timeout", s.cfg.RestTimeout).Msg(semLogContext)
	}

	if s.cfg.RetryCount != 0 {
		s.restClient.SetRetryCount(s.cfg.RetryCount)
		log.Trace().Int("rest-retry-count", s.cfg.RetryCount).Msg(semLogContext)
	}

	if s.cfg.RetryWaitTime != 0 {
		s.restClient.SetRetryWaitTime(s.cfg.RetryWaitTime)
		log.Trace().Dur("rest-wait-time", s.cfg.RetryWaitTime).Msg(semLogContext)
	}

	if s.cfg.RetryMaxWaitTime != 0 {
		s.restClient.SetRetryMaxWaitTime(s.cfg.RetryMaxWaitTime)
		log.Trace().Dur("rest-max-wait-time", s.cfg.RetryMaxWaitTime).Msg(semLogContext)
	}

	if len(s.cfg.RetryOnHttpError) > 0 {
		s.restClient.AddRetryCondition(retryCondition(s.cfg.RetryOnHttpError))
		log.Trace().Interface("rest-retry on error", s.cfg.RetryOnHttpError).Msg(semLogContext)
	}

	if s.cfg.SkipVerify {
		s.restClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return s
}

func retryCondition(errorsList []int) resty.RetryConditionFunc {
	return func(resp *resty.Response, err error) bool {

		const semLogContext = "restclient.NewClient"

		if len(errorsList) == 0 || err != nil {
			log.Trace().Err(err).Msg(semLogContext + " retry condition satisifed")
			return true
		}

		sc := resp.StatusCode()
		for i := 0; i < len(errorsList); i++ {
			if sc == errorsList[i] {
				log.Trace().Int("http-status", sc).Msg(semLogContext + " retry condition satisifed")
				return true
			}
		}

		log.Trace().Int("http-status", sc).Msg(semLogContext + " retry condition NOT satisifed")
		return false
	}
}

func (s *Client) Close() {
	if s.span != nil && s.spanOwned {
		s.span.Finish()
	}
}

func (s *Client) Execute(opName string, reqId string, lraId string, reqDef *Request) (*Entry, error) {

	now := time.Now()
	e := &Entry{
		Comment:         reqId,
		StartedDateTime: now.Format("2006-01-02T15:04:05.999999999Z07:00"),
		StartDateTimeTm: now,
		Request:         reqDef,
	}

	var reqSpanName string
	if s.cfg.TraceRequestName != "" {
		reqSpanName = strings.Replace(s.cfg.TraceRequestName, RequestTraceNameOpNamePlaceHolder, opName, 1)
		reqSpanName = strings.Replace(reqSpanName, RequestTraceNameRequestIdPlaceHolder, opName, 1)
	} else {
		reqSpanName = strings.Join([]string{opName, reqId}, "_")
	}
	reqSpan := startSpan(s.span, reqSpanName)
	defer reqSpan.Finish()

	// reqDef.Headers = append(reqDef.Headers, NameValuePair{Name: "Accept", Value: "application/json"})
	req := s.getRequestWithSpan(reqDef, reqSpan)

	var resp *resty.Response
	var err error

	u := reqDef.URL
	switch reqDef.Method {
	case http.MethodGet:
		resp, err = req.Get(u)
	case http.MethodHead:
		resp, err = req.Head(u)
	case http.MethodPost:
		resp, err = req.Post(u)
	case http.MethodPut:
		resp, err = req.Put(u)
	}

	s.setSpanTags(reqSpan, opName, reqId, lraId, u, reqDef.Method, resp.StatusCode(), err)

	var r *Response
	if err == nil {
		r = &Response{
			Status:      resp.StatusCode(),
			HTTPVersion: "1.1",
			StatusText:  resp.Status(),
			HeadersSize: -1,
			BodySize:    resp.Size(),
			Cookies:     []Cookie{},
			Content: &Content{
				MimeType: resp.Header().Get("Content-type"),
				Size:     resp.Size(),
				Data:     resp.Body(),
			},
		}

		for n, _ := range resp.Header() {
			r.Headers = append(r.Headers, NameValuePair{Name: n, Value: resp.Header().Get(n)})
		}
	} else {
		sc, st := DetectStatusCodeStatusTextFromError(resp.StatusCode(), err)
		err = util.NewError(strconv.Itoa(sc), err)
		r = NewResponse(sc, st, "text/plain", []byte(err.Error()), nil)
	}

	if e.StartedDateTime != "" {
		elapsed := time.Since(e.StartDateTimeTm)
		e.Time = float64(elapsed.Milliseconds())
	}

	e.Timings = &Timings{
		Blocked: -1,
		DNS:     -1,
		Connect: -1,
		Send:    -1,
		Wait:    e.Time,
		Receive: -1,
		Ssl:     -1,
	}

	e.Response = r
	return e, err
	// return resp.StatusCode(), resp.Body(), resp.Header(), err
}

func (s *Client) getRequestWithSpan(reqDef *Request, reqSpan opentracing.Span) *resty.Request {

	req := s.restClient.R()
	// Transmit the span's TraceContext as HTTP headers on our outbound request.
	_ = opentracing.GlobalTracer().Inject(reqSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))

	switch reqDef.Method {
	case http.MethodGet:
	case http.MethodHead:
	case http.MethodPost:
		if reqDef.HasBody() {
			req = req.SetBody(reqDef.PostData.Data)
		}

	case http.MethodPut:
		if reqDef.HasBody() {
			req = req.SetBody(reqDef.PostData.Data)
		}
	}

	// Setting first the default headers.
	for _, h := range s.cfg.Headers {
		req.SetHeader(h.Name, h.Value)
	}

	// Setting more specific headers next
	for _, h := range reqDef.Headers {
		req.SetHeader(h.Name, h.Value)
	}

	return req
}

/*
func (s *Client) getRequestSpan() opentracing.Span {

	var reqSpan opentracing.Span

	if s.span != nil {
		parentCtx := s.span.Context()
		reqSpan = opentracing.StartSpan(
			s.cfg.TraceOpName,
			opentracing.ChildOf(parentCtx),
		)
	} else {
		reqSpan = opentracing.StartSpan(
			s.cfg.TraceOpName,
		)
	}

	return reqSpan
}
*/

func startSpan(parentSpan opentracing.Span, spanName string) opentracing.Span {

	var span opentracing.Span

	if parentSpan != nil {
		parentCtx := parentSpan.Context()
		span = opentracing.StartSpan(
			spanName,
			opentracing.ChildOf(parentCtx),
		)
	} else {
		span = opentracing.StartSpan(
			spanName,
		)
	}

	return span
}

func (s *Client) setSpanTags(reqSpan opentracing.Span, opName, reqId, lraId, endpoint, method string, statusCode int, err error) {

	reqSpan.SetTag(util.HttpUrlTraceTag, endpoint)
	reqSpan.SetTag(util.HttpMethodTraceTag, method)
	reqSpan.SetTag(util.HttStatusCodeTraceTag, statusCode)

	if opName != "" {
		reqSpan.SetTag(OpNameTraceTag, opName)
	}

	if lraId != "" {
		reqSpan.SetTag(LraHttpContextTraceTag, lraId)
	}

	if reqId != "" {
		reqSpan.SetTag(RequestIdTraceTag, reqId)
	}

	if err != nil {
		reqSpan.SetTag("error", err.Error())
	}
}

func DetectStatusCodeStatusTextFromError(c int, err error) (int, string) {
	if c != 0 {
		return c, http.StatusText(http.StatusRequestTimeout)
	}

	if os.IsTimeout(err) {
		return http.StatusRequestTimeout, http.StatusText(http.StatusRequestTimeout)
	}

	rc := http.StatusInternalServerError
	rt := http.StatusText(rc)

	switch t := err.(type) {
	case *url.Error:
		rc = http.StatusServiceUnavailable
		rt = http.StatusText(rc)
		if t1, ok := t.Err.(*net.OpError); ok {
			switch t1.Op {
			case "dial":
				rt = "Unknown host"
			case "read":
				rt = "Connection refused"
			}
		}

	case *net.OpError:
		rc = http.StatusServiceUnavailable
		rt = http.StatusText(rc)
		switch t.Op {
		case "dial":
			rt = "Unknown host"
		case "read":
			rt = "Connection refused"
		}

	case syscall.Errno:
		rc = http.StatusServiceUnavailable
		rt = http.StatusText(rc)
		if t == syscall.ECONNREFUSED {
			rt = "Connection refused"
		}
	}

	return rc, rt
}
