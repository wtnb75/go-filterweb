package filterweb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os/exec"

	"github.com/go-viper/mapstructure/v2"
)

// single command
type CommandConfig struct {
	Filter
	KeepEnvs         bool
	InputContentType string
	ContentType      string
	Dir              string
	Env              map[string]string
	Args             []string
}

func (cc *CommandConfig) Name() string {
	return "command"
}

func (cc *CommandConfig) New() Filter {
	return &CommandConfig{}
}

func (cc *CommandConfig) Accepts() []string {
	return []string{}
}

func (cc *CommandConfig) Prep(config Config, data Data) error {
	// defaults
	cc.KeepEnvs = false
	cc.ContentType = "text/plain"
	err := mapstructure.Decode(config.Params, cc)
	if err != nil {
		return err
	}
	// mandatory
	if len(cc.Args) == 0 {
		slog.Error("command filter requires 'args' parameter")
		return ErrMissingParams
	}
	return nil
}

func (cc *CommandConfig) Process(data Data) (Data, error) {
	res := Data{ContentType: cc.ContentType}
	cmd := exec.Command(cc.Args[0], cc.Args[1:]...)
	if cc.Dir != "" {
		cmd.Dir = cc.Dir
	}
	if !cc.KeepEnvs {
		cmd.Env = []string{}
	}
	for k, v := range cc.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		slog.Error("stdinpipe", "error", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("stdoutpipe", "error", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		slog.Error("stderrpipe", "error", err)
	}
	stdoutbuf := bytes.Buffer{}
	stderrbuf := bytes.Buffer{}
	go func() {
		io.Copy(&stdoutbuf, stdout)
		stdout.Close()
	}()
	go func() {
		io.Copy(&stderrbuf, stderr)
		stderr.Close()
	}()
	go func() {
		if cc.InputContentType != "" {
			bbuf, err := EncodeContentType(cc.InputContentType, data.Data)
			if err != nil {
				slog.Error("encode failed", "contenttype", cc.InputContentType, "error", err)
			} else {
				stdin.Write(bbuf)
			}
		} else if bbuf, ok := data.Data.([]byte); ok {
			stdin.Write(bbuf)
		} else if sbuf, ok := data.Data.(string); ok {
			stdin.Write([]byte(sbuf))
		} else {
			enc := json.NewEncoder(stdin)
			if err := enc.Encode(data.Data); err != nil {
				slog.Error("encode error", "error", err)
			}
		}
	}()
	err = cmd.Run()
	if err != nil {
		slog.Error("command failed", "error", err, "stdout", stdoutbuf.String(), "stderr", stderrbuf.String())
		return res, err
	}
	res.Data, err = DecodeContentType(cc.ContentType, stdoutbuf.Bytes())
	return res, err
}

func (cc *CommandConfig) Post(config Config, data Data) error {
	return nil
}

func init() {
	RegisterFilter(&CommandConfig{})
}
