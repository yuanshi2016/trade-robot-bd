package service

import (
	"context"
	"fmt"
	"log"
	"time"
	"trade-robot-bd/api/response"
	pb "trade-robot-bd/api/usercenter/v1"
	walletpb "trade-robot-bd/api/wallet/v1"
	"trade-robot-bd/app/usercenter-svc/internal/dao"
	"trade-robot-bd/app/usercenter-svc/internal/model"
	"trade-robot-bd/libs/bcrypt2"
	"trade-robot-bd/libs/env"
	"trade-robot-bd/libs/logger"
	"trade-robot-bd/libs/message"
	validate_code "trade-robot-bd/libs/validate-code"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	pb.UnimplementedUserServer
	dao       *dao.Dao
	WalletSrv walletpb.WalletClient
}

func NewUserService(r *etcd.Registry) *UserService {
	if r == nil {
		return &UserService{
			dao: dao.New(),
		}
	}
	// w, err := r.GetService(context.Background(), env.WalletSrvName)
	// en := w[0].Endpoints[1]
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithDiscovery(r),
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%v", env.WalletSrvName)),
	)
	if err != nil {
		log.Fatalf("Wallet服务发现失败%v", err.Error())
	}
	// conn.Connect()
	wall := walletpb.NewWalletClient(conn)
	// dd, err := wall.CreateWallet(context.Background(), &walletpb.UidReq{UserId: "1"})
	return &UserService{
		dao:       dao.New(),
		WalletSrv: wall,
	}
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	var resp *pb.LoginResp
	user := s.dao.GetWqUserBaseByPhone(req.Phone)
	if user == nil {
		return nil, response.NewLoginPasswordOrPhoneErrMsg(errID)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, response.NewLoginPasswordOrPhoneErrMsg(errID)
	}
	resp.UserId = user.UserID
	resp.Phone = user.Phone
	resp.Name = user.Name
	resp.Avatar = user.Avatar
	resp.InvitationCode = user.InvitationCode
	resp.LastLoginAt = timestamppb.New(user.LastLogin)
	resp.LoginCount = int32(user.LoginCount)
	return &pb.LoginResp{}, nil
}
func (s *UserService) SendValidateCode(ctx context.Context, req *pb.ValidateCodeReq) (*pb.ValidateCodeResp, error) {
	var resp *pb.ValidateCodeResp
	//ipCount := limitReq.GetReqCount(req.Phone)
	//if ipCount > 0 {
	//	return response.NewLoginReqFreqErrMsg(errID)
	//}
	//if err := limitReq.SetReqCount(req.Phone, ipCount+1); err != nil {
	//	logger.Errorf("limitReq.SetReqCount %v %s", err, req.Phone)
	//	return response.NewInternalServerErrMsg(errID)
	//}
	count, err := validate_code.CheckCount(req.Phone)
	if err != nil { // 禁止调用次数超出
		return nil, response.NewLoginReqMaxErrMsg(errID)
	}
	code := validate_code.GenValidateCode(6)
	logger.Infof("手机验证码 %s", code)
	timeout := 10 * time.Minute
	if err := validate_code.SaveValidateCode(req.Phone, code, count+1, timeout); err != nil {
		logger.Errorf("SendValidateCode 保存手机验证码到redis错误 :%v", err)
		return nil, response.NewInternalServerErrMsg(errID)
	}
	err = message.SendMsg(req.Phone, code)
	if err != nil {
		logger.Errorf("send code server err :%v", err)
		return nil, response.NewInternalServerErrMsg(errID)
	}
	resp.Code = code
	return &pb.ValidateCodeResp{}, nil
}

func (s *UserService) Register(ctx context.Context, req *pb.RegisterReq) (*emptypb.Empty, error) {
	vcode, err := validate_code.GetValidateCode(req.Phone)
	if err != nil {
		logger.Warnf("GetValidateCode 查验证码失败 %v", err)
		// return nil, response.NewValidateCodeExpireErrMsg("验证码查找失败")
	}
	if req.ValidateCode != vcode {
		logger.Warnf("GetValidateCode [%v]验证码错误 %v", req.Phone, err)
		// return nil, response.NewValidateCodeErrMsg(errID)
	}
	var userMaster *model.WqUserBase
	if req.InvitationCode != "" {
		//主人用户 发出邀请码的用户
		userMaster = s.dao.GetWqUserBaseByInCode(req.InvitationCode)
		if userMaster == nil {
			return nil, response.NewInvitationCodeErrMsg(errID)
		}
	}
	if user := s.dao.GetWqUserBaseByPhone(req.Phone); user != nil {
		return nil, response.NewPhoneHasRegisterErrMsg(errID)
	}
	dbclt := s.dao.BeginTran()
	user := model.NewWqUserBase(req.GetPhone(), req.GetPassword())
	if err := s.dao.CreateWqUserBase(dbclt, user); err != nil {
		dbclt.Rollback()
		return nil, response.NewUserCreateErrMsg(errID)
	}
	//创建钱包
	aaa, err := s.WalletSrv.CreateWallet(context.Background(), &walletpb.UidReq{UserId: user.UserID})
	if err != nil {
		logger.Warnf("用户%s 注册时,创建钱包失败 %v", user.UserID, err)
	}
	log.Fatalln(aaa, err)
	if userMaster == nil {
		dbclt.Commit()
		return nil, nil
	}
	if err := s.dao.CreateWqUserInvite(dbclt, userMaster.UserID, user.UserID); err != nil {
		dbclt.Rollback()
		return nil, response.NewUserCreateErrMsg(errID)
	}
	//给发出邀请的用户增加ifc
	if err := s.AddIfcBalance(userMaster.UserID, user.UserID, "", "register", 5.0); err != nil {
		logger.Warnf("注册时增加邀请码用户的ifc失败 uid %s userMasterId %s err %v", user.UserID, userMaster.UserID, err)
	}
	dbclt.Commit()
	//delete validate_code
	validate_code.DeleteValidateCode(req.Phone)
	return nil, nil
}

