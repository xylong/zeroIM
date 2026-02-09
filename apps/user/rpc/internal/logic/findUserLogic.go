package logic

import (
	"context"
	"gorm.io/gen"
	"zeroIM/apps/user/models"

	"github.com/jinzhu/copier"

	"zeroIM/apps/user/rpc/internal/svc"
	"zeroIM/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	var (
		users []*models.User
		err   error
	)

	u := l.svcCtx.Dao.User.WithContext(l.ctx).Debug() // IUserDo（负责查），debug查看 sql
	f := l.svcCtx.Dao.User                            // *UserDo（负责字段）
	conds := make([]gen.Condition, 0)

	if in.Phone != "" {
		conds = append(conds, f.Phone.Eq(in.Phone))
	}
	if in.Name != "" {
		conds = append(conds, f.Nickname.Like("%"+in.Name+"%"))
	}
	if len(in.Ids) > 0 {
		conds = append(conds, f.ID.In(in.Ids...))
	}
	if len(conds) > 0 {
		u = u.Where(conds...)
	}
	if users, err = u.Find(); err != nil {
		return nil, err
	}

	var resp []*user.UserEntity
	copier.Copy(&resp, users)

	return &user.FindUserResp{
		User: resp,
	}, nil
}
