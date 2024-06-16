package errs

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "runtime"
    "strings"
    "time"

    "store-api/pkg/data/constant"
)

var (
    dunno      = []byte("???")
    centerDot  = []byte("·")
    dot        = []byte(".")
    slash      = []byte("/")
    pathPrefix = "/go/src/store-api"
)

func GetCurrentFileAndLine(index int) string {
    _, file, line, ok := runtime.Caller(index)
    if !ok {
        return ""
    }
    return fmt.Sprintf("%s:%d", file, line)
}

func GetFileAndFuncForLogrus() string {
    buf := new(bytes.Buffer)
    var lines [][]byte

    pc, file, line, ok := runtime.Caller(9)
    if !ok {
        return ""
    }
    fmt.Fprintf(buf, "%s:%d\n", file, line)
    if file != "" {
        data, err := ioutil.ReadFile(file)
        if err != nil {
            return ""
        }
        lines = bytes.Split(data, []byte{'\n'})
    }
    fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))

    return strings.ReplaceAll(buf.String(), pathPrefix, "...")
}

func Stack(skip int) string {
    stack, _ := StackAndFile(skip)
    return stack
}

// stack returns a nicely formatted stack frame, skipping skip frames.
//
func StackAndFile(skip int) (string, string) {
    buf := new(bytes.Buffer) // the returned data
    // As we loop, we open files and read them. These variables record the currently
    // loaded file.
    var lines [][]byte
    var lastFile string
    var firstFile string
    for i := skip; ; i++ { // Skip the expected number of frames

        pc, file, line, ok := runtime.Caller(i)
        if !ok {
            break
        }
        if firstFile == "" {
            firstFile = strings.ReplaceAll(fmt.Sprintf("%s:%d", file, line), pathPrefix, "...")
        }

        // ---- Skip un-necessary trace
        if strings.Contains(file, "net/http/server.go") {
            /**
            Often

            C:/Program Files/Go/src/net/http/server.go:2046 (0x3efdae)
            	HandlerFunc.ServeHTTP: f(w, r)
            C:/Users/Will/Desktop/code/yii2/tnt/vendor/github.com/gorilla/mux/mux.go:210 (0x51b64e)
            	(*Router).ServeHTTP: handler.ServeHTTP(w, req)
            C:/Users/Will/Desktop/code/yii2/tnt/vendor/gopkg.in/DataDog/dd-trace-go.v1/contrib/internal/httputil/trace.go:57 (0xad1109)
            	TraceAndServe: httpinstr.WrapHandler(h, span).ServeHTTP(cfg.ResponseWriter, cfg.Request.WithContext(ctx))
            C:/Users/Will/Desktop/code/yii2/tnt/vendor/gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux/mux.go:119 (0xadb00b)
            	(*Router).ServeHTTP: httputil.TraceAndServe(r.Router, &httputil.TraceConfig{
            C:/Program Files/Go/src/net/http/server.go:2878 (0x3f331a)
            	serverHandler.ServeHTTP: handler.ServeHTTP(rw, req)
            C:/Program Files/Go/src/net/http/server.go:1929 (0x3eee87)
            	(*conn).serve: serverHandler{c.server}.ServeHTTP(w, w.req)
            C:/Program Files/Go/src/runtime/asm_amd64.s:1581 (0x1a8f00)
            	goexit: BYTE	$0x90	// NOP
            */
            break
        }

        // Print this much at least.  If we can't find the source, it won't show.
        fmt.Fprintf(buf, "%s:%d\n", file, line)
        if file != lastFile {
            data, err := ioutil.ReadFile(file)
            if err != nil {
                continue
            }
            lines = bytes.Split(data, []byte{'\n'})
            lastFile = file
        }
        fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))

    }
    return strings.ReplaceAll(buf.String(), pathPrefix, "..."), firstFile
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
    n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
    if n < 0 || n >= len(lines) {
        return dunno
    }
    return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
    fn := runtime.FuncForPC(pc)
    if fn == nil {
        return dunno
    }
    name := []byte(fn.Name())
    // The name includes the path name to the package, which is unnecessary
    // since the file name is already included.  Plus, it has center dots.
    // That is, we see
    //	runtime/debug.*T·ptrmethod
    // and want
    //	*T.ptrmethod
    // Also the package path might contains dot (e.g. code.google.com/...),
    // so first eliminate the path prefix
    if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
        name = name[lastSlash+1:]
    }
    if period := bytes.Index(name, dot); period >= 0 {
        name = name[period+1:]
    }
    name = bytes.Replace(name, centerDot, dot, -1)
    return name
}

// timeFormat returns a customized time string for logger.
func timeFormat(t time.Time) string {
    return t.Format(constant.DefaultDatetimeLayout)
}
