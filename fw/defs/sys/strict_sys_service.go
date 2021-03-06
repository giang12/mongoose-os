// Code generated by clubbygen.
// GENERATED FILE DO NOT EDIT
// +build clubby_strict

package sys

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cesanta.com/clubby"
	"cesanta.com/clubby/endpoint"
	"cesanta.com/clubby/frame"
	"cesanta.com/common/go/ourjson"
	"cesanta.com/common/go/ourtrace"
	"github.com/cesanta/errors"
	"golang.org/x/net/trace"

	"github.com/cesanta/ucl"
	"github.com/cesanta/validate-json/schema"
	"github.com/golang/glog"
)

var _ = bytes.MinRead
var _ = fmt.Errorf
var emptyMessage = ourjson.RawMessage{}
var _ = ourtrace.New
var _ = trace.New

const ServiceID = "http://mongoose-iot.com/fwSys"

type RebootArgs struct {
	Delay_ms *int64 `json:"delay_ms,omitempty"`
}

type Service interface {
	Reboot(ctx context.Context, args *RebootArgs) error
}

type Instance interface {
	Call(context.Context, string, *frame.Command) (*frame.Response, error)
	TraceCall(context.Context, string, *frame.Command) (context.Context, trace.Trace, func(*error))
}

type _validators struct {
	// This comment prevents gofmt from aligning types in the struct.
	RebootArgs *schema.Validator
}

var (
	validators     *_validators
	validatorsOnce sync.Once
)

func initValidators() {
	validators = &_validators{}

	loader := schema.NewLoader()

	service, err := ucl.Parse(bytes.NewBuffer(_ServiceDefinition))
	if err != nil {
		panic(err)
	}
	// Patch up shortcuts to be proper schemas.
	for _, v := range service.(*ucl.Object).Find("methods").(*ucl.Object).Value {
		if s, ok := v.(*ucl.Object).Find("result").(*ucl.String); ok {
			for kk := range v.(*ucl.Object).Value {
				if kk.Value == "result" {
					v.(*ucl.Object).Value[kk] = &ucl.Object{
						Value: map[ucl.Key]ucl.Value{
							ucl.Key{Value: "type"}: s,
						},
					}
				}
			}
		}
		if v.(*ucl.Object).Find("args") == nil {
			continue
		}
		args := v.(*ucl.Object).Find("args").(*ucl.Object)
		for kk, vv := range args.Value {
			if s, ok := vv.(*ucl.String); ok {
				args.Value[kk] = &ucl.Object{
					Value: map[ucl.Key]ucl.Value{
						ucl.Key{Value: "type"}: s,
					},
				}
			}
		}
	}
	var s *ucl.Object
	_ = s // avoid unused var error
	s = &ucl.Object{
		Value: map[ucl.Key]ucl.Value{
			ucl.Key{Value: "properties"}: service.(*ucl.Object).Find("methods").(*ucl.Object).Find("Reboot").(*ucl.Object).Find("args"),
			ucl.Key{Value: "type"}:       &ucl.String{Value: "object"},
		},
	}
	if req, found := service.(*ucl.Object).Find("methods").(*ucl.Object).Find("Reboot").(*ucl.Object).Lookup("required_args"); found {
		s.Value[ucl.Key{Value: "required"}] = req
	}
	validators.RebootArgs, err = schema.NewValidator(s, loader)
	if err != nil {
		panic(err)
	}
}

func NewClient(i Instance, addr string) Service {
	validatorsOnce.Do(initValidators)
	return &_Client{i: i, addr: addr}
}

type _Client struct {
	i    Instance
	addr string
}

func (c *_Client) Reboot(pctx context.Context, args *RebootArgs) (err error) {
	cmd := &frame.Command{
		Cmd: "Sys.Reboot",
	}
	ctx, tr, finish := c.i.TraceCall(pctx, c.addr, cmd)
	defer finish(&err)
	_ = tr

	tr.LazyPrintf("args: %s", ourjson.LazyJSON(&args))
	cmd.Args = ourjson.DelayMarshaling(args)
	b, err := cmd.Args.MarshalJSON()
	if err != nil {
		glog.Errorf("Failed to marshal args as JSON: %+v", err)
	} else {
		v, err := ucl.Parse(bytes.NewReader(b))
		if err != nil {
			glog.Errorf("Failed to parse just serialized JSON value %q: %+v", string(b), err)
		} else {
			if err := validators.RebootArgs.Validate(v); err != nil {
				glog.Warningf("Sending invalid args for Reboot: %+v", err)
				return errors.Annotatef(err, "invalid args for Reboot")
			}
		}
	}
	resp, err := c.i.Call(ctx, c.addr, cmd)
	if err != nil {
		return errors.Trace(err)
	}
	if resp.Status != 0 {
		return errors.Trace(&endpoint.ErrorResponse{Status: resp.Status, Msg: resp.StatusMsg})
	}
	return nil
}

func RegisterService(i *clubby.Instance, impl Service) error {
	validatorsOnce.Do(initValidators)
	s := &_Server{impl}
	i.RegisterCommandHandler("Sys.Reboot", s.Reboot)
	i.RegisterService(ServiceID, _ServiceDefinition)
	return nil
}

type _Server struct {
	impl Service
}

func (s *_Server) Reboot(ctx context.Context, src string, cmd *frame.Command) (interface{}, error) {
	b, err := cmd.Args.MarshalJSON()
	if err != nil {
		glog.Errorf("Failed to marshal args as JSON: %+v", err)
	} else {
		if v, err := ucl.Parse(bytes.NewReader(b)); err != nil {
			glog.Errorf("Failed to parse valid JSON value %q: %+v", string(b), err)
		} else {
			if err := validators.RebootArgs.Validate(v); err != nil {
				glog.Warningf("Got invalid args for Reboot: %+v", err)
				return nil, errors.Annotatef(err, "invalid args for Reboot")
			}
		}
	}
	var args RebootArgs
	if len(cmd.Args) > 0 {
		if err := cmd.Args.UnmarshalInto(&args); err != nil {
			return nil, errors.Annotatef(err, "unmarshaling args")
		}
	}
	return nil, s.impl.Reboot(ctx, &args)
}

var _ServiceDefinition = json.RawMessage([]byte(`{
  "methods": {
    "Reboot": {
      "args": {
        "delay_ms": {
          "doc": "Optional delay, in milliseconds. Default is 100.",
          "type": "integer"
        }
      },
      "doc": "Reboot the device"
    }
  },
  "name": "Sys",
  "namespace": "http://mongoose-iot.com/fw"
}`))
