package util

const (
	// 巨量营销刷新Refresh Token https://open.oceanengine.com/labels/7/docs/1696710506097679
	RefreshToken = "https://api.oceanengine.com/open_api/oauth2/refresh_token/"

	// 代理商管理账户列表 https://open.oceanengine.com/labels/7/docs/1696710516003852
	APIAdvertiserSelect = "https://api.oceanengine.com/open_api/2/agent/advertiser/select/"

	// 投放账户信息查询 https://open.oceanengine.com/labels/7/docs/1809915654787136
	APIAdvertiserInfo = "https://api.oceanengine.com/open_api/2/agent/advertiser_info/query/"

	// 【代理返点】创建下载任务-返点(非素材)&激励 https://open.oceanengine.com/labels/7/docs/1812970751420483
	APIRebateDownloadCreateTask = "https://api.oceanengine.com/open_api/2/file/rebate/rebate_download/create_task/"

	// 【代理返点】数据报表-【代理】2025返点明点化素材数据-创建下载任务
	APICreateMingdianDownloadTask = "https://api.oceanengine.com/open_api/2/file/rebate/common_download/create_task/"

	// 【代理返点】数据报表-【代理】2025返点明点化素材数据-查询下载任务
	APIQueryMingdianDownloadTask = "https://api.oceanengine.com/open_api/2/file/rebate/common_download/get_download_task_list/"

	// 【代理返点】数据报表-【代理】2025返点明点化素材数据-下载任务结果
	APIDownloadMingdianData = "https://api.oceanengine.com/open_api/2/file/rebate/common_download/download_file/"

	// 代理商账户批量复制 https://open.oceanengine.com/labels/7/docs/1775097300794371
	APIAdvertisAccountCopy = "https://api.oceanengine.com/open_api/2/agent/advertiser/copy/"

	// 代理商消耗报表数据 https://open.oceanengine.com/labels/7/docs/1784979080790218
	APICostReportListQuery = "https://api.oceanengine.com/open_api/2/agent/adv/cost_report/list/query/"

	// 转账-获取最大可转余额（代理） https://open.oceanengine.com/labels/7/docs/1789754975045699
	APICanTransferBalanceGet = "https://api.oceanengine.com/open_api/v3.0/cg_transfer/query_can_transfer_balance/"

	// 转账-查询账户转账余额（代理） https://open.oceanengine.com/labels/7/docs/1789754859486282
	APITransferBalanceGet = "https://api.oceanengine.com/open_api/v3.0/cg_transfer/query_transfer_balance/"

	// 转账-发起转账（代理） https://open.oceanengine.com/labels/7/docs/1789755060558916
	APITransferCreate = "https://api.oceanengine.com/open_api/v3.0/cg_transfer/create_transfer/"

	// 转账-查询转账单信息（代理） https://open.oceanengine.com/labels/7/docs/1789755120706634
	APITransferDetailGet = "https://api.oceanengine.com/open_api/v3.0/cg_transfer/query_transfer_detail/"

	// 数据报表-【自定义报表】 https://open.oceanengine.com/labels/7/docs/1741387668314126
	APICustomReportData = "https://api.oceanengine.com/open_api/v3.0/report/custom/get/"

	//【代理商】上传自产首发素材至方舟（搬运治理） https://open.oceanengine.com/labels/7/docs/1792582253929536
	APIUploadVideo = "https://api.oceanengine.com/open_api/2/file/video/agent/"

	//代理商创建前测任务 https://open.oceanengine.com/labels/7/docs/1816970745502732
	APICreateDiagnosisTask = "https://api.oceanengine.com/open_api/2/diagnosis_task/agent/create/"

	//代理商轮询前测任务结果 https://open.oceanengine.com/labels/7/docs/1816970934411355
	APIDiagnosisTaskGet = "https://api.oceanengine.com/open_api/2/diagnosis_task/agent/get/"

	//营销素材预审 https://open.oceanengine.com/labels/7/docs/1823745911787786
	APIOpenMaterialAudit = "https://api.oceanengine.com/open_api/v3.0/security/open_material_audit/"

	//营销素材预审结果 https://open.oceanengine.com/labels/7/docs/1825271224960218
	APIAuditResultsGet = "https://api.oceanengine.com/open_api/v3.0/security/audit_results/"

	//获取单元列表 https://open.oceanengine.com/labels/7/docs/1741028841006095
	APIPromotionListGet = "https://api.oceanengine.com/open_api/v3.0/promotion/list/"

	//获取本地推单元列表 https://open.oceanengine.com/labels/37/docs/1808147672950851
	APILocalPromotionListGet = "https://api.oceanengine.com/open_api/v3.0/local/promotion/list/"

	//获取千川账户下计划列表（不含创意）https://open.oceanengine.com/labels/12/docs/1697467558690816
	APIQianchuanAdListGet = "https://ad.oceanengine.com/open_api/v1.0/qianchuan/ad/get/"

	//AD获取视频素材 https://open.oceanengine.com/labels/7/docs/1696710601820172
	APIFileVideoGet = "https://api.oceanengine.com/open_api/2/file/video/get/"

	//获取千川素材库视频 https://open.oceanengine.com/labels/12/docs/1739309912219663
	APIQianchuanVideoGet = "https://ad.oceanengine.com/open_api/v1.0/qianchuan/video/get/"

	//本地推素材库视频 https://open.oceanengine.com/labels/37/docs/1808613640441882
	APILocalFileVideoGet = "https://api.oceanengine.com/open_api/v3.0/local/file/video/get/"

	// 查询账户累计积分 https://open.oceanengine.com/labels/7/docs/1807434247414986
	APIAccountScoreGet = "https://api.oceanengine.com/open_api/v3.0/security/score_total/get/"

	// 查询违规积分明细 https://open.oceanengine.com/labels/7/docs/1807434338681868
	APIViolationScoreDetailGet = "https://api.oceanengine.com/open_api/v3.0/security/score_violation_event/get/"

	// 查看积分处置详情 https://open.oceanengine.com/labels/7/docs/1807434405738569
	APIScoreDisposalInfoGet = "https://api.oceanengine.com/open_api/v3.0/security/score_disposal_info/get/"

	// 【代理商】查询营销违规信息 https://open.oceanengine.com/labels/7/docs/1790052406659072
	APIRiskPromotionListGet = "https://api.oceanengine.com/open_api/2/agent/query/risk_promotion_list/"

	// 【代理激励】激励政策信息查询 https://open.oceanengine.com/labels/12/docs/1869409753087043
	APIIncentivePolicyBaseInfoGet = "https://api.oceanengine.com/open_api/2/file/incentive_policy_base_info/get/"

	// 批量更新营销素材启用状态 https://open.oceanengine.com/labels/7/docs/1755355780973568
	APIPromotionMaterialStatusUpdate = "https://api.oceanengine.com/open_api/v3.0/material/status/update/"

	// 删除营销下素材 https://open.oceanengine.com/labels/7/docs/1797183832412380
	APIPromotionMaterialDelete = "https://api.oceanengine.com/open_api/v3.0/promotion/material/delete/"

	// 批量更新营销启用状态 https://open.oceanengine.com/labels/7/docs/1741031308559364
	APIPromotionStatusUpdate = "https://api.oceanengine.com/open_api/v3.0/promotion/status/update/"

	// 批量更新营销出价 https://open.oceanengine.com/labels/7/docs/1741031138305028
	APIPromotionBidUpdate = "https://api.oceanengine.com/open_api/v3.0/promotion/bid/update/"

	// 批量修改深度出价 https://open.oceanengine.com/labels/7/docs/1755355890182159
	APIPromotionDeepBidUpdate = "https://api.oceanengine.com/open_api/v3.0/promotion/deepbid/update/"

	// 批量更新营销预算 https://open.oceanengine.com/labels/7/docs/1741030872454148
	APIPromotionBudgetUpdate = "https://api.oceanengine.com/open_api/v3.0/promotion/budget/update/"

	// 批量更新项目预算 https://open.oceanengine.com/labels/7/docs/1755353873798155
	APIProjectBudgetUpdate = "https://api.oceanengine.com/open_api/v3.0/project/budget/update/"

	// 批量更新项目状态 https://open.oceanengine.com/labels/7/docs/1740941413906432
	APIProjectStatusUpdate = "https://api.oceanengine.com/open_api/v3.0/project/status/update/"

	// 更新账户日预算 https://open.oceanengine.com/labels/7/docs/1696710531631116
	APIAdvertiserBudgetUpdate = "https://ad.oceanengine.com/open_api/2/advertiser/update/budget/"

	// 获取账户日预算 https://open.oceanengine.com/labels/7/docs/1696710531128335
	APIAdvertiserBudgetGet = "https://ad.oceanengine.com/open_api/2/advertiser/budget/get/"
)
