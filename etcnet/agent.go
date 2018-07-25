package etcnet

// func (a *Agent) reject(reason uint8) {
// 	if a.ws != nil {
// 		rejectW(a.ws, reason)
// 		return
// 	}

// 	a.l.reject(reason)
// }

// func (a *Agent) handleInit(head *etc.Buffer) {
// 	flags := head.ReadByte()
// 		initialized := head.ReadDate()
// 		signkey := head.ReadBytes(32)
// 		sessionkey := head.ReadBoxKey()
// 		signature := head.ReadBytes(64)

// 		tnow := time.Now()

// 		testData := etc.NewBuffer()
// 		testData.WriteDate(initialized)
// 		testData.Write(sessionkey[:])

// 		if a.l.c.KeyCheck != nil {
// 			ok := a.l.c.KeyCheck(strings.ToUpper(hex.EncodeToString(signkey[:])))
// 			if !ok {
// 				a.reject(REJECT_SIGNING)
// 				return
// 			}
// 		}

// 		ok := ed25519.Verify(signkey, testData.Bytes(), signature)
// 		if !ok {
// 			rejectW(c, REJECT_SIGNING)
// 			return
// 		}

// 		if tnow.Sub(initialized) > ((MAX_CLOCK_DIFFERENCE) * time.Millisecond) {
// 			rejectW(c, REJECT_CLOCK)
// 			return
// 		}

// 		ag := new(Agent)
// 		ag.ws = c
// 		ag.id = signkey
// 		ag.sessionPeerKey = sessionkey
// 		ag.sessionKey = new([32]byte)
// 		curve25519.ScalarMult(ag.sessionKey, ag.sessionPeerKey, ag.sessionPrivateKey)

// 		if ag.ws != nil {
// 			l.agents.Store(etc.GenerateRandomUUID().String(), ag)
// 		}
// }
