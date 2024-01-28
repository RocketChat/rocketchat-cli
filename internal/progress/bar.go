package progress

import (
	"fmt"
	"io"
	"math"
	"os"
	"sync"

	"golang.org/x/sys/unix"
)

const extraCharacters = "[] []"

// pbmanager manages all progressbars
type pbmanager struct {
	mut sync.Mutex
	// current is a flawed idea
	current      int
	counter      int
	progressBars int
	writer       io.Writer
}

var _m = &pbmanager{current: 0, counter: -1, mut: sync.Mutex{}, writer: os.Stdout}

type progressBar struct {
	order int

	writer   io.Writer
	label    string
	maxBytes int64

	barWidth int

	//state
	bytesWritten   int64
	columnsWritten int
	_progressbar   []byte
	_manager       *pbmanager
}

type ProgressBarOptions struct {
	BarWidth float64
}

func (p *progressBar) Write(data []byte) (int, error) {
	n, err := p.writer.Write(data)
	if err != nil {
		return 0, err
	}

	if n == 0 {
		return 0, nil
	}

	p._manager.mut.Lock()
	defer p._manager.mut.Unlock()

	if p._manager.current > p.order {
		fmt.Fprintf(p._manager.writer, "\033[%dF", p._manager.current-p.order) // go up by x
		p._manager.current = p.order
	} else if p._manager.current < p.order {
		for i := p._manager.current; i < p.order; i++ {
			fmt.Fprintf(p._manager.writer, "\n")
		}
		p._manager.current = p.order
	}

	fmt.Fprintf(p._manager.writer, "\r[%s] [%s]", p.label, p.buildprogressbar(n))

	if p.bytesWritten == p.maxBytes {
		p._manager.counter--
		if p._manager.counter == -1 {
			for i := p.order; i < p._manager.progressBars; i++ {
				p._manager.writer.Write([]byte{10})
			}
		}
	}

	return n, err
}

func (p *progressBar) buildprogressbar(bytes int) string {
	p.bytesWritten += int64(bytes)

	fill := int(math.Floor(float64(p.barWidth) * float64(p.bytesWritten) / float64(p.maxBytes)))

	for i := p.columnsWritten; i < fill; i += 1 {
		p._progressbar[i] = '='
	}

	p.columnsWritten = fill

	return string(p._progressbar)
}

func NewWriteProgressBar(label string, maxBytes int64, writer io.Writer, opts *ProgressBarOptions) (io.Writer, error) {
	size, err := unix.IoctlGetWinsize(0, unix.TIOCGWINSZ)
	if err != nil {
		return nil, err
	}

	_m.mut.Lock()
	defer _m.mut.Unlock()

	p := &progressBar{writer: writer, label: label, maxBytes: maxBytes, _manager: _m}

	_m.counter += 1
	p.order = _m.counter
	_m.progressBars++

	columns := int(size.Col)

	if opts != nil {
		columns = int(math.Floor(float64(columns) * float64(opts.BarWidth) / 100.0))
	}

	p.barWidth = int(columns) - len(label) - len(extraCharacters)

	p._progressbar = make([]byte, p.barWidth)
	for i := range p._progressbar {
		p._progressbar[i] = ' '
	}

	return p, nil
}
