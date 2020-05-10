// +build darwin

package watch

/*
#include <CoreServices/CoreServices.h>
typedef void (*CFRunLoopPerformCallBack)(void*);
void source(uintptr_t);
void notify(uintptr_t, uintptr_t, size_t, uintptr_t, uintptr_t, uintptr_t);
static FSEventStreamRef EventStreamCreate(
	FSEventStreamContext * context,
	uintptr_t info, CFArrayRef paths,
	FSEventStreamEventId since,
	CFTimeInterval latency,
	FSEventStreamCreateFlags flags
) {
	context->info = (void*) info;
	return FSEventStreamCreate(NULL, (FSEventStreamCallback) notify, context, paths, since, latency, flags); //nolint: lll
}
#cgo LDFLAGS: -framework CoreServices
*/
import "C"

import (
	"context"
	"errors"
	"io"
	"runtime"
	"sync"
	"unsafe"

	"github.com/tvastar/gotools/pkg/handles"
)

var hTable handles.Table //nolint: gochecknoglobals

// DirFSEvents implements watching a directory (and its descendants)
// using FSEvents.
func DirFSEvents(dir string) Stream {
	return &fse{dir: dir}
}

type fse struct {
	dir     string
	ref     C.FSEventStreamRef
	handle  uintptr
	runloop C.CFRunLoopRef
	ch      chan string
	closed  chan struct{}
}

func (f *fse) NextPath(ctx context.Context) (string, error) {
	if f.ch == nil {
		f.ch = make(chan string)
		f.closed = make(chan struct{})
		if err := f.init(); err != nil {
			f.ch = nil
			return "", err
		}
	}

	select {
	case <-f.closed:
		return "", io.EOF
	case s := <-f.ch:
		return s, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (f *fse) Close() error {
	if f.ch == nil {
		return nil
	}
	f.stop()
	close(f.closed)
	return nil
}

func (f *fse) init() error {
	var wg sync.WaitGroup
	h := hTable.Add(&wg)
	ctx := &C.CFRunLoopSourceContext{
		perform: (C.CFRunLoopPerformCallBack)(C.source),
		info:    unsafe.Pointer(h),
	}
	source := C.CFRunLoopSourceCreate(C.kCFAllocatorDefault, 0, ctx)
	go func() {
		runtime.LockOSThread()
		f.runloop = C.CFRunLoopGetCurrent()
		C.CFRunLoopAddSource(f.runloop, source, C.kCFRunLoopDefaultMode)
		C.CFRunLoopRun()
	}()
	wg.Wait()
	hTable.Delete(h)

	f.handle = hTable.Add(f)
	return f.listen()
}

func (f *fse) listen() error {
	dir := C.CFStringCreateWithCStringNoCopy(
		C.kCFAllocatorDefault,
		C.CString(f.dir),
		C.kCFStringEncodingUTF8,
		C.kCFAllocatorDefault,
	)
	path := C.CFArrayCreate(
		C.kCFAllocatorDefault,
		(*unsafe.Pointer)(unsafe.Pointer(&dir)),
		1,
		nil, //nolint: gocritic
	)

	var nilRef C.FSEventStreamRef
	ref := f.createEventStream(C.uintptr_t(f.handle), path)
	if ref == nilRef {
		return errors.New("could not create stream")
	}

	C.FSEventStreamScheduleWithRunLoop(ref, f.runloop, C.kCFRunLoopDefaultMode)
	if C.FSEventStreamStart(ref) == C.Boolean(0) {
		C.FSEventStreamInvalidate(ref)
		return errors.New("could not create stream")
	}
	C.CFRunLoopWakeUp(f.runloop)
	f.ref = ref
	return nil
}

func (f *fse) createEventStream(info C.uintptr_t, path C.CFArrayRef) C.FSEventStreamRef {
	ctx := C.FSEventStreamContext{}
	since := C.FSEventsGetCurrentEventId()
	var latency C.CFTimeInterval
	flags := C.FSEventStreamCreateFlags(C.kFSEventStreamCreateFlagFileEvents | C.kFSEventStreamCreateFlagNoDefer)

	//nolint: gocritic
	return C.EventStreamCreate(&ctx, C.uintptr_t(f.handle), path, since, latency, flags)
}

func (f *fse) notify(paths []string) {
	for _, path := range paths {
		select {
		case <-f.closed:
		case f.ch <- path:
		}
	}
}

func (f *fse) stop() {
	var nilref C.FSEventStreamRef
	if f.ref == nilref {
		return
	}
	C.FSEventStreamStop(f.ref)
	C.FSEventStreamInvalidate(f.ref)
	C.CFRunLoopWakeUp(f.runloop)
	f.ref = nilref
	hTable.Delete(f.handle)
}

//export source
func source(info uintptr) {
	if v, ok := hTable.Get(info); ok {
		v.(*sync.WaitGroup).Done()
	}
}

//export notify
func notify(_, info uintptr, n C.size_t, cpaths, flags, ids uintptr) {
	const offchar = unsafe.Sizeof((*C.char)(nil))
	paths := make([]string, 0, int(n))
	for i := uintptr(0); i < uintptr(n); i++ {
		path := C.GoString(*(**C.char)(unsafe.Pointer(cpaths + i*offchar)))
		paths = append(paths, path)
	}
	if v, ok := hTable.Get(info); ok {
		v.(*fse).notify(paths)
	}
}