func (s *UserService) ResetPassword(ctx context.Context, req *pb.ChangePasswordReq) (*emptypb.Empty, error) {
	user := s.dao.GetWqUserBaseByUID(req.UserId)
	if user == nil {
		return nil, response.NewUserNotFoundErrMsg(errID)
	}
	userUpdateField := &model.WqUserBase{UserID: user.UserID, Password: bcrypt2.CryptPassword(req.Password),
		UpdatedAt: time.Now()}
	if err := s.dao.UpdateWqUserBaseByUID(userUpdateField); err != nil {
		return nil, response.NewUserSetPassErrMsg(errID)
	}
	return nil, nil
}

func (s *UserService) ForgetPassword(ctx context.Context, req *pb.ForgetPasswordReq) (*emptypb.Empty, error) {
	vcode, err := validate_code.GetValidateCode(req.Phone)
	if err != nil {
		logger.Warnf("GetValidateCode 查验证码失败 %v key %s", err)
		return nil, response.NewValidateCodeExpireErrMsg(errID)
	}
	if req.ValidateCode != vcode {
		return nil, response.NewValidateCodeErrMsg(errID)
	}
	user, err := s.GetUserInfoByPhone(req.Phone)
	if err != nil {
		return nil, response.NewUserNotFoundErrMsg(errID)
	}
	changePassReq := &pb.ChangePasswordReq{
		UserId:          user.UserID,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}
	if _, err := s.ResetPassword(context.Background(), changePassReq); err != nil {
		return nil, response.NewUserSetPassErrMsg(errID)
	}
	validate_code.DeleteValidateCode(req.Phone)
	return nil, nil
}
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*emptypb.Empty, error) {
	user := s.dao.GetWqUserBaseByUID(req.GetUserId())
	if user == nil {
		return nil, response.NewUserNotFoundErrMsg(errID)
	}
	switch {
	case req.Avatar != "":
		user.Avatar = req.Avatar
	case req.Name != "":
		user.Name = req.Name
	}
	user.UpdatedAt = time.Now()
	if err := s.dao.UpdateWqUserBaseByUID(user); err != nil {
		return nil, response.NewUserUpdateBaseErrMsg(errID)
	}
	return nil, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, req *pb.UserInfoReq) (*pb.LoginResp, error) {
	var resp = new(pb.LoginResp)
	user := s.dao.GetWqUserBaseByUID(req.UserId)
	if user == nil {
		return nil, response.NewUserNotFoundErrMsg(errID)
	}
	resp.UserId = user.UserID
	resp.Name = user.Name
	resp.Avatar = user.Avatar
	resp.Phone = user.Phone
	resp.InvitationCode = user.InvitationCode
	resp.LoginCount = int32(user.LoginCount)
	resp.LastLoginAt, _ = ptypes.TimestampProto(user.LastLogin)
	return resp, nil
}

func (s *UserService) GetAllUserInfo(ctx context.Context, req *emptypb.Empty) (*pb.AllUserInfoResp, error) {
	var resp = new(pb.AllUserInfoResp)
	users := s.dao.GetAllUsers()
	for _, user := range users {
		lastLogin, _ := ptypes.TimestampProto(user.LastLogin)
		resp.UserInfo = append(resp.UserInfo, &pb.LoginResp{
			UserId:         user.UserID,
			InvitationCode: user.InvitationCode,
			Name:           user.Name,
			Avatar:         user.Avatar,
			Phone:          user.Phone,
			LoginCount:     int32(user.LoginCount),
			LastLoginAt:    lastLogin,
		})
	}
	return resp, nil
}

func (s *UserService) GetUserMasterByInViteUser(ctx context.Context, req *pb.GetUserMasterReq) (*pb.UserMasterResp, error) {
	data := s.dao.GetUserMasterByInUserId(req.InviteUid)
	if data == nil {
		return nil, response.NewDataNotFound(errID, "没有找到邀请数据")
	}
	var resp = new(pb.UserMasterResp)
	resp.UserMasterId = data.UserID
	return resp, nil
}
