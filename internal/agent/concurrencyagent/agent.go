package concurrencyagent

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"runtime"
	"sync"
	"ya-metrics/internal/agent/runableagent"
)

type ConcurrencyAgent struct {
	RateLimit uint
	mu        sync.Mutex
	counter   uint
	reqCh     chan *http.Request
	logger    *zap.SugaredLogger
	client    *http.Client
}

func New(logger *zap.SugaredLogger, client *http.Client, rateLimit uint) *ConcurrencyAgent {
	return &ConcurrencyAgent{
		RateLimit: rateLimit,
		counter:   rateLimit,
		reqCh:     make(chan *http.Request, rateLimit),
		logger:    logger,
		client:    client,
	}
}

func (c *ConcurrencyAgent) Run(ctx context.Context, srvrAddr string, pCount int64, reportIntervalSec int, secretKey string) {
	go c.requestFactory(ctx, srvrAddr, pCount, reportIntervalSec, secretKey)

	go func() {
		countIncomeReq := 0
		for req := range c.reqCh {
			c.logger.Infof("Num gorutine:%d", runtime.NumGoroutine())
			c.logger.Infof("chan len: %d", len(c.reqCh))
			countIncomeReq++
			c.mu.Lock()
			if c.counter == 0 {
				c.counter = c.RateLimit - 1
			}
			go c.send(ctx, req, int(c.counter))
			c.mu.Unlock()
			c.logger.Infof("count requests: %d", countIncomeReq)
		}
	}()
}

func (c *ConcurrencyAgent) send(ctx context.Context, req *http.Request, order int) {
	select {
	default:
	case <-ctx.Done():
		c.logger.Infof("worker № %d stopped by closing context", order)
		return
	}
	resp, err := c.client.Do(req)
	if err != nil && !errors.Is(err, io.EOF) {
		c.logger.Errorf("err on send request:%s", err)
		return
	}
	//TODO: ВОПРОС : тут падает vet test, но действительно ли нужно закрывать тело ответа, если оно не используется?
	defer resp.Body.Close()
	c.mu.Lock()
	if c.counter < c.RateLimit {
		c.counter++
	}
	c.mu.Unlock()
}

func (c *ConcurrencyAgent) requestFactory(ctx context.Context, srvrAddr string, pCount int64, reportIntervalSec int, secretKey string) {
	compressJSONAgent := runableagent.CompressJSONAgent{
		SecretKey: secretKey,
		SendCh:    c.reqCh,
		Logger:    c.logger,
	}
	SimpleAgent := runableagent.SimpleAgent{
		SecretKey: secretKey,
		SendCh:    c.reqCh,
		Logger:    c.logger,
	}

	for {
		select {
		default:
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				compressJSONAgent.SendMetrics(srvrAddr, pCount, reportIntervalSec)
			}()
			wg.Add(1)
			go func() {
				defer wg.Done()
				SimpleAgent.SendMetrics(srvrAddr, pCount, reportIntervalSec)
			}()
			wg.Wait()
		case <-ctx.Done():
			c.logger.Info("requestFactory stopped by context")
			return
		}
	}
}
