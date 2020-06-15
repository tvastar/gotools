package script

import (
	"context"
	"io"
	"os"
)

type pipe struct {
	tasks            []Task
	writers, readers []*os.File
}

func (p pipe) Start(ctx context.Context) error {
	for kk := range p.tasks[1:] {
		r, w, err := os.Pipe()
		if err != nil {
			return err
		}
		p.writers[kk] = w
		p.readers[kk] = r
		p.tasks[kk] = p.tasks[kk].Stdout(w)
		p.tasks[kk+1] = p.tasks[kk+1].Stdin(r)
	}

	for kk := range p.tasks {
		if err := p.tasks[kk].Start(ctx); err != nil {
			if kk > 0 {
				p.readers[kk-1].Close()
				p.readers[kk-1] = nil
			}
			return err
		}
	}
	return nil
}

func (p pipe) Wait(ctx context.Context) error {
	var err error
	for kk, t := range p.tasks {
		if err2 := t.Wait(ctx); err2 != nil {
			err = err2
		}
		if kk < len(p.writers) {
			p.writers[kk].Close()
		}
		if kk > 0 {
			p.readers[kk-1].Close()
		}
	}
	return err
}

func (p pipe) Stdin(r io.Reader) Task {
	writers := make([]*os.File, len(p.tasks)-1)
	readers := make([]*os.File, len(p.tasks)-1)
	result := pipe{p.tasks, writers, readers}
	for kk := range p.tasks {
		result.tasks[kk] = p.tasks[kk]
		if kk == 0 {
			result.tasks[kk] = p.tasks[kk].Stdin(r)
		}
	}
	return result
}

func (p pipe) Stdout(w io.Writer) Task {
	writers := make([]*os.File, len(p.tasks)-1)
	readers := make([]*os.File, len(p.tasks)-1)
	result := pipe{p.tasks, writers, readers}
	for kk := range p.tasks {
		result.tasks[kk] = p.tasks[kk]
		if kk == len(p.tasks)-1 {
			result.tasks[kk] = p.tasks[kk].Stdout(w)
		}
	}
	return result
}
