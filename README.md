### BlockEmulator Handbook (Extended Version)

**English (primary)** | [简体中文](README.zh-CN.md)

This is an extended version of the official handbook of BlockEmulator (version 1.0) developed by HuangLab (School of Software Engineering, Sun Yat-sen University). The original official handbook and Chinese-version handbook are available at the [original repository](https://github.com/HuangLab-SYSU/block-emulator/tree/main).  

#### My Contributions (Based on BlockEmulator v1.0)
1. **Consensus Protocol Extension**: Implemented the Kaspa DAG (Directed Acyclic Graph) consensus protocol at the blockchain sharding layer of BlockEmulator, expanding the platform’s support for DAG-based sharding mechanisms (complementary to the existing PBFT, Relay, and BrokerChain protocols).  
2. **Code Refactoring & Optimization**: Restructured core modules related to sharding consensus and transaction processing in the original codebase to adapt to the characteristics of Kaspa DAG, including transaction pool management, block packaging logic, and cross-shard transaction verification.  
3. **Experimental Validation**: Completed end-to-end experimental verification of the integrated Kaspa DAG consensus protocol, ensuring the modified BlockEmulator can stably run sharding experiments with DAG-based consensus and output valid experimental metrics (e.g., throughput, latency, transaction confirmation rate).  

---
below is from original handbook:

### Introduction
#### 1. Background of BlockEmulator
Initiated by HuangLab (a research group in the School of Software Engineering, Sun Yat-sen University, China), *BlockEmulator is a blockchain testbed that enables researchers to verify their proposed new protocols and mechanisms. It supports multiple consensus protocols, particularly the blockchain sharding mechanism.  

The main purpose of this testbed is to help users (researchers, students, etc.) quickly verify their own blockchain consensus protocols and blockchain-sharding protocols.  

BlockEmulator is designed as an experimental platform that adopts a lightweight system architecture. It simplifies the implementation of industrial-class blockchains since BlockEmulator only implements the core functions of a blockchain, including the transaction pool, block packaging, consensus protocols, and on-chain transaction storage. It also supports common consensus protocols, such as Practical Byzantine Fault Tolerance (PBFT).  

In particular, BlockEmulator offers the system-level design and implementation for blockchain-sharding mechanisms. For example, the cross-shard transaction mechanisms implemented by BlockEmulator include the following two representative solutions, i.e., i) Relay transaction mechanism proposed by Monoxide (NSDI'2019), and the BrokerChain protocol (INFOCOM'2022) PDF.  

BlockEmulator is oriented toward blockchain researchers. It provides a blockchain experimental platform for quickly implementing their own algorithms, protocols, and mechanisms. It also offers very helpful functions to help researchers collect experimental data, facilitating their plotting experimental figures.  

#### 2. Official Technical Paper & Citation
To provide an official handbook for BlockEmulator, we have written a technical paper titled "BlockEmulator: An Emulator Enabling to Test Blockchain Sharding Protocols" arXiv page. Please cite our paper if you use BlockEmulator as an experiment tool in your own paper.  

```bibtex
@article{huang2023blockemulator,
   title={BlockEmulator: An Emulator Enabling to Test Blockchain Sharding Protocols},
   author={Huang, Huawei and Ye, Guang and Chen, Qinde and Yin, Zhaokang and Luo, Xiaofei and Lin, Jianru and Li, Taotao and Yang, Qinglin and Zheng, Zibin},
   journal={arXiv preprint arXiv:2311.03612},
   year={2023}
}
```

#### 3. Related Work
The following papers from HuangLab's publications have adopted BlockEmulator as an experimental tool.  

- BrokerChain: A Cross-Shard Blockchain Protocol for Account/Balance-based State Sharding (published at INFOCOM 2022) [PDF]  
- Achieving Scalability and Load Balance across Blockchain Shards for State Sharding (published at SRDS 2022) [PDF]  
- tMPT: Reconfiguration across Blockchain Shards via Trimmed Merkle Patricia Trie (published at IWQoS 2023) [PDF]  
- MVCom: Scheduling Most Valuable Committees for the Large-Scale Sharded Blockchain (published at ICDCS 2021) [PDF]  

#### 4. Highlights of BlockEmulator
- **Lightweight**: BlockEmulator is a lightweight testbed platform for blockchain experiments.  
- **Fast Configuration**: BlockEmulator enables users to set up their environments quickly and supports remote deployment on the Cloud.  
- **Customization**: BlockEmulator is implemented using Go, which supports users' customization and modification.  
- **Easy to Conduct Experiments**: BlockEmulator supports the replay of historical transactions of mainstream blockchains (such as Ethereum). It can automatically yield experimental log files. Using those log files, researchers can interpret plenty of metrics such as system throughput, confirmation latency of transactions, the queueing of the transaction pool, etc. This function is very useful for researchers and students to facilitate their experimental data collection and plotting of experimental charts.  

#### 5. Getting Started to Use BlockEmulator
Quickly get started with BlockEmulator through the following document. Please refer to: BlockEmulator help document. More details of the user guidebook can be found at [Handbook of BlockEmulator (Eng.)](url) or [Handbook of BlockEmulator (Chn.)](url).  

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
