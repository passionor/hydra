package cron

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/utility"
)

const (
	unstarted = 1
	pause     = 2
	running   = 4
)

//Processor cron管理程序，用于管理多个任务的执行，暂停，恢复，动态添加，移除
type Processor struct {
	*dispatcher.Engine
	lock      sync.Mutex
	done      bool
	closeChan chan struct{}
	length    int
	index     int
	span      time.Duration
	slots     []cmap.ConcurrentMap //time slots
	startTime time.Time
	status    int
}

//NewProcessor 创建processor
func NewProcessor() (p *Processor) {
	p = &Processor{
		status:    unstarted,
		closeChan: make(chan struct{}),
		span:      time.Second,
		length:    60,
		startTime: time.Now(),
	}
	p.Engine = dispatcher.New()
	p.Engine.Use(middleware.Recovery().DispFunc(CRON))
	p.Engine.Use(middleware.Logging().DispFunc())
	p.Engine.Use(middleware.Trace().DispFunc()) //跟踪信息
	p.Engine.Use(middleware.Delay().DispFunc()) //

	p.slots = make([]cmap.ConcurrentMap, p.length, p.length)
	for i := 0; i < p.length; i++ {
		p.slots[i] = cmap.New(2)
	}
	return p
}

//Start 所有任务
func (s *Processor) Start() error {
START:
	for {
		select {
		case <-s.closeChan:
			break START
		case <-time.After(s.span):
			s.execute()
		}
	}
	return nil
}

//Add 添加任务
func (s *Processor) Add(ts ...*task.Task) (err error) {
	for _, t := range ts {
		if t.Disable {
			s.Remove(t.GetUNQ())
			continue
		}

		//处理单次执行服务
		if ok, err := s.onceHandle(t); ok {
			return err
		}

		task, err := NewCronTask(t)
		if err != nil {
			return fmt.Errorf("构建cron.task失败:%v", err)
		}

		if !s.Engine.Find(task.GetService()) {
			s.Engine.Handle(task.GetMethod(), task.GetService(), middleware.ExecuteHandler(task.Service).DispFunc(CRON))
		}
		if _, _, err := s.add(task); err != nil {
			return err
		}
	}
	return

}
func (s *Processor) onceHandle(t *task.Task) (bool, error) {
	if t.Cron != "@once" && t.Cron != "@now" {
		return false, nil
	}
	if s.status == unstarted {
		select {
		case <-time.After(time.Second):
			if s.status == unstarted {
				return true, nil
			}
		}
	}
	if s.status == pause{
		return true,nil

	//注册服务
	task, _ := NewCronTask(t)

	//注册服务
	if !s.Engine.Find(task.GetService()) {
		s.Engine.Handle(task.GetMethod(), task.GetService(), middleware.ExecuteHandler(task.Service).DispFunc(CRON))
	}

	//异步执行
	go s.Engine.HandleRequest(task)
	return true, nil
}

func (s *Processor) add(task *CronTask) (offset int, round int, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.done {
		return -1, -1, nil
	}
	now := time.Now()
	nextTime := task.NextTime(now)
	if nextTime.Sub(now) < 0 {
		return -1, -1, errors.New("next time less than now.1")
	}
	offset, round = s.getOffset(now, nextTime)
	if offset < 0 || round < 0 {
		return -1, -1, errors.New("next time less than now.2")
	}
	task.Round.Update(round)
	s.slots[offset].Set(utility.GetGUID(), task)
	return
}

//Remove 移除服务
func (s *Processor) Remove(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, slot := range s.slots {
		slot.RemoveIterCb(func(k string, value interface{}) bool {
			task := value.(*CronTask)
			task.Disable = true
			return task.GetName() == name
		})
	}
}

//Pause 暂停所有任务
func (s *Processor) Pause() (bool, error) {
	if s.status != pause {
		s.status = pause
		return true, nil
	}
	return false, nil
}

//Resume 恢复所有任务
func (s *Processor) Resume() (bool, error) {
	if s.status != running {
		s.status = running
		return true, nil
	}
	return false, nil
}

//Close 退出
func (s *Processor) Close() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.done {
		s.done = true
		close(s.closeChan)
	}
}

//-------------------------------------内部处理------------------------------------

func (s *Processor) getOffset(now time.Time, next time.Time) (pos int, circle int) {
	d := next.Sub(now) //剩余时间
	delaySeconds := int(d/1e9) + 1
	intervalSeconds := int(s.span.Seconds())
	circle = int(delaySeconds / intervalSeconds / s.length)
	pos = int(s.index+delaySeconds/intervalSeconds) % s.length
	return
}

func (s *Processor) execute() {
	s.startTime = time.Now()
	s.lock.Lock()
	defer s.lock.Unlock()
	s.index = (s.index + 1) % s.length
	current := s.slots[s.index]
	current.RemoveIterCb(func(k string, value interface{}) bool {
		task := value.(*CronTask)
		task.Round.Reduce()
		if task.Round.Get() <= 0 {
			go s.handle(task)
			return true
		}
		return false
	})
}
func (s *Processor) handle(task *CronTask) error {
	if s.done || task.Disable {
		return nil
	}
	if s.status == running {
		task.Counter.Increase()
		s.Engine.HandleRequest(task) //触发服务引擎进行业务处理
	}
	_, _, err := s.add(task)
	if err != nil {
		return err
	}
	return nil

}
