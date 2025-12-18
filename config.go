package main

type Contract byte

const WrappedQContract Contract = 0
const V3CoreContract Contract = 1
const MultiCallContract Contract = 2
const ProxyAdminContract Contract = 3
const TickLensContract Contract = 4
const NftDescriptorLibraryContract Contract = 5
const NftPositionDescriptorContract Contract = 6
const TransparentProxyContract Contract = 7
const NonfungiblePositionManagerContract Contract = 8
const V3MigratorContract Contract = 9
const V3StakerContract Contract = 10
const QuoterV2Contract Contract = 11
const V2CoreContract Contract = 12
const V3SwapRouterContract Contract = 13

type DeploymentSettings struct {
	ContractId   Contract
	ContractName string
	Gas          uint64
}

var (
	v2CoreContract                     = DeploymentSettings{ContractId: V2CoreContract, ContractName: "V2CoreContract", Gas: 6000000}
	wqContract                         = DeploymentSettings{ContractId: WrappedQContract, ContractName: "WrappedQContract", Gas: 6000000}
	v3CoreContract                     = DeploymentSettings{ContractId: V3CoreContract, ContractName: "V3CoreContract", Gas: 6000000}
	multiCallContract                  = DeploymentSettings{ContractId: MultiCallContract, ContractName: "MultiCallContract", Gas: 6000000}
	proxyAdminContract                 = DeploymentSettings{ContractId: ProxyAdminContract, ContractName: "ProxyAdminContract", Gas: 6000000}
	tickLensContract                   = DeploymentSettings{ContractId: TickLensContract, ContractName: "TickLensContract", Gas: 6000000}
	nftDescriptorLibraryContract       = DeploymentSettings{ContractId: NftDescriptorLibraryContract, ContractName: "NftDescriptorLibraryContract", Gas: 6000000}
	nftPositionDescriptorContract      = DeploymentSettings{ContractId: NftPositionDescriptorContract, ContractName: "NftPositionDescriptorContract", Gas: 6000000}
	transparentProxyContract           = DeploymentSettings{ContractId: TransparentProxyContract, ContractName: "TransparentProxyContract", Gas: 6000000}
	nonfungiblePositionManagerContract = DeploymentSettings{ContractId: NonfungiblePositionManagerContract, ContractName: "NonfungiblePositionManagerContract", Gas: 6000000}
	v3MigratorContract                 = DeploymentSettings{ContractId: V3MigratorContract, ContractName: "V3MigratorContract", Gas: 6000000}
	v3StakerContract                   = DeploymentSettings{ContractId: V3StakerContract, ContractName: "V3StakerContract", Gas: 6000000}
	quoterV2Contract                   = DeploymentSettings{ContractId: QuoterV2Contract, ContractName: "QuoterV2Contract", Gas: 6000000}
	v3SwapRouterContract               = DeploymentSettings{ContractId: V3SwapRouterContract, ContractName: "V3SwapRouterContract", Gas: 6000000}

	ContractMap = map[Contract]DeploymentSettings{
		WrappedQContract:                   wqContract,
		V2CoreContract:                     v2CoreContract,
		V3CoreContract:                     v3CoreContract,
		MultiCallContract:                  multiCallContract,
		ProxyAdminContract:                 proxyAdminContract,
		TickLensContract:                   tickLensContract,
		NftDescriptorLibraryContract:       nftDescriptorLibraryContract,
		NftPositionDescriptorContract:      nftPositionDescriptorContract,
		TransparentProxyContract:           transparentProxyContract,
		NonfungiblePositionManagerContract: nonfungiblePositionManagerContract,
		V3MigratorContract:                 v3MigratorContract,
		V3StakerContract:                   v3StakerContract,
		QuoterV2Contract:                   quoterV2Contract,
		V3SwapRouterContract:               v3SwapRouterContract,
	}
)
