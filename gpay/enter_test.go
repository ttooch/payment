package gpay

import (
	"testing"
	"fmt"
	"io/ioutil"
)

func TestEnter_Handle(t *testing.T) {
	en := new(Enter)

	pfxData,_ := ioutil.ReadFile("test.pfx")

	en.InitBaseConfig(&BaseConfig{
		AgentNo:"zl@cs.sh.cn",
		Key:"ZD7D0ZM7SC8CVV6Q",
		PfxData:pfxData,
		CertPassWord:"12345678",
	})

	picInfo := new (PicInfo)
	picInfo.LegalPersonIdFrontPic="/zl@cs.sh.cn/2018-07-25/1.png"  //法人身份正面照
	picInfo.LegalPersonIdOppositePic="/zl@cs.sh.cn/2018-07-25/1.png" //法人证件号反面
	picInfo.LegalPersonBanKCardPic="/zl@cs.sh.cn/2018-07-25/1.png" //法人银行卡图片
	picInfo.LegalPersonIdFrontPic="/zl@cs.sh.cn/2018-07-25/1.png"  //法人身份正面照
	picInfo.OperatorIdDeadlineFrontPic="/zl@cs.sh.cn/2018-07-25/1.png" //经办人身份证正面
	picInfo.OperatorIdDeadlineOppositePic="/zl@cs.sh.cn/2018-07-25/1.png" //经办人身份证反面
	picInfo.MerchantDoorHeadPic="/zl@cs.sh.cn/2018-07-25/1.png"  //商户门头照片
	picInfo.MerchantFrontPic="/zl@cs.sh.cn/2018-07-25/1.png" //商户门脸照片
	picInfo.MerchantInsidePic="/zl@cs.sh.cn/2018-07-25/1.png" //商户内饰照片
	//picInfo.NoSealAgreement="/zl@cs.sh.cn/2018-07-25/1.png" //协议 - 未盖章
	//picInfo.SealAgreement="/zl@cs.sh.cn/2018-07-25/1.png" //协议 - 已盖章
	picInfo.ContractConfirm="/zl@cs.sh.cn/2018-07-25/1.png" //合同确认图片

	settleInfo := new (SettleInfo)
	settleInfo.SettleAccountType="2"
	settleInfo.AccountName="周超"
	settleInfo.OpenBankName="中国民生银行股份有限公司"
	settleInfo.BranchBankName="上海杨浦支行"
	settleInfo.BankCardNo="6216910206331032"
	settleInfo.SettleCycle="T+24"

	feeInfoList := []*FeeInfo{
		{
			FeeRate:"3",
			FeeType:"0",
			ProductType:"20010002-1001-106",
		},
		{
			FeeRate:"3",
			FeeType:"0",
			ProductType:"20010002-1001-107",
		},
		{
			FeeRate:"3",
			FeeType:"0",
			ProductType:"20010002-1001-108",
		},
		{
			FeeRate:"3",
			FeeType:"0",
			ProductType:"20010002-1001-204",
		},
		{
			FeeRate:"3",
			FeeType:"0",
			ProductType:"20010002-1001-206",
		},
	}

	ret, err := en.Handle(&EnterConf{
		LoginNo:"15923126443",
		AccountType:"1",//开户类型
		MerchantName:"哟兔（信息）科技有限公司",//商户名称
		MerchantType:"9",//商户类型
		MccCode:"5331",//mcc码
		MerchantMail:"zhouc@ttouch.com.cn",//商户邮箱
		SendMail:"1",//是否发送邮件
		PhoneNo:"15923126443",//联系人手机号
		//LicenseProvince:"500000",
		//LicenseCity:"500100",
		//LicenseDistrict:"500103",

		LegalPersonName:"周超",//法人名称
		LegalPersonId:"500234199208102174",//法人身份证号
		LegalPersonIdDeadline:"2023-01-15",//法人身份证有效期
		OperatorName:"周超",//经办人姓名
		OperatorId:"500234199208102174",//经办身份证号
		OperatorIdDeadline:"2023-01-15",//经办人身份证有效期
		MerchantContactsName:"周超",//商户联系人
		MerchantAddress:"重庆渝中区大坪英利国际1号",//商户地址
		ServicePhone:"15923126443",//客服电话
		MerchantShortName:"呦点便利",//商户简称
		PicInfo:picInfo,//图片信息
		SettleInfo:settleInfo,//结算信息
		FeeInfo:feeInfoList,//费率信息
		FeeStartTime:"2018-10-01 21:20:00",//费率生效开始时间
		FeeEndTime:"2020-10-02 12:20:00",//费率生效结束时间
		NotifyUrl:"http://www.baidu.com",
	})

	fmt.Println(err)
	fmt.Println(ret)
}