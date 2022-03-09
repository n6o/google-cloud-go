package testutil

import (
	"context"
	"github.com/golang/protobuf/proto"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	"log"
)

type InMemInstanceAdminServer interface {
	instancepb.InstanceAdminServer
	Stop()
	Resps() []proto.Message
	SetResps([]proto.Message)
	Reqs() []proto.Message
	SetReqs([]proto.Message)
	SetErr(error)
}
type inMemInstanceAdminServer struct {
	instancepb.InstanceAdminServer
	reqs  []proto.Message
	err   error
	resps []proto.Message
}

func gologoo__NewInMemInstanceAdminServer_b747d06f230b85267a2a9e269dfef1d5() InMemInstanceAdminServer {
	res := &inMemInstanceAdminServer{}
	return res
}
func (s *inMemInstanceAdminServer) gologoo__GetInstance_b747d06f230b85267a2a9e269dfef1d5(ctx context.Context, req *instancepb.GetInstanceRequest) (*instancepb.Instance, error) {
	s.reqs = append(s.reqs, req)
	if s.err != nil {
		defer func() {
			s.err = nil
		}()
		return nil, s.err
	}
	return s.resps[0].(*instancepb.Instance), nil
}
func (s *inMemInstanceAdminServer) gologoo__Stop_b747d06f230b85267a2a9e269dfef1d5() {
}
func (s *inMemInstanceAdminServer) gologoo__Resps_b747d06f230b85267a2a9e269dfef1d5() []proto.Message {
	return s.resps
}
func (s *inMemInstanceAdminServer) gologoo__SetResps_b747d06f230b85267a2a9e269dfef1d5(resps []proto.Message) {
	s.resps = resps
}
func (s *inMemInstanceAdminServer) gologoo__Reqs_b747d06f230b85267a2a9e269dfef1d5() []proto.Message {
	return s.reqs
}
func (s *inMemInstanceAdminServer) gologoo__SetReqs_b747d06f230b85267a2a9e269dfef1d5(reqs []proto.Message) {
	s.reqs = reqs
}
func (s *inMemInstanceAdminServer) gologoo__SetErr_b747d06f230b85267a2a9e269dfef1d5(err error) {
	s.err = err
}
func NewInMemInstanceAdminServer() InMemInstanceAdminServer {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NewInMemInstanceAdminServer_b747d06f230b85267a2a9e269dfef1d5")
	log.Printf("Input : (none)\n")
	r0 := gologoo__NewInMemInstanceAdminServer_b747d06f230b85267a2a9e269dfef1d5()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemInstanceAdminServer) GetInstance(ctx context.Context, req *instancepb.GetInstanceRequest) (*instancepb.Instance, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetInstance_b747d06f230b85267a2a9e269dfef1d5")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__GetInstance_b747d06f230b85267a2a9e269dfef1d5(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *inMemInstanceAdminServer) Stop() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Stop_b747d06f230b85267a2a9e269dfef1d5")
	log.Printf("Input : (none)\n")
	s.gologoo__Stop_b747d06f230b85267a2a9e269dfef1d5()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemInstanceAdminServer) Resps() []proto.Message {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Resps_b747d06f230b85267a2a9e269dfef1d5")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__Resps_b747d06f230b85267a2a9e269dfef1d5()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemInstanceAdminServer) SetResps(resps []proto.Message) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SetResps_b747d06f230b85267a2a9e269dfef1d5")
	log.Printf("Input : %v\n", resps)
	s.gologoo__SetResps_b747d06f230b85267a2a9e269dfef1d5(resps)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemInstanceAdminServer) Reqs() []proto.Message {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Reqs_b747d06f230b85267a2a9e269dfef1d5")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__Reqs_b747d06f230b85267a2a9e269dfef1d5()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemInstanceAdminServer) SetReqs(reqs []proto.Message) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SetReqs_b747d06f230b85267a2a9e269dfef1d5")
	log.Printf("Input : %v\n", reqs)
	s.gologoo__SetReqs_b747d06f230b85267a2a9e269dfef1d5(reqs)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemInstanceAdminServer) SetErr(err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SetErr_b747d06f230b85267a2a9e269dfef1d5")
	log.Printf("Input : %v\n", err)
	s.gologoo__SetErr_b747d06f230b85267a2a9e269dfef1d5(err)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
