# BlockEmulator 手册（扩展版）

[English (primary)](README.md) | **简体中文**

这是 HuangLab（中山大学软件工程学院）开发的 BlockEmulator 1.0 官方手册的扩展版本。官方英文手册和中文手册见[上游仓库](https://github.com/HuangLab-SYSU/block-emulator/tree/main)。

## 本项目的扩展（基于 BlockEmulator 1.0）

1. **共识协议扩展：**在 BlockEmulator 的区块链分片层集成 Kaspa DAG（有向无环图）共识协议，补充现有 PBFT、Relay 和 BrokerChain 的实验能力。
2. **代码重构与优化：**调整分片共识和交易处理核心模块，包括交易池管理、区块打包与跨分片交易验证，以适应 Kaspa DAG。
3. **实验验证：**原项目说明记录了集成协议的端到端实验验证，可输出吞吐量、延迟与交易确认率等指标。

## 上游手册简介

BlockEmulator 是面向区块链研究者的轻量级实验平台，用于验证新共识协议、分片机制和跨分片交易方法。它实现交易池、区块打包、共识协议和链上交易存储等核心功能，支持 PBFT，以及 Monoxide 的 Relay 和 BrokerChain 等跨分片机制。

平台使用 Go，支持配置与扩展，也可回放以太坊等区块链的历史交易并输出实验日志，供研究者分析吞吐量、交易确认延迟和交易池排队情况。

## 技术论文与引用

使用上游 BlockEmulator 作为实验工具时，请引用其论文。完整英文手册、参考文献和贡献者名单保留在 [英文主版](README.md)。

```bibtex
@article{huang2023blockemulator,
   title={BlockEmulator: An Emulator Enabling to Test Blockchain Sharding Protocols},
   author={Huang, Huawei and Ye, Guang and Chen, Qinde and Yin, Zhaokang and Luo, Xiaofei and Lin, Jianru and Li, Taotao and Yang, Qinglin and Zheng, Zibin},
   journal={arXiv preprint arXiv:2311.03612},
   year={2023}
}
```

## 上游参考文献与贡献者

以下保留上游手册的文献条目和贡献者信息，姓名与原文引文保持英文。

### Reference
[NSDI'2019 Monoxide] J. Wang and H. Wang, “Monoxide: Scale out blockchains with asynchronous consensus zones,” in 16th USENIX Symposium on Networked Systems Design and Implementation (NSDI 19), 2019, pp. 95–112.  

[INFOCOM'2022 BrokerChain] H. Huang, X. Peng, J. Zhan, S. Zhang, Y. Lin, Z. Zheng, and S. Guo, “Brokerchain: A cross-shard blockchain protocol for account/balance-based state sharding,” in Proc. of IEEE Conference on Computer Communications (INFOCOM’22), 2022, pp. 1–10.  

[SRDS'2022 Achieving] C. Li, H. Huang, Y. Zhao, X. Peng, R. Yang, Z. Zheng, and S. Guo, “Achieving scalability and load balance across blockchain shards for state sharding,” in Proc. of 2022 41st International Symposium on Reliable Distributed Systems (SRDS’22), 2022, pp. 284–294.  

### Contributors
- **Huawei Huang**: Associate professor at Sun Yat-sen University, IEEE senior member, Ph.D. from The University of Aizu (Japan). Research interests: blockchain, Web3, metaverse, distributed computing/protocols.  
- **Zibin Zheng (IEEE Fellow)**: Professor at Sun Yat-Sen University, Ph.D. from the Chinese University of Hong Kong (2012). Research interests: service computing, blockchain, cloud computing.  
- **Jianru Lin**: Research scientist engineer, technical mentor of HuangLab. Research interests: decentralized systems, smart contract languages, virtual machines.  
- **Taotao Li**: Postdoctoral researcher at Sun Yat-sen University, Ph.D. in cyberspace security (Institute of Information Engineering, CAS, 2022). Research interests: blockchain theory and application, cryptography.  
- **Xiaofei Luo**: Postdoctoral researcher at Sun Yat-sen University, Ph.D. from the University of Aizu (2023). Research interests: blockchain, payment channel networks, reinforcement learning.  
- **Qinglin Yang (IEEE Member)**: Research fellow, Ph.D. from the University of Aizu (2021). Research interests: intelligent edge cloud computing, federated learning privacy protection, Web3.  
- **Qinde Chen**: Ph.D. student at Sun Yat-sen University, joined HuangLab in 2022.  
- **Guang Ye**: Master’s student at Sun Yat-sen University, joined HuangLab in 2021 (undergraduate intern).  
- **Zhaokang Yin**: Master’s student at Sun Yat-sen University, joined HuangLab in 2022.  
- **Xiaowen Peng**: M.Sc. from Sun Yat-sen University (2023), joined HuangLab in 2021.  
- **Yue Lin**: Master’s student at Sun Yat-sen University, joined HuangLab in 2021.  
- **Miaoyong Xu**: Master’s student at Sun Yat-sen University, joined HuangLab in 2021.  
- **Junhao Wu**: Undergraduate intern at Sun Yat-sen University, joined HuangLab in 2022.  
