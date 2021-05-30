package service

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/LotteWong/giotto-gateway-core/constants"
	"github.com/LotteWong/giotto-gateway-core/dao/mysql"
	"github.com/LotteWong/giotto-gateway-core/load_balance"
	"github.com/LotteWong/giotto-gateway-core/models/po"
	"github.com/e421083458/golang_common/lib"
)

var lbService *LbService

type LbService struct {
	LoadBalanceMap   map[string]*po.LoadBalanceDetail
	LoadBalanceSlice []*po.LoadBalanceDetail
	RWLock           sync.RWMutex

	loadBalanceOperator *mysql.LoadBalanceOperator
}

func NewLbService() *LbService {
	service := &LbService{
		LoadBalanceMap:      map[string]*po.LoadBalanceDetail{},
		LoadBalanceSlice:    []*po.LoadBalanceDetail{},
		RWLock:              sync.RWMutex{},
		loadBalanceOperator: mysql.NewLoadBalanceOperator(),
	}
	return service
}

func GetLbService() *LbService {
	if lbService == nil {
		lbService = NewLbService()
	}
	return lbService
}

func (s *LbService) GetLbWithConfForSvc(svc *po.ServiceDetail) (load_balance.LoadBalance, error) {
	// hit in cache, use cache data
	for _, lb := range s.LoadBalanceSlice {
		if lb.ServiceName == svc.Info.ServiceName {
			return lb.LoadBalancer, nil
		}
	}

	// miss in cache, new a load balance with config
	activeIps := svc.LoadBalance.GetEnabledIpList()
	weights := svc.LoadBalance.GetWeightList()
	ipWeightMap := map[string]int{}
	for idx, ip := range activeIps {
		weight, err := strconv.Atoi(weights[idx])
		if err != nil {
			return nil, err
		}
		ipWeightMap[ip] = weight
	}

	var scheme string
	switch svc.Info.ServiceType {
	case constants.ServiceTypeHttp:
		switch svc.HttpRule.NeedHttps {
		case constants.Enable:
			scheme = "https://"
		case constants.Disable:
			scheme = "http://"
		default:
			scheme = "http://"
		}
	case constants.ServiceTypeTcp:
		scheme = ""
	case constants.ServiceTypeGrpc:
		scheme = ""
	default:
		scheme = ""
	}
	format := fmt.Sprintf("%s%s", scheme, "%s")

	service := svc.Info.ServiceName
	tag := fmt.Sprintf("%d", svc.Info.Id)

	var conf load_balance.LoadBalanceConf
	if lib.GetBoolConf("base.consul.enable") {
		conf = load_balance.NewServerSvcDiscoveryLbConf(activeIps, ipWeightMap, format, service, tag)
	} else {
		conf = load_balance.NewClientSvcDiscoveryLbConf(activeIps, ipWeightMap, format)
	}
	lbr := load_balance.LoadBalanceWithConfFactory(load_balance.LbType(svc.LoadBalance.RoundType), conf)

	// miss in cache, write back to cache
	lb := &po.LoadBalanceDetail{
		LoadBalancer: lbr,
		ServiceName:  svc.Info.ServiceName,
	}

	s.LoadBalanceSlice = append(s.LoadBalanceSlice, lb)
	s.RWLock.Lock()
	defer s.RWLock.Unlock()
	s.LoadBalanceMap[svc.Info.ServiceName] = lb

	return lbr, nil
}
