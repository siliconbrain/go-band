package band

import "sync"

// New returns a new, initialized Band instance
func New() *Band {
	return &Band{
		disbanded:  make(latch),
		disbanding: make(latch),
	}
}

// Band represents a group of members that can be joined and left, disbanded with a reason, and followed until disbanding
type Band struct {
	disbanded  latch
	disbanding latch
	members    int
	mutex      sync.Mutex
	reason     interface{}
}

// Collab joins the band, executes the specified function, and leaves the band (even if the function panics)
func (b *Band) Collab(fn func()) {
	defer b.Join().Leave()
	fn()
}

// Disband initiates band break-up with the specified reason
func (b *Band) Disband(reason interface{}) {
	if !open(b.disbanded) {
		return
	}
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if open(b.disbanding) {
		b.reason = reason
		close(b.disbanding)
		if b.members == 0 {
			close(b.disbanded)
		}
	}
}

// Disbanded returns a channel that is closed when the band has finished disbanding (i.e. all members left)
func (b *Band) Disbanded() WaitChannel {
	return b.disbanded
}

// Disbanding returns a channel that is closed when band break-up has been initiated
func (b *Band) Disbanding() WaitChannel {
	return b.disbanding
}

// Follow waits for the band to disband and returns with the reason for disbanding
func (b *Band) Follow() (reason interface{}) {
	<-b.disbanded
	return b.reason // reason can be read without synchronization since it won't change after the band has disbanded
}

// Join joins the band as a member and returns a "token of membership"
func (b *Band) Join() *Membership {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.members++
	return &Membership{
		band: b,
	}
}

func (b *Band) leave() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.members--
	if b.members == 0 && !open(b.disbanding) && open(b.disbanded) {
		close(b.disbanded)
	}
}

// Membership should be used to leave the band
type Membership struct {
	band *Band
	once sync.Once
}

// Leave leaves the band associated with this membership
func (m *Membership) Leave() {
	m.once.Do(m.band.leave)
}

type WaitChannel = <-chan void

type latch = chan void
type void = struct{}

func open(l WaitChannel) bool {
	select {
	case <-l:
		return false
	default:
		return true
	}
}
