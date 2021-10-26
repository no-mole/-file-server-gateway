package grpc_pool

import (
	"context"
	"errors"
	"sync"
	"time"

	"google.golang.org/grpc/connectivity"
)

const (
	defaultMaxStreamsPerConn  = 100
	defaultConnIdleSeconds    = float64(time.Minute / time.Second)
	defaultCycleMonitorTicker = 5 * time.Second
	defaultChannelCap         = 20
)

var ClosedErr = errors.New("grpc pool has closed")

type Pool interface {
	Get() (Conn, error)
	Restore(Conn)
	Close()
}

type pool struct {
	sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	ticker    *time.Ticker
	builder   Builder
	pool      *stack
	storage   map[string]*rpcConn
	toRecycle map[string]time.Time
	reStorage chan *rpcConn
	buffer    chan *rpcConn
	activity  chan byte
	closed    chan bool
}

// Option optional configs
type Option func(*option)

type option struct {
}

func NewPool(builder Builder, options ...Option) (Pool, error) {
	if builder == nil {
		return nil, errors.New("builder is null")
	}
	opt := new(option)
	for _, f := range options {
		f(opt)
	}

	stack := new(stack)
	conn := newConn(builder)
	stack.Push(conn)

	ctx, cancel := context.WithCancel(context.Background())

	pool := &pool{
		ctx:       ctx,
		cancel:    cancel,
		ticker:    time.NewTicker(defaultCycleMonitorTicker),
		builder:   builder,
		pool:      stack,
		storage:   make(map[string]*rpcConn),
		toRecycle: make(map[string]time.Time),
		reStorage: make(chan *rpcConn, defaultChannelCap),
		buffer:    make(chan *rpcConn, defaultChannelCap),
		activity:  make(chan byte, defaultChannelCap),
		closed:    make(chan bool),
	}
	pool.toRecycle[conn.id] = time.Now()
	go pool.Hold()
	return pool, nil
}

func (p *pool) Get() (Conn, error) {
	for {
		select {
		case <-p.ctx.Done():
			return nil, ClosedErr
		default:
			p.activity <- 1
			if rpcConn := <-p.buffer; rpcConn.conn.GetState() != connectivity.Shutdown {
				p.Add(1)
				return rpcConn, nil
			}
		}
	}
}

func (p *pool) Restore(c Conn) {
	if c == nil {
		return
	}
	p.reStorage <- c.(*rpcConn)
	p.Done()
}

func (p *pool) Close() {
	p.cancel()
	p.ticker.Stop()

	p.Wait()
	close(p.closed)

	//关闭正在回收的conn
	for id := range p.toRecycle {
		p.pool.Remove(id).conn.Close()
	}
	for !p.pool.Empty() {
		p.pool.Pop().conn.Close()
	}
}

func (p *pool) Hold() {
	for {
		select {
		case <-p.closed:
			return
		case <-p.ticker.C:
			for id, ts := range p.toRecycle {
				//超过conn的最大连接时间
				if time.Since(ts).Seconds() > defaultConnIdleSeconds {
					conn := p.pool.Remove(id)
					conn.conn.Close()
					delete(p.toRecycle, id)
				}
			}
		case conn := <-p.reStorage:
			if c, ok := p.storage[conn.id]; ok {
				//释放一个conn
				delete(p.storage, c.id)
				c.streams--
				p.pool.Push(c)
				continue
			}
			if conn.streams >= 1 {
				if conn.streams--; conn.streams == 0 {
					p.toRecycle[conn.id] = time.Now()
				}
			}
		case <-p.activity:
		GET:
			conn := p.pool.Peek()
			if conn == nil {
				newConn := newConn(p.builder)
				p.pool.Push(newConn)
				goto GET
			}

			if conn.streams+1 <= defaultMaxStreamsPerConn {
				if conn.streams++; conn.streams == 1 {
					delete(p.toRecycle, conn.id)
				}
				goto PUT
			}
			p.storage[conn.id] = p.pool.Pop()
			if !p.pool.Empty() {
				goto GET
			}
			conn = newConn(p.builder)
			conn.streams++
			p.pool.Push(conn)
		PUT:
			p.buffer <- conn
		}
	}
}
