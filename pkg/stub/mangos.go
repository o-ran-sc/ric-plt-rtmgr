/*
==================================================================================
  Copyright (c) 2019 AT&T Intellectual Property.
  Copyright (c) 2019 Nokia

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

   This source code is part of the near-RT RIC (RAN Intelligent Controller)
   platform project (RICP).

==================================================================================
*/
/*
  Mnemonic:	mangos.go
  Abstract:
  Date:		3 May 2019
*/

package stub

import "errors"

type MangosMessage struct {
	Header []byte
	Body   []byte
	Pipe   MangosPipe
	bbuf   []byte
	hbuf   []byte
	bsize  int
	pool   interface{}
}

type MangosProtocolInfo struct {
	Self     uint16
	Peer     uint16
	SelfName string
	PeerName string
}

// Mangos Listener Stub

type MangosListener struct {
}

func (l MangosListener) Listen() error {
	return nil
}

func (l MangosListener) Close() error {
	return nil
}

func (l MangosListener) Address() string {
	return ""
}

func (l MangosListener) SetOption(s string, i interface{}) error {
	return nil
}

func (l MangosListener) GetOption(s string) (interface{}, error) {
	return nil, nil
}

// Mangos Dialer Stub

type MangosDialer struct {
}

func (d MangosDialer) Open() error {
	return nil
}

func (d MangosDialer) Close() error {
	return nil
}

func (d MangosDialer) Address() string {
	return ""
}

func (d MangosDialer) SetOption(s string, i interface{}) error {
	return nil
}

func (d MangosDialer) GetOption(s string) (interface{}, error) {
	return nil, nil
}

// Mangos Context Stub

type MangosContext struct {
}

func (c MangosContext) Close() error {
	return nil
}

func (c MangosContext) SetOption(s string, i interface{}) error {
	return nil
}

func (c MangosContext) GetOption(s string) (interface{}, error) {
	return nil, nil
}

func (c MangosContext) Send(b []byte) error {
	return nil
}

func (c MangosContext) Recv() ([]byte, error) {
	return make([]byte, 0), nil
}

func (c MangosContext) SendMsg(*MangosMessage) error {
	return nil
}

func (c MangosContext) RecvMsg() (*MangosMessage, error) {
	return nil, nil
}

// Mangos Pipe Stub

type MangosPipe struct {
}

func (p MangosPipe) ID() uint32 {
	return 0
}

func (p MangosPipe) Listener() MangosListener {
	return MangosListener{}
}

func (p MangosPipe) Dialer() MangosDialer {
	return MangosDialer{}
}

func (p MangosPipe) Close() error {
	return nil
}

func (p MangosPipe) Address() string {
	return ""
}

func (p MangosPipe) GetOption(s string) (interface{}, error) {
	return nil, nil
}

// Mangos PipeEventHook Stub

type PipeEventHook func(int, MangosPipe)

// Mangos Socket Stub

type MangosSocket struct {
	GenerateSocketCloseError  bool
	GenerateSocketSendError   bool
	GenerateSocketDialError   bool
	GenerateSocketListenError bool
}

func (s MangosSocket) Info() MangosProtocolInfo {
	return MangosProtocolInfo{}
}

func (s MangosSocket) Close() error {
	if s.GenerateSocketCloseError {
		return errors.New("stub generated Socket Close error")
	}
	return nil
}

func (s MangosSocket) Send(b []byte) error {
	if s.GenerateSocketSendError {
		return errors.New("stub generated Socket Send error")
	}
	return nil
}

func (s MangosSocket) Recv() ([]byte, error) {
	return make([]byte, 0), nil
}

func (s MangosSocket) SendMsg(*MangosMessage) error {
	return nil
}

func (s MangosSocket) RecvMsg() (*MangosMessage, error) {
	return nil, nil
}

func (s MangosSocket) Dial(t string) error {
	if s.GenerateSocketDialError {
		return errors.New("stub generated Socket Dial error")
	}
	return nil
}

func (s MangosSocket) DialOptions(t string, m map[string]interface{}) error {
	if err := s.Dial(t); err != nil {
		return err
	}
	return nil
}

func (s MangosSocket) NewDialer(t string, m map[string]interface{}) (MangosDialer, error) {
	return MangosDialer{}, nil
}

func (s MangosSocket) Listen(t string) error {
	if s.GenerateSocketListenError {
		return errors.New("stub generated Socket Listen error")
	}
	return nil
}

func (s MangosSocket) ListenOptions(t string, m map[string]interface{}) error {
	return nil
}

func (s MangosSocket) NewListener(t string, m map[string]interface{}) (MangosListener, error) {
	return MangosListener{}, nil
}

func (s MangosSocket) SetOption(t string, i interface{}) error {
	return nil
}

func (s MangosSocket) GetOption(t string) (interface{}, error) {
	return nil, nil
}

func (s MangosSocket) OpenContext() (MangosContext, error) {
	return MangosContext{}, nil
}

func (s MangosSocket) SetPipeEventHook(p PipeEventHook) PipeEventHook {
	return nil
}

// Mangos ProtocolPipe Stub

type MangosProtocolPipe struct {
}

func (p MangosProtocolPipe) ID() uint32 {
	return 0
}

func (p MangosProtocolPipe) Close() error {
	return nil
}

func (p MangosProtocolPipe) SendMsg(m *MangosMessage) error {
	return nil
}

func (p MangosProtocolPipe) RecvMsg() *MangosMessage {
	return nil
}

// Mangos ProtocolContext Stub

type MangosProtocolContext struct {
}

func (p MangosProtocolContext) Close() error {
	return nil
}

func (p MangosProtocolContext) SendMsg(m *MangosMessage) error {
	return nil
}

func (p MangosProtocolContext) RecvMsg() (*MangosMessage, error) {
	return nil, nil
}

func (p MangosProtocolContext) GetOption(s string) (interface{}, error) {
	return nil, nil
}

func (p MangosProtocolContext) SetOption(s string, i interface{}) error {
	return nil
}

// Mangos ProtocolBase Stub

type MangosProtocolBase struct {
	MangosProtocolContext
}

func (p MangosProtocolBase) Info() MangosProtocolInfo {
	return MangosProtocolInfo{}
}

func (p MangosProtocolBase) AddPipe(t MangosProtocolPipe) error {
	return nil
}

func (p MangosProtocolBase) RemovePipe(MangosProtocolPipe) {

}

func (p MangosProtocolBase) OpenContext() (MangosProtocolContext, error) {
	return MangosProtocolContext{}, nil
}
